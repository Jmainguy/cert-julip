module github.com/jmainguy/cert-julip

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.9 // indirect
	github.com/Azure/go-autorest v11.5.2+incompatible // indirect
	github.com/coreos/prometheus-operator v0.26.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-openapi/spec v0.18.0 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gophercloud/gophercloud v0.0.0-20190318015731-ff9851476e98 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/jetstack/cert-manager v0.8.1
	github.com/openshift/api v3.9.0+incompatible
	github.com/operator-framework/operator-sdk v0.8.2-0.20190522220659-031d71ef8154
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.27.1
	k8s.io/apimachinery v0.28.0-alpha.0
	k8s.io/client-go v2.0.0+incompatible
	k8s.io/code-generator v0.27.1
	k8s.io/gengo v0.0.0-20230306165830-ab3349d207d4
	k8s.io/kube-openapi v0.0.0-20230327201221-f5883ff37f0c
	sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools v0.11.3
)

// Pinned to kubernetes-1.13.1
replace (
	k8s.io/api => k8s.io/api v0.27.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.27.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.27.1
	k8s.io/client-go => k8s.io/client-go v0.27.1
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.8.1
	k8s.io/code-generator => k8s.io/code-generator v0.27.1
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20230327201221-f5883ff37f0c
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.11.3
)

go 1.13
