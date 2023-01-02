FROM golang:1.19 as build

WORKDIR /go/src/app
COPY pkg/ pkg/
COPY vendor/ vendor/
COPY go.* .
COPY main.go main.go

# RUN go test

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]