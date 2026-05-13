output "namespace" {
  value = kubernetes_namespace.go_assistant.metadata[0].name
}

output "app_service" {
  value = "go-assistant-service"
}
