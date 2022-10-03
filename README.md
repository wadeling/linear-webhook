# Linear webhook

this is a linear webhook service which could recv linear webhook msg and notify to feishu

## Features

- support deliver txt msg to feishu
- support recv multi linear-workspaces webhook
  - linear workspaces should user same webhook url
  - config multi workspace api key in config.yaml
- limit: currently only supports the use of mobile phone number to at someone in Feishu messages

## Usage

### build image
```asciidoc
./docker_build.sh
```

### run webhook service

- run as http service
```asciidoc
docker run -d --name linear-webhook -p51723:6500 \
-v $(pwd)/config/staff.json:/root/staff.json \
-v $(pwd)/config/config.yaml:/root/config.yaml \
linear/linear-webhook:latest -config=/root/config.yaml
```

- run as https service
```asciidoc
docker run -d --name linear-webhook -p51723:6500 \
-v $(pwd)/cert/my.crt:/root/cert/my.crt -v $(pwd)/cert/my.key:/root/cert/my.key \
-v $(pwd)/config/staff.json:/root/staff.json \
-v $(pwd)/config/config.yaml:/root/config.yaml \
linear/linear-webhook:latest -config=/root/config.yaml
```
config.yaml parameter:
```yaml
server:
  listen_addr: ":6500"
https:
  certificate: "/root/cert/my.crt"
  private_key: "/root/cert/my.key"
linear:
  api_addr: "https://api.linear.app/graphql"
  api_keys:
    -
      #linear workspace name
      workspace: workspace-bugs
      #api key in this workspace
      api_key: lin_api_xxx
    -
      workspace: workspace-test
      api_key: lin_api_xxx
feishu:
  #feishu robot webhook addr
  webhook_url: "https://open.feishu.cn/open-apis/bot/v2/hook/bcxxxx"
  #feishu robot app id
  app_id: cli_xxx
  #feishu robot app secret
  app_secret: 7nzxxxG
  #staff info which contain username and mobile number
  staff_file: /root/staff.json
```

staff_file format：
```json
[
    {"user_name":"tom","mobile":"123456"},
    {"user_name":"jerry","mobile":"123456"}
]
```
