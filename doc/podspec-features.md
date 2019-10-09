# PodSpec and Container API implementation

The following table provides an overview of the current implementation of the [`PodSpec` v1 core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#podspec-v1-core) and [`Container` v1 core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#container-v1-core) API.

- Some properties have implementation specific characteristics, *marked with an asterisk (\*)*
- Properties without CRI implications are *marked with a dash (-)*.
- Properties who are CRI related but not yet taken into account are *marked with a question mark (?)*

| `PodSpec` property  | In LXE implemented | Notes | Related LXC config |
| -- | -- | -- | -- |
| `activeDeadlineSeconds` | - | _not CRI related_ |  |
| `affinity` | - | _not CRI related_ |  |
| `automountServiceAccountToken` | yes | implicitly provided with [`CRI Mounts`](https://github.com/kubernetes/kubernetes/blob/release-1.12/pkg/kubelet/apis/cri/runtime/v1alpha2/api.pb.go#L1835) |  |
| `containers` | yes* | only one container per pod currently, see [FAQ](development-preview-faq.md) | the lxc containers |
| `dnsConfig` | yes | see `dnsPolicy` | |
| `dnsPolicy` | yes | kubelet does all the work and provides the target settings |  |
| `hostAliases` | yes | kubelet does all the work and provides the hosts file as CRI Mount |  |
| `hostIPC` | ? |  |  |
| `hostNetwork` | yes* | if false LXE calls [CNI](https://github.com/containernetworking/cni/blob/master/SPEC.md#network-configuration) | if true then `config.raw.lxc.include` to a file containing `lxc.net.0.type=none` |
| `hostPID` | ? |  |  |
| `hostname` | yes* | providing hostname using cloud-init vendor-data, see [FAQ](development-preview-faq.md) | unfortunately in LXD the container name *is* the hostname, so providing via `config.user.vendor-data` |
| `imagePullSecrets` | ? | authentication to LXD servers are different than to docker, see `container.image` |  |
| `initContainers` | ? |  |  |
| `nodeName` | - | _not CRI related_ |  |
| `nodeSelector` | - | _not CRI related_ |  |
| `priority` | - | _not CRI related_ |  |
| `priorityClassName` | - | _not CRI related_ |  |
| `readinessGates` | - | _not CRI related_ |  |
| `restartPolicy` | - | _not CRI related_ |  |
| `runtimeClassName` | - | _not CRI related_ |  |
| `schedulerName` | - | _not CRI related_ |  |
| `securityContext` | incomplete* |  |  |
| `serviceAccount` | - | _not CRI related_ |  |
| `serviceAccountName` | - | _not CRI related_ |  |
| `shareProcessNamespace` | ? |  |  |
| `subdomain` | - | _Not CRI related_ |  |
| `terminationGracePeriodSeconds` | - | _Not CRI related_ |  |
| `tolerations` | - | _Not CRI related_ |  |
| `volumes` | yes | only `container.volumeMounts` are relevant for CRI |  |

| `Container` property  | In LXE implemented | Notes | Related LXC config |
| -- | -- | -- | -- |
| `args` | no* | see below `command` |  |
| `command` | no* | lxc containers with lxd have no entrypoint-like option, can be differently provided with cloud-init user-data, see [FAQ](development-preview-faq.md) | `config.user.user-data` |
| `env` | yes* | there are some additional reserved fields for cloud-init: `env.meta-data`, `env.network-config`, `env.user-data` | `config.environment.*` |
| `envFrom` | yes | kubelet does all the work and are merged with `env` |  |
| `image` | yes* | only lxc images, see [FAQ](development-preview-faq.md) | the container image |
| `imagePullPolicy` | yes | kubelet decides itself when to pull the image through CRI |  |
| `lifecycle` | - | _not CRI related_ |  |
| `livenessProbe` | - | _not CRI related_ |  |
| `name` | yes |  |  |
| `ports` | yes |  | `config.devices.*.type=proxy` |
| `readinessProbe` | - | _not CRI related_ |  |
| `resources` | yes | see [limits.md](limits.md) | `config.limits.*` |
| `securityContext` | incomplete* | yet only `securityContext.privileged` | `config.security.privileged` |
| `stdin` | ? |  |  |
| `stdinOnce` | ? |  |  |
| `terminationMessagePath` | ? |  |  |
| `terminationMessagePolicy` | ? |  |  |
| `tty` | ? |  |  |
| `volumeDevices` | yes | with [`CRI Devices`](https://github.com/kubernetes/kubernetes/blob/release-1.12/pkg/kubelet/apis/cri/runtime/v1alpha2/api.pb.go#L1837) | `config.devices.*.type=block` |
| `volumeMounts` | yes | with [`CRI Mounts`](https://github.com/kubernetes/kubernetes/blob/release-1.12/pkg/kubelet/apis/cri/runtime/v1alpha2/api.pb.go#L1835) | `config.devices.*.type=disk` |
| `workingDir` | ? |  |  |
