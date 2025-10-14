# kcp Production Quickstart Guide

kcp is a Kubernetes-like control plane that enables multi-tenant and
multi-cluster scenarios. It is designed to be run in production environments,
but setting it up requires several components to work together.

Every configuration described here uses 2 shards (root and alpha). This is the recommended
approach so you can adapt to multi-shard operations from the start.

## Certificate Management Scenarios

We cover three different certificate management approaches, each with specific use cases:

### Scenario 1: Self-Signed Certificates (kcp-columbus namespace)
**Best for**: Development, testing, or closed internal environments

- **Certificate Strategy**: All certificates are self-signed using an internal CA
- **Access Pattern**: Only front-proxy is publicly accessible, shards are private
- **Use Case**: Closed system where controllers run within the same cluster as kcp
- **Trade-off**: Requires adding the internal CA to client trust stores
- **Network**: Simple single-cluster deployment

### Scenario 2: External Certificates with In-Cluster Issuer (kcp-vespucci namespace)  
**Best for**: Production environments with external controller access

- **Certificate Strategy**: Front-proxy uses external certificates (Let's Encrypt), shards use internal certificates
- **Access Pattern**: Front-proxy and shards are publicly accessible
- **Use Case**: External integrations where controllers run outside the kcp cluster
- **Trade-off**: Mixed certificate trust model - external for front-proxy, internal for shards
- **Network**: Public access to all components with proper DNS setup

### Scenario 3: Dual Front-Proxy with Edge Re-encryption (kcp-comer namespace)
**Best for**: Production with CDN/edge services and certificate authority restrictions

- **Certificate Strategy**: Two front-proxy instances - public (CloudFlare-managed) and internal (self-signed)
- **Access Pattern**: Public front-proxy for external access, internal for direct cluster communication
- **Use Case**: When certificate authentication doesn't work with TLS termination at the edge
- **Trade-off**: More complex setup but better security and CDN integration
- **Network**: Edge re-encryption setup with CloudFlare or similar services

### Future Scenarios (Not Yet Implemented)

4. **Multi-Cluster & Multi-Region Public**: Distributed kcp across multiple Kubernetes clusters with public access
5. **Multi-Cluster & Multi-Region Private**: Distributed kcp with VPN/Service Mesh connectivity

## Prerequisites

Before setting up any scenario, you must first deploy the shared components that all configurations depend on.

## Setup Instructions

### Step 1: Install Shared Components
Follow the instructions in [shared components](0-shared.md) to set up:
- etcd-druid operator for database storage
- cert-manager for certificate management  
- kcp-operator for kcp lifecycle management
- OIDC provider (dex) for authentication
- DNS configuration

### Step 2: Choose and Deploy Your Scenario

#### Scenario 1: Self-Signed Certificates
Follow the instructions in [kcp-columbus (self-signed certificates)](kcp-columbus-self-signed.md) for a development or internal environment setup.

#### Scenario 2: External Certificates  
Follow the instructions in [kcp-vespucci (external certificates)](kcp-vespucci-external-certs.md) for production with external controller access.

#### Scenario 3: Dual Front-Proxy with Edge Re-encryption
Follow the instructions in [kcp-comer (dual front-proxy)](kcp-comer-dual-frontproxy.md) for production with CDN integration.

## Decision Matrix

| Requirement | Scenario 1<br>(Columbus) | Scenario 2<br>(Vespucci) | Scenario 3<br>(Comer) |
|-------------|---------------------------|---------------------------|-------------------------|
| **Environment** | Dev/Test/Internal | Production | Production + CDN |
| **External Controllers** | No | Yes | Yes |
| **Certificate Trust** | Manual CA trust | Automatic (Let's Encrypt) | Mixed (Edge + Internal) |
| **CDN Integration** | No | Limited | Full Support |
| **Setup Complexity** | Simple | Moderate | Complex |
| **Security Level** | Medium | High | Highest |

## Troubleshooting

### Common Errors

#### HTTP/2 Stream Error
**Error**: `(92) HTTP/2 stream 1 was not closed cleanly: INTERNAL_ERROR (err 2)`

**Cause**: You are trying to access a URL that is not served by the front-proxy or kcp.

**Solution**: Verify the URL is correct. This error is not related to TLS or certificates.

#### Certificate Trust Issues
**Error**: Certificate verification failures when using kubectl

**Solutions**:
- **Scenario 1 (Columbus)**: Add the internal CA to your system's trust store
- **Scenario 2 (Vespucci)**: Ensure Let's Encrypt root CA is in your kubeconfig
- **Scenario 3 (Comer)**: Use OIDC authentication instead of certificate-based auth

#### DNS Resolution Problems
**Error**: Cannot resolve kcp domain names

**Solution**: Verify DNS records point to the correct LoadBalancer IPs:
```bash
kubectl get svc -n <namespace> frontproxy-front-proxy
nslookup <your-kcp-domain>
```

#### OIDC Authentication Failures  
**Error**: Authentication redirects fail or return errors

**Solutions**:
1. Verify OIDC provider is running and accessible
2. Check that the client secret matches in both kcp and OIDC provider configurations
3. Ensure redirect URLs are properly configured