module github.com/mjm/pi-tools

go 1.14

require (
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/etherlabsio/healthcheck v0.0.0-20191224061800-dd3d2fd8c3f6
	github.com/golang-migrate/migrate/v4 v4.13.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.3
	github.com/google/go-github/v33 v33.0.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hako/durafmt v0.0.0-20200710122514-c0fb7b4da026
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/jonboulle/clockwork v0.2.2
	github.com/lib/pq v1.8.0
	github.com/mdlayher/unifi v0.0.0-20180604180305-be73af8d7922
	github.com/mjm/graphql-go v0.0.0-20200213070811-48021cabfccf
	github.com/rs/cors v1.7.0 // indirect
	github.com/segmentio/ksuid v1.0.3
	github.com/stretchr/testify v1.6.1
	github.com/zserge/hid v0.0.0-20190124175232-e1626f1782f3
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.14.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.14.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.14.0
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.14.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.14.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
	k8s.io/api v0.19.1
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v0.19.1
	k8s.io/kubernetes v1.19.4
	sigs.k8s.io/sig-storage-lib-external-provisioner/v6 v6.1.0
	zombiezen.com/go/postgrestest v1.0.0
)

replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/apiserver => k8s.io/apiserver v0.19.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.0
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.0
	k8s.io/code-generator => k8s.io/code-generator v0.19.0
	k8s.io/component-base => k8s.io/component-base v0.19.0
	k8s.io/cri-api => k8s.io/cri-api v0.19.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.0
	k8s.io/kubectl => k8s.io/kubectl v0.19.0
	k8s.io/kubelet => k8s.io/kubelet v0.19.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.0
	k8s.io/metrics => k8s.io/metrics v0.19.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.0
)
