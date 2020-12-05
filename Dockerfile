FROM golang:1.15-alpine3.12 as build
WORKDIR /go/src/github.com/orisano/targd
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN go build -o /bin/targd

FROM alpine:3.12
COPY --from=build /bin/targd /bin/
ENTRYPOINT ["/bin/targd"]

