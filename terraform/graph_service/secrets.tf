resource "kubernetes_secret" "graph_service_secret" {
  metadata {
    name = local.graph_service_secret
  }

  data = {
    REDIS_CRED = var.redis_cred
    NEO4J_CRED = var.neo4j_cred
    S2AG_KEY   = var.ss_api_key
  }
}