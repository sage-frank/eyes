FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件
RUN go build -o eyes .

###################
# 接下来创建一个小镜像
###################
FROM debian:buster-slim
#FROM centos:7.9.2009
COPY ./wait-for.sh /
COPY ./templates /templates
COPY ./config /config/

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/eyes /

RUN set -eux; apt-get update; apt-get install -y --no-install-recommends netcat;  chmod 755 wait-for.sh

RUN set -eux; chmod 755 wait-for.sh

# 声明服务端口
EXPOSE 8081

# 需要运行的命令
ENTRYPOINT ["/eyes", "--env=prod"]