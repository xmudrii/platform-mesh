# 1 Cluster 1 Region. Two front-proxy public, shards - public.

This document describes how to set up a production-grade kcp environment in a single
kubernetes cluster. Throughout this guide, we will install kcp and every dependency into
the namespace `kcp-comer` for deployments.

The main difference with this setup from `kcp-columbus` is that here we will use
external certificates for the kcp front-proxy and shards. And two front-proxy instances,
one with CloudFlare managed certificate for public access and one with internal certificate
for private access.

In this scenario we will use CloudFlare Google WE1 Certificate issuer. Donwload CA 
for extended trust from: https://ssl-tools.net/subjects/b9bed5f1a61e40b24196b0c29e7e1a9d8bfcb520 

```bash
curl -L -o google-we1.pem https://ssl-tools.net/certificates/108fbf794e18ec5347a414e4370cc4506c297ab2.pem 
```

```yaml
kubectl create secret generic google-we1-ca --from-file=tls.crt=google-we1.pem -n kcp-comer
```
  

1. Create certificates for etcd

```yaml
kubectl create namespace kcp-comer
kubectl apply -f kcp/assets/kcp-comer/certificate-etcd.yaml
```

2. Create etcd cluster

```yaml
kubectl apply -f kcp/assets/kcp-comer/etcd-druid-root.yaml
kubectl apply -f kcp/assets/kcp-comer/etcd-druid-alpha.yaml
```

3. Create kcp certificates for system components
```yaml
kubectl apply -f kcp/assets/kcp-comer/certificate-kcp.yaml
```

4. Create kcp deployment

```yaml
kubectl apply -f kcp/assets/kcp-comer/kcp-root-shard.yaml
kubectl apply -f kcp/assets/kcp-comer/kcp-alpha-shard.yaml
kubectl apply -f kcp/assets/kcp-comer/kcp-front-proxy.yaml
kubectl apply -f kcp/assets/kcp-comer/kcp-front-proxy-internal.yaml
```

At this point the certificate is missing for the front-proxy, so kcp will not start properly.
Get the Service LoadBalancer IP or domain name and create a certificate for the front-proxy and make sure DNS is 
pointing to the LoadBalancer IP.

```bash
kubectl get svc -n kcp-comer frontproxy-front-proxy
```

Do the same for the internal front-proxy, shards:

```bash
kubectl get svc -n kcp-comer frontproxy-internal-front-proxy
kubectl get svc -n kcp-comer root-shard
kubectl get svc -n kcp-comer alpha-shard
```

5. Create kubeconfig for admin user

```yaml
kubectl apply -f kcp/assets/kcp-comer/kubeconfig-kcp-admin.yaml
kubectl apply -f kcp/assets/kcp-comer/kubeconfig-kcp-admin-internal.yaml
```

In this scenario, we informs shards to use private front-proxy for internal communication,
but due to this reason if public front-proxy is down or DNS is not resolving, kubeconfig will not work.

In this case we need to make sure few thigs are configured properly on CloudFlare side:
1. api.comer is set to proxy
2. There is rule deployed to `Rewrite port to 6443`
3. SSL tab has `Add custom certificate authority (CA)` uploaded so CloudFlare can trust the
   certificate used by the front-proxy.
``` 
kubectl get secret -n kcp-comer root-ca -o jsonpath='{.data.ca\.crt}' | base64 -d 
```


6. Extract kubeconfig and check access.

```bash
kubectl get secret -n kcp-comer kcp-admin-frontproxy -o jsonpath='{.data.kubeconfig}' | base64 -d > kcp-admin-kubeconfig-comer.yaml
```

Now this would NOT work. When using certificate re-encryption on CloudFlare side,
the certificate presented to the client is CloudFlare's certificate, not the one
we have configured on the front-proxy. So if one want to use external certificate 
we should use OIDC authentication instead.

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

Now via external shard:

```bash
KUBECONFIG=kcp-admin-kubeconfig-comer.yaml kubectl get shards
NAME    REGION   URL                                               EXTERNAL URL                                   AGE
alpha            https://alpha.comer.genericcontrolplane.io:6443   https://api.comer.genericcontrolplane.io:443   16m
root             https://root.comer.genericcontrolplane.io:6443    https://api.comer.genericcontrolplane.io:443   16m
```

Now via internal shard:
```bash
kubectl get secret -n kcp-comer kcp-admin-frontproxy-internal -o jsonpath='{.data.kubeconfig}' | base64 -d > kcp-admin-kubeconfig-comer-internal.yaml
KUBECONFIG=kcp-admin-kubeconfig-comer-internal.yaml kubectl get shards
NAME    REGION   URL                                                  EXTERNAL URL                                       AGE
```

The core difference with the `kcp-columbus` setup is that here we are using external certificates
for the front-proxy and shards, so you can access kcp without any additional configuration.
Internal certificates are also used for shards, but they are publicly accessible, so
controllers can run outside of the cluster as well, as long as they have the right configuration.
