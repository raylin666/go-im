# API Interface Documentation (HTTP / WebSocket)

> 该项目为 IM 服务, HTTP 和 WebSocket 共用同个端口, 路由通过 Nginx 转发只对应的协议模块。

### Nginx 配置

```shell
# 处理 HTTP
upstream go-http
{
    server 127.0.0.1:10010 weight=1 max_fails=2 fail_timeout=10s;
    keepalive 16;
}

# 处理 WebSocket
upstream go-ws
{
    server 127.0.0.1:10010 weight=1 max_fails=2 fail_timeout=10s;
    keepalive 16;
}

server {
        listen       80;
        server_name  im.docker;

        location / {
            proxy_set_header Host $host;
            proxy_pass http://go-ws/ws;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
            proxy_set_header Connection "";
            proxy_redirect off;
            proxy_intercept_errors on;
            client_max_body_size 10m;
        }

        location /api
        {
            proxy_set_header Host $host;
            proxy_pass http://go-http/api;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            proxy_redirect off;
            proxy_intercept_errors on;
            client_max_body_size 30m;
        }
}
```

### HTTP 接口文档

> HTTP 地址: `http://im.docker/api`

### WebSocket 接口文档

> WebSocket 地址: `ws://im.docker`

###### 请求消息协议

> 格式参考示例: `{"seq": "2938372621", "event": "ping", "data": {}}`

| 字段值   | 字段类型      | 是否必须 | 字段描述                                               |
|-------|-----------|------|----------------------------------------------------|
| seq   | string    | 是    | 消息唯一ID                                             |
| event | string    | 是    | 消息事件, 该值请参考 [请求消息事件](#message_event_req) , 否则为无效事件 |
| data  | interface | 否    | JSON 数据包                                           |

###### 请求消息事件
<a id="message_event_req"></a>

| 事件名称 | 事件内容 | 事件描述                |
|------|------|---------------------|
| ping | 无    | 发送 Socket PING 心跳检测 |

###### 响应消息协议

> 格式参考示例: `{"seq":"2938372621","event":"ping","response":{"code":200,"message":"OK","data":"PONG"}}`

| 字段值      | 字段类型                                       | 是否必须 | 字段描述                                         |
|----------|--------------------------------------------|------|----------------------------------------------|
| seq      | string                                     | 是    | 消息唯一ID                                       |
| event    | string                                     | 是    | 消息事件, 该值请参考 [响应消息事件](#message_event_resp) 解析 |
| response | [MessageResponse](#struct_MessageResponse) | 是    | 消息内容                                         |

###### 响应消息事件
<a id="message_event_resp"></a>

| 事件名称 | 事件内容         | 事件描述                |
|------|--------------|---------------------|
| ping | 固定字符串 "PONE" | 正常 Socket PING 心跳状态 |


###### 响应状态码
<a id="message_responseCode"></a>

| CODE | MESSAGE |
|------|---------|
| 200  | 响应成功    |

### 数据结构

###### [WebSocket] 消息响应 MessageResponse
<a id="struct_MessageResponse"></a>

| 字段值     | 字段类型      | 是否必须 | 字段描述                                        |
|---------|-----------|------|---------------------------------------------|
| code    | uint32    | 是    | 响应状态码, 该值请参考 [响应状态码](#message_responseCode) |
| message | string    | 是    | 响应描述                                        |
| data    | interface | 是    | JSON 数据包                                    |


