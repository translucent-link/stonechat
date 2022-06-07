#! /bin/sh
GOOS=linux GOARCH=amd64 go build -o stonechat.linux . 
docker build --platform=linux/amd64 -t translucentlink/stonechat:$1 .
docker push translucentlink/stonechat:$1