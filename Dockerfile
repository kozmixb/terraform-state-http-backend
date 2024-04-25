########## BUILDER ##########
FROM golang:1.22 as builder

COPY . /src
WORKDIR /src

# pre-copy/cache go.mod for pre-downloading dependencies and 
# only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download && go mod verify

# Build statically linked file and strip debug information
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static -s -w" -v -o app

########## RESULT ##########
FROM alpine:latest

COPY --from=builder /src/app /app

VOLUME [ "/storage" ]
EXPOSE 8080

ENTRYPOINT  ["/app"]
