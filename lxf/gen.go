package lxf // import "github.com/automaticserver/lxe/lxf"

// lxdfakes
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o lxdfakes/fake_operation.go github.com/lxc/lxd/client.Operation
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o lxdfakes/fake_remote_operation.go github.com/lxc/lxd/client.RemoteOperation
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o lxdfakes/fake_server.go github.com/lxc/lxd/client.Server
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o lxdfakes/fake_image_server.go github.com/lxc/lxd/client.ImageServer
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o lxdfakes/fake_container_server.go github.com/lxc/lxd/client.ContainerServer
