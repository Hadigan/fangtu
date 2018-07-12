## 安装golang

下载二进制文件包

```
wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
```

解压

```
tar -C /usr/local -xzf go1.10.3.linux-amd64.tar.gz
```

修改环境变量

```
vim /etc/profile
```

添加

```
export GOPATH=/root/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin/:$GOPATH/bin/
```

查看是否安装正确```go version```

## 安装docker

