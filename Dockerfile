# bilicoin service Dockerfile
# version 1.1.1
# author r3inbowari
FROM golang:1.17.2 as builder
LABEL stage="builder"

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

# use netgo
ENV CGO_ENABLED=0

RUN chmod 777 build.sh
RUN  ./build.sh

RUN mv ./build/bilicoin_linux_amd64_v1.1.1 ./build/bilicoin
RUN chmod 777 ./build/bilicoin

FROM alpine

WORKDIR /app

COPY --from=builder /app/build/bilicoin .

ENTRYPOINT ["./bilicoin", "-a"]
