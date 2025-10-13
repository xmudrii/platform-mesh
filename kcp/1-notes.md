# 1 Cluster 1 Region. Front-proxy public, shards - private.

This document describes how to set up a production-grade kcp environment in a single
kubernetes cluster. Throughout this guide, we will install kcp and every dependency into
the namespace `kcp-columbus` for deployments. 

In this guide we will NOT use external certificates, so everything will be self-signed
by the CA we configure. This means any external user will have to trust our CA to be able
to access kcp.

1. Create certificates for etcd

```yaml
kubectl create namespace kcp-columbus
kubectl apply -f kcp/assets/kcp-columbus/certificate-etcd.yaml
```

2. Create etcd cluster

```yaml
kubectl apply -f kcp/assets/kcp-columbus/etcd-druid-root.yaml
kubectl apply -f kcp/assets/kcp-columbus/etcd-druid-alpha.yaml
```

3. Create kcp certificates for system components
```yaml
kubectl apply -f kcp/assets/kcp-columbus/certificate-kcp.yaml
```

4. Create kcp deployment

```yaml
kubectl apply -f kcp/assets/kcp-columbus/kcp-root-shard.yaml
kubectl apply -f kcp/assets/kcp-columbus/kcp-alpha-shard.yaml
kubectl apply -f kcp/assets/kcp-columbus/kcp-front-proxy.yaml
```

At this point the certificate is missing for the front-proxy, so kcp will not start properly.
Get the Service LoadBalancer IP or domain name and create a certificate for the front-proxy and make sure DNS is 
pointing to the LoadBalancer IP.

```bash
kubectl get svc -n kcp-columbus frontproxy-front-proxy
```

and check if certificate is issued:

```bash
kubectl get certificate -n kcp-columbus root-frontproxy-server -o yaml
```

5. Create kubeconfig for admin user

```yaml
kubectl apply -f kcp/assets/kcp-columbus/kubeconfig-kcp-admin.yaml
```

6. Extract kubeconfig and check access

```bash
kubectl get secret -n kcp-columbus kcp-admin-frontproxy -o jsonpath='{.data.kubeconfig}' | base64 -d > kcp-admin-kubeconfig-columbus.yaml
KUBECONFIG=kcp-admin-kubeconfig-columbus.yaml kubectl get shards
NAME    REGION   URL                                                           EXTERNAL URL                                       AGE
alpha            https://alpha-shard-kcp.kcp-columbus.svc.cluster.local:6443   https://api.columbus.genericcontrolplane.io:6443   14m
root             https://root-kcp.kcp-columbus.svc.cluster.local:6443          https://api.columbus.genericcontrolplane.io:6443   14m
```


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
  --exec-arg=--oidc-redirect-url=http://127.0.0.1:8000/ \
  --exec-arg=--oidc-client-secret=Z2Fyc2lha2FsYmlzdmFuZGVuekWplCg==

kubectl config set-context --current --user=oidc
```

At this point, if you try to access kcp with `kubectl get shards`, you should be redirected to GitHub for authentication
but get an error because the user does not have permissions yet.
