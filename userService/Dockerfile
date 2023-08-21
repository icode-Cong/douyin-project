FROM golang:1.20-alpine

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /app

# 将代码复制到容器中

COPY . .
RUN go build -o userService ./cmd

EXPOSE 8082

CMD [ "./userService" ]