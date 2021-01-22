# Ambassador Example

1. Install kind `https://kind.sigs.k8s.io/`
2. Create a cluster
```
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
```
3. Set up Ambassador
```
kubectl apply -f https://github.com/datawire/ambassador-operator/releases/latest/download/ambassador-operator-crds.yaml
kubectl apply -n ambassador -f https://github.com/datawire/ambassador-operator/releases/latest/download/ambassador-operator-kind.yaml
kubectl wait --timeout=180s -n ambassador --for=condition=deployed ambassadorinstallations/ambassador
```
4. Set up services with ingress
```
kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml
```
5. Add annotation so that Ambassador will detect the ingress
```
kubectl annotate ingress example-ingress kubernetes.io/ingress.class=ambassador
```
6. Test services
```
curl localhost/foo
curl localhost/bar
```
