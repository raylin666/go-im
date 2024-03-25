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

        # WebSocket 连接路由
	    location / {
            proxy_set_header Host $host;
            proxy_pass http://go-ws/app/ws;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
            proxy_set_header Connection "";
            proxy_redirect off;
            proxy_intercept_errors on;
            client_max_body_size 10m;
        }

        # HTTP Web API 路由
        location /app {
            proxy_set_header Host $host;
            proxy_pass http://go-http/app/api;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
            proxy_redirect off;
            proxy_intercept_errors on;
            client_max_body_size 30m;
        }

        # HTTP RPC API 路由
        location /api {
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

### 公共请求参数

##### Header 鉴权 (非必须 - 需要鉴权的接口文档会说明)

<a id="http_header_auth"></a>

| 参数名称    | 是否必须 | 示例                                                                                                                               | 说明                       |
|---------|------|----------------------------------------------------------------------------------------------------------------------------------|--------------------------|
| key     | 是    | 2475506427                                                                                                                       | 应用ID                     |
| secret  | 是    | 938d3afb72664b1b8652cd71c2a4d0f0                                                                                                 | 应用密钥                     |
| usersig | 是    | baafe63b4c11e0ec5ede99efd0df8795d0d26a9af4d1bebcc7671450d56a9124cff0f9c12a4200dc3780ff235f9e4a91d94cfcd87b3fa0937e43ba4766d5b0c5 | 授权账号用户签名TOKEN (应用下的账号体系) |

### HTTP Web API 接口文档

<i> HTTP 地址: `http://im.docker/app` </i>

### HTTP RPC API 服务接口文档

<i> HTTP 地址: `http://im.docker/api` </i>

##### 基础模块

###### 项目心跳接口

> 请求方式: `GET`
>
> 请求URL: `/api/heartbeat`

响应内容

| 名称      | 类型     | 是否必须 | 示例   | 描述                    |
|---------|--------|------|------|-----------------------|
| message | string | 是    | PONE | 固定 "PONG" 值, 表示项目正常运行 |

##### 应用管理模块

###### 创建应用接口

> 请求方式: `POST`
>
> 请求URL: `/api/manager/create`

请求参数

| 名称         | 类型        | 是否必须 | 示例                   | 校验规则                      | 描述           |
|------------|-----------|------|----------------------|---------------------------|--------------|
| ident      | string    | 是    | raylin666            | 6-50字以内,必须是字母开头,由字母数字和.组成 | 唯一标识, 用来标识来源 |
| name       | string    | 是    | 正式环境                 | 2-30字以内                   | 应用名称         |
| expired_at | timestamp | 是    | 2099-12-31T23:59:59Z | 必须大于当前时间                  | 过期时间         |
| status     | enum      | 否    | 1                    | 必须为 0,1,2 数字              | 应用状态         |

响应内容

| 名称         | 类型     | 是否必须 | 示例                               | 描述                              |
|------------|--------|------|----------------------------------|---------------------------------|
| id         | string | 是    | 1                                | 自增ID(无实质业务性质)                   |
| ident      | string | 是    | raylin666                        | 唯一标识, 用来标识来源                    |
| name       | string | 是    | 正式环境                             | 应用名称                            |
| key        | string | 是    | 403227602                        | 应用KEY                           |
| secret     | string | 是    | 1808c3d2a764499eb2924e70731f76d5 | 应用密钥                            |
| status     | string | 是    | OPEN                             | 应用状态 CLOSE:停用 OPEN:启用 FREEZE:冻结 |
| expired_at | string | 是    | 2099-12-11T23:59:59Z             | 过期时间                            |
| created_at | string | 是    | 2024-03-17T13:27:02.384994424Z   | 创建时间                            |

##### 账号管理模块

###### 创建账号接口

必须 Header 鉴权登录, 参数请参考 [Header 鉴权](#http_header_auth)

> 请求方式: `POST`
>
> 请求URL: `/api/account/create`

请求参数

| 名称       | 类型     | 是否必须 | 示例                                                                                | 校验规则               | 描述                 |
|----------|--------|------|-----------------------------------------------------------------------------------|--------------------|--------------------|
| user_id  | string | 是    | 91283746167                                                                       | 1-30字以内,只能是字母或数字组成 | 唯一标识, 用户ID         |
| username | string | 是    | 小明                                                                                | 1-30字以内            | 用户名称               |
| avatar   | string | 是    | https://inews.gtimg.com/om_bt/OUlYolgW1-7XmNSQchf24_Vqcs5KWNun2D4ospT3Wn8xYAA/641 | 必须是有效URL地址         | 用户头像               |
| is_admin | bool   | 否    | false                                                                             |                    | 是否管理员(后端推送消息必须该权限) |

响应内容

| 名称         | 类型     | 是否必须 | 示例                                                                                | 描述                 |
|------------|--------|------|-----------------------------------------------------------------------------------|--------------------|
| user_id    | string | 是    | 91283746167                                                                       | 唯一标识, 用户ID         |
| username   | string | 是    | 小明                                                                                | 用户名称               |
| avatar     | string | 是    | https://inews.gtimg.com/om_bt/OUlYolgW1-7XmNSQchf24_Vqcs5KWNun2D4ospT3Wn8xYAA/641 | 用户头像               |
| is_admin   | bool   | 是    | false                                                                             | 是否管理员(后端推送消息必须该权限) |
| created_at | string | 是    | 2024-03-17T13:27:02.384994424Z                                                    | 创建时间               |

### WebSocket 接口文档

<i> WebSocket 地址 (key、secret 为应用授权值): `ws://im.docker?key=2465562260&secret=92f57e94c7af48b1af8980ef2b843b5b` </i>

###### 请求消息协议

> 格式参考示例: `{"seq": "2938372621", "event": "ping", "data": {}}`

| 字段值   | 字段类型      | 是否必须 | 字段描述                                                              |
|-------|-----------|------|-------------------------------------------------------------------|
| seq   | string    | 是    | 消息唯一ID                                                            |
| event | string    | 是    | 消息事件, 该值请参考 [请求消息事件](#message_event_req) , 否则为无效事件                |
| data  | interface | 否    | JSON 数据包, 不同事件设定不同的数据包, 具体参考 [请求消息事件](#message_event_req) 的事件数据结构 |

###### 请求消息事件

<a id="message_event_req"></a>

| 事件名称        | 事件数据结构                                                 | 事件描述                   |
|-------------|--------------------------------------------------------|------------------------|
| ping        | 无                                                      | 发送 Socket PING 心跳检测    |
| login       | [WebSocketLoginRequest](#struct_WebSocketLoginRequest) | 账号登录 (必须登录完成后才能进行用户事件) |
| logout      | 无                                                      | 账号登出                   |
| loginStatus | 无                                                      | 获取账号登录状态               |

###### 响应消息协议

> 格式参考示例: `{"seq":"2938372621","event":"ping","response":{"code":200,"message":"OK","data":"PONG"}}`

| 字段值      | 字段类型                                                         | 是否必须 | 字段描述                                         |
|----------|--------------------------------------------------------------|------|----------------------------------------------|
| seq      | string                                                       | 是    | 消息唯一ID                                       |
| event    | string                                                       | 是    | 消息事件, 该值请参考 [响应消息事件](#message_event_resp) 解析 |
| response | [WebSocketMessageResponse](#struct_WebSocketMessageResponse) | 是    | 消息内容                                         |

###### 响应消息事件

<a id="message_event_resp"></a>

| 事件名称        | 事件内容                                                                 | 事件描述                |
|-------------|----------------------------------------------------------------------|---------------------|
| ping        | 固定字符串 "PONE"                                                         | 正常 Socket PING 心跳状态 |
| login       | [WebSocketLoginResponse](#struct_WebSocketLoginResponse)             | 账号登录                |
| logout      | [WebSocketLogoutResponse](#struct_WebSocketLogoutResponse)            | 账号登出                |
| loginStatus | [WebSocketLoginStatusResponse](#struct_WebSocketLoginStatusResponse) | 账号登录状态              |

###### 响应状态码

<a id="message_WebSocketResponseCode"></a>

| CODE | MESSAGE |
|------|---------|
| 200  | 响应成功    |

### 数据结构

###### [WebSocket] 消息响应 WebSocketMessageResponse

<a id="struct_WebSocketMessageResponse"></a>

| 字段值     | 字段类型      | 是否必须 | 字段描述                                                 |
|---------|-----------|------|------------------------------------------------------|
| code    | uint32    | 是    | 响应状态码, 该值请参考 [响应状态码](#message_WebSocketResponseCode) |
| message | string    | 是    | 响应描述                                                 |
| data    | interface | 是    | JSON 数据包                                             |

###### [WebSocket] - [请求事件: login] 账号登录 WebSocketLoginRequest

<a id="struct_WebSocketLoginRequest"></a>

| 字段值     | 字段类型      | 是否必须 | 字段描述 |
|---------|-----------|------|------|
| user_id | string    | 是    | 用户ID |
| usersig | string    | 是    | 用户签名 |

###### [WebSocket] - [响应事件: login] 账号登录 WebSocketLoginResponse

<a id="struct_WebSocketLoginResponse"></a>

| 字段值              | 字段类型      | 是否必须 | 字段描述                                                                             |
|------------------|-----------|------|----------------------------------------------------------------------------------|
| user_id          | string    | 是    | 用户ID                                                                             |
| username         | string    | 是    | 用户名称                                                                             |
| avatar           | string    | 是    | 用户头像                                                                             |
| is_admin         | bool      | 是    | 是否管理员                                                                            |
| status           | string    | 是    | 在线状态 离线(Offline) 在线(Online)                                                      |
| first_login_time | time.Time | 是    | 用户首次登录时间                                                                         |
| last_login_time  | time.Time | 是    | 用户最后登录时间                                                                         |
| last_login_ip    | string    | 是    | 用户最后登录IP                                                                         |
| repeat_login    | bool      | 是    | 在此之前用户是否已登录状态 (如果用户此前未登录将返回false); 如果返回true, 表示服务端并不会更新登录信息, 只是把上次登录成功的信息返回给客户端。 |

###### [WebSocket] - [响应事件: logout] 账号登出 WebSocketLogoutResponse

<a id="struct_WebSocketLogoutResponse"></a>

| 字段值         | 字段类型      | 是否必须 | 字段描述   |
|-------------|-----------|------|--------|
| user_id     | string    | 是    | 用户ID   |
| logout_time | time.Time | 是    | 用户登出时间 |

###### [WebSocket] - [响应事件: loginStatus] 账号登录状态 WebSocketLoginStatusResponse

<a id="struct_WebSocketLoginStatusResponse"></a>

| 字段值     | 字段类型   | 是否必须 | 字段描述   |
|---------|--------|------|--------|
| user_id | string | 是    | 用户ID, 未登录返回空字符串   |
| status  | string | 是    | 登录状态 已登录(Login) 未登录(Logout) |

