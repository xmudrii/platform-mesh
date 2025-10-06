# TODO

## Documentation Status

### KCP Production Setup
- ✅ **Scenario 1**: 1 Cluster & 1 Region. Front-proxy public, shards private (kcp-columbus) - [1-notes.md](1-notes.md)
- ✅ **Scenario 2**: 1 Cluster & 1 Region. Public with certificates (kcp-vespucci) - [2-notes.md](2-notes.md)
- ✅ **Shared components**: Common setup requirements - [0-shared.md](0-shared.md)
- ✅ **Main guide**: Overview and quickstart - [README.md](README.md)

### Related Platform Mesh Documentation
- ✅ **Account Model ADR**: Architecture decision for account model implementation - [../adr/adr-01-account-model.md](../adr/adr-01-account-model.md)
- ✅ **KCP & MFP Integration POC**: Proof of concept for micro frontend platform integration - [../poc/kcp-poc-mfp.md](../poc/kcp-poc-mfp.md)

## Outstanding Issues
1. Etcd Druid: https://github.com/gardener/etcd-druid/issues/1176
2. kcp-operator: https://github.com/kcp-dev/kcp-operator/issues/99 
3. kcp-operator - good-to-have: https://github.com/kcp-dev/kcp-operator/issues/100 
4. kcp (good to have): https://github.com/kcp-dev/kcp/issues/3613
5. kcp-operator CA bundle support: https://github.com/kcp-dev/kcp-operator/issues/103
6. How to run when OIDC custom CA needs to be provided.
7. http://github.com/kcp-dev/kcp-operator/issues/67 Add kubeconfig status.

## Future Scenarios (TBD)
- **Scenario 3**: n Clusters & n Regions. Public
- **Scenario 4**: n Clusters & n Regions. Private (VPN/Direct Connect/Service Mesh)
- **Scenario 5**: n Clusters & n Regions. Filtered & private