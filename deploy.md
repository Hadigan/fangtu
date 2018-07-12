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

安装包管理器支持https的包

```
apt-get update &&
apt-get install apt-transport-https ca-certificates curl software-properties-common
```
添加docker官方GPG key

```
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```

添加repository

```
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```

安装

```
apt-get update && apt-get install docker-ce 
```

## 安装docker-compose

下载docker-compose二进制文件

```
curl -L https://github.com/docker/compose/releases/download/1.21.2/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
```

添加执行权限

```
chmod +x /usr/local/bin/docker-compose
```

##