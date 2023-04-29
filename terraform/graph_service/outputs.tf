output "graph_service_endpoint" {
    value = "${local.graph_service_label}:${local.graph_service_port}"
}