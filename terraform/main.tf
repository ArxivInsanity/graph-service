terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.52.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.1"
    }
  }
  cloud {
    hostname     = "app.terraform.io"
    organization = "Arxiv-Insanity"
    workspaces {
      name = "graph-service"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region
  zone    = var.zone
}

data "terraform_remote_state" "gke" {
  backend = "remote"
  config = {
    organization = "Arxiv-Insanity"
    workspaces = {
      name = "app-infra"
    }
  }
}

data "google_client_config" "default" {}

data "google_container_cluster" "my_cluster" {
  name     = data.terraform_remote_state.gke.outputs.gke_outputs.cluster_name
  location = data.terraform_remote_state.gke.outputs.gke_outputs.location
}

provider "kubernetes" {
  host = "https://${data.terraform_remote_state.gke.outputs.gke_outputs.cluster_host}"

  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate)
}

module "graph_service" {
  source              = "./graph_service"
  graph_service_image = var.graph_service_image
  ss_api_key          = var.ss_api_key
  redis_cred          = var.redis_cred
  neo4j_cred          = var.neo4j_cred
}
