module github.com/StreetEasy/literalconfigmap

go 1.14

require (
	sigs.k8s.io/kustomize/api v0.3.2
	sigs.k8s.io/yaml v1.1.0
)

exclude (
	sigs.k8s.io/kustomize/api v0.2.0
)
