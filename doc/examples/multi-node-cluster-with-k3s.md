# Using LXE in a multinode kubernetes cluster created with k3s and LXD

In this step by step example we create a multinode kubernetes cluster with the help of [k3s](https://k3s.io/) running on [LXD VMs](https://documentation.ubuntu.com/lxd/latest/explanation/instances/). One of the nodes will later be configured to use LXE as kubelet's container runtime that will run kubernetes pods on LXD. You can also use other means of VMs as long as they meet the [requirements for LXD](https://documentation.ubuntu.com/lxd/latest/requirements/) and the [requirements for k3s](https://docs.k3s.io/installation/requirements). In the current example we use Ubuntu noble as the host system with 4 CPU cores and 4GB of RAM.

## Preparations

Either already have 2 separate VMS ready or create those VMs with one host VM using LXD:

### Install LXD on the host and create 2 VMs

On the host [install or update LXD](https://documentation.ubuntu.com/lxd/latest/installing/):

```
root@host:~# snap install lxd --channel=5.0/stable
lxd (5.0/stable) 5.0.4-497fe1e from Canonical✓ installed

root@host:~# snap refresh lxd
snap "lxd" has no updates available

root@host:~# lxc version
Client version: 5.0.4
Server version: 5.0.4
```

And initialize a [minimal LXD setup](https://documentation.ubuntu.com/lxd/latest/howto/initialize/#minimal-setup):

```
root@host:~# lxd init --minimal
```

We are now going to create two Ubuntu noble VMs:

```
root@host:~# lxc launch ubuntu:noble --vm node1 -c limits.cpu=2 -c limits.memory=2GiB
Creating node1
Starting node1

root@host:~# lxc launch ubuntu:noble --vm node2 -c limits.cpu=2 -c limits.memory=2GiB
Creating node2
Starting node2
```

We should now have 2 running VMs (interfaces might take a bit until they're online):

```
root@host:~# lxc list
+-------+---------+------------------------+-------------------------------------------------+-----------------+-----------+
| NAME  |  STATE  |          IPV4          |                      IPV6                       |      TYPE       | SNAPSHOTS |
+-------+---------+------------------------+-------------------------------------------------+-----------------+-----------+
| node1 | RUNNING | 10.22.106.8 (enp5s0)   | fd42:679e:b2b4:3dea:216:3eff:fe5c:d348 (enp5s0) | VIRTUAL-MACHINE | 0         |
+-------+---------+------------------------+-------------------------------------------------+-----------------+-----------+
| node2 | RUNNING | 10.22.106.141 (enp5s0) | fd42:679e:b2b4:3dea:216:3eff:fea9:d46f (enp5s0) | VIRTUAL-MACHINE | 0         |
+-------+---------+------------------------+-------------------------------------------------+-----------------+-----------+
```

And we should be able to exec into them once the [LXD agent](https://documentation.ubuntu.com/lxd/stable-5.21/guest-os-compatibility/#lxd-agent) is running:

```
root@host:~# lxc exec node1 -- bash -l
root@node1:~# exit
exit

root@host:~# lxc exec node2 -- bash -l
root@node2:~# exit
exit
```

Make sure networking is working by e.g. updating the system:

```
root@node1:~# apt-get update && apt-get dist-upgrade
[...]

root@node2:~# apt-get update && apt-get dist-upgrade
[...]
```

### Install k3s on both nodes to create a multinode kuberentes cluster

We [install a minimal k3s setup](https://docs.k3s.io/quick-start) on both of these nodes to form a cluster. The first node will be setup as the master node:

```
root@node1:~# curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.32.4+k3s1" sh -
[INFO]  Using v1.32.4+k3s1 as release
[INFO]  Downloading hash https://github.com/k3s-io/k3s/releases/download/v1.32.4+k3s1/sha256sum-amd64.txt
[INFO]  Downloading binary https://github.com/k3s-io/k3s/releases/download/v1.32.4+k3s1/k3s
[INFO]  Verifying binary download
[INFO]  Installing k3s to /usr/local/bin/k3s
[INFO]  Skipping installation of SELinux RPM
[INFO]  Creating /usr/local/bin/kubectl symlink to k3s
[INFO]  Creating /usr/local/bin/crictl symlink to k3s
[INFO]  Creating /usr/local/bin/ctr symlink to k3s
[INFO]  Creating killall script /usr/local/bin/k3s-killall.sh
[INFO]  Creating uninstall script /usr/local/bin/k3s-uninstall.sh
[INFO]  env: Creating environment file /etc/systemd/system/k3s.service.env
[INFO]  systemd: Creating service file /etc/systemd/system/k3s.service
[INFO]  systemd: Enabling k3s unit
Created symlink /etc/systemd/system/multi-user.target.wants/k3s.service → /etc/systemd/system/k3s.service.
[INFO]  systemd: Starting k3s
```

Wait and check that all pods are running:

```
root@node1:~# watch kubectl get pods -A
Every 2.0s: kubectl get pods -A                                                                               node1: ...

NAMESPACE     NAME                                      READY   STATUS      RESTARTS   AGE
kube-system   coredns-697968c856-xxssm                  1/1     Running     0          46s
kube-system   helm-install-traefik-crd-jntf7            0/1     Completed   0          47s
kube-system   helm-install-traefik-mkgds                0/1     Completed   1          47s
kube-system   local-path-provisioner-774c6665dc-fd9qr   1/1     Running     0          46s
kube-system   metrics-server-6f4c6675d5-2lrl8           1/1     Running     0          46s
kube-system   svclb-traefik-20713c58-btsvz              2/2     Running     0          27s
kube-system   traefik-c98fdf6fb-jcfxt                   1/1     Running     0          27s
```

You can also [get the generated kubeconfig](https://docs.k3s.io/cluster-access) and use it directly with kubectl.

Next, obtain the node token so we can add the second node later:

```
root@node1:~# cat /var/lib/rancher/k3s/server/node-token
K1029a66d5894033508ee58745bbf603b9389c159082f9c4766f76320d05b7eef61::server:9888e64abbf745d4aa376c8462fa2619
```

Install k3s on the second node as an additional node (LXD offers dns domain lxd over dhcp so we can use `node1.lxd` as the hostname):

```
root@node2:~# curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="v1.32.4+k3s1" K3S_URL=https://node1.lxd:6443 K3S_TOKEN=K1029a66d5894033508ee58745bbf603b9389c159082f9c4766f76320d05b7eef61::server:9888e64abbf745d4aa376c8462fa2619 sh -
[INFO]  Using v1.32.4+k3s1 as release
[INFO]  Downloading hash https://github.com/k3s-io/k3s/releases/download/v1.32.4+k3s1/sha256sum-amd64.txt
[INFO]  Downloading binary https://github.com/k3s-io/k3s/releases/download/v1.32.4+k3s1/k3s
[INFO]  Verifying binary download
[INFO]  Installing k3s to /usr/local/bin/k3s
[INFO]  Skipping installation of SELinux RPM
[INFO]  Creating /usr/local/bin/kubectl symlink to k3s
[INFO]  Creating /usr/local/bin/crictl symlink to k3s
[INFO]  Creating /usr/local/bin/ctr symlink to k3s
[INFO]  Creating killall script /usr/local/bin/k3s-killall.sh
[INFO]  Creating uninstall script /usr/local/bin/k3s-agent-uninstall.sh
[INFO]  env: Creating environment file /etc/systemd/system/k3s-agent.service.env
[INFO]  systemd: Creating service file /etc/systemd/system/k3s-agent.service
[INFO]  systemd: Enabling k3s-agent unit
Created symlink /etc/systemd/system/multi-user.target.wants/k3s-agent.service → /etc/systemd/system/k3s-agent.service.
[INFO]  systemd: Starting k3s-agent
```

Check on the first node that the second node has been added:

```
root@node1:~# kubectl get nodes
NAME    STATUS   ROLES                  AGE     VERSION
node1   Ready    control-plane,master   2m19s   v1.32.4+k3s1
node2   Ready    <none>                 21s     v1.32.4+k3s1
```

[Taint](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) node2 as [LXD doesn't support OCI images](https://documentation.ubuntu.com/lxd/latest/explanation/instances/#application-containers-vs-system-containers) so that those pods don't get scheduled on this node:

```
root@node1:~# kubectl taint nodes node2 containerruntime=lxd:NoSchedule
node/node2 tainted

root@node1:~# kubectl taint nodes node2 containerruntime=lxd:NoExecute
node/node2 tainted
```

[Label](https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/) node2 so pods can have a constraint to be scheduled to such a node with this label:

```
root@node1:~# kubectl label nodes node2 containerruntime=lxd
node/node2 labeled
```

All current pods should be only running on the first node:

```
root@node1:~# kubectl get pods -A -o wide
NAMESPACE     NAME                                      READY   STATUS      RESTARTS   AGE     IP          NODE    NOMINATED NODE   READINESS GATES
kube-system   coredns-697968c856-xxssm                  1/1     Running     0          4m3s    10.42.0.2   node1   <none>           <none>
kube-system   helm-install-traefik-crd-jntf7            0/1     Completed   0          4m4s    10.42.0.5   node1   <none>           <none>
kube-system   helm-install-traefik-mkgds                0/1     Completed   1          4m4s    10.42.0.6   node1   <none>           <none>
kube-system   local-path-provisioner-774c6665dc-fd9qr   1/1     Running     0          4m3s    10.42.0.4   node1   <none>           <none>
kube-system   metrics-server-6f4c6675d5-2lrl8           1/1     Running     0          4m3s    10.42.0.3   node1   <none>           <none>
kube-system   svclb-traefik-20713c58-btsvz              2/2     Running     0          3m44s   10.42.0.7   node1   <none>           <none>
kube-system   traefik-c98fdf6fb-jcfxt                   1/1     Running     0          3m44s   10.42.0.8   node1   <none>           <none>
```

[Install or update LXD](https://linuxcontainers.org/lxd/getting-started-cli/#installing-a-package) here as well:

```
root@node2:~# snap install lxd --channel=5.0/stable
lxd (5.0/stable) 5.0.4-497fe1e from Canonical✓ installed

root@node2:~# snap refresh lxd
snap "lxd" has no updates available

root@node2:~# lxc version
Client version: 5.0.4
Server version: 5.0.4
```

And initialize a [minimal LXD setup](https://linuxcontainers.org/lxd/getting-started-cli/#initial-configuration) here as well:

```
root@node2:~# lxd init --minimal
```

## Install LXE

The easiest way to install LXE is by using the `go install` command (see [installation instructions](https://github.com/automaticserver/lxe#installing-lxe-from-source) for the "traditional" way):

```
root@node2:~# apt-get install golang
[...]

root@node2:~# go version
go version go1.22.2 linux/amd64

root@node2:~# go install github.com/automaticserver/lxe/cmd/lxe@latest
go: downloading ...
[...]
```

Make sure you have run an lxc command at least once as the lxc command will generate the default remotes file, you'll get an extra output in the beginning. If you followed the steps, this should've happened in `lxc version` above.

> If this is your first time running LXD on this machine ...

There should be a remote config file now in `/root/snap/lxd/common/config/config.yml` (Here we use the same remote config file in LXE. This way any additions to this file benefit both the root user and LXE):

```
root@node2:~# cat /root/snap/lxd/common/config/config.yml 
default-remote: local
remotes:
  local:
    addr: unix://
    public: false
aliases: {}
```

We're also going to reuse the default LXD bridge `lxdbr0` (for now) which LXD has created during initialisation.

```
root@node2:~# lxc network list
+--------+----------+---------+-----------------+---------------------------+-------------+---------+---------+
|  NAME  |   TYPE   | MANAGED |      IPV4       |           IPV6            | DESCRIPTION | USED BY |  STATE  |
+--------+----------+---------+-----------------+---------------------------+-------------+---------+---------+
| cni0   | bridge   | NO      |                 |                           |             | 0       |         |
+--------+----------+---------+-----------------+---------------------------+-------------+---------+---------+
| enp5s0 | physical | NO      |                 |                           |             | 0       |         |
+--------+----------+---------+-----------------+---------------------------+-------------+---------+---------+
| lxdbr0 | bridge   | YES     | 10.240.106.1/24 | fd42:8bd7:6585:3356::1/64 |             | 1       | CREATED |
+--------+----------+---------+-----------------+---------------------------+-------------+---------+---------+
```

Start lxe with the obtained informations above and leave the terminal open:

```
root@node2:~# /root/go/bin/lxe --log-level info
WARNING[05-18|16:58:30.422] starting lxe...                               builddate=undef buildnumber=undef gitcommit=undef gittreestate=undef packagename=undef version=0.0.0
INFO   [05-18|16:58:30.452] Connected to LXD                              lxdsocket=/var/snap/lxd/common/lxd/unix.socket
INFO   [05-18|16:58:30.543] started lxe CRI shim                          socket=/run/lxe.sock
INFO   [05-18|16:58:30.544] started streaming server                      baseurl="http://localhost:44124" endpoint="localhost:44124"
```

## Set LXE as runtime

In a new terminal session login again to the second node and now we configure the k3s agent to use some specific kubelet flags so it uses the LXE socket:

```
root@node2:~# mkdir /etc/rancher/k3s

root@node2:~# cat <<EOF >/etc/rancher/k3s/config.yaml
kubelet-arg:
  - "container-runtime-endpoint=unix:///run/lxe.sock"
  - "containerd="
EOF
```

Restart the k3s agent service and we should see the arguments above are passed to the kubelet:

```
root@node2:~# systemctl restart k3s-agent

root@node2:~# journalctl -u k3s-agent | grep "Running kubelet"
[...]
May 07 19:49:43 node2 k3s[6386]: time="2025-05-07T19:49:43Z" level=info msg="Running kubelet [...] --container-runtime-endpoint=unix:///run/lxe.sock --containerd= [...]
```

If we check with kubectl (from the first node) we can also see that the node is now running on LXE:

```
root@node1:~# kubectl get node node2 -o wide
NAME    STATUS   ROLES    AGE   VERSION        INTERNAL-IP     EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
node2   Ready    <none>   17m   v1.32.4+k3s1   10.22.106.141   <none>        Ubuntu 24.04.2 LTS   6.8.0-59-generic   lxe://0.0.0
```

## Launch the first pod

We are ready to launch the first pod on LXD ([using a scheduling hint](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/) to only assign this pod to nodes with matching labels):

```
root@node1:~# cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: ubuntu
  namespace: default
spec:
  containers:
  - name: ubuntu
    image: ubuntu/noble
  nodeSelector:
    containerruntime: lxd
  tolerations:
  - key: containerruntime
    operator: Exists
EOF

pod/ubuntu created
```

You'll see that kubelet is sending commands to LXE in its output:

```
INFO   [05-07|19:51:35.133] run pod                                       namespace=default podname=ubuntu poduid=15e74145-474e-4deb-a909-770f06c0f0d1
INFO   [05-07|19:51:35.177] run pod successful                            namespace=default podid=u7eosd7tz6l5b37p podname=ubuntu poduid=15e74145-474e-4deb-a909-770f06c0f0d1
INFO   [05-07|19:51:45.251] create container                              attempt=0 containername=ubuntu image=5ff364a309c64c01bd06b87ba8f17728b3ff59339fc8561a40cb778dfb5931cc podid=u7eosd7tz6l5b37p
INFO   [05-07|19:51:55.542] create container successful                   attempt=0 containername=ubuntu image=5ff364a309c64c01bd06b87ba8f17728b3ff59339fc8561a40cb778dfb5931cc podid=u7eosd7tz6l5b37p
INFO   [05-07|19:51:55.544] start container                               containerid=uh2rbciu2w5ciw6b
INFO   [05-07|19:51:55.874] event detected                                containerid=uh2rbciu2w5ciw6b event=instance-started
INFO   [05-07|19:51:55.900] start container successful                    containerid=uh2rbciu2w5ciw6b
```

And kubernetes sees the pod as running:

```
root@node1:~# kubectl get pods -o wide
NAME     READY   STATUS    RESTARTS      AGE     IP               NODE    NOMINATED NODE   READINESS GATES
ubuntu   1/1     Running   0             3m31s   10.240.106.169   node2   <none>           <none>
```

We can now exec into the pod and do stuff:

```
root@node1:~# kubectl exec -it ubuntu -- bash -l
root@ubuntu:~# uptime
 19:57:42 up 1 min,  0 user,  load average: 0.13, 0.35, 0.24
root@ubuntu:~# apt-get update && apt-get dist-upgrade
[...]
root@ubuntu:~# exit
exit
```

## What's next?

This example took a little shortcut in networking by using the default lxd bridge. Node administrators usually want to hook up these containers into the kubernetes networking setup using CNI, e.g. kubeproxy, flannel or whatever they use.

LXE supports CNI by definining `--network-mode cni`. Ubuntu and Debian offer officially the package `containernetworking-plugins` for the standard CNI plugins. You either install and setup the appropriate networking services on the second node itself or create appropriate pods/daemonsets with either the finished prepared image or using cloud-init to install and setup those services. Using the former variant requires the pods to want host networking capability and for that LXE needs [a little file](../../fixtures/hostnetwork.conf) **persisted** for the `--hostnetwork-file` argument, e.g. save it to `/var/lib/lxe/hostnetwork.conf`. See `lxe --help` for more info and further configuration options.
