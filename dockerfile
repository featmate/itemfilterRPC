# 构造可执行文件
FROM --platform=$TARGETPLATFORM golang:alpine as build_bin
ENV GO111MODULE=on
# 停用cgo
ENV CGO_ENABLED=0
WORKDIR /code
COPY go.mod /code/go.mod
COPY go.sum /code/go.sum
# 添加源文件
COPY filterstore /code/filterstore
COPY itemfilterRPC_pb /code/itemfilterRPC_pb
COPY itemfilterRPC_serv /code/itemfilterRPC_serv
COPY metadatastore /code/metadatastore
COPY main.go /code/main.go
RUN go build -ldflags "-s -w" -o itemfilterRPC-go main.go

# 使用upx压缩可执行文件
FROM --platform=$TARGETPLATFORM alpine:3.11 as upx
WORKDIR /code
# 安装upx
RUN apk update && apk add --no-cache upx && rm -rf /var/cache/apk/*
COPY --from=build_bin /code/itemfilterRPC-go .
RUN upx --best --lzma -o itemfilterRPC itemfilterRPC-go

# 编译获得grpc-health-probe
FROM --platform=$TARGETPLATFORM golang:buster as build_grpc-health-probe
ENV GO111MODULE=on
# 停用cgo
ENV CGO_ENABLED=0
# 安装grpc-health-probe
RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

# 使用压缩过的可执行文件构造镜像
FROM --platform=$TARGETPLATFORM alpine:3.16.2 as build_img
# 打包镜像
COPY --from=build_grpc-health-probe /go/bin/grpc-health-probe .
RUN chmod +x /grpc-health-probe
COPY --from=upx /code/itemfilterRPC .
RUN chmod +x /itemfilterRPC
EXPOSE 5000
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/grpc-health-probe","-addr=:5000" ]
ENTRYPOINT [ "/itemfilterRPC"]