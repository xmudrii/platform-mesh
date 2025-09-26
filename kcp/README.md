# kcp Production Quickstart Guide

This document describes how to setup a production-grade kcp environment. Throughout
this guide, we will install kcp and every dependency into the namespace `columbo`.

## Day 1

### 1 – Etcd Druid

Since kcp requires etcd, we first have to setup Etcd Druid to create and manage
etcd clusters. Installation involves cloning the etcd-druid repository and using
its scripting to set everything up.

```shell
# setup your kubeconfig for the target cluster
export KUBECONFIG=<path to your kubeconfig file>

# clone etcd-druid
git clone https://github.com/gardener/etcd-druid

# ensure etcd-druid's scripting can find the kubeconfig
mkdir -p etcd-druid/hack/kind
cp "$KUBECONFIG" etcd-druid/hack/kind/kubeconfig

# deploy etcd-druid
make -C etcd-druid deploy
```

### 2 – cert-manager

To secure the traffic between kcp and etcd, TLS certificates are used. These can
be provisioned using cert-manager, which is thankfully very easy to install:

```shell
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

### 3 – Prepare etcd PKI

Before setting up the actual etcd cluster, we need to provision the certificates
that it should use. This is not done by Etcd Druid itself, but there are Helm
charts that can automate the creation of the necessary `Certificate` objects.

Before proceeding, create a new `etcd-certs-values.yaml` to configure what
etcd clusters the certificates are meant for:

```yaml
# This should be identical to the name of the Etcd object we're creating later.
etcdName: etcd-test

# 3 is the default; if you plan on using more members in you etcd cluster,
# increase this (usually the max member count is 5 for etcd clusters).
replicas: 3
```

With this, we can now install the Helm chart and let it create the certificates:

```shell
helm repo add hajowieland https://charts.wieland.tech
helm repo update

helm upgrade \
  --install \
  --namespace columbo \
  --create-namespace \
  --values etcd-certs-values.yaml \
  my-etcd-certs hajowieland/etcd-druid-certs
```

You can check the progress by running `kubectl -n columbo get certs`. Once all
certificates are ready, you can proceed with the next step.

### 4 – Setup etcd

Now it's finally time to bootstrap our etcd cluster. This needs to be done for
every kcp shard (incl. the root shard) individually (i.e. kcp shards must not
use the same etcd cluster among them).

A simple example etcd with TLS can be provisioned by creating this YAML file
first and then applying it:

```yaml
apiVersion: druid.gardener.cloud/v1alpha1
kind: Etcd
metadata:
  name: etcd-test
  namespace: columbo
  labels:
    app: etcd-statefulset
    role: test
spec:
  replicas: 3

  etcd:
    metrics: basic
    defragmentationSchedule: "0 */24 * * *"
    resources:
      limits: { cpu: 500m, memory: 1Gi }
      requests: { cpu: 100m, memory: 200Mi }
    clientPort: 2379
    serverPort: 2380
    quota: 8Gi

    # configure the certificates we just created

    clientUrlTls:
      tlsCASecretRef:     { name: "etcd-test-etcd-ca-tls" }
      serverTLSSecretRef: { name: "etcd-test-etcd-server-tls" }
      clientTLSSecretRef: { name: "etcd-test-etcd-client-tls" }

    peerUrlTls:
      tlsCASecretRef:     { name: "etcd-test-etcd-peer-ca-tls" }
      serverTLSSecretRef: { name: "etcd-test-etcd-peer-tls" }
      clientTLSSecretRef: { name: "etcd-test-etcd-peer-tls" }

  backup:
    port: 8080
    fullSnapshotSchedule: "0 */24 * * *"
    resources:
      limits: { cpu: 200m, memory: 1Gi }
      requests: { cpu: 23m, memory: 128Mi }
    garbageCollectionPolicy: Exponential
    garbageCollectionPeriod: 43200s
    deltaSnapshotPeriod: 300s
    deltaSnapshotMemoryLimit: 1Gi
    compression:
      enabled: false
      policy: "gzip"
    leaderElection:
      reelectionPeriod: 5s
      etcdConnectionTimeout: 5s

    # configure the certificates we just created

    tls:
      tlsCASecretRef:     { name: "etcd-test-etcd-backup-restore-ca-tls" }
      serverTLSSecretRef: { name: "etcd-test-etcd-backup-restore-server-tls" }
      clientTLSSecretRef: { name: "etcd-test-etcd-backup-restore-client-tls" }

  sharedConfig:
    autoCompactionMode: periodic
    autoCompactionRetention: "30m"

  annotations:
    app: etcd-statefulset
    role: test
  labels:
    app: etcd-statefulset
    role: test
```

Apply it using `kubectl`:

```shell
kubectl apply --filename etcd.yaml
```

You can watch the progress of the rollout by running `watch kubectl -n columbo get etcd`.

### 5 – Setup kcp-operator

The kcp project offers a Kubernetes operator to manage kcp setups.

```shell
helm repo add kcp https://kcp-dev.github.io/helm-chartsc
helm repo update

helm upgrade \
  --install \
  --namespace kcp-operator \
  --create-namespace \
  kcp-operator kcp/kcp-operator
```

Installing the Helm chart will make the operator resources, like `RootShard`, `FrontProxy`
and `Kubeconfig` available in your cluster.

### 6 – Setup kcp

Now it's finally time to put all the pieces together and create a kcp environment.
To ensure no part of the system incorrectly assumes just a single kcp shard, it's
recommended to always start with a sharded setup.

First we create a new Issuer that will be responsible for signing all certificates
used by kcp:

```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: selfsigned
  namespace: columbo
spec:
  selfSigned: {}
```

Just `kubectl apply` the file above.

Next we can bootstrap our kcp root shard. `kubectl apply` this manifest:

```yaml
apiVersion: operator.kcp.io/v1alpha1
kind: RootShard
metadata:
  name: rooty
  namespace: columbo
spec:
  external:
    # replace the hostname with the external DNS name for your kcp instance
    hostname: example.operator.kcp.io
    port: 6443
  certificates:
    # this references the issuer created above
    issuerRef:
      group: cert-manager.io
      kind: Issuer
      name: selfsigned
  cache:
    embedded:
      # kcp comes with a cache server accessible to all shards,
      # in this case it is fine to enable the embedded instance
      enabled: true
  etcd:
    endpoints:
      - https://etcd-test-client.columbo.svc.cluster.local:2379
    tlsConfig:
      secretRef:
        name: etcd-test-etcd-client-tls
```

Every kcp setup needs at least one front-proxy that can also be applied the
same way:

```yaml
apiVersion: operator.kcp.io/v1alpha1
kind: FrontProxy
metadata:
  name: main-proxy
  namespace: collumbo
spec:
  rootShard:
    ref:
      name: rooty
  externalHostname: example.operator.kcp.io
```
