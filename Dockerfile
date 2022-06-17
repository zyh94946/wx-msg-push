FROM alpine:3.16

WORKDIR /usr/src/app

COPY ./build/linux-amd64/wx-msg-push .

ENTRYPOINT ["./wx-msg-push"]