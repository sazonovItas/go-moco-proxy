ARG GO_VERSION=1.22
FROM golang:${GO_VERSION} as build

WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  go mod download -x

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 GOARCH=${TARGETARCH} GOOS=${TARGETOS} go build \
  -o /bin/moco-proxy .

FROM alpine:3.20.0 AS development

COPY --from=build /bin/moco-proxy /bin/

ENTRYPOINT [ "/bin/moco-proxy" ]
