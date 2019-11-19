
# LXE

<img src="fixtures/logo/logo_lxe_150.png" align="right" title="LXE Logo">

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

[![Build Status](https://img.shields.io/travis/automaticserver/lxe/master)](https://travis-ci.org/automaticserver/lxe)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/automaticserver/lxe)](https://github.com/automaticserver/lxe/releases)
[![GitHub](https://img.shields.io/github/license/automaticserver/lxe?color=lightgrey)](https://github.com/automaticserver/lxe/blob/master/COPYING)
[![Gitter](https://img.shields.io/gitter/room/automaticserver/lxe?color=blueviolet)](https://gitter.im/automaticserver-lxe)

LXE is a shim of the Kubernetes [Container Runtime Interface](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-node/container-runtime-interface.md) for LXD.
This project is currently under heavy development, expect incompatible changes.

## Requirements

You need to have LXD >= 3.3 installed, which packages are officially only available [via snap](https://linuxcontainers.org/lxd/getting-started-cli/#snap-package-archlinux-debian-fedora-opensuse-and-ubuntu). A LXD built by source is also sufficient.

## Installing LXE from packages

There are only official builds right now. Migration of our internal pipeline to a public location in progress.

## Installing LXE from source

LXE uses [Go Modules](https://github.com/golang/go/wiki/Modules) so the minimum Go version required is 1.11. Clone this repo to your wished location. If you checked it out within `$GOPATH` set `GO111MODULE=on`.

### Building

Build this project using the following command, which will give you the binary in `./bin/`

```bash
make build
```

## Getting started

### LXD prerequisites

Please follow these steps carefully. Some parameters and arguments depend on whether you installed lxd by source or via snap.

Make sure [that you have LXD running](https://github.com/lxc/lxd#machine-setup) and your default profile *only includes the root device and no interfaces*, since LXE organizes the networking and so interface names could interfere. Here's an example default profile:

```yaml
# lxc profile show default
config: {}
description: Default LXD profile
devices:
  root:
    path: /
    pool: default
    type: disk
name: default
used_by: []
```

Also make sure the LXD-client's remote configuration file exists (e.g. by running `lxc list` once), you'll need that later.

- if you built LXD by source, this file is located in `~/.config/lxc/config.yml` (LXE will guess this automatically by default)
- if you installed LXD via snap, the file is located in `~/snap/lxd/current/.config/lxc/config.yml`
- or you wrote that configration file on a location of your choice

LXE can be run as a non-privileged user, so give it [access to lxd's socket](https://linuxcontainers.org/lxd/getting-started-cli/#access-control).

### Running LXE

#### Parameters

The most important LXE options are the following:

```c
      --lxd-remote-config string    Path to the LXD remote config (guessed by default)
      --lxd-socket string           LXD's unix socket (default "/var/lib/lxd/unix.socket")
      --network-plugin string       The network plugin to use. '' is the standard network plugin and manages a lxd bridge 'lxebr0'. 'cni' uses kubernetes cni tools to attach interfaces.
      --socket string               The unix socket under which LXE will expose its service to Kubernetes (default "/var/run/lxe.sock")
```

You may need to provide the LXD socket path:

- if you built LXD by source, the socket is located in `/var/lib/lxd/unix.socket` (which is also default in LXE)
- if you installed LXD via snap, the socket is located in `/var/snap/lxd/common/lxd/unix.socket`

For all options, consider looking into `lxe --help`.

#### Starting the daemon

You might want to use `--verbose` for some feedback, otherwise the daemon is pretty silent when no errors occur. Warning: `--debug` is *very* verbose.

- if you built LXD by source, `lxe --verbose`
- if you installed LXD via snap, `lxe --lxd-socket /var/snap/lxd/common/lxd/unix.socket --lxd-remote-config ~/snap/lxd/current/.config/lxc/config.yml --verbose`

You should be greeted with:

```bash
INFO[10-03|19:02:07] Connected to LXD via "/var/lib/lxd/unix.socket" 
INFO[10-03|19:02:07] Starting streaming server on :44124 
INFO[10-03|19:02:07] Started LXE/0.1.21.gc4ee124.dirty CRI shim on UNIX socket "/var/run/lxe.sock" 
```

#### Configure Kubelet to use LXE

Now that you have LXE running on your system you can define the LXE socket as CRI endpoint in kubelet. You'll have to define the following options `--container-runtime=remote` and `--container-runtime-endpoint=unix:///var/run/lxe.sock` and your kubelet should be able to connect to your LXE socket.

## Bug reports

Bug reports can be filed at the [github issue tracker](https://github.com/automaticserver/lxe/issues/new)

## Contributing

Contribution guidelines are not yet defined.

## Documentation / FAQ

[A lot of options are missing and not yet implemented](doc/podspec-features.md) from the Kubernetes PodSpec.
Limitations and decisions of the current state are described in the [development preview FAQ](/doc/development-preview-faq.md).
