########## BUILDER ##########
FROM golang:1.26.4-alpine3.23 AS builder

WORKDIR /src
RUN apk add --no-cache build-base

# pre-copy/cache go.mod for pre-downloading dependencies and 
# only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build statically linked file and strip debug information.
# CGO is required by github.com/mattn/go-sqlite3.
RUN CGO_ENABLED=1 go build -tags "sqlite_omit_load_extension" -ldflags="-linkmode external -extldflags '-static' -s -w" -v -o /tmp/terraform-state-http-backend

########## RESULT ##########
FROM alpine:3.23

COPY --from=builder /tmp/terraform-state-http-backend /app

RUN addgroup -S app \
	&& adduser -S -D -H -h /nonexistent -s /sbin/nologin -G app app \
	&& mkdir -p /storage \
	&& chown app:app /storage \
	&& chmod 0750 /storage

VOLUME [ "/storage" ]
EXPOSE 8080

USER app:app

ENTRYPOINT ["/app"]
