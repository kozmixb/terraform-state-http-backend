########## BUILDER ##########
FROM golang:1.22-alpine AS builder

WORKDIR /src
RUN apk add --no-cache build-base

# pre-copy/cache go.mod for pre-downloading dependencies and 
# only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build statically linked file and strip debug information.
# CGO is required by github.com/mattn/go-sqlite3.
RUN CGO_ENABLED=1 go build -tags "sqlite_omit_load_extension" -ldflags="-linkmode external -extldflags '-static' -s -w" -v -o app

########## RESULT ##########
FROM alpine:latest

COPY --from=builder /src/app /app

VOLUME [ "/storage" ]
EXPOSE 8080

ENTRYPOINT  ["/app"]
