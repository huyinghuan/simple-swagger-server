### swagger 服务

一键部署 swagger 服务

###  运行

```bash
go run main.go
# 或者在release 下载编译好的 linux amd64 的二进制文件
chmod +x swagger-server.bin
./swagger-server.bin
```

默认服务启动在`8888`端口，默认读取进程所在的目录的`docs`目录下的所有`*.json`文件作为`swagger`服务配置文件

可以通过以下命令：

```
--port 8080 #指定服务端口
--docs mydocs #指定swagger服务配置文件目录
--auth password #指定swagger服务的访问密码
```

### 命令行上传
    
```bash
curl -X POST -H "Content-Type: multipart/form-data" -F "file=@user.json" http://ip:port/api/upload
```

如果有认证密码:

```bash
curl -X POST -H "Content-Type: multipart/form-data" -F "file=@user.json" "http://ip:port/api/upload?token=password"
```
