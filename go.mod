module github.com/automaticserver/lxe

go 1.14

require (
	github.com/containernetworking/cni v0.8.0
	github.com/dionysius/errand v1.0.0
	github.com/docker/docker v1.13.1
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20200710112657-153946a5f232 // indirect
	github.com/flosch/pongo2 v0.0.0-20200529170236-5abacdfa4915 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golangci/golangci-lint v1.28.1
	github.com/golangci/revgrep v0.0.0-20180812185044-276a5c0a1039 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/juju/errors v0.0.0-20200330140219-3fe23663418f
	github.com/juju/go4 v0.0.0-20160222163258-40d72ab9641a // indirect
	github.com/juju/loggo v0.0.0-20200526014432-9ce3a2e09b5e // indirect
	github.com/juju/persistent-cookiejar v0.0.0-20171026135701-d5e5a8405ef9 // indirect
	github.com/juju/testing v0.0.0-20200706033705-4c23f9c453cd // indirect
	github.com/juju/webbrowser v1.0.0 // indirect
	github.com/lxc/lxd v0.0.0-20200825183131-2deb2bfbbce1
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.3
	github.com/mmcloughlin/professor v0.0.0-20170922221822-6b97112ab8b3
	github.com/opencontainers/runtime-spec v1.0.2
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/shurcooL/go v0.0.0-20191216061654-b114cc39af9f // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae
	google.golang.org/genproto v0.0.0-20200711021454-869866162049 // indirect
	google.golang.org/grpc v1.30.0
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/httprequest.v1 v1.2.1 // indirect
	gopkg.in/ini.v1 v1.52.0 // indirect
	gopkg.in/macaroon-bakery.v2 v2.2.0 // indirect
	gopkg.in/retry.v1 v1.0.3 // indirect
	gopkg.in/robfig/cron.v2 v2.0.0-20150107220207-be2e0b0deed5 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.16.15
	k8s.io/client-go v0.16.15
	k8s.io/cri-api v0.0.0
	k8s.io/kubernetes v1.18.19
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19
	mvdan.cc/unparam v0.0.0-20191111180625-960b1ec0f2c2 // indirect
	sourcegraph.com/sqs/pbtypes v1.0.0 // indirect
)

replace k8s.io/kubernetes => k8s.io/kubernetes v1.16.15

replace k8s.io/api => k8s.io/api v0.16.15

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.15

replace k8s.io/apimachinery => k8s.io/apimachinery v0.16.16-rc.0

replace k8s.io/apiserver => k8s.io/apiserver v0.16.15

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.15

replace k8s.io/client-go => k8s.io/client-go v0.16.15

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.15

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.15

replace k8s.io/code-generator => k8s.io/code-generator v0.16.16-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.16.15

replace k8s.io/cri-api => k8s.io/cri-api v0.16.16-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.15

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.15

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.15

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.15

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.15

replace k8s.io/kubectl => k8s.io/kubectl v0.16.15

replace k8s.io/kubelet => k8s.io/kubelet v0.16.15

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.15

replace k8s.io/metrics => k8s.io/metrics v0.16.15

replace k8s.io/node-api => k8s.io/node-api v0.16.15

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.15

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.16.15

replace k8s.io/sample-controller => k8s.io/sample-controller v0.16.15

replace vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
