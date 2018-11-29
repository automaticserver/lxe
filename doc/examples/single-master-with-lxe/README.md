# Single Master with LXE

This example provides the installation of a Kubernetes master with LXD as runtime using LXE. You need a fresh ubuntu xenial server and the script might not be failure proof. The installation is similiar what `kubeadm` would do. Only try this in a test environment! It will:

- remove existing lxd and lxc packages
- install lxd via snap
- enable kubernetes apt
- install lxe via deb using latest github release
- install cri-tools and configuration for lxe
- install and configure kubelet and create certificates and credentials
- create pod manifests for various required containers for kubernetes (networking, proxy, dns, etc.)

To run the script in one command:

```
git clone https://github.com/automaticserver/lxe.git; cd lxe/doc/examples/single-master-with-lxe/; ./setup.sh
```