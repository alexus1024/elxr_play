



resource "aws_ecr_repository" "this" {
  for_each = toset(var.repositories)

  name                 = each.key
  image_tag_mutability = "MUTABLE"
  force_delete         = true
}
