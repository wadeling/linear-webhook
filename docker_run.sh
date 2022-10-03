#!/bin/bash

notify="https://open.feishu.cn/open-apis/bot/v2/hook/bc848220-7379-4ce3-b95e-236cffa94dec"
certPath="/root/go/src/github.com/wadeling/linear-webhook/cert"

docker stop linear-webhook
docker rm linear-webhook

docker run -d --name linear-webhook -p6500:6500 wade23/linear-webhook:latest -listen=":6500" --notify-url=$notify
#docker run -d --name linear-webhook -p6500:6500 -v $(pwd)/cert/my.crt:/root/cert/my.crt -v $(pwd)/cert/my.key:/root/cert/my.key \
#	wade23/linear-webhook:latest -listen=":6500" --notify-url=$notify -linear-addr="" -linear-api-key="" -tls=true -cert-file=/root/cert/my.crt -key-file=/root/cert/my.key
