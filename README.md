# Cosine Platform Golden Path

Independent proof-of-concept inspired by Cosine's public Senior Platform / Infrastructure role.

Cosine's role describes a platform engineer owning Kubernetes, CI/CD, observability, networking,
security, enterprise/on-prem deployment patterns, SLOs, and developer self-service. This repository
turns that into a small but concrete **platform golden path**.

It does **not** claim knowledge of Cosine's private infrastructure.

## What it demonstrates

A service team supplies one declarative workload contract. The platform layer evaluates whether the
service is ready to ship and generates the Kubernetes/rollout expectations around it.

The contract captures:

- workload identity
- environment and tenant isolation
- SLOs
- resource requests/limits
- autoscaling bounds
- disruption budgets
- health probes
- rollout strategy
- observability requirements
- secret source
- egress policy
- backup/restore expectations
- GPU requirements where relevant

The gate rejects unsafe production workloads before they reach the cluster.

## Core idea

```text
Product / Research team
        |
        v
 workload contract
        |
        v
 platformctl gate
   |      |      |
 SLOs   security rollout
   |      |      |
   +--- golden path ---> Helm / EKS / GitOps
                    |
                    +--> evidence.json
```

## Example

```bash
go test ./...
go run ./cmd/platformctl -contract examples/agent-api-prod.json
```

A safe example should return `allowed: true`.

## Why this maps to Cosine

The public role calls for Kubernetes/EKS ownership, Helm, infrastructure as code, progressive rollouts, rollback, observability, SLOs, incident readiness, least-privilege access, enterprise/on-prem patterns, and internal paved roads. This proof-of-concept focuses directly on those responsibilities.

## Repository layout

```text
cmd/platformctl/          admission/readiness CLI
internal/gate/            production-readiness policy engine
terraform/modules/eks/    reference EKS platform module
helm/cosine-service/      reusable Helm golden path
examples/                 workload contracts
docs/                     architecture and roadmap
.github/workflows/        CI + platform contract validation
```

## Production-readiness rules

For production workloads the gate requires CPU and memory requests/limits, at least two replicas for stateless workloads, readiness/liveness probes, PodDisruptionBudget, bounded autoscaling, explicit rollout and rollback, availability/latency SLOs, metrics/logs/traces, managed secrets, explicit egress policy, tenant/environment isolation, backup/restore for stateful workloads, and no public administrative ingress.

## Enterprise / on-prem pattern

A production evolution would support Cosine-managed AWS, customer-managed cloud, and air-gapped/on-prem deployments behind one workload contract.

## Reliability model

The platform should make failure boring: SLOs are part of the service contract, rollout policy is selected before deploy, canary steps have measurable promotion criteria, failed health/SLO checks block promotion, rollback is explicit and testable, and incident metadata can be generated from the same contract.

## Disclaimer

Independent engineering prototype based only on public Cosine role/product information.
