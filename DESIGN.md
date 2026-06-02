# Starling — design

Starling is a Kubernetes control plane for inference serving. It is built from
the substrate up rather than assembled from off-the-shelf operators: a custom
resource and a controller that reconciles model-serving workloads, and a
scheduler-framework plugin that admits a replica set as a gang. It runs on a
production-shaped stack — Terraform-provisioned EKS, Argo CD delivery, and
Karpenter-managed spot capacity.

The name: a starling murmuration moves as one coordinated unit. Gang
scheduling enforces the same all-or-nothing collective motion on a distributed
workload — every replica is placed together, or none is.

## Architecture

### `InferenceService` CRD + controller (`operator/`)

A custom resource declares a model image and a replica count. A controller
written directly on `client-go` — shared informer, rate-limited work queue,
level-triggered reconciliation — drives an owned Deployment and Service toward
that declaration, converges them when they drift, and reports readiness back
through the status subresource.

The controller deliberately does **not** use Kubebuilder / controller-runtime.
That framework hides the informer, the work queue, and the reconcile contract
behind a manager; Starling keeps them explicit, because the control loop is the
substance of the project rather than an implementation detail to delegate. The
resource shape mirrors KServe's `InferenceService`, kept minimal: the
control-plane behavior is identical whether the served pod runs vLLM or a
placeholder server, so the served image is intentionally not the focus.

### Gang scheduler (`scheduler/`)

An out-of-tree scheduler-framework plugin whose `Permit` stage holds every pod
of a replica set until all members are simultaneously schedulable, then
releases them together — and rejects the whole gang on timeout. A distributed
deployment therefore never half-starts: N−1 replicas never sit Running while a
peer waits for capacity that will not arrive. This is the mechanism behind
Kueue's gang admission and Volcano's gang scheduling.

### Delivery (`infra/`)

Terraform provisions an EKS cluster with no NAT gateway and a Karpenter
`NodePool` of spot instances that scales to zero when idle. Argo CD syncs the
operator manifests from this repository (GitOps), so cluster state tracks the
repository rather than imperative `kubectl` edits.

## Roadmap

| Milestone | Scope | Status |
|---|---|---|
| Reconcile core | CRD + `client-go` controller; owned Deployment/Service; level-triggered convergence; status subresource | in progress |
| EKS delivery | Terraform EKS (no NAT), Karpenter spot, ECR, Argo CD GitOps | planned |
| Queue-depth autoscaling | Scale on inflight concurrency, not CPU; scale-to-zero | planned |
| Gang scheduler | `Permit`-stage all-or-nothing admission of a replica set | planned |
| Debugging notes | Reconcile hot-loop, missed event, Permit deadlock — symptom to root cause | planned |
| Production mapping | Each component mapped to its production counterpart, honest about omissions | planned |
| GPU | g4dn spot, NVIDIA device plugin, a tensor-parallel gang | optional |

## Relation to production systems

Each component has a production counterpart; Starling implements the minimum of
each to make the mechanism legible end to end.

| Starling | Production counterpart |
|---|---|
| `InferenceService` CRD + controller | KServe |
| Gang admission | Kueue, Volcano |
| Queue-depth autoscaling | KEDA, Knative |
| Spot node provisioning | Karpenter |
| GitOps delivery | Argo CD |

## Substrate and cost

Development runs on a local `kind` cluster; the AWS milestones run on EKS in a
concentrated burst, spot-backed and torn down between sessions. See
[HARDWARE.md](HARDWARE.md).
