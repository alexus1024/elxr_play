# Workstation Setup

Tools and dependencies installed on WSL2 (Ubuntu) for this project.

## Go

```bash
# Go toolchain (1.24+)
# https://go.dev/doc/install

# Swag - Swagger doc generator
go install github.com/swaggo/swag/cmd/swag@latest
# Ensure ~/go/bin is in PATH
```

## Docker

```bash
# Docker Desktop (WSL2 backend) or Docker Engine
# Needed for building container images
```

## Terraform

```bash
# Option A: HashiCorp APT repo
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt-get update && sudo apt-get install terraform

# Option B: Direct binary
curl -fsSL https://releases.hashicorp.com/terraform/1.10.5/terraform_1.10.5_linux_amd64.zip -o /tmp/tf.zip
unzip /tmp/tf.zip -d /tmp && sudo mv /tmp/terraform /usr/local/bin/
```

## AWS CLI

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o /tmp/awscliv2.zip
unzip /tmp/awscliv2.zip -d /tmp
sudo /tmp/aws/install

# Configure credentials
aws configure
# Access Key ID, Secret Access Key, Region (ca-central-1), Output (json)
```

## kubectl

```bash
# Installed via AWS CLI plugin or standalone
# Configured for EKS:
aws eks update-kubeconfig --name elxr --region ca-central-1
```
