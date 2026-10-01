# A small Cloud Run v2 root: runtime SA, Artifact Registry, one write-only secret, the service.
# Every block is the provider documentation's minimal form; references/resources-gcp.md has the URLs.

resource "google_project_service" "apis" {
  for_each           = toset(["run.googleapis.com", "artifactregistry.googleapis.com", "secretmanager.googleapis.com"])
  service            = each.value
  disable_on_destroy = false
}

resource "google_service_account" "run" {
  account_id   = "${var.app}-run"
  display_name = "${var.app} Cloud Run runtime"
}

resource "google_artifact_registry_repository" "app" {
  location               = var.region
  repository_id          = var.repo
  format                 = "DOCKER"
  cleanup_policy_dry_run = false
  docker_config {
    immutable_tags = true # the git-SHA tag can never be moved to another image
  }
  cleanup_policies {
    id     = "delete-untagged"
    action = "DELETE"
    condition {
      tag_state = "UNTAGGED"
    }
  }
  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret" "app" {
  secret_id = "${var.app}-app-secret"
  replication {
    auto {}
  }
  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "app" {
  secret                 = google_secret_manager_secret.app.id
  secret_data_wo         = var.app_secret
  secret_data_wo_version = "1"
}

resource "google_secret_manager_secret_iam_member" "run_reads_app" {
  project   = google_secret_manager_secret.app.project
  secret_id = google_secret_manager_secret.app.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.run.member
}

resource "google_cloud_run_v2_service" "app" {
  name                 = var.app
  location             = var.region
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = var.public

  template {
    service_account = google_service_account.run.email
    scaling {
      max_instance_count = 10
    }
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project}/${var.repo}/${var.app}:${var.git_sha}"
      env {
        name = "APP_SECRET"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.app.secret_id
            version = google_secret_manager_secret_version.app.version # pinned to the version this root manages, never a floating alias (cloud-run-deploy pins too)
          }
        }
      }
      startup_probe {
        http_get {
          path = "/healthz"
        }
      }
    }
  }


  lifecycle {
    # cloud-run-deploy owns the revisions: the image tag of each deploy, the canary split, a rollback.
    # Terraform shapes the service and never resets those on the next apply.
    ignore_changes = [template[0].containers[0].image, traffic, client, client_version]
  }

  depends_on = [
    google_secret_manager_secret_iam_member.run_reads_app,
    google_secret_manager_secret_version.app, # the revision reads the version at start; create it first
    google_project_service.apis,
  ]
}
