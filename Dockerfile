# syntax=docker/dockerfile:experimental

FROM golang:1.14.2

ARG KUSTOMIZE_REPO=https://github.com/kubernetes-sigs/kustomize.git
ARG KUSTOMIZE_VERSION

RUN git clone ${KUSTOMIZE_REPO} /go/src/sigs.k8s.io/kustomize \
  && cd /go/src/sigs.k8s.io/kustomize \
  && git checkout kustomize/${KUSTOMIZE_VERSION}

WORKDIR /go/src/sigs.k8s.io/kustomize/kustomize

ENV GO_PKG=sigs.k8s.io/kustomize/api/provenance
RUN go build -o /usr/local/bin/kustomize \
  -ldflags "-s \
    -X ${GO_PKG}.version=${KUSTOMIZE_VERSION} \
    -X ${GO_PKG}.gitCommit=`git rev-parse HEAD` \
    -X ${GO_PKG}.buildDate=`date -u +'%Y-%m-%dT%H:%M:%SZ'`"

COPY plugin/devops.streeteasy.com /go/src/sigs.k8s.io/kustomize/plugin/devops.streeteasy.com

RUN cd /go/src/sigs.k8s.io/kustomize/plugin/devops.streeteasy.com/v1/envkeysecret \
  && go build -buildmode=plugin -o $HOME/.config/kustomize/plugin/devops.streeteasy.com/v1/envkeysecret/EnvKeySecret.so

RUN cd /go/src/sigs.k8s.io/kustomize/plugin/devops.streeteasy.com/v1/envkeyconfigmap \
  && go build -buildmode=plugin -o $HOME/.config/kustomize/plugin/devops.streeteasy.com/v1/envkeyconfigmap/EnvKeyConfigMap.so

RUN cd /go/src/sigs.k8s.io/kustomize/plugin/devops.streeteasy.com/v1/literalsecret \
  && go build -buildmode=plugin -o $HOME/.config/kustomize/plugin/devops.streeteasy.com/v1/literalsecret/LiteralSecret.so

RUN cd /go/src/sigs.k8s.io/kustomize/plugin/devops.streeteasy.com/v1/literalconfigmap \
  && go build -buildmode=plugin -o $HOME/.config/kustomize/plugin/devops.streeteasy.com/v1/literalconfigmap/LiteralConfigMap.so

ENTRYPOINT ["/usr/local/bin/kustomize"]
