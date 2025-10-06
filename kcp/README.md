# kcp Production Quickstart Guide

kcp is a Kubernetes-like control plane that enables multi-tenant and
multi-cluster scenarios. It is designed to be run in production environments,
but setting it up requires several components to work together.

Every configuration we will describe here will have 2 shards. This is the recommended
way so you can adjust for multi-shard operations from the start, including the root shard.

In general we will cover a few different configuration modes:

1. 1 Cluster & 1 Region. Front-proxy public with self signed certificates, shards - private.

  This is the most simple configuration. In this mode, you have a single kcp cluster deployed in a single kubernetes cluster.
  Only the front-proxy is publicly accessible, and shards are private. This means it's a closed system, where only the front-proxy
  is accessible from outside. In this mode every controller should be running within the same cluster as kcp
  to be able to access shards. This has a caveat that the certificate of the front-proxy is internally signed, so you need to
  add it to your local trust store to be able to access it. 

2. 1 Cluster & 1 Region. With external certificates. Front-proxy public, shards - semi-private.

  In this mode, you have a single kcp cluster deployed in a single kubernetes cluster.
  The front-proxy and shards are all in the same cluster and publicly accessible.
  This means all components have public IPs and DNS name, but only the front-proxy has
  public certificates. The shards have internally signed certificates, but are publicly accessible.

  This allows for external integrations, like running controllers outside of the kcp
  deployment. And due to everything being public, it is easy to get started with.

3. (TBD) n Clusters & n Regions. Public.

  In this mode, you have a single kcp cluster deployed across multiple kubernetes clusters.
  The front-proxy has instances in different clusters for HA replication and shards are 
  in different clusters and publicly accessible. This requires distributed deployment configuration
  management, like GitOps to orchestrate the deployments, certificates, etc.

4. (TBD) n Clusters & n Regions. Private. (Using VPN or Direct Connect or Service Mesh)

  In this mode it's similar to the previous one, but only the front-proxy is publicly accessible.
  This means shards should be operating in private networks, and the front-proxy should be able to
  access them. This requires a more complex networking setup, but is more secure. 
  The caveat is that if external integrations are needed, they need to be able to access the private 
  network as well to be able to operate within the shards network.

This document describes how to set up a production-grade kcp environment. Throughout
this guide, we will install kcp and every dependency into the namespace `columbo`.

## Shared components

Follow the instructions in [shared components](0-shared.md) to set up shared components.

## Scenario 1: 1 Cluster & 1 Region. Front-proxy public, shards - private.

Follow the instructions in [1-cluster-1-region-frontproxy-public-shards-private](1-notes.md) to set up this scenario.

## Scenario 2: 1 Cluster & 1 Region. Public.

Follow the instructions in [1-cluster-1-region-public](2-notes.md) to set up this scenario.

# FAQ:

1. I got error `(92) HTTP/2 stream 1 was not closed cleanly: INTERNAL_ERROR (err 2)`
 
  This most likely means you are trying to access a URL which is not served by front-proxy & kcp.
  This error is a "red herring" and is not related to TLS or certificates. Check your URL.