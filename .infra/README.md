# Infra

## Create k8s cluster

```sh
export PROJECT_ID=redislabs-redisvpc-dev-238506
gcloud container clusters create-auto redis-ai-workshop-dec-2025 \
    --location=europe-west1 \
    --project=$PROJECT_ID
```

## Install Ollama

The default configuration uses NVIDIA L4 GPUs. Make sure your GKE Autopilot cluster has access to GPU resources.

```sh
helm repo add otwld https://helm.otwld.com/
helm repo update
helm install ollama otwld/ollama -f ./helm/ollama/values.yaml
```

### Get Public IP Address

The Ollama service is exposed via LoadBalancer. Get the external IP:

```sh
# Wait for external IP to be assigned (may take 1-2 minutes)
kubectl get svc ollama -w

# Or get it directly once assigned
export OLLAMA_IP=$(kubectl get svc ollama -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "Ollama is available at: http://${OLLAMA_IP}:11434"
```

### Test Ollama

Once you have the external IP, test it:

```sh
# List available models
curl http://${OLLAMA_IP}:11434/api/tags

# Generate a response
curl http://${OLLAMA_IP}:11434/api/generate -d '{
  "model": "llama3.2:3b",
  "prompt": "Why is the sky blue?",
  "stream": false
}'

# Or use the ollama CLI (if installed locally)
export OLLAMA_HOST=http://${OLLAMA_IP}:11434
ollama list
ollama run llama3.2:3b "Why is the sky blue?"
```

### Restrict Access to Specific IP

To restrict access to only your current IP address:

```sh
# Get your current IP
MY_IP=$(curl -s https://api.ipify.org)
echo "Your IP: ${MY_IP}"

# Patch the service to restrict access
kubectl patch svc ollama -p "{\"spec\":{\"loadBalancerSourceRanges\":[\"${MY_IP}/32\"]}}"

# Verify the restriction
kubectl get svc ollama -o jsonpath='{.spec.loadBalancerSourceRanges}'
```

**Note:** The Helm chart doesn't persist `loadBalancerSourceRanges`, so you'll need to reapply this patch after each `helm upgrade`.

## Install AISHE

AISHE requires Ollama to be running in the cluster. Make sure you've installed Ollama first (see above).

```sh
# Install AISHE
helm install aishe ./helm/aishe
```

### Configuration

The AISHE Helm chart is configured to connect to Ollama via the internal Kubernetes service:

- **Ollama Host:** `http://ollama.default.svc.cluster.local:11434`
- **Ollama Model:** `llama3.2:3b`

### Verify AISHE Installation

```sh
# Check pod status
kubectl get pods -l app.kubernetes.io/name=aishe

# Check logs
kubectl logs -l app.kubernetes.io/name=aishe -f

# Test the health endpoint
kubectl port-forward svc/aishe 8000:8000
curl http://localhost:8000/health
```

### Get AISHE Public IP Address

The AISHE service is exposed via LoadBalancer. Get the external IP:

```sh
# Wait for external IP to be assigned (may take 1-2 minutes)
kubectl get svc aishe -w

# Or get it directly once assigned
export AISHE_IP=$(kubectl get svc aishe -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "AISHE is available at: http://${AISHE_IP}:8000"
```

### Test AISHE API

Once you have the external IP, test it:

```sh
# Test the root endpoint
curl http://${AISHE_IP}:8000/

# Ask a question
curl http://${AISHE_IP}:8000/api/v1/ask -X POST \
  -H "Content-Type: application/json" \
  -d '{"question": "What is Kubernetes?"}'

# Ask another question
curl http://${AISHE_IP}:8000/api/v1/ask -X POST \
  -H "Content-Type: application/json" \
  -d '{"question": "Explain Redis in simple terms"}'
```

### Cleanup

```sh
# Uninstall AISHE
helm uninstall aishe

# Uninstall Ollama
helm uninstall ollama

# Delete the cluster (if needed)
gcloud container clusters delete redis-ai-workshop-dec-2025 \
    --location=europe-west1 \
    --project=$PROJECT_ID
```

