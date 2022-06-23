module github.com/ory/oathkeeper-maester

go 1.16

require (
	github.com/avast/retry-go v2.4.1+incompatible
	github.com/bitly/go-simplejson v0.5.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/go-logr/logr v0.4.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.10.2
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.12.2
	google.golang.org/protobuf => google.golang.org/protobuf v1.28.0
)
