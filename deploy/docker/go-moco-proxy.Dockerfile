ARG GO_VERSION=1.22
FROM golang:${GO_VERSION} as build

WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  go mod download -x

ARG TARGETOS=linux
ARG TARGETARCH=amd64

ARG VERSION
ARG COMMIT
ARG DATE
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 GOARCH=${TARGETARCH} GOOS=${TARGETOS} go build \
  -ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.Commit=${COMMIT}' -X 'main.Date=${DATE}'" -o /bin/go-moco .

FROM build AS test

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  go test -v -count=1 ./... -race

FROM alpine:3.20.0 AS development

COPY --from=build /bin/go-moco /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/go-moco" ]

FROM alpine:3.20.0 AS release

RUN --mount=type=cache,target=/var/cache/apk/ \
  apk --update add \
  ca-certificates \
  tzdata \
  && \
  update-ca-certificates

ARG UID=10001
RUN adduser -H -D \
  --uid "${UID}" appuser
USER appuser

COPY --from=build /bin/go-moco /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/go-moco" ]
