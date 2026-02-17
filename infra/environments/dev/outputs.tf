output "cluster_name" {
  value = module.eks.cluster_name
}

output "ecr_urls" {
  value = module.ecr.repository_urls
}
