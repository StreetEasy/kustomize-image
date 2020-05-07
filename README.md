# kustomize-image

```
docker build \
  --build-arg KUSTOMIZE_REPO=https://github.com/kubernetes-sigs/kustomize.git \
  --build-arg KUSTOMIZE_VERSION=v3.5.4 \
  -t zillownyc/kustomize:3.5.4 \
  -t zillownyc/kustomize:3.5 \
  -t zillownyc/kustomize:3 \
  -t latest .
```
