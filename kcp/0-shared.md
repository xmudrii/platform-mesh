# Shared components

For any scenarios to work, there are some shared components that need to be
deployed. These components are required for kcp to operate in any of the
deployment modes, like `etcd`, `cert-manager`, `kcp-operator`, `oidc`, etc.
This section describes how to deploy these shared components.


If you find files with names containing a `.template` suffix, it means that you need to
customize them for your own setup. Please copy them to a new file without the
`.template` suffix and edit them as needed.

For all certificates we are going to use `cert-manager` and `letsencrypt` as the issuer.
All DNS names are managed externally, by creating ServiceType LoadBalancer and
adding records to your DNS provider. One could automate this using `external-dns`,
but this is out of scope for these documents. 

Due to this we use the externalissuer configuration with api 
challenge. One could use an ingress controller and automate this.

```yaml
    solvers:
    - dns01:
        cloudflare:
          email: email@example.com
          apiKeySecretRef:
            name: cloudflare-api-key-secret
            key: api-key
```

1. ETCD

You can run any etcd compatible cluster. For production-grade setup, we
recommend using [etcd-druid](TBC)

```bash
helm install etcd-druid oci://europe-docker.pkg.dev/gardener-project/releases/charts/gardener/etcd-druid \
  --namespace etcd-druid \
  --create-namespace \
  --version v0.32.0
```

HACK: The etcd-druid chart does not install CRDs, so we need to install them manually:
see: https://github.com/gardener/etcd-druid/issues/1185

```bash
kubectl apply -f https://gist.githubusercontent.com/mjudeikis/94382cc47c6a8611a2a1b85e20ec8380/raw/953478d8605bcb3419de2a6958147a5fca932b20/etcdcopybackupstasks.druid.gardener.cloud
kubectl apply -f https://gist.githubusercontent.com/mjudeikis/94382cc47c6a8611a2a1b85e20ec8380/raw/953478d8605bcb3419de2a6958147a5fca932b20/etcds.druid.gardener.cloud
```

2. Cert-manager

We need cert-manager for both etcd and kcp. Install it using the following commands:

```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm upgrade \
  --install \
  --namespace cert-manager \
  --create-namespace \
  --version v1.18.2 \
  --set crds.enabled=true \
  --atomic \
  cert-manager jetstack/cert-manager
```

Add ClusterIssuer for letsencrypt:

```yaml
kubectl apply -f kcp/assets/cert-manager/cluster-issuer.yaml
```

We use direct integration with Cloudflare for DNS01 challenge. If you use different DNS provider, you need to change the issuer configuration accordingly.

Follow [cloud-flare-dns01](https://cert-manager.io/docs/configuration/acme/dns01/cloudflare/#api-keys) for more details.

Create secret with Cloudflare API key:

```yaml
kubectl apply -f kcp/assets/cert-manager/cloudflare-secret.yaml
```

3. Prepare etcd PKI

Before setting up the actual etcd cluster, we need to provision the certificates
that it should use. This is not done by Etcd Druid itself, but we are going to use Cert-Manager:

These issuers will live in the `cert-manager` namespace. One can change the namespace if 
needed by re-configuring cert-manager.

Check issue for more details: https://github.com/gardener/etcd-druid/issues/1187

```yaml
kubectl apply -f kcp/assets/etcd-druid/certificate-etcd-issuer.yaml
```

4. Install kcp-operator:

```bash
helm repo add kcp https://kcp-dev.github.io/helm-charts

# remove --image and --image.repository once we fix base url issue.
# Once this is merged and publiced: https://github.com/kcp-dev/kcp-operator/pull/101
helm upgrade --install kcp-operator kcp/kcp-operator --create-namespace --namespace kcp-operator --set image.tag=v36 --set image.repository=ghcr.io/mjudeikis/kcp-operator

# hack: update CRDs for BaseURL support.
kubectl apply -f https://raw.githubusercontent.com/kcp-dev/kcp-operator/2b2d13960d3c660d2a7ebaa74469a1e98146aec5/config/crd/bases/operator.kcp.io_frontproxies.yaml
kubectl apply -f https://raw.githubusercontent.com/kcp-dev/kcp-operator/2b2d13960d3c660d2a7ebaa74469a1e98146aec5/config/crd/bases/operator.kcp.io_kubeconfigs.yaml
kubectl apply -f https://raw.githubusercontent.com/kcp-dev/kcp-operator/2b2d13960d3c660d2a7ebaa74469a1e98146aec5/config/crd/bases/operator.kcp.io_rootshards.yaml
kubectl apply -f https://raw.githubusercontent.com/kcp-dev/kcp-operator/2b2d13960d3c660d2a7ebaa74469a1e98146aec5/config/crd/bases/operator.kcp.io_shards.yaml
```


5. OIDC Provider

If you have an existing OIDC provider, you can skip this step. We are going to use dex as the OIDC provider
for simplicity. And we will use postgres as the backend for dex. 

5.1. Install postgres operator

```bash
kubectl create namespace oidc
```

```bash
kubectl apply --server-side -f \
  https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.26/releases/cnpg-1.26.0.yaml
```

Create a cluster and database for dex:

```yaml
kubectl apply -f kcp/assets/oidc-dex/postgres-cluster.yaml
kubectl apply -f kcp/assets/oidc-dex/postgres-database.yaml
```

Create certificate request for dex to allow cert-manager to issue certificates:

```yaml
kubectl apply -f kcp/assets/oidc-dex/certificate-dns.yaml
```

Wait until certificate is issued and install dex itself:

```bash
helm repo add dex https://charts.dexidp.io

helm upgrade -i \
    --create-namespace \
    --namespace oidc \
    dex dex/dex \
    -f kcp/assets/oidc-dex/values.yaml  
```

5. DNS

For this to work we will require custom domains for `front-proxy` and 2 shards.
For simplicity we are using the most simple and most reliable option: 
   
  External DNS + Service type LoadBalancer.

We configure the external DNS provider with 3 domains:
```
# front-proxy domain
api.columbus.genericcontrolplane.io

At this point our cluster is ready to host kcp and we can proceed to the next step of
installing kcp itself.
