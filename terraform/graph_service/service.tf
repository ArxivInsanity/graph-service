resource "kubernetes_service" "graph_service_service" {
  metadata {
    name = local.graph_service_label
  }
  spec {
    selector = {
      App = local.graph_deployment_label
    }
    port {
      port        = local.graph_service_port
      target_port = local.graph_service_port
    }

    type = "NodePort"
  }
}
