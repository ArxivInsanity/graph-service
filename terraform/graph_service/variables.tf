variable "graph_service_image" {
  description = "The docker image for graph service application that should be deployed in kubernetes pod"
}
variable "ss_api_key" {
  description = "This semantic scholar api key to be used by graph service"
}
variable "redis_cred" {
  description = "The credentials to connect to redis"
}
variable "neo4j_cred" {
  description = "The credentials to connect to neo4j"
}