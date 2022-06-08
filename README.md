
# LXE

<!-- markdownlint-disable-next-line MD033 -->
<img src="fixtures/logo/logo_lxe_150.png" align="right" title="LXE Logo">

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

[![Build Status](https://github.com/automaticserver/lxe/workflows/build/badge.svg?branch=master)](https://github.com/automaticserver/lxe/actions?query=workflow:build+branch:master)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/automaticserver/lxe)](https://github.com/automaticserver/lxe/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/automaticserver/lxe)](https://goreportcard.com/report/github.com/automaticserver/lxe)
[![GitHub](https://img.shields.io/github/license/automaticserver/lxe?color=lightgrey)](https://github.com/automaticserver/lxe/blob/master/COPYING)
[![Gitter](https://img.shields.io/gitter/room/automaticserver/lxe?color=blueviolet)](https://gitter.im/automaticserver-lxe)

LXE is a shim of the Kubernetes [Container Runtime Interface](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-node/container-runtime-interface.md) for LXD.
This project is currently under heavy development, expect incompatible changes.

## Requirements

You need to have LXD installed, which packages are officially only available [via snap](https://linuxcontainers.org/lxd/getting-started-cli/#snap-package-archlinux-debian-fedora-opensuse-and-ubuntu). Debian is working on a LXD [deb package](https://wiki.debian.org/LXD), other distros might be as well. A LXD built by source is also supported.

## Installing LXE from packages

There are only manual builds right now, see [releases page](https://github.com/automaticserver/lxe/releases) for available builds.

## Getting started

### LXD prerequisites

Please follow these steps carefully. Some parameters and arguments depend on how you installed LXD.

Make sure [that you have LXD running](https://github.com/lxc/lxd#machine-setup) and the LXD-client's remote configuration file exists (e.g. by running `lxc list` once), LXE will need that later.

### Running LXE

LXE can be run as a non-privileged user, so give it [access to lxd's socket](https://linuxcontainers.org/lxd/getting-started-cli/#access-control). When using the network-plugin cni root permissions are required.

#### Parameters

The most important LXE options are the following:

```cmd
      --lxd-remote-config string    Path to the LXD remote config. (guessed by default)
      --lxd-socket string           Path of the socket where LXD provides it's API. (guessed by default)
      --network-plugin string       The network plugin to use. 'bridge' manages the lxd bridge defined in --bridge-name. 'cni' uses kubernetes cni tools to attach interfaces using configuration defined in --cni-conf-dir (default "bridge")
      --socket string               Path of the socket where it should provide the runtime and image service to kubelet. (default "/run/lxe.sock")
```

We recommend to use CNI as the network plugin as it offers more flexibility and integration to [common kubernetes network setups](https://kubernetes.io/docs/concepts/cluster-administration/networking/). But for sure you can use the currently default network plugin, which uses lxd's integrated networking, and build kubernetes cluster networking around it.

The CNI plugin is selected by passing the `--network-plugin=cni` option. The CNI configuration is read from within `--cni-conf-dir` (default /etc/cni/net.d) and uses that file to set up each pod’s network. The CNI configuration file must match the [CNI specification](https://github.com/containernetworking/cni/blob/master/SPEC.md#network-configuration), and any required CNI plugins referenced by the configuration must be present in `--cni-bin-dir` (default /opt/cni/bin).

If there are multiple CNI configuration files in the directory, the first configuration file by name in lexicographic order is used. Keep in mind you can also chain several plugins using a conflist file. Example configuration `/etc/cni/net.d/10-mynet.conf`:

```json
{
  "cniVersion": "0.3.1",
  "name": "mynet",
  "type": "bridge",
  "bridge": "cni0",
  "isGateway": true,
  "ipMasq": true,
  "ipam": {
    "type": "host-local",
    "ranges": [
      [
        {
          "subnet": "10.22.0.0/16",
          "rangeStart": "10.22.0.50",
          "rangeEnd": "10.22.0.100",
          "gateway": "10.22.0.1"
        }
      ]
    ],
    "routes": [
      {
        "dst": "0.0.0.0/0"
      }
    ]
  }
}
```

For all options, consider looking into `lxe --help`.

#### Starting the daemon

You might want to use `--log-level info` for some feedback, otherwise the daemon is pretty silent when no warnings or errors occur:

`lxe --network-plugin cni --log-level info`

You should be greeted with:

```bash
WARNING[07-20|17:11:18.042] starting lxe...                               packagename=lxe version=0.0.0 gitcommit=... gittreestate=dirty buildnumber=undef builddate="..."
INFO   [07-20|17:11:18.051] Connected to LXD                              lxd-socket=/var/lib/lxd/unix.socket
INFO   [07-20|17:11:18.068] Started lxe CRI shim                          socket=/run/lxe.sock
INFO   [07-20|17:11:18.068] Started streaming server                      endpoint=":44124" baseurl="http://10.249.100.169:44124"
```

#### Configuration options

Aside from command-line parameters, LXE also supports configuration files in various file formats (YAML, JSON, TOML, etc.) and environment variables. All parameters are automatically detected depending on rules for each configuaration type:

- Environment variables are always upper case and the dashes are underscores: E.g. `--lxd-socket` is `LXD_SOCKET`
- The keys in config files are split on every dash and the remaining part of the key continues as sub-object: E.g. in JSON `--lxd-socket` is `{"lxd": {"socket": ...}}`

You can also print the currently loaded configuration with `lxd config show yaml` (or any other supported extension). _Right now this only shows what configuration would be loaded if you would run lxe like that, it does not retrieve the current running configuration in case you have already an lxe process running_!

```toml
$ lxe config show toml
config = ""
socket = "/run/lxe.sock"

[bridge]
  name = "lxdbr0"

  [bridge.dhcp]
    range = ""

...
```

If no config argument is provided via command-line or environment, it will automatically look for a file `lxe.<ext>` in `~/.local/lxe/` and `/etc/lxe/` (The folder `lxe` can be customized with the make variable `PACKAGE_NAME`)

You can also combine all these variants. Command-line parameters have precedence over environment variables, which have precedence over configuration file settings, those in turn have precedence over defaults. Please be aware you can't set the config variable in a config file, it has no effect. Once a variable is set in any way, even if empty, the default is overridden.

### Configure Kubelet to use LXE

Now that you have LXE running on your system you can define the LXE socket as CRI endpoint in kubelet. You'll have to define the following options `--container-runtime=remote` and `--container-runtime-endpoint=unix:///run/lxe.sock` and your kubelet should be able to connect to your LXE socket.

## Installing LXE from source

LXE requires golang 1.18.

### Quick Installation

Simply run:

```bash
go install github.com/automaticserver/lxe/cmd/lxe@latest
```

And the binary will be located in `$GOPATH/bin`

### Building & Tests

Clone this repo to your wished location. Build this project using the following command, which will give you the binary in `./bin/`

```bash
make build
```

You can also run the program directly using:

```bash
make run lxe -- --log-level info --network-plugin cni
```

There are also tests and lints available.

```bash
make test
make lint
```

For testing with [`critest` from cri-tools](https://github.com/kubernetes-sigs/cri-tools/) have a look in the [testing documentation](doc/testing.md).

To list all available make targets have a look at help target `make help`.

## Examples

Please have a look at the [multi node cluster guide](doc/examples/multi-node-cluster-with-k3s.md)

## Bug reports

Bug reports can be filed at the [github issue tracker](https://github.com/automaticserver/lxe/issues/new)

## Documentation / FAQ

[A lot of options are missing and not yet implemented](doc/podspec-features.md) from the Kubernetes PodSpec.
Limitations and decisions of the current state are described in the [development preview FAQ](/doc/development-preview-faq.md).
