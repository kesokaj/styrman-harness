terraform {
  required_version = "~> 1.16"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 8.0"
    }
  }
  # State in GCS once the bucket exists (SKILL.md § Bootstrap); literal values only:
  # backend "gcs" {
  #   bucket = "PROJECT-tfstate"
  #   prefix = "APP"
  # }
}

provider "google" {
  project = var.project
  region  = var.region
}
