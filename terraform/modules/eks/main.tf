variable "cluster_name" { type = string }
variable "private_endpoint" { type = bool default = true }
variable "system_min_size" { type = number default = 2 }
variable "app_min_size" { type = number default = 3 }

# Reference skeleton: intentionally provider-neutral at module boundaries.
# Production implementation would pin the AWS/EKS module version and expose
# only the platform capabilities product teams actually need.

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  name               = var.cluster_name
  kubernetes_version = "1.33"

  endpoint_public_access  = !var.private_endpoint
  endpoint_private_access = true

  enable_irsa = true

  eks_managed_node_groups = {
    system = {
      min_size       = var.system_min_size
      desired_size   = var.system_min_size
      max_size       = 6
      instance_types = ["m7i.large"]
      labels         = { workload = "system" }
    }

    application = {
      min_size       = var.app_min_size
      desired_size   = var.app_min_size
      max_size       = 30
      instance_types = ["m7i.xlarge"]
      labels         = { workload = "application" }
    }
  }
}
