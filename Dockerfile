FROM golang:1.20-alpine as build

WORKDIR /go/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app ./cmd/api/main.go

FROM busybox AS busybox

ARG BUSYBOX_VERSION=1.31.0-i686-uclibc
ADD https://busybox.net/downloads/binaries/$BUSYBOX_VERSION/busybox_WGET /wget
RUN chmod a+x /wget

FROM gcr.io/distroless/static-debian11

COPY --from=busybox /wget /usr/bin/wget
COPY --from=build /go/src/app/config/config.yml /config/config.yml
COPY --from=build /go/src/app/data /data
COPY --from=build /go/bin/app /

CMD ["/app"]