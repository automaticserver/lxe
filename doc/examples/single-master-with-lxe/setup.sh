#!/bin/bash

set -ex

#cluster_domain="cluster.local"
#cluster_servicesubnet="10.96.0.0/12"
#cluster_podsubnet="10.244.0.0/16"

# remove existing lxd and lxc packages
apt-get purge lxd* lxc* liblxc* -y

# install lxd via snap
snap install lxd
cat preseed-lxd-init.yaml | lxd init --preseed

# enable kubernetes apt
apt-get install -y apt-transport-https curl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
apt-get update

# install lxe via deb using latest github release
apt-get install kubernetes-cni socat ebtables -y
rm lxe*.debian-lxd-snap.deb
wget https://github.com/automaticserver/lxe/releases/download/v0.1.2/lxe_0.1.2.debian-lxd-snap.deb
dpkg -i lxe*.debian-lxd-snap.deb
systemctl status lxe | cat

# copy various predefined files
apt-get install facter -y
cp -r files/* /
sed -i "s@\[\[ipaddress\]\]@$(facter ipaddress_ens3)@g" /etc/kubernetes/kubeadm.conf

# install cri-tools and configuration for lxe
apt-get install cri-tools
crictl version

# install and configure kubelet and create certificates and credentials
apt-get install kubeadm kubelet kubectl -o Dpkg::Options::="--force-confold" --force-yes -y
kubeadm alpha phase certs all --config /etc/kubernetes/kubeadm.conf
kubeadm alpha phase kubeconfig all --config /etc/kubernetes/kubeadm.conf
mkdir -p ~/.kube
cp -n /etc/kubernetes/admin.conf ~/.kube/config

modprobe ip_vs ip_vs_rr ip_vs_wrr ip_vs_sh nf_conntrack_ipv4
modprobe br_netfilter
sysctl net.bridge.bridge-nf-call-iptables=1
kubectl completion bash > /etc/bash_completion.d/kubectl
systemctl daemon-reload
systemctl enable kubelet.service

# create pod manifests for various required containers for kubernetes (networking, proxy, dns, etc.)
# TODO

# simple pod for access
#kubectl apply -f pods/nginx.yaml