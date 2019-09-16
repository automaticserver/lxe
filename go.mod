module github.com/lxc/lxe

require (
	code.cloudfoundry.org/systemcerts v0.0.0-20180917154049-ca00b2f806f2 // indirect
	github.com/containernetworking/cni v0.7.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/docker v1.13.1
	// Repo does not have any tags at all
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20181111060418-2ce16c963a8a // indirect
	github.com/emicklei/go-restful v2.8.0+incompatible // indirect
	github.com/flosch/pongo2 v0.0.0-20180809100617-24195e6d38b0 // indirect
	github.com/frankban/quicktest v1.1.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.2.0 // indirect
	// Repo does not have any tags at all
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/juju/errors v0.0.0-20181118221551-089d3ea4e4d5 // indirect
	github.com/juju/go4 v0.0.0-20160222163258-40d72ab9641a // indirect
	github.com/juju/persistent-cookiejar v0.0.0-20171026135701-d5e5a8405ef9 // indirect
	github.com/juju/retry v0.0.0-20180821225755-9058e192b216 // indirect
	github.com/juju/testing v0.0.0-20180920084828-472a3e8b2073 // indirect
	github.com/juju/utils v0.0.0-20180820210520-bf9cc5bdd62d // indirect
	github.com/juju/webbrowser v0.0.0-20180907093207-efb9432b2bcb // indirect
	github.com/julienschmidt/httprouter v1.2.0 // indirect
	// Last commit of LXD v3.3
	github.com/lxc/lxd v0.0.0-20181220183431-fba7538d485e
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	// Repo do not use SemVer (Tags are "1.0.3", not "v1.0.3"). This commit is tagged with 1.0.3
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	// Only commit in this repo
	github.com/rogpeppe/fastuuid v0.0.0-20150106093220-6724a57986af // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.2.2
	// Do not have SemVer tags - latest commits at time of writing
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9 // indirect
	golang.org/x/net v0.0.0-20181220203305-927f97764cc3
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890 // indirect
	golang.org/x/sync v0.0.0-20181108010431-42b317875d0f // indirect
	golang.org/x/sys v0.0.0-20181221135038-a79f1b190785 // indirect
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/appengine v1.4.0 // indirect
	// Does not have SemVer tags - latest commit at time of writing
	google.golang.org/genproto v0.0.0-20181221010529-a1fde7408246 // indirect
	google.golang.org/grpc v1.17.0
	gopkg.in/errgo.v1 v1.0.0 // indirect
	gopkg.in/httprequest.v1 v1.1.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/macaroon-bakery.v2 v2.1.0 // indirect
	gopkg.in/macaroon.v2 v2.0.0 // indirect
	gopkg.in/retry.v1 v1.0.2 // indirect
	gopkg.in/robfig/cron.v2 v2.0.0-20150107220207-be2e0b0deed5 // indirect
	gopkg.in/yaml.v2 v2.2.2
	// These imports do not use SemVer (Tags are "kubernetes-1.12.1")
	k8s.io/api v0.0.0-20181130031204-d04500c8c3dd // indirect
	k8s.io/apimachinery v0.0.0-20181220065808-98853ca904e8
	k8s.io/apiserver v0.0.0-20181220070914-ce7b605bead3 // indirect
	// We need to use the master branch as long as commit a6d1c60475b25ad is not in a tag
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v0.1.0 // indirect
	k8s.io/kubernetes v1.14.1
	// Repo does not have any tags at all
	k8s.io/utils v0.0.0-20181115163542-0d26856f57b3 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)
