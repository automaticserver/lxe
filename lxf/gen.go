package lxf

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// lxdfakes
//counterfeiter:generate -o lxdfakes/fake_operation.go github.com/lxc/lxd/client.Operation
//counterfeiter:generate -o lxdfakes/fake_remote_operation.go github.com/lxc/lxd/client.RemoteOperation
//counterfeiter:generate -o lxdfakes/fake_server.go github.com/lxc/lxd/client.Server
//counterfeiter:generate -o lxdfakes/fake_image_server.go github.com/lxc/lxd/client.ImageServer
//counterfeiter:generate -o lxdfakes/fake_container_server.go github.com/lxc/lxd/client.ContainerServer
