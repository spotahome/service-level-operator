# Delivery

In deploy/manifests there are example manifests to deploy this operator.

- Set the correct namespaces on the manifests.
- Set the correct namespace on the service account.
- If you are using [prometheus-operator] check `deploy/manifests/prometheus.yaml` and edit accordingly.
- Image is set to `latest`, this is only the example, it's a bad practice to not use versioned applications.

[prometheus-operator]: https://github.com/coreos/prometheus-operator
