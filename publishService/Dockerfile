# 使用 Golang 镜像作为基础镜像
FROM golang:1.20-alpine

# 更换为阿里云镜像源并更新
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update

# 设置环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 安装 FFmpeg 和其他依赖
RUN apk add --no-cache ffmpeg

# 移动到工作目录：/app
WORKDIR /app

# 将代码复制到容器中
COPY . .

# 编译您的 Golang 项目
RUN go build -o publishService ./cmd

# 暴露端口
EXPOSE 8083

# 定义启动命令，启动您的 Golang 微服务
CMD [ "./publishService" ]