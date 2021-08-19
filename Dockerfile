FROM golang:1.17-alpine3.14 as build
WORKDIR /go/src/github.com/orisano/targd
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN go build -o /bin/targd

FROM alpine:3.14
COPY --from=build /bin/targd /bin/
ENTRYPOINT ["/bin/targd"]

