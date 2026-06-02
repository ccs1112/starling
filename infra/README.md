# infra

Terraform for the EKS cluster and the Argo CD manifests that deliver the
operator. Lands with the **EKS delivery** milestone (see `DESIGN.md`). Until
then, work runs on a local `kind` cluster — see `HARDWARE.md`.

Planned contents:

- `*.tf` — EKS cluster (no NAT gateway), Karpenter `NodePool` (spot,
  scale-to-zero), ECR repository, IRSA for the controller.
- `argocd/` — the Argo CD `Application` that syncs `operator/config/` from
  this repository.

State files and `*.tfvars` are gitignored. Never commit `terraform.tfstate`.
