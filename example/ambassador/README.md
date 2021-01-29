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
3. Get cluster info
```
kubectl cluster-info --context kind-kind
```
4. Set up Ambassador
```
kubectl apply -f https://github.com/datawire/ambassador-operator/releases/latest/download/ambassador-operator-crds.yaml
kubectl apply -n ambassador -f https://github.com/datawire/ambassador-operator/releases/latest/download/ambassador-operator-kind.yaml
kubectl wait --timeout=180s -n ambassador --for=condition=deployed ambassadorinstallations/ambassador
```
5. Set up services with ingress
```
kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml
```
6. Add annotation so that Ambassador will detect the ingress
```
kubectl annotate ingress example-ingress kubernetes.io/ingress.class=ambassador
```
7. Set up zipkin
```
kubectl apply -f zipkin.yaml
```
8. Restart ambassador (needed to pick up tracing configuration)
```
kubectl rollout restart deployment -n ambassador
```
9. Set up authentication
```
kubectl apply -f authz.yaml
kubectl apply -f authz-ambassador.yaml
```
10. Test services
```
curl -H "Content-Type: application/json" --data '{"abc": 123}' -v localhost/bar
curl -H "Content-Type: application/json" --data '{"abc": 123}' -v localhost/foo
```
11. View results in zipkin:
```
kubectl get pods
```
Returns something like
```
NAME                            READY   STATUS    RESTARTS   AGE
bar-app                         1/1     Running   0          28m
example-auth-79746f86c9-zc6js   1/1     Running   0          27m
foo-app                         1/1     Running   0          28m
zipkin-68bdbf69f6-9gwp5         1/1     Running   0          27m
```
Make zipkin accessible
```
kubectl port-forward zipkin-68bdbf69f6-9gwp5 9411
```
In browser, go to
```
localhost:9411
```
