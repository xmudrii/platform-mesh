# Scenario 3: Dual Front-Proxy with Edge Re-encryption (kcp-comer)

This scenario deploys kcp with dual front-proxy instances to handle edge re-encryption scenarios, typically used with CDN services like CloudFlare. This configuration is ideal for production environments requiring CDN integration where certificate authentication doesn't work with TLS termination at the edge.

## Architecture Overview

- **Certificate Strategy**: Two front-proxy instances - public (CloudFlare-managed) and internal (self-signed)
- **Access Pattern**: Public front-proxy for external access, internal for direct cluster communication  
- **Namespace**: `kcp-comer`
- **Use Case**: Production with CDN integration and edge re-encryption requirements

## Key Differences from Other Scenarios

- **Dual Front-Proxy**: Public instance for CDN, internal for direct access
- **Edge Re-encryption**: CDN terminates external TLS, re-encrypts to internal front-proxy
- **CloudFlare Integration**: Requires CloudFlare-specific configuration for certificate trust
- **Mixed Authentication**: OIDC recommended for external access due to certificate chain complexity

## Prerequisites

Before starting, ensure you have completed the [shared components setup](0-shared.md).

## CloudFlare Certificate Authority Setup

This scenario uses CloudFlare's edge certificates. Download the required CA certificate for extended trust:

**Important**: Verify this URL and certificate source with CloudFlare documentation before use in production.

```bash
curl -L -o google-we1.pem https://ssl-tools.net/certificates/108fbf794e18ec5347a414e4370cc4506c297ab2.pem
```

Create a secret with the CA certificate:

```bash
kubectl create namespace kcp-comer
kubectl create secret generic google-we1-ca --from-file=tls.crt=google-we1.pem -n kcp-comer
```
  

## Deployment Steps

### 1. Create ETCD Certificates

Configure certificates for etcd clusters:

```bash
kubectl apply -f kcp/assets/kcp-comer/certificate-etcd.yaml
```

### 2. Deploy ETCD Clusters

Deploy the etcd clusters for both root and alpha shards:

```bash
kubectl apply -f kcp/assets/kcp-comer/etcd-druid-root.yaml
kubectl apply -f kcp/assets/kcp-comer/etcd-druid-alpha.yaml
```

### 3. Configure KCP System Certificates

Create the certificate authorities and issuers for kcp components:

```bash
kubectl apply -f kcp/assets/kcp-comer/certificate-kcp.yaml
```

### 4. Deploy KCP Components

Deploy the kcp shards and both front-proxy instances:

```bash
kubectl apply -f kcp/assets/kcp-comer/kcp-root-shard.yaml
kubectl apply -f kcp/assets/kcp-comer/kcp-alpha-shard.yaml
kubectl apply -f kcp/assets/kcp-comer/kcp-front-proxy.yaml
kubectl apply -f kcp/assets/kcp-comer/kcp-front-proxy-internal.yaml
```

### 5. Configure DNS for All Components

All components need DNS configuration. Get the LoadBalancer IPs for each service:

1. **Public Front-Proxy**:
   ```bash
   kubectl get svc -n kcp-comer frontproxy-front-proxy
   ```

2. **Internal Front-Proxy and Shards**:
   ```bash
   kubectl get svc -n kcp-comer frontproxy-internal-front-proxy
   kubectl get svc -n kcp-comer root-shard
   kubectl get svc -n kcp-comer alpha-shard
   ```

Create DNS A records for all services pointing to their respective LoadBalancer IPs.

### 6. Create Admin Kubeconfigs

Generate kubeconfigs for both public and internal access:

```bash
kubectl apply -f kcp/assets/kcp-comer/kubeconfig-kcp-admin.yaml
kubectl apply -f kcp/assets/kcp-comer/kubeconfig-kcp-admin-internal.yaml
```

**Important**: In this scenario, shards use the private front-proxy for internal communication. If the public front-proxy is down or DNS isn't resolving, the public kubeconfig will not work.

### 7. Configure CloudFlare Settings

For the public front-proxy to work with CloudFlare, configure the following in your CloudFlare dashboard:

1. **Set `api.comer` to "Proxied"** (orange cloud icon)
2. **Add a Page Rule**: "Rewrite port to 6443" for the API domain
3. **Upload Custom CA** in the SSL/TLS tab so CloudFlare can trust the internal front-proxy certificate:

   ```bash
   kubectl get secret -n kcp-comer root-ca -o jsonpath='{.data.ca\.crt}' | base64 -d
   ```

### 8. Test Access

Extract the public kubeconfig:

```bash
kubectl get secret -n kcp-comer kcp-admin-frontproxy \
  -o jsonpath='{.data.kubeconfig}' | base64 -d > kcp-admin-kubeconfig-comer.yaml
```

**Certificate Authentication Limitation**: Due to CloudFlare's certificate re-encryption, the certificate presented to clients is CloudFlare's certificate, not the internal front-proxy certificate. Therefore, certificate-based authentication through the public front-proxy will not work directly.

**Recommended Approach**: Use OIDC authentication for external access through the public front-proxy.

## OIDC Authentication Setup

Due to the certificate re-encryption complexity, OIDC authentication is the recommended approach for external access:

### Install kubectl OIDC Plugin

```bash
# macOS
brew install int128/kubelogin/kubelogin

# For other platforms, see: https://github.com/int128/kubelogin
```

### Configure OIDC Credentials

```bash
kubectl config set-credentials oidc \
  --exec-api-version=client.authentication.k8s.io/v1beta1 \
  --exec-command=kubectl \
  --exec-arg=oidc-login \
  --exec-arg=get-token \
  --exec-arg=--oidc-issuer-url="https://auth.genericcontrolplane.io" \
  --exec-arg=--oidc-client-id="platform-mesh" \
  --exec-arg=--oidc-extra-scope="email" \
  --exec-arg=--oidc-client-secret=<YOUR_CLIENT_SECRET>

kubectl config set-context --current --user=oidc
```

### Verify Public Access

Test access via the public front-proxy with OIDC:

```bash
KUBECONFIG=kcp-admin-kubeconfig-comer.yaml kubectl get shards
```

Expected output:
```
NAME    REGION   URL                                               EXTERNAL URL                                   AGE
alpha            https://alpha.comer.genericcontrolplane.io:6443   https://api.comer.genericcontrolplane.io:443   16m
root             https://root.comer.genericcontrolplane.io:6443    https://api.comer.genericcontrolplane.io:443   16m
```

### Verify Internal Access

Extract and test the internal kubeconfig:

```bash
kubectl get secret -n kcp-comer kcp-admin-frontproxy-internal \
  -o jsonpath='{.data.kubeconfig}' | base64 -d > kcp-admin-kubeconfig-comer-internal.yaml

KUBECONFIG=kcp-admin-kubeconfig-comer-internal.yaml kubectl get shards
```

## Advantages of This Setup

- **CDN Integration**: Full CloudFlare integration with edge optimization
- **Dual Access Patterns**: Public for external users, internal for cluster services
- **High Security**: Edge re-encryption with certificate validation at multiple layers
- **Flexible Authentication**: OIDC for external access, certificates for internal access

## Use Cases

This scenario is particularly well-suited for:
- Production environments requiring CDN acceleration
- Organizations using CloudFlare for DDoS protection and traffic management  
- Scenarios where certificate authentication conflicts with edge TLS termination
- Mixed access patterns (external users + internal controllers)

## Next Steps

Your dual front-proxy kcp deployment is complete. You can now:
- Configure CloudFlare rules for traffic management and security
- Set up monitoring for both front-proxy instances
- Deploy controllers using the appropriate access method (internal vs external)

## Troubleshooting

**CloudFlare Certificate Trust**: Ensure the internal CA is uploaded to CloudFlare's SSL settings.

**Port Configuration**: Verify CloudFlare page rules correctly rewrite traffic to port 6443.

**DNS Propagation**: Allow time for DNS changes to propagate before testing access.
