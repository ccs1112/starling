#!/usr/bin/env bash
# verify.sh — the local inner loop.
#
# Builds the controller from the current working tree, loads it into the
# kind cluster, applies the CRD and manifests, and asserts the change you
# just made actually works. Picks up working-tree state — no commit
# required. Asserts grow over time; early on, the build is the only check,
# since later steps are skipped when their inputs do not exist yet.
#
# Env overrides:
#   STARLING_CLUSTER  kind cluster name      (default starling)
#   STARLING_IMG      controller image:tag   (default starling/operator:dev)
#
# On EKS, Argo CD syncs the manifests from git instead; this is the fast
# local loop.

set -euo pipefail

CLUSTER="${STARLING_CLUSTER:-starling}"
IMG="${STARLING_IMG:-starling/operator:dev}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$REPO_ROOT"

# --- 1. Compile. The floor for every change: the code builds.
if [[ -f operator/go.mod ]]; then
  echo ">> go build ./... (operator)"
  ( cd operator && go build ./... )
  echo "ok: operator compiles"
else
  echo "verify.sh: no operator/go.mod yet" >&2
  exit 1
fi

# --- 2. Is a kind cluster up? If not, building is all we can check.
if ! kind get clusters 2>/dev/null | grep -qx "$CLUSTER"; then
  echo "note: kind cluster '$CLUSTER' not found — skipping cluster asserts."
  echo "      create it with: kind create cluster --name $CLUSTER"
  echo "=== verify.sh OK (build only) ==="
  exit 0
fi
kubectl config use-context "kind-$CLUSTER" >/dev/null

# --- 3. Apply the CRD if it exists.
if [[ -d operator/config/crd ]]; then
  echo ">> kubectl apply -f operator/config/crd/"
  kubectl apply -f operator/config/crd/
  if ! kubectl get crd inferenceservices.inference.noetics.dev >/dev/null 2>&1; then
    echo "FAIL: InferenceService CRD did not register" >&2
    exit 1
  fi
  echo "ok: InferenceService CRD registered"
fi

# --- 4. Build + load the controller image if there is a Dockerfile.
if [[ -f operator/Dockerfile ]]; then
  echo ">> docker build + kind load ($IMG)"
  docker build -t "$IMG" operator/
  kind load docker-image "$IMG" --name "$CLUSTER"
  echo "ok: image loaded into kind"
fi

# --- 5. Apply controller manifests + roll out if present.
if [[ -d operator/config/manager ]]; then
  echo ">> kubectl apply -f operator/config/manager/"
  kubectl apply -f operator/config/manager/
  if kubectl get deploy -n starling operator >/dev/null 2>&1; then
    kubectl rollout status -n starling deploy/operator --timeout=60s
    echo "ok: controller rolled out"
  fi
fi

echo "=== verify.sh OK ==="
