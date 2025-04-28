package third_party

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ../fakes/lxd/client/fake_operation.go github.com/canonical/lxd/client.Operation
//counterfeiter:generate -o ../fakes/lxd/client/fake_remote_operation.go github.com/canonical/lxd/client.RemoteOperation
//counterfeiter:generate -o ../fakes/lxd/client/fake_server.go github.com/canonical/lxd/client.Server
//counterfeiter:generate -o ../fakes/lxd/client/fake_image_server.go github.com/canonical/lxd/client.ImageServer
//counterfeiter:generate -o ../fakes/lxd/client/fake_instance_server.go github.com/canonical/lxd/client.InstanceServer

//counterfeiter:generate -o ../fakes/lxe/lxf/fake_client.go github.com/automaticserver/lxe/lxf.Client

//counterfeiter:generate -o ../fakes/containernetworking/libcni/fake_cni.go github.com/containernetworking/cni/libcni.CNI
