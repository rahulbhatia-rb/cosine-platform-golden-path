# Architecture

## Cluster topology

A production fleet should separate system workloads, observability/control-plane add-ons, application workloads, and optional GPU workloads.

Isolation can be achieved with node groups, taints/tolerations, namespaces, network policy, and separate clusters/accounts where the blast-radius or customer boundary warrants it.

## GitOps

Recommended flow:

```text
PR -> tests -> image build -> scan/sign -> contract gate
   -> environment repo -> Argo CD -> progressive rollout
   -> SLO evaluation -> promote / rollback
```

## SLOs

Service SLOs should be declarative inputs to dashboards, alerts, rollout promotion, and incident review. Burn-rate alerting is preferable to static threshold spam.

## Enterprise deployments

The control plane should distinguish managed SaaS, customer cloud, and disconnected/on-prem. The same application release should produce different delivery artifacts without requiring a fork of the application architecture.
