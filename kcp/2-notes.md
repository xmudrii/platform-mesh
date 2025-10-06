# 1 Cluster 1 Region. Front-proxy public, shards - public.

This document describes how to set up a production-grade kcp environment in a single
kubernetes cluster. Throughout this guide, we will install kcp and every dependency into
the namespace `kcp-vespucci` for deployments.

The main difference with this setup from `kcp-columbus` is that here we will use
external certificates for the kcp front-proxy and shards. One could add external certificates to the
`kcp-columbus` setup as well, but for simplicity we will use self-signed certificates there.


1. Create certificates for etcd

```yaml
kubectl create namespace kcp-vespucci
kubectl apply -f kcp/assets/kcp-vespucci/certificate-etcd.yaml
```

2. Create etcd cluster

```yaml
kubectl apply -f kcp/assets/kcp-vespucci/etcd-druid-root.yaml
kubectl apply -f kcp/assets/kcp-vespucci/etcd-druid-alpha.yaml
```

3. Create kcp certificates for system components
```yaml
kubectl apply -f kcp/assets/kcp-vespucci/certificate-kcp.yaml
```

4. Create kcp deployment

```yaml
kubectl apply -f kcp/assets/kcp-vespucci/kcp-root-shard.yaml
kubectl apply -f kcp/assets/kcp-vespucci/kcp-alpha-shard.yaml
kubectl apply -f kcp/assets/kcp-vespucci/kcp-front-proxy.yaml
```

At this point the certificate is missing for the front-proxy, so kcp will not start properly.
Get the Service LoadBalancer IP or domain name and create a certificate for the front-proxy and make sure DNS is 
pointing to the LoadBalancer IP.

```bash
kubectl get svc -n kcp-vespucci frontproxy-front-proxy
```

and check if certificate is issued:

```bash
kubectl get certificate -n kcp-vespucci root-frontproxy-server -o yaml
```

5. Create kubeconfig for admin user

```yaml
kubectl apply -f kcp/assets/kcp-vespucci/kubeconfig-kcp-admin.yaml
```

6. Extract kubeconfig and check access


```bash
kubectl get secret -n kcp-vespucci kcp-admin-frontproxy -o jsonpath='{.data.kubeconfig}' | base64 -d > kcp-admin-kubeconfig-vespucci.yaml
KUBECONFIG=kcp-admin-kubeconfig-vespucci.yaml kubectl get shards                                                                                                                                    10:20:41
NAME    REGION   URL                                                  EXTERNAL URL                                       AGE
alpha            https://alpha.vespucci.genericcontrolplane.io:6443   https://api.vespucci.genericcontrolplane.io:6443   2d19h
root             https://root.vespucci.genericcontrolplane.io:6443    https://api.vespucci.genericcontrolplane.io:6443   2d20h
```

The core difference with the `kcp-columbus` setup is that here we are using external certificates
for the front-proxy and shards, so you can access kcp without any additional configuration.
Internal certificates are also used for shards, but they are publicly accessible, so
controllers can run outside of the cluster as well, as long as they have the right configuration.

# Optional: Test OIDC authentication
Our system has OIDC enabled, so we can try using it with `kcp` and the `kubectl oidc-login` plugin.
```bash
brew install int128/kubelogin/kubelogin 
kubectl config set-credentials oidc \
  --exec-api-version=client.authentication.k8s.io/v1beta1 \
  --exec-command=kubectl \
  --exec-arg=oidc-login \
  --exec-arg=get-token \
  --exec-arg=--oidc-issuer-url="https://auth.genericcontrolplane.io" \
  --exec-arg=--oidc-client-id="platform-mesh" \
  --exec-arg=--oidc-extra-scope="email" \
  --exec-arg=--oidc-client-secret=Z2Fyc2lha2FsYmlzdmFuZGVuekWplCg==

kubectl config set-context --current --user=oidc
```

At this point, if you try to access kcp with `kubectl get shards`, you should be redirected to GitHub for authentication
but get an error because the user does not have permissions yet.
