# Resource requests and limits

This section describes the resource limits you can set for a container, the disk or networking.

## Container

### Kubernetes keywords

Kuberenetes has [specific resource limits section in the podspec](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container) which lxe will apply through CRI to [LXD container configuration limits](https://lxd.readthedocs.io/en/latest/containers/#keyvalue-configuration). Also see [LXD's CPU limits](https://lxd.readthedocs.io/en/latest/containers/#cpu-limits) and [Kubernetes' How pod limits are run](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#how-pods-with-resource-limits-are-run).

| Kubernetes resource keyword                   | LXD container configuration keyword | Translation Notes                                                                                                       |
|-----------------------------------------------|-------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| `spec.containers[].resources.requests.cpu`    | - (not used)                        | [LXD issue](https://github.com/lxc/lxd/issues/6231)                                                                     |
| `spec.containers[].resources.limits.cpu`      | `limits.cpu.allowance`              | Translated into allowed cpu time usage. E.g. Kuberentes cpu limit of `1.5` or `1500m` cpu will result to `150ms/100ms`. |
| `spec.containers[].resources.requests.memory` | - (not used)                        | -                                                                                                                       |
| `spec.containers[].resources.limits.memory`   | `limits.memory`                     | -                                                                                                                       |
