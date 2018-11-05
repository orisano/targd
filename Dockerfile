FROM golang:1.11-alpine3.8 as vendor
WORKDIR /go/src/github.com/orisano/targd
RUN apk add --no-cache git
RUN wget -O- https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only && rm -rf /go/pkg/dep/sources

FROM golang:1.11-alpine3.8 as build
WORKDIR /go/src/github.com/orisano/targd
RUN apk add --no-cache gcc musl-dev
COPY --from=vendor /go/src/github.com/orisano/targd .
COPY . .
RUN go build -o /bin/targd

FROM alpine:3.8
COPY --from=build /bin/targd /bin/
ENTRYPOINT ["/bin/targd"]

