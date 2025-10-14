# Scenario 2: External Certificates with In-Cluster Issuer (kcp-vespucci)

This scenario deploys kcp in a single Kubernetes cluster using external certificates (Let's Encrypt) for the front-proxy and in-cluster certificate issuers for shards. This configuration is ideal for production environments where external controllers need access to kcp.

## Architecture Overview

- **Certificate Strategy**: Front-proxy uses external certificates (Let's Encrypt), shards use internal certificates
- **Access Pattern**: Front-proxy and shards are publicly accessible
- **Namespace**: `kcp-vespucci`
- **Use Case**: Production with external controller access and automatic certificate management

## Key Differences from Scenario 1 (Columbus)

- Uses Let's Encrypt for front-proxy certificates (automatic renewal)
- Shards are publicly accessible with internal certificates
- External controllers can run outside the cluster
- Mixed certificate trust model requiring CA bundle configuration

## Prerequisites

Before starting, ensure you have completed the [shared components setup](0-shared.md).

## Certificate Authority Setup

Since we're using Let's Encrypt, we need to ensure kubectl has the proper CA data in kubeconfigs. While browsers and operating systems typically trust Let's Encrypt by default, kubectl requires explicit CA configuration.

### Download Let's Encrypt Root CA

```bash
curl -L -o isrgrootx1.pem https://letsencrypt.org/certs/isrgrootx1.pem
```

### Create CA Secret

Create a secret containing the Let's Encrypt CA certificate:

```bash
kubectl create namespace kcp-vespucci
kubectl create secret generic letsencrypt-ca --from-file=tls.crt=isrgrootx1.pem -n kcp-vespucci
```


## Deployment Steps

### 1. Create ETCD Certificates

Configure certificates for etcd clusters:

```bash
kubectl apply -f kcp/assets/kcp-vespucci/certificate-etcd.yaml
```

### 2. Deploy ETCD Clusters

Deploy the etcd clusters for both root and alpha shards:

```bash
kubectl apply -f kcp/assets/kcp-vespucci/etcd-druid-root.yaml
kubectl apply -f kcp/assets/kcp-vespucci/etcd-druid-alpha.yaml
```

### 3. Configure KCP System Certificates

Create the certificate authorities and issuers for kcp components:

```bash
kubectl apply -f kcp/assets/kcp-vespucci/certificate-kcp.yaml
```

### 4. Deploy KCP Components

Deploy the kcp shards and front-proxy:

```bash
kubectl apply -f kcp/assets/kcp-vespucci/kcp-root-shard.yaml
kubectl apply -f kcp/assets/kcp-vespucci/kcp-alpha-shard.yaml
kubectl apply -f kcp/assets/kcp-vespucci/kcp-front-proxy.yaml
```

### 5. Configure DNS and Certificates

The front-proxy certificate depends on DNS resolution. Configure DNS and verify certificate issuance:

1. **Get the LoadBalancer IP**:
   ```bash
   kubectl get svc -n kcp-vespucci frontproxy-front-proxy
   ```

2. **Create DNS A record** pointing `api.vespucci.genericcontrolplane.io` to the LoadBalancer IP

3. **Verify certificate issuance** (this may take a few minutes):
   ```bash
   kubectl get certificate -n kcp-vespucci root-frontproxy-server -o yaml
   ```

### 6. Create Admin Kubeconfig

Generate the kubeconfig for administrative access:

```bash
kubectl apply -f kcp/assets/kcp-vespucci/kubeconfig-kcp-admin.yaml
```

### 7. Test Access

Extract the kubeconfig and verify kcp functionality:

```bash
kubectl get secret -n kcp-vespucci kcp-admin-frontproxy \
  -o jsonpath='{.data.kubeconfig}' | base64 -d > kcp-admin-kubeconfig-vespucci.yaml

KUBECONFIG=kcp-admin-kubeconfig-vespucci.yaml kubectl get shards
```

Expected output:
```
NAME    REGION   URL                                                  EXTERNAL URL                                       AGE
alpha            https://alpha.vespucci.genericcontrolplane.io:6443   https://api.vespucci.genericcontrolplane.io:6443   2d19h
root             https://root.vespucci.genericcontrolplane.io:6443    https://api.vespucci.genericcontrolplane.io:6443   2d20h
```

## Advantages of This Setup

- **Automatic Certificate Renewal**: Let's Encrypt certificates renew automatically
- **External Controller Support**: Controllers can run outside the cluster with public shard access
- **No Manual Certificate Trust**: Browsers and most tools automatically trust Let's Encrypt certificates
- **Mixed Trust Model**: Combines the security of external certificates with internal cluster communication

## Optional: OIDC Authentication

If you deployed the optional OIDC provider (Dex) in the shared components, you can configure OIDC authentication.

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

**Important**: Replace `<YOUR_CLIENT_SECRET>` with your actual OIDC client secret.

### Test OIDC Authentication

When you run `kubectl get shards`, you should be redirected to your OIDC provider for authentication. Users will need appropriate RBAC permissions to access resources.

## Next Steps

Your kcp deployment with external certificates is complete. You can now:
- Configure external controllers to connect to public shard endpoints  
- Set up RBAC for user and service account access
- Deploy multi-tenant workloads across shards

## Troubleshooting

**Certificate Issuance**: If Let's Encrypt certificates fail to issue, verify DNS propagation and that CloudFlare API credentials are correct.

**Mixed Certificate Trust**: Some clients may need CA bundle configuration to trust both external and internal certificates.
