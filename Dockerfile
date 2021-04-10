FROM golang:1.16-alpine3.13 as build
WORKDIR /go/src/github.com/orisano/targd
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN go build -o /bin/targd

FROM alpine:3.13
COPY --from=build /bin/targd /bin/
ENTRYPOINT ["/bin/targd"]

