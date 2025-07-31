#!/bin/bash
set -e

# Start Minikube if not running
if ! minikube status &>/dev/null; then
  echo "Starting Minikube..."
  minikube start
fi

# Set Docker to use Minikube's Docker daemon
eval $(minikube docker-env)

# Build the Docker image
echo "Building Docker image..."
docker build -t accio:latest .

# Apply Kubernetes manifests
echo "Applying Kubernetes manifests..."
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/pvc.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Wait for deployment to be ready
echo "Waiting for deployment to be ready..."
kubectl rollout status deployment/accio

# Get the URL
echo "Accio is available at: $(minikube service accio --url)"