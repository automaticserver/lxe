# Development Preview FAQ

During development we found various technical and non-technical issues and reasonings in different interfaces. Let's try to break them down and open them for discussions.

## Images (LXC vs OCI [docker])

### "LXD only does system containers"

LXD [does not offer to create containers from OCI images](https://discuss.linuxcontainers.org/t/using-oci-templates-in-lxd/1911) (follow the thread). LXC itself can. According to the last comment you might be able to work around this issue, possibly by providing a remote for these images, or maybe using the underlying LXC itself.

But here lies a **conceptual conflict between LXD and Kubernetes**. Kubernetes wants only application containers while LXD is designed to only run system containers. For now, LXE supports only but any provided image from a (public) LXD remote.

### Image name format

Kubernetes is too focused on the OCI (or docker?) image name format and communications. It applies [default tags](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/images/image_manager.go#L95), provides only [docker specific credentials](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/container/runtime.go#L140) if they are defined and [validates the image name](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/images/image_manager.go#L150) to [docker grammar](https://github.com/docker/distribution/blob/master/reference/reference.go#L4). So this leaves us with a conflict, that even the image name format of a normal LXD image (e.g. `images:ubuntu/trusty`) is not allowed.

So this leaves us with an awkward workaround to have to specify the following grammar and rules:

```txt
reference          := name [ ":" tag ]
name               := remote-name '/' path-component ['/' path-component]*
remote-name        := remote-component ['.' remote-component]*
remote-component   := /([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9])/
path-component     := alpha-numeric [separator alpha-numeric]*
alpha-numeric      := /[a-z0-9]+/
separator          := /[_.]|[-]*/
```

- `tag` gets ignored, they don't exist in this form in lxc
- `remote-name` must already exist on the host
- interpreted `reference` is matched agains the `lxc image alias` (an alias is used to find the pulled image on the `local` remote)

#### Examples of the image name interpretation

| Pod.Container[].Image  | Kubelet validator | Kubelet interpretes | LXE interpretes | equals LXE alias | in LXC syntax |
| -- | -- | -- | -- | -- | -- |
| busybox | docker.io/library/busybox | busybox:latest | busybox | ???/busybox | [invalid] |
| busybox:other | docker.io/library/busybox:other | busybox:other | busybox | ???/busybox | [invalid] |
| hub.example.io/busybox:other | hub.example.io/busybox:other | hub.example.io/busybox:other | hub.example.io/busybox | hub.example.io/busybox | hub.example.io:busybox |
| hub.example.io/someuser/images/busybox:other | hub.example.io/someuser/images/busybox:other | hub.example.io/someuser/images/busybox:other | hub.example.io/someuser/images/busybox | hub.example.io/someuser/images/busybox | hub.example.io:someuser/images/busybox |
| images/ubuntu/14.04 | docker.io/images/ubuntu/14.04 | images/ubuntu/14.04:latest | images/ubuntu/14.04 | images/ubuntu/14.04 | images:ubuntu/14.04 |
| missingremote/example/ubuntu/14.04 | docker.io/missingremote/example/ubuntu/14.04 | missingremote/example/ubuntu/14.04:latest | missingremote/example/ubuntu/14.04 | missingremote/example/ubuntu/14.04 | [notfound] |

## Environment variables

Environment variables defined in the ContainerSpec of the PodSpec are passed to the [lxd container config](https://lxd.readthedocs.io/en/latest/containers/) as `config.environment.*`, which are passed to the init process of the container (see `cat /proc/1/environ`) and usually the init system does not forward these. In systemd, you could use [PassEnvironment](https://www.freedesktop.org/software/systemd/man/systemd.exec.html#PassEnvironment=) to make these visible for your unit.

## TBD

- only one container per pod (for now)
- cloud-init user-data instead of `PodSpec`'s `command` and `args`
- container kind and lifecycle, exited = shutdown
- Supported networking types and its implications
- Kubernetes' critest
- LXE specific `PodSpec` additions
- Examples / LXC images for kube binaries
