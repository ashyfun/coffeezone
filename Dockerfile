FROM golang:1.20-alpine3.18

RUN apk add --update --no-cache chromium curl tar

WORKDIR /tmp

ARG MIGRATE_VERSION=4.16.2
ARG MIGRATE_OS=linux
ARG MIGRATE_ARCH=amd64

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v$MIGRATE_VERSION/migrate.$MIGRATE_OS-$MIGRATE_ARCH.tar.gz | \
    tar xvz && \
    mv -f ./migrate /usr/bin

WORKDIR /usr/src/coffeezone

COPY . .

RUN go build -o ./build/scraper ./scraper && \
    go build -o ./build/rest ./rest && \
    for bin in ./build/*; do ln -s "$(realpath $bin)" "/usr/bin/${bin##*/}"; done
