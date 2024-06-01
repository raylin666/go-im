## 概述

> 即时通信服务, `HTTP` 和 `WebSocket` 共用同个端口, 路由通过 `Nginx` 转发只对应的协议模块。

## NGINX 配置

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

## 服务启动

这里假设您已经算是一名 `goer` 了，所以不过多描述基础部署，只需要如下几步即可完成服务启动。

##### 下载服务项目源代码

> `git clone git@github.com:raylin666/go-im.git`

##### Make 能帮你做很多事

> `make init`

> `make generate`

> `make wire`

##### 运行项目

> `make run`

##### 访问项目

> 访问 WebSocket 地址: `ws://127.0.0.1:10010`

> 访问 RPC API 接口地址: `http://127.0.0.1:10010/api`

> 访问 HTTP Web API 接口地址: `http://127.0.0.1:10010/app`


## 注意要点

> 调用IM服务时都必须鉴权才能正常连接。

#### 阅读推荐

服务架构：https://github.com/raylin666/go-im

服务文档：https://github.com/raylin666/doc-goim