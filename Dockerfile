########## BUILDER ##########
FROM --platform=$BUILDPLATFORM golang:1.26.4-alpine3.23 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# pre-copy/cache go.mod for pre-downloading dependencies and 
# only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN mkdir -p /tmp/storage \
	&& CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-s -w" -trimpath -v -o /tmp/terraform-state-http-backend

########## RESULT ##########
FROM scratch

COPY --from=builder --chown=65532:65532 /tmp/terraform-state-http-backend /app
COPY --from=builder --chown=65532:65532 /tmp/storage /storage

VOLUME [ "/storage" ]
EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/app"]
