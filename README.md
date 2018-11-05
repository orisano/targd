# targd
targd is multistage build helper. targd outputs dependent stage of the specified target stage. (includes target stage)

## Installation
```bash
go get github.com/orisano/targd
```

## How to use
```bash
$ targd
  -f string
    	Dockerfile path (default "Dockerfile")
  -target string
    	target stage name (required)
```
```bash
$ targd -target test | docker build -f - .
```
or
```bash
$ docker run -v $(PWD):/host orisano/targd -f /host/Dockerfile -target production
```

## Example
```dockerfile
# ./Dockerfile
FROM alpine:3.8 as builder
COPY ./depfile ./lockfile ./
RUN install_cmd
COPY . .
RUN build_cmd

FROM builder as test
RUN test_cmd

FROM alpine:3.8 as runtime
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
CMD ["--help"]
```
```bash
$ targd -target test
FROM alpine:3.8 as builder
COPY ./depfile ./lockfile ./
RUN install_cmd
COPY . .
RUN build_cmd

FROM builder as test
RUN test_cmd

```
```bash
$ targd -target runtime
FROM alpine:3.8 as builder
COPY ./depfile ./lockfile ./
RUN install_cmd
COPY . .
RUN build_cmd

FROM alpine:3.8 as runtime
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
CMD ["--help"]

```

## Author
Nao YONASHIRO (@orisano)

## License
MIT
