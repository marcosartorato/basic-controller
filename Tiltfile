# --- Kubernetes manifests to apply ---
k8s_yaml([
    'k3d/crd/greetings.yaml',
    'k3d/rbac/role.yaml',
    'k3d/rbac/binding.yaml',
    'k3d/deploy/operator.yaml',  # Deployment uses image: localhost:5000/greeting-operator:latest
])

# --- image build + push (fast inner loop) ---
docker_build(
    'k3d-tilt-registry:5000/greeting-operator',
    context='.',
    dockerfile='Dockerfile',
)

# Name your workload resource so Tilt shows logs/port-forwards nicely
k8s_resource('greeting-operator', port_forwards=['8080:8080', '8081:8081'])
