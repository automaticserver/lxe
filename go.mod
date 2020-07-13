module github.com/automaticserver/lxe

go 1.14

require (
	emperror.dev/errors v0.7.0
	github.com/containernetworking/cni v0.8.0
	github.com/docker/docker v1.13.1
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20200710112657-153946a5f232 // indirect
	github.com/emicklei/go-restful v2.13.0+incompatible // indirect
	github.com/flosch/pongo2 v0.0.0-20200529170236-5abacdfa4915 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/golangci/golangci-lint v1.28.3
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/juju/errors v0.0.0-20200330140219-3fe23663418f
	github.com/juju/go4 v0.0.0-20160222163258-40d72ab9641a // indirect
	github.com/juju/loggo v0.0.0-20200526014432-9ce3a2e09b5e // indirect
	github.com/juju/persistent-cookiejar v0.0.0-20171026135701-d5e5a8405ef9 // indirect
	github.com/juju/testing v0.0.0-20200706033705-4c23f9c453cd // indirect
	github.com/juju/webbrowser v1.0.0 // indirect
	github.com/lxc/lxd v0.0.0-20200713101704-ccf7026d1616
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.3
	github.com/opencontainers/runtime-spec v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae
	google.golang.org/genproto v0.0.0-20200711021454-869866162049 // indirect
	google.golang.org/grpc v1.30.0
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/httprequest.v1 v1.2.1 // indirect
	gopkg.in/macaroon-bakery.v2 v2.2.0 // indirect
	gopkg.in/retry.v1 v1.0.3 // indirect
	gopkg.in/robfig/cron.v2 v2.0.0-20150107220207-be2e0b0deed5 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.5
	k8s.io/apiserver v0.18.5 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kubernetes v1.15.12
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19
)

replace k8s.io/api => k8s.io/api v0.15.12

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.15.12

replace k8s.io/apimachinery => k8s.io/apimachinery v0.15.13-beta.0

replace k8s.io/apiserver => k8s.io/apiserver v0.15.12

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.15.12

replace k8s.io/client-go => k8s.io/client-go v0.15.12

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.15.12

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.15.12

replace k8s.io/code-generator => k8s.io/code-generator v0.15.13-beta.0

replace k8s.io/component-base => k8s.io/component-base v0.15.12

replace k8s.io/cri-api => k8s.io/cri-api v0.15.13-beta.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.15.12

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.15.12

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.15.12

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.15.12

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.15.12

replace k8s.io/kubectl => k8s.io/kubectl v0.15.13-beta.0

replace k8s.io/kubelet => k8s.io/kubelet v0.15.12

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.15.12

replace k8s.io/metrics => k8s.io/metrics v0.15.12

replace k8s.io/node-api => k8s.io/node-api v0.15.12

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.15.12

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.15.12

replace k8s.io/sample-controller => k8s.io/sample-controller v0.15.12
