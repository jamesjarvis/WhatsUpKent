# WhatsUpKent deployment

This is just for reference.

To deploy from scratch, start with a kubernetes cluster (Because I'm cheap I use k3s on a small VPS).

Once it is set up, deploy the dgraph service:

```bash
kubectl create -f deployments/dgraph.yaml
```

Then, set up the service for the API (As this is the only one which is going to be actually exposed to the outside)

```bash
kubectl apply -f deployments/service-api.yaml
```

Then, create the initial deployments, and verify they have successfully been created:

```bash
kubectl apply -f deployments/deployment-scraper.yaml
kubectl rollout status deployment whatsupkent-scraper
```

```bash
kubectl apply -f deployments/deployment-api.yaml
kubectl rollout status deployment whatsupkent-api
```

To redeploy, just run the `kubectl apply` commands again, with the updated configuration files, or to just fetch the latest container:
`kubectl rollout restart deployment xyz`

To retrieve the logs (to check if the containers are functioning correctly), run:

```bash
kubectl get pods
```

To retrieve the current pods and their names, then run:

```bash
kubectl logs pod/pod_name_x
```

To retrieve the logs
