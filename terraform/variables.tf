variable "region" {
  type        = string
  default     = "us-east-2"
  description = "AWS region"
}

variable "cluster_dns" {
  type    = string
  default = "talk-k8s-autoscaling.example.com"
}

variable "cluster_name" {
  type    = string
  default = "education-eks-autoscaling"
}

variable "autoscaler_k8s_namespace" {
  type    = string
  default = "autoscaler"
}

variable "k8s_service_account_name" {
  type    = string
  default = "cluster-autoscaler-aws-cluster-autoscaler"
}

variable "cluster_identity_oidc_issuer_arn" {
  type    = string
  default = "none"
}