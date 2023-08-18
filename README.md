### swagger 服务

一键部署 swagger 服务

###  运行

```
go run main.go
```

默认服务启动在`8888`端口，默认读取进程所在的目录的`docs`目录下的所有`*.json`文件作为`swagger`服务配置文件

可以通过以下命令：

```
--port 8080 #指定服务端口
--docs mydocs #指定swagger服务配置文件目录
```
