# Starling

> A Kubernetes control plane for GPU inference — a custom operator and an
> all-or-nothing gang scheduler, built on raw client-go.

Built from the substrate up rather than assembled from off-the-shelf
operators: an `InferenceService` custom resource and a controller that
reconciles model-serving workloads directly on `client-go`, plus a
scheduler-framework plugin that admits a replica set as a gang. Deployed the
way production clusters run it — Terraform-provisioned EKS, Argo CD delivery,
Karpenter spot capacity.

A starling murmuration moves as one coordinated unit; gang scheduling enforces
the same all-or-nothing collective motion on a distributed workload.

## Components

- `operator/` — the `InferenceService` CRD and its controller (Go, `client-go`)
- `scheduler/` — the out-of-tree gang-scheduling plugin (Go)
- `infra/` — Terraform for the EKS cluster and the Argo CD manifests

See [DESIGN.md](DESIGN.md) for the architecture and roadmap, and
[HARDWARE.md](HARDWARE.md) for the deployment substrate and cost model.

## Quickstart (local)

```sh
kind create cluster --name starling
cd operator && go build ./...
../scripts/verify.sh        # build, load into kind, apply, assert
```

EKS provisioning and Argo CD delivery live under `infra/`.

## License

Apache-2.0. See [LICENSE](LICENSE).
