FROM --platform=amd64 docker.io/library/alpine:latest

WORKDIR /usr/local

RUN apk update && \
    apk add git && \
    wget -O golang.tar.gz \
        https://go.dev/dl/go1.23.1.linux-amd64.tar.gz && \
    tar -xzf golang.tar.gz && \
    rm -f golang.tar.gz && \
    git clone https://github.com/estoneman/crawly.git

WORKDIR /usr/local/crawly

RUN /usr/local/go/bin/go build -o crawly cmd/crawly/main.go

ENTRYPOINT [ "./crawly" ]
CMD [ "https://blog.boot.dev/", "20", "100" ]
