FROM golang:1.14-alpine3.11 as build
WORKDIR /go/src/github.com/orisano/targd
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN go build -o /bin/targd

FROM alpine:3.11
COPY --from=build /bin/targd /bin/
ENTRYPOINT ["/bin/targd"]

