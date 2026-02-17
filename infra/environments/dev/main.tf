terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

module "vpc" {
  source         = "../../modules/vpc"
  name           = var.name
  azs            = ["${var.region}a", "${var.region}b"]
  public_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
}

module "eks" {
  source     = "../../modules/eks"
  name       = var.name
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.public_subnet_ids
}

module "ecr" {
  source       = "../../modules/ecr"
  repositories = ["api"]
}
