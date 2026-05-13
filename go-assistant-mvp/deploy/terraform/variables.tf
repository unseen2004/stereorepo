variable "kubeconfig_path" {
  type    = string
  default = "~/.kube/config"
}

variable "app_image" {
  type    = string
  default = "go-assistant:latest"
}

variable "namespace" {
  type    = string
  default = "go-assistant"
}
