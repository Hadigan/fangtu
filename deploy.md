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

## 拉取docker镜像

下载配置文件

```
git clone https://github.com/Hadigan/fangtu.git
```

获取fabric二进制文件以及docker镜像

```
./fangtu/pull-images.sh 
```

可以看见在当前文件夹下多了一个bin文件夹
把bin文件夹内容移到$GOPATH/bin/下，这样子就不用在修改环境变量了
如果$GOPATH/bin文件夹不存在就手动创建

```
mv ./bin/* /root/go/bin/
```

## ca-server配置

caserver的管理员账户密码在fangtu/config/docker-compose-caserver.yaml中指定了，请在第一次启动caserver的容器之前进行修改。否则只能通过ca-client进行修改。默认的caserver的管理员账户为==admin== ，密码为 ==caserver==。

启动ca-server容器,在/root/workspaces/下执行操作

```
docker-compose -f ./fangtu/config/docker-compose-caserver.yaml up -d
```







