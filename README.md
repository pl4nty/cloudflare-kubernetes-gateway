# Cloudflare Kubernetes Gateway

Manage Kubernetes ingress traffic with [Cloudflare Tunnels](https://developers.cloudflare.com/tunnel/) via the [Gateway API](https://gateway-api.sigs.k8s.io/).

## Prerequisites

- [Gateway API](https://gateway-api.sigs.k8s.io/guides/getting-started/introduction/#installing-gateway-api) v1.x

## Getting Started

1. Install cloudflare-kubernetes-gateway with Kustomize:
   <!-- x-release-please-start-version -->
   ```sh
   kubectl apply -k github.com/pl4nty/cloudflare-kubernetes-gateway//config/default?ref=v0.9.0
   ```
   <!-- x-release-please-end -->
2. [Create a Cloudflare API token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) with these permissions:
   - For User Tokens, **Account > Cloudflare Tunnel > Edit** and **Zone > DNS > Edit**
   - For Account Tokens, **Entire Account > Cloudflare One / Zero Trust > Argo Tunnel (Legacy) > Edit** and **All Domains > DNS & Zones > DNS > Edit**
3. Create a Secret with your Cloudflare [Account ID](https://developers.cloudflare.com/fundamentals/setup/find-account-and-zone-ids/) and API token:
   <details>
   <summary>Secret manifest</summary>

   ```yaml
   # kubectl create secret generic -n cloudflare-gateway cloudflare-gateway-token --from-literal=ACCOUNT_ID=<account-id> --from-literal=TOKEN=<api-token>
   apiVersion: v1
   kind: Secret
   metadata:
     name: cloudflare-gateway-token
     namespace: cloudflare-gateway
   type: Opaque
   stringData:
     ACCOUNT_ID: <account-id>
     TOKEN: <api-token>
   ```

   </details>
4. Create a GatewayClass for this controller, referencing the Secret:
   <details>
   <summary>GatewayClass manifest</summary>

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
       name: cloudflare-gateway-token
   ```

   </details>
5. Create a Gateway from the GatewayClass; each Gateway corresponds to a Cloudflare Tunnel.
   In most cases you should only need one per cluster.
   <details>
   <summary>Gateway manifest</summary>

   ```yaml
   apiVersion: gateway.networking.k8s.io/v1
   kind: Gateway
   metadata:
     name: gateway
     namespace: cloudflare-gateway
   spec:
     gatewayClassName: cloudflare
     # listeners are not used but need to be present
     listeners:
     - name: http
       protocol: HTTP
       port: 80
   ```
   </details>
6. Manage traffic with [HTTPRoutes](https://gateway-api.sigs.k8s.io/guides/http-routing/).
   <details>
   <summary>HTTPRoute example manifest</summary>

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

   </details>
7. (optional) Install Prometheus ServiceMonitors to collect controller and cloudflared metrics:
   <!-- x-release-please-start-version -->
   ```sh
   kubectl apply -k github.com/pl4nty/cloudflare-kubernetes-gateway//config/prometheus?ref=v0.9.0
   ```
   <!-- x-release-please-end -->

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

## Configuring cloudflared

By default, a [Cloudflare Tunnel client](https://github.com/cloudflare/cloudflared) (cloudflared) runs for each Gateway, as a Deployment in the Gateway's namespace.
Additional clients can be deployed ([guide](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/deploy-tunnels/deployment-guides/))
to customise parameters that aren't exposed in the Gateway config,
and traffic will be load-balanced between them and the built-in client.

The internal cloudflared can be configured with a ConfigMap referenced from the Gateway:
<details>
<summary>Gateway manifest</summary>

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: gateway
  namespace: cloudflare-gateway
spec:
  gatewayClassName: cloudflare
  listeners:
  - name: http
    protocol: HTTP
    port: 80
  infrastructure:
    parametersRef:
      group: ""
      kind: ConfigMap
      name: gateway
```

</details>
<details>
<summary>ConfigMap manifest</summary>

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gateway
  namespace: cloudflare-gateway
data:
  # Disables the internal cloudflared deployment entirely. Separate clients must be deployed
  disableDeployment: "true" # string

  # The following are literal strings in yaml format that are passed directly through to the deployment spec
  # Values are examples, not the defaults

  # DeploymentSpec.replicas
  replicas: "1" # string
  # PodSpec.nodeSelector
  nodeSelector: |
    disktype: ssd
  # PodSpec.affinity
  affinity: |
    nodeAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: disktype
            operator: In
            values:
            - ssd
  # PodSpec.tolerations
  tolerations: |
    - key: "key1"
      operator: "Equal"
      value: "value1"
      effect: "NoSchedule"
  # PodSpec.containers[gateway].resources
  resources: |
    requests:
      memory: "64Mi"
      cpu: "250m"
    limits:
      memory: "128Mi"
      cpu: "500m"
```

</details>

See also:
- [Assign Pods to Nodes](https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/)
- [Assign Pods to Nodes using Node Affinity](https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes-using-node-affinity/)
- [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/)
