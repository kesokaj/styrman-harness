#!/usr/bin/env bash
# setup-eval.sh — Setup and install the harness evaluation suite and GCP infrastructure.
#
# Usage:
#   ./setup-eval.sh [OPTIONS]
#
# Options:
#   --gcp                  Configure and provision GCP resources (Vertex AI & Cloud Storage)
#   --project <PROJECT_ID> GCP Project ID (default: active gcloud project)
#   --region <REGION>      GCP Region (default: active gcloud compute/region or europe-west1)
#   --bucket <BUCKET>      GCS Bucket name for eval run bundles (default: <PROJECT>-harness-eval)
#   --deploy-ui            Build and deploy the visualizer Web UI to Cloud Run
#   --install-bin          Install eval-harness to $HOME/.agents/bin/
#   --dry-run              Display the planned actions without making any changes
#   -h, --help             Show this help message
#
# Examples:
#   ./setup-eval.sh                      # Local build and install only
#   ./setup-eval.sh --gcp                # Local build + provision GCP Vertex AI & Cloud Storage
#   ./setup-eval.sh --gcp --deploy-ui    # Local build + provision GCP + deploy Web UI to Cloud Run

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Defaults
ENABLE_GCP=false
DEPLOY_UI=false
INSTALL_BIN=false
DRY_RUN=false
PROJECT_ID=""
REGION=""
BUCKET_NAME=""
VERTEX_MODEL="gemini-1.5-pro"

# Color helpers
info() { printf "\033[34m==>\033[0m %s\n" "$*"; }
success() { printf "\033[32m✔\033[0m %s\n" "$*"; }
warn() { printf "\033[33m⚠\033[0m %s\n" "$*"; }
error() { printf "\033[31m✖\033[0m %s\n" "$*" >&2; }
die() { error "$*"; exit 1; }

usage() {
  cat <<'EOF'
setup-eval.sh — Setup and install the harness evaluation suite and GCP infrastructure

Usage:
  ./setup-eval.sh [OPTIONS]

Options:
  --gcp                  Configure and provision GCP resources (Vertex AI & Cloud Storage)
  --project <PROJECT_ID> GCP Project ID (default: active gcloud project)
  --region <REGION>      GCP Region (default: active gcloud compute/region or europe-west1)
  --bucket <BUCKET>      GCS Bucket name for eval run bundles (default: <PROJECT>-harness-eval)
  --deploy-ui            Build and deploy the visualizer Web UI to Cloud Run
  --install-bin          Install eval-harness to $HOME/.agents/bin/
  --dry-run              Display the planned actions without making any changes
  -h, --help             Show this help message
EOF
  exit 0
}

# Parse CLI arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --gcp)
      ENABLE_GCP=true
      shift
      ;;
    --project)
      [[ -n "${2:-}" ]] || die "Flag --project requires a value"
      PROJECT_ID="$2"
      ENABLE_GCP=true
      shift 2
      ;;
    --region)
      [[ -n "${2:-}" ]] || die "Flag --region requires a value"
      REGION="$2"
      shift 2
      ;;
    --bucket)
      [[ -n "${2:-}" ]] || die "Flag --bucket requires a value"
      BUCKET_NAME="$2"
      ENABLE_GCP=true
      shift 2
      ;;
    --deploy-ui)
      DEPLOY_UI=true
      ENABLE_GCP=true
      shift
      ;;
    --install-bin)
      INSTALL_BIN=true
      shift
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    -h|--help)
      usage
      ;;
    *)
      die "Unknown argument: $1 (see --help)"
      ;;
  esac
done

info "Checking local development prerequisites..."

# 1. Verify Go toolchain
if ! command -v go >/dev/null 2>&1; then
  die "Go is not installed or not on PATH. Please install Go (>= 1.22)."
fi
GO_VERSION=$(go version | awk '{print $3}')
success "Found Go toolchain ($GO_VERSION)"

# 2. Check Git
if ! command -v git >/dev/null 2>&1; then
  die "git is not installed or not on PATH."
fi
success "Found git ($(git --version))"

# 3. Check Agent CLIs (advisory only)
CLAUDE_AVAILABLE=false
AGY_AVAILABLE=false
if command -v claude >/dev/null 2>&1; then
  CLAUDE_AVAILABLE=true
  success "Found Claude Code CLI (claude)"
else
  warn "claude CLI not found on PATH (will be required for Claude eval runs)"
fi

if command -v agy >/dev/null 2>&1; then
  AGY_AVAILABLE=true
  success "Found Antigravity CLI (agy)"
else
  warn "agy CLI not found on PATH (will be required for Antigravity eval runs)"
fi

# 4. Build local eval binary
info "Building eval-harness binary..."
BIN_DIR="${SCRIPT_DIR}/bin"
mkdir -p "$BIN_DIR"
if [ "$DRY_RUN" = true ]; then
  info "[dry-run] Would build: go build -o ${BIN_DIR}/eval-harness ./cmd/eval"
else
  go build -o "${BIN_DIR}/eval-harness" ./cmd/eval
  success "Built ${BIN_DIR}/eval-harness"
fi

if [ "$INSTALL_BIN" = true ]; then
  AGENT_BIN="${HOME}/.agents/bin"
  info "Installing binary to ${AGENT_BIN}/eval-harness..."
  if [ "$DRY_RUN" = true ]; then
    info "[dry-run] Would copy: cp ${BIN_DIR}/eval-harness ${AGENT_BIN}/eval-harness"
  else
    mkdir -p "$AGENT_BIN"
    cp "${BIN_DIR}/eval-harness" "${AGENT_BIN}/eval-harness"
    chmod +x "${AGENT_BIN}/eval-harness"
    success "Installed to ${AGENT_BIN}/eval-harness"
  fi
fi

# 5. GCP Configuration and Infrastructure Provisioning
if [ "$ENABLE_GCP" = true ]; then
  info "Configuring Google Cloud Platform (GCP) resources..."

  if ! command -v gcloud >/dev/null 2>&1; then
    die "gcloud CLI is required for GCP setup but is not installed or not on PATH."
  fi

  # Resolve Project ID
  if [ -z "$PROJECT_ID" ]; then
    PROJECT_ID=$(gcloud config get-value project 2>/dev/null || true)
    if [ -z "$PROJECT_ID" ] || [ "$PROJECT_ID" = "(unset)" ]; then
      # Fallback to reading gcloud config files directly
      ACTIVE_CONF="default"
      if [ -f "${HOME}/.config/gcloud/active_config" ]; then
        ACTIVE_CONF=$(tr -d '[:space:]' < "${HOME}/.config/gcloud/active_config")
      fi
      CONF_FILE="${HOME}/.config/gcloud/configurations/config_${ACTIVE_CONF}"
      if [ -f "$CONF_FILE" ]; then
        PROJECT_ID=$(grep -E '^\s*project\s*=' "$CONF_FILE" | head -n 1 | awk -F'=' '{print $2}' | tr -d '[:space:]' || true)
      fi
    fi
    if [ -z "$PROJECT_ID" ] || [ "$PROJECT_ID" = "(unset)" ]; then
      die "No active GCP project configured. Pass --project <ID> or run: gcloud config set project <ID>"
    fi
  fi
  success "GCP Project: ${PROJECT_ID}"

  # Resolve Region
  if [ -z "$REGION" ]; then
    REGION=$(gcloud config get-value compute/region 2>/dev/null || true)
    if [ -z "$REGION" ] || [ "$REGION" = "(unset)" ]; then
      REGION="europe-west1"
    fi
  fi
  success "GCP Region: ${REGION}"

  # Resolve Bucket Name
  if [ -z "$BUCKET_NAME" ]; then
    BUCKET_NAME="${PROJECT_ID}-harness-eval"
  fi
  success "GCS Bucket: gs://${BUCKET_NAME}"

  # Check Application Default Credentials (ADC)
  info "Verifying Application Default Credentials (ADC)..."
  if ! gcloud auth application-default print-access-token >/dev/null 2>&1; then
    warn "ADC is not set up or has expired. Run 'gcloud auth application-default login' so the Go eval tool can authenticate with Vertex AI."
  else
    success "Application Default Credentials valid."
  fi

  # Enable Required APIs
  APIS=(
    "aiplatform.googleapis.com"
    "storage.googleapis.com"
  )
  if [ "$DEPLOY_UI" = true ]; then
    APIS+=(
      "run.googleapis.com"
      "artifactregistry.googleapis.com"
    )
  fi

  info "Enabling required GCP APIs: ${APIS[*]}..."
  if [ "$DRY_RUN" = true ]; then
    info "[dry-run] Would enable APIs: ${APIS[*]}"
  else
    gcloud services enable "${APIS[@]}" --project="${PROJECT_ID}"
    success "GCP APIs enabled."
  fi

  # Create GCS Bucket if not already existing
  info "Checking GCS Bucket: gs://${BUCKET_NAME}..."
  if [ "$DRY_RUN" = true ]; then
    info "[dry-run] Would ensure bucket exists: gs://${BUCKET_NAME}"
  else
    if gcloud storage buckets describe "gs://${BUCKET_NAME}" >/dev/null 2>&1; then
      success "GCS Bucket gs://${BUCKET_NAME} already exists."
    else
      info "Creating GCS Bucket gs://${BUCKET_NAME} in ${REGION}..."
      gcloud storage buckets create "gs://${BUCKET_NAME}" \
        --project="${PROJECT_ID}" \
        --location="${REGION}" \
        --uniform-bucket-level-access
      success "Created GCS Bucket gs://${BUCKET_NAME}."
    fi
  fi

  # Write configuration file
  CONFIG_DIR="${SCRIPT_DIR}/.agents"
  mkdir -p "$CONFIG_DIR"
  CONFIG_FILE="${CONFIG_DIR}/eval-config.json"
  info "Writing configuration to ${CONFIG_FILE}..."
  if [ "$DRY_RUN" = true ]; then
    info "[dry-run] Would write config: ${CONFIG_FILE}"
  else
    cat > "$CONFIG_FILE" <<EOF
{
  "projectId": "${PROJECT_ID}",
  "region": "${REGION}",
  "gcsBucket": "${BUCKET_NAME}",
  "vertexModel": "${VERTEX_MODEL}",
  "localLedger": ".agents/docs/evals/eval-ledger.jsonl"
}
EOF
    success "Saved configuration to ${CONFIG_FILE}"
  fi
fi

# 6. Deploy Web UI to Cloud Run (if requested)
if [ "$DEPLOY_UI" = true ]; then
  info "Deploying Web Visualizer UI to Google Cloud Run..."

  if [ ! -d "${SCRIPT_DIR}/web" ]; then
    warn "web/ directory does not exist yet. Run this with --deploy-ui after the web UI has been initialized."
  else
    SERVICE_NAME="harness-eval-ui"
    REPO_NAME="harness-eval"
    IMAGE_TAG=$(git rev-parse --short HEAD 2>/dev/null || echo "v1")
    IMAGE_PATH="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO_NAME}/${SERVICE_NAME}:${IMAGE_TAG}"

    info "Ensuring Artifact Registry repository exists: ${REPO_NAME}..."
    if [ "$DRY_RUN" = true ]; then
      info "[dry-run] Would create repository and deploy Cloud Run service"
    else
      if ! gcloud artifacts repositories describe "$REPO_NAME" --location="$REGION" --project="$PROJECT_ID" >/dev/null 2>&1; then
        gcloud artifacts repositories create "$REPO_NAME" \
          --repository-format=docker \
          --location="$REGION" \
          --project="$PROJECT_ID" \
          --description="Harness Eval Visualizer Images"
      fi

      info "Configuring Docker auth for ${REGION}-docker.pkg.dev..."
      gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet

      info "Building and pushing container image: ${IMAGE_PATH}..."
      # Building multi-platform or AMD64 for Cloud Run
      docker buildx build --platform linux/amd64 -t "${IMAGE_PATH}" --push -f "${SCRIPT_DIR}/web/Dockerfile" "${SCRIPT_DIR}/web"

      info "Deploying service to Cloud Run: ${SERVICE_NAME}..."
      gcloud run deploy "$SERVICE_NAME" \
        --image "$IMAGE_PATH" \
        --region "$REGION" \
        --project "$PROJECT_ID" \
        --platform managed \
        --port 8080 \
        --allow-unauthenticated \
        --set-env-vars "GCS_BUCKET=${BUCKET_NAME},GCP_PROJECT=${PROJECT_ID}"

      URL=$(gcloud run services describe "$SERVICE_NAME" --region "$REGION" --project "$PROJECT_ID" --format 'value(status.url)')
      success "Web Visualizer deployed successfully: ${URL}"
    fi
  fi
fi

# 7. Final Summary
echo ""
printf "\033[32m==================================================\033[0m\n"
printf "\033[32m  Harness Eval Setup Complete!                    \033[0m\n"
printf "\033[32m==================================================\033[0m\n"
echo ""
echo "Next steps:"
echo "  1. List available evaluation tasks:"
echo "     ./bin/eval-harness list"
echo ""
echo "  2. Run an evaluation task (e.g. greenfield Go CLI):"
echo "     ./bin/eval-harness run --task gf-go-cli"
echo ""
echo "  3. Run with baseline A/B comparison (Harnessed vs Vanilla):"
echo "     ./bin/eval-harness run --task gf-go-cli --baseline"
echo ""
if [ "$ENABLE_GCP" = true ]; then
  echo "  4. GCP Integration status:"
  echo "     Project:    ${PROJECT_ID}"
  echo "     GCS Bucket: gs://${BUCKET_NAME}"
  echo "     LLM Judge:  Vertex AI (${VERTEX_MODEL})"
  echo ""
fi
echo "  5. To launch the Web UI locally:"
echo "     cd web && npm install && npm run dev"
echo ""
