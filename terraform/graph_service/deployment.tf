resource "kubernetes_deployment" "graph_service_deployment" {
  depends_on = [kubernetes_secret.graph_service_secret]
  metadata {
    name = local.graph_deployment_label
    labels = {
      App = local.graph_deployment_label
    }
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        App = local.graph_deployment_label
      }
    }
    template {
      metadata {
        labels = {
          App = local.graph_deployment_label
        }
      }
      spec {
        container {
          image             = var.graph_service_image
          name              = local.graph_service_label
          image_pull_policy = "Always"

          port {
            container_port = local.graph_service_port
          }

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
          env_from {
            secret_ref {
              name = local.graph_service_secret
            }
          }
          readiness_probe {
            http_get {
              path = "/"
              port = local.graph_service_port
            }
          }
          liveness_probe {
            http_get {
              path = "/"
              port = local.graph_service_port
            }
          }
        }
      }
    }
  }
}
