# Using LXE in a multinode kubernetes cluster created with k3s and LXD

In this step by step example we create a multinode kubernetes cluster with the help of [k3s](https://k3s.io/) running on [LXD VMs](https://linuxcontainers.org/lxd/docs/master/virtual-machines/). One of the nodes will later be configured to use LXE as kubelet's container runtime that will run kubernetes pods on LXD. You can also use other means of VMs as long as they meet the [requirements for LXD](https://linuxcontainers.org/lxd/docs/master/requirements/) and the [requirements for k3s](https://rancher.com/docs/k3s/latest/en/installation/installation-requirements/). In the current example we use Ubuntu jammy as the host system with 4 CPU cores and 4GB of RAM.

## Preparations

Either already have 2 separate VMS ready or create those VMs with one host VM using LXD:

### Install LXD on the host and create 2 VMs

On the host [install or update LXD](https://linuxcontainers.org/lxd/getting-started-cli/#installing-a-package):

```
root@host:~# snap install lxd
snap "lxd" is already installed, see 'snap help refresh'

root@host:~# snap refresh lxd
snap "lxd" has no updates available

root@host:~# lxc version
Client version: 5.0.0
Server version: 5.0.0
```

And initialize a [minimal LXD setup](https://linuxcontainers.org/lxd/getting-started-cli/#initial-configuration):

```
root@host:~# lxd init --minimal
```

We are now going to create Ubuntu VMs which [require an extra step to access it after](https://discuss.linuxcontainers.org/t/running-virtual-machines-with-lxd-4-0/7519#extra-steps-for-official-ubuntu-images-6). The following command prepares an instance config that has a few steps stitched together and will configure the root user with the plain-text password `secret`:

```
root@host:~# cat <<EOF >lxd-ubuntu-vm-config.yaml
config:
  user.user-data: |
    #cloud-config
    users:
      - name: root
        lock_passwd: false
        hashed_passwd: '$1$SaltSalt$YhgRYajLPrYevs14poKBQ0'
devices:
  config:
    source: cloud-init:config
    type: disk
EOF
```

And create 2 Ubuntu jammy VMs:

```
root@host:~# lxc launch ubuntu:jammy --vm node1 < lxd-ubuntu-vm-config.yaml 
Creating node1
Starting node1

root@host:~# lxc launch ubuntu:jammy --vm node2 < lxd-ubuntu-vm-config.yaml 
Creating node2
Starting node2
```

We should now have 2 running VMs (interfaces might take a bit until they're online):

```
root@host:~# lxc list
+-------+---------+-----------------------+-------------------------------------------------+-----------------+-----------+
| NAME  |  STATE  |         IPV4          |                      IPV6                       |      TYPE       | SNAPSHOTS |
+-------+---------+-----------------------+-------------------------------------------------+-----------------+-----------+
| node1 | RUNNING | 10.4.147.50 (enp5s0)  | fd42:44cd:69a4:e5d2:216:3eff:fe35:20 (enp5s0)   | VIRTUAL-MACHINE | 0         |
+-------+---------+-----------------------+-------------------------------------------------+-----------------+-----------+
| node2 | RUNNING | 10.4.147.154 (enp5s0) | fd42:44cd:69a4:e5d2:216:3eff:fe49:2c9c (enp5s0) | VIRTUAL-MACHINE | 0         |
+-------+---------+-----------------------+-------------------------------------------------+-----------------+-----------+
```

And we should be able to exec into them once the LXD VM agent is running (Alternatively login from the console using `lxc console <nodename>` using user `root` with password `secret`):

```
root@host:~# lxc exec node1 bash
root@node1:~# exit
exit

root@host:~# lxc exec node2 bash
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

We [install a minimal k3s setup](https://rancher.com/docs/k3s/latest/en/quick-start/) on both of these nodes to form a cluster. The first node will be setup as the master node:

```
root@node1:~# curl -sfL https://get.k3s.io | sh -
[INFO]  Finding release for channel stable
[INFO]  Using v1.23.6+k3s1 as release
[INFO]  Downloading hash https://github.com/k3s-io/k3s/releases/download/v1.23.6+k3s1/sha256sum-amd64.txt
[INFO]  Downloading binary https://github.com/k3s-io/k3s/releases/download/v1.23.6+k3s1/k3s
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
root@node1:~# kubectl get pods -A
NAMESPACE     NAME                                      READY   STATUS      RESTARTS   AGE
kube-system   local-path-provisioner-6c79684f77-7wmr6   1/1     Running     0          2m9s
kube-system   coredns-d76bd69b-tpqbr                    1/1     Running     0          2m9s
kube-system   metrics-server-7cd5fcb6b7-kghqm           1/1     Running     0          2m9s
kube-system   helm-install-traefik-crd-vlgvw            0/1     Completed   0          2m10s
kube-system   helm-install-traefik-btjjt                0/1     Completed   2          2m10s
kube-system   svclb-traefik-c4tjr                       2/2     Running     0          31s
kube-system   traefik-df4ff85d6-jnh9r                   1/1     Running     0          32s
```

You can also [get the generated kubeconfig](https://rancher.com/docs/k3s/latest/en/cluster-access/) and use it directly with kubectl.

Next, obtain the node token so we can add the second node later:

```
root@node1:~# cat /var/lib/rancher/k3s/server/node-token
K1056ffce43bdcab0acdad527b8cfb1d8a00afb4e3a4e01eed2dcb54acdc075c84d::server:61d8b0d8abcbb0279c71e28cdc7eb1fe
```

Install k3s on the second node as an additional node (LXD offers dns domain lxd over dhcp so we can use `node1.lxd` as the hostname):

```
root@node2:~# curl -sfL https://get.k3s.io | K3S_URL=https://node1.lxd:6443 K3S_TOKEN=K1056ffce43bdcab0acdad527b8cfb1d8a00afb4e3a4e01eed2dcb54acdc075c84d::server:61d8b0d8abcbb0279c71e28cdc7eb1fe sh -
[INFO]  Finding release for channel stable
[INFO]  Using v1.23.6+k3s1 as release
[INFO]  Downloading hash https://github.com/k3s-io/k3s/releases/download/v1.23.6+k3s1/sha256sum-amd64.txt
[INFO]  Downloading binary https://github.com/k3s-io/k3s/releases/download/v1.23.6+k3s1/k3s
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
node1   Ready    control-plane,master   7m51s   v1.23.6+k3s1
node2   Ready    <none>                 4s      v1.23.6+k3s1
```

[Taint](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/) node2 as [LXD doesn't support OCI images](https://discuss.linuxcontainers.org/t/using-oci-templates-in-lxd/1911/4) (like docker uses) so that those pods don't get scheduled on this node:

```
root@node1:~# kubectl taint nodes node2 containerruntime=lxd:NoSchedule
node/node2 tainted

root@node1:~# kubectl taint nodes node2 containerruntime=lxd:NoExecute
node/node2 tainted
```

All current pods should be only running on the first node:

```
root@node1:~# kubectl get pods -A -o wide
NAMESPACE     NAME                                      READY   STATUS      RESTARTS   AGE   IP          NODE    NOMINATED NODE   READINESS GATES
kube-system   local-path-provisioner-6c79684f77-7wmr6   1/1     Running     0          17m   10.42.0.5   node1   <none>           <none>
kube-system   coredns-d76bd69b-tpqbr                    1/1     Running     0          17m   10.42.0.6   node1   <none>           <none>
kube-system   metrics-server-7cd5fcb6b7-kghqm           1/1     Running     0          17m   10.42.0.2   node1   <none>           <none>
kube-system   helm-install-traefik-crd-vlgvw            0/1     Completed   0          17m   10.42.0.4   node1   <none>           <none>
kube-system   helm-install-traefik-btjjt                0/1     Completed   2          17m   10.42.0.3   node1   <none>           <none>
kube-system   svclb-traefik-c4tjr                       2/2     Running     0          16m   10.42.0.8   node1   <none>           <none>
kube-system   traefik-df4ff85d6-jnh9r                   1/1     Running     0          16m   10.42.0.7   node1   <none>           <none>
```

[Install or update LXD](https://linuxcontainers.org/lxd/getting-started-cli/#installing-a-package) here as well:

```
root@node2:~# snap install lxd
snap "lxd" is already installed, see 'snap help refresh'

root@node2:~# snap refresh lxd
snap "lxd" has no updates available

root@node2:~# lxc version
Client version: 5.0.0
Server version: 5.0.0
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
go version go1.18.1 linux/amd64

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
  images:
    addr: https://images.linuxcontainers.org
    protocol: simplestreams
    public: true
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
| lxdbr0 | bridge   | YES     | 10.108.115.1/24 | fd42:5694:8d95:4f1c::1/64 |             | 1       | CREATED |
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
  - "container-runtime=remote"
  - "container-runtime-endpoint=unix:///run/lxe.sock"
  - "containerd="
EOF
```

Restart the k3s agent service and we should see the arguments above are passed to the kubelet:

```
root@node2:~# systemctl restart k3s-agent

root@node2:~# journalctl -u k3s-agent | grep "Running kubelet"
[...]
May 18 19:57:23 node2 k3s[15774]: time="2022-05-18T19:57:23Z" level=info msg="Running kubelet --address=0.0.0.0 --anonymous-auth=false --authentication-token-webhook=true --authorization-mode=Webhook --cgroup-driver=cgroupfs --client-ca-file=/var/lib/rancher/k3s/agent/client-ca.crt --cloud-provider=external --cluster-dns=10.43.0.10 --cluster-domain=cluster.local --cni-bin-dir=/var/lib/rancher/k3s/data/8c2b0191f6e36ec6f3cb68e2302fcc4be850c6db31ec5f8a74e4b3be403101d8/bin --cni-conf-dir=/var/lib/rancher/k3s/agent/etc/cni/net.d --container-runtime=remote --container-runtime-endpoint=unix:///run/lxe.sock --containerd= --eviction-hard=imagefs.available<5%,nodefs.available<5% --eviction-minimum-reclaim=imagefs.available=10%,nodefs.available=10% --fail-swap-on=false --healthz-bind-address=127.0.0.1 --hostname-override=node2 --kubeconfig=/var/lib/rancher/k3s/agent/kubelet.kubeconfig --node-labels= --pod-manifest-path=/var/lib/rancher/k3s/agent/pod-manifests --read-only-port=0 --resolv-conf=/tmp/k3s-resolv.conf --serialize-image-pulls=false --tls-cert-file=/var/lib/rancher/k3s/agent/serving-kubelet.crt --tls-private-key-file=/var/lib/rancher/k3s/agent/serving-kubelet.key"
```

If we check with kubectl (from the first node) we can also see that the node is now running on LXE:

```
root@node1:~# kubectl get node node2 -o wide
NAME    STATUS   ROLES    AGE     VERSION        INTERNAL-IP    EXTERNAL-IP   OS-IMAGE           KERNEL-VERSION    CONTAINER-RUNTIME
node2   Ready    <none>   8m10s   v1.23.6+k3s1   10.4.147.106   <none>        Ubuntu 22.04 LTS   5.15.0-1008-kvm   lxe://0.0.0
```

## Launch the first pod

We are ready to launch the first pod on LXD:

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
    image: ubuntu/jammy
  tolerations:
  - key: containerruntime
    operator: Exists
EOF

pod/ubuntu created
```

You'll see that kubelet is sending commands to LXE in its output:

```
INFO   [06-08|11:08:48.064] run pod                                       namespace=default podname=ubuntu poduid=ede51976-e988-4ae9-a9c0-6f0c5d0b26c1
INFO   [06-08|11:08:48.387] run pod successful                            namespace=default podid=uk7h64y23ycwf7hs podname=ubuntu poduid=ede51976-e988-4ae9-a9c0-6f0c5d0b26c1
INFO   [06-08|11:09:05.410] create container                              attempt=0 containername=ubuntu image=c73fb1ddeb3ba971b230e79565817cd5a8e6053bfa9526afe19cd10e3008f895 podid=uk7h64y23ycwf7hs
INFO   [06-08|11:09:51.655] create container successful                   attempt=0 containername=ubuntu image=c73fb1ddeb3ba971b230e79565817cd5a8e6053bfa9526afe19cd10e3008f895 podid=uk7h64y23ycwf7hs
INFO   [06-08|11:09:51.697] start container                               containerid=uetrgnza3i2rugcr
INFO   [06-08|11:09:55.016] start container successful                    containerid=uetrgnza3i2rugcr
```

And kubernetes sees the pod as running:

```
root@node1:~# kubectl get pod ubuntu -o wide
NAME     READY   STATUS    RESTARTS   AGE     IP              NODE    NOMINATED NODE   READINESS GATES
ubuntu   1/1     Running   0          2m10s   10.108.115.23   node2   <none>           <none>
```

We can now exec into the pod and do stuff:

```
root@node1:~# kubectl exec -it ubuntu -- bash
root@ubuntu:~# uptime
 11:12:07 up 2 min,  0 users,  load average: 0.60, 0.95, 0.69
root@ubuntu:~# apt-get update && apt-get dist-upgrade
[...]
root@ubuntu:~# exit
exit
```

## What's next?

This example took a little shortcut in networking by using the default lxd bridge. Node administrators usually want to hook up these containers into the kubernetes networking setup using CNI, kubeproxy, flannel or whatever they use.

LXE supports CNI by definining `--network-mode cni`. Ubuntu and Debian offer officially the package `containernetworking-plugins` for the standard CNI plugins. You either install and setup the appropriate networking services on the second node itself or create appropriate pods/daemonsets with either the finished prepared image or using cloud-init to install and setup those services. Using the former variant requires the pods to want host networking capability and for that LXE needs [a little file](../../fixtures/hostnetwork.conf) **persisted** for the `--hostnetwork-file` argument, e.g. save it to `/var/lib/lxe/hostnetwork.conf`. See `lxe --help` for more info and further configuration options.
