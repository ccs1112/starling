# Substrate and cost

Two substrates, by design. The control-plane work — the reconcile loop, the
autoscaler, the gang scheduler — runs identically anywhere, so it runs on a
**free local `kind` cluster**. Only the AWS milestones need **EKS**, and only
in a concentrated burst.

## Local (most of the time, $0)

```sh
kind create cluster --name starling
kubectl config use-context kind-starling
```

`kind` (Kubernetes-in-Docker) provides a real API server, scheduler, and
controller-manager on a laptop. The controller, the autoscaler, and a second
scheduler all run here. Tear down with `kind delete cluster --name starling`.

## Cloud (the EKS burst)

The constraint that picks the approach: **the EKS control plane bills $0.10/hr
and cannot be paused — only deleted.** Left up for a month that is ~$72. So the
EKS milestones run in **one concentrated block**: create the cluster once, use
it daily that week, destroy it at the end.

### The shape

**EKS in `ap-northeast-1` (Tokyo), provisioned by Terraform, Karpenter-managed
spot nodes, no NAT gateway.** Develop on `kind`, validate here.

- **Terraform over `eksctl`**: Terraform is the IaC tool the target stack uses,
  and the state is itself part of the artifact. Provisioning is a milestone, not
  a prerequisite glossed over.
- **Karpenter spot nodes, scale-to-zero**: with no `InferenceService` applied,
  the node pool drains to zero and only the control plane bills. Spot is ~70%
  cheaper and the workloads are interruptible.
- **No NAT gateway (public subnets)**: a NAT gateway is a silent ~$32/mo plus
  data charges and buys nothing here. This is the most common surprise on an EKS
  bill — designed out from the start.
- **Tokyo (`ap-northeast-1`)**: lowest RTT for `kubectl`/`scp`.

### Spin-up

```sh
cd infra && terraform init && terraform apply       # ~15 min for the control plane
aws eks update-kubeconfig --name starling --region ap-northeast-1
kubectl get nodes                                   # Karpenter brings nodes on demand
```

Argo CD then syncs the `operator/` manifests from this repository.

## Lifecycle

- **`kind delete`** when done with local work — the containers hold RAM.
- **`terraform destroy`** at the end of an EKS session — an EKS control plane
  cannot be stopped, so destroy the cluster rather than let $0.10/hr run
  overnight. Recreate next session (~15 min), which is why the EKS milestones
  are batched into one block rather than spread thin.

```sh
cd infra && terraform destroy
```

## Cost ballpark

- Local `kind`: **$0**, most of the time.
- EKS control plane: $0.10/hr; one focused week, destroyed nightly, ≈ **~$4–17**.
- Karpenter spot nodes during sessions: **~$5**.
- ECR + gateway + misc: **~$5**.
- **Total well under $50.** The failure mode is not overspend but a forgotten
  NAT gateway or an un-destroyed control plane. Check Cost Explorer at the end
  of each EKS day.

## Dependency cooldown

Never install a package version published less than 7 days ago. When `go mod
tidy` or `terraform init` resolves versions, pin to the newest release that has
been public for at least 7 days; take security/CVE patches immediately. Applies
to Go modules, Terraform providers, and Helm charts (Karpenter, KEDA, Argo CD)
alike.

## Sanity checks (local)

```sh
kubectl cluster-info --context kind-starling
kubectl api-resources | grep -i inference     # after the CRD is applied
go version
```
