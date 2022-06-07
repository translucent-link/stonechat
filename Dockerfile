FROM alpine:3
WORKDIR /
COPY stonechat.linux stonechat

CMD ["/stonechat"]

EXPOSE 8080
