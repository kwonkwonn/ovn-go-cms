FROM golang:1.24-alpine AS goinit

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /src
RUN git clone https://github.com/kwonkwonn/ovn-go-cms.git

WORKDIR /src/ovn-go-cms
RUN go build  -o ovn-go-cms .

FROM gcr.io/distroless/static:nonroot
COPY --from=goinit /src/ovn-go-cms/ovn-go-cms /usr/local/bin/ovn-go-cms

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/ovn-go-cms"]
