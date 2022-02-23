package network

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ./libcnifake/cni.go github.com/containernetworking/cni/libcni.CNI
