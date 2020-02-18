FROM golang:1.13 as build-env
WORKDIR /go/src/app
ADD main.go /go/src/app
RUN go get -d -v ./...
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base
LABEL version="0.1.0"
LABEL repository="https://github.com/EnriqueTejeda/github-akamai-purge"
LABEL homepage="https://github.com/EnriqueTejeda/github-akamai-purge"
LABEL maintainer="Luis Enrique Tejeda Rodriguez"
COPY --from=build-env /go/bin/app /
CMD ["/app"]