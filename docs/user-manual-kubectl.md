# Install Mondoo Operator with kubectl

The following steps sets up the mondoo operator using kubectl and a manifest file.

## Preconditions:

From manifests file

- kubectl with admin role

## Deployment of Operator using Manifests

//TODO where to fetch manifests from

1. GET operator manifests

```bash
kubectl apply -f mondoo-operator-manifests.yaml
```

2. Configure the Mondoo secret:

- Download service account from [Mondooo](https://mondoo.com)
- Convert json to yaml via:

```bash
yq e -P creds.json > creds.yaml
```

- Store service account as a secret in the mondoo namespace via:

```bash
kubectl create secret generic mondoo-client --namespace mondoo-operator-system --from-file=config=creds.yaml
```

Once the secret is configure, we configure the operator to define the scan targets:

3. Create `mondoo-config.yaml`

```yaml
apiVersion: k8s.mondoo.com/v1alpha1
kind: MondooAuditConfig
metadata:
  name: mondoo-client
  namespace: mondoo-operator-system
spec:
  workloads:
    enable: true
    serviceAccount: mondoo-operator-workload
  nodes:
    enable: true
  mondooSecretRef: mondoo-client
```

4. Apply the configuration via:

```bash
kubectl apply -f mondoo-config.yaml
```

## FAQ

**I do not see the service running, only the operator?**

First check that the CRD is properly registered with the operator:

```bash
kubectl get crd
NAME                           CREATED AT
mondooauditconfigs.k8s.mondoo.com   2022-01-14T14:07:28Z
```

Then make sure a configuration for the mondoo client is deployed:

```bash
kubectl get mondooauditconfigs
NAME                  AGE
mondooauditconfig-sample   2m44s
```
