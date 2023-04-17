FROM golang:1.20-alpine as build

WORKDIR /go/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app ./cmd/api/main.go

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/src/app/config/config.yml /config/config.yml
COPY --from=build /go/src/app/data /data
COPY --from=build /go/bin/app /

CMD ["/app"]