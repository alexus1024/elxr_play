# RPi5 Manual Setup

One-time steps required after flashing and booting the Pi.

## kubectl
k3s ships with kubectl built-in, no separate install needed.

## kubeconfig
Make kubeconfig available to all tools (helm, kubectl, etc.):
```bash
mkdir -p ~/.kube && ln -s /etc/rancher/k3s/k3s.yaml ~/.kube/config
```

## Helm
```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

Add NATS chart repo:
```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update
```
