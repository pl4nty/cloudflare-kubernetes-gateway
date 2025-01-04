# Cloudflare Kubernetes Gateway

Manage Kubernetes ingress traffic with Cloudflare Tunnels via the [Gateway API](https://gateway-api.sigs.k8s.io/).

## Getting Started

1. Install v1 or later of the Gateway API CRDs: `kubectl apply -k github.com/kubernetes-sigs/gateway-api//config/crd?ref=v1.0.0`
2. Install cloudflare-kubernetes-gateway: `kubectl apply -k github.com/pl4nty/cloudflare-kubernetes-gateway//config/default?ref=v0.8.0` <!-- x-release-please-version -->
3. [Find your Cloudflare account ID](https://developers.cloudflare.com/fundamentals/setup/find-account-and-zone-ids/)
4. [Create a Cloudflare API token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) with the Account Cloudflare Tunnel Edit and Zone DNS Edit permissions
5. Use them to create a Secret: `kubectl create secret -n cloudflare-gateway generic cloudflare --from-literal=ACCOUNT_ID=your-account-id --from-literal=TOKEN=your-token`
6. Create a file containing your GatewayClass, then apply it with `kubectl apply -f file.yaml`:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: cloudflare
spec:
  controllerName: github.com/pl4nty/cloudflare-kubernetes-gateway
  parametersRef:
    group: ""
    kind: Secret
    namespace: cloudflare-gateway
    name: cloudflare
```

7. [Create a Gateway and HTTPRoute(s)](https://gateway-api.sigs.k8s.io/guides/http-routing/) to start managing traffic! For example:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway
  namespace: cloudflare-gateway
spec:
  gatewayClassName: cloudflare
  listeners:
  - protocol: HTTP
    port: 80
    name: http
```

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: example-route
  namespace: default
spec:
  parentRefs:
  - name: gateway
    namespace: cloudflare-gateway
  hostnames:
  - example.com
  rules:
  - backendRefs:
    - name: example-service
      port: 80
```

8. (optional) Install Prometheus ServiceMonitors to collect controller and cloudflared metrics: `kubectl apply -k github.com/pl4nty/cloudflare-kubernetes-gateway//config/prometheus?ref=v0.8.0` <!-- x-release-please-version -->

## Features

The v1 Core spec is not yet supported, as some features (eg header-based routing) aren't available with Tunnels. The following features are supported:

* HTTPRoute hostname and path matching
* HTTPRoute Service backendRefs without filtering or weighting
* Gateway gatewayClassName and listeners only
* GatewayClass Core fields

> [!WARNING]
> Currently, DNS records are not deleted when route hostnames are modified or when routes are deleted.
> Requests to orphaned hostnames respond with an HTTP 404 Not Found, rather than a DNS lookup failure.
> For more details, see [#206](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/206).

<!-- * HTTPRoute Gateway parentRefs, without sectionName
* HTTPRoute hostnames, but not listener filtering or precedence
* HTTPRoute rule path match only
* HTTPRoute backendRefs without filtering or weighting
* Gateway gatewayClassName, listeners aren't validated
* GatewayClass Core fields -->

## Standalone cloudflared

By default, a [Cloudflare Tunnel client](https://github.com/cloudflare/cloudflared) (cloudflared) runs for each Gateway, as a Deployment in the Gateway's namespace.
Additional clients can be deployed ([guide](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/deploy-tunnels/deployment-guides/)) to customise parameters like replicas or tolerations, and traffic will be load-balanced between them and the built-in client.
To disable the built-in Deployment and only use standalone clients:

1. Create a ConfigMap: `kubectl create configmap -n cloudflare-gateway gateway --from-literal=disableDeployment=true`
2. Reference it from the gateway:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway
  namespace: cloudflare-gateway
spec:
  gatewayClassName: cloudflare
  listeners:
  - protocol: HTTP
    port: 80
    name: http
  infrastructure:
    parametersRef:
      group: ""
      kind: ConfigMap
      name: gateway
```
