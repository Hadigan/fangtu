# 1. 基础环境准备

## 1.1. 安装golang

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

## 1.2. 安装docker

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

## 1.3. 安装docker-compose

下载docker-compose二进制文件

```
curl -L https://github.com/docker/compose/releases/download/1.21.2/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
```

添加执行权限

```
chmod +x /usr/local/bin/docker-compose
```

## 1.4. 拉取docker镜像

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

# 2. ca-server配置
这一节的操作都在ca服务器上完成

## 2.1. 启动ca服务

将我们需要的配置文件拷出来

```
root@fabric-60-21:~/workspaces# cp -r ./fangtu/ca-server ./
```

caserver的管理员账户密码在./ca-server/docker-compose-caserver.yaml中指定了，请在第一次启动caserver的容器之前进行修改。否则只能通过ca-client进行修改。默认的caserver的管理员账户为==admin== ，密码为 ==caserver==。

启动ca-server容器,在```/root/workspaces/ca-server/```下执行操作

```
root@fabric-60-21:~/workspaces/ca-server# docker-compose -f docker-compose-caserver.yaml up -d
```

如果报错，请删除之前启动的名为fabric-ca-server的容器

启动成功之后 ```/root/workspaces/ca-server/fabric-ca-server/```下就是caserver的配置文件以及数据库以及msp。

利用我们ca服务默认管理员账户admin和密码caserver生成admin账户的凭证

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client enroll -u http://admin:caserver@127.0.0.1:7054 -H ./fabric-ca-files/admin
```
## 2.2. 删除默认的组织结构
下面的命令中的-H参数代表我们连接ca服务所使用的用户

可以看到初始时联盟组织结构如下：

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin affiliation list
affiliation: .
   affiliation: org2
      affiliation: org2.department1
   affiliation: org1
      affiliation: org1.department1
      affiliation: org1.department2
```

我们现在要删除这些组织结构，然后创建我们自己的组织

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin  affiliation remove --force  org1
```
```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin  affiliation remove --force  org2
```
## 2.3. 创建自己的组织
创建自己的组织

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin  affiliation add  com
```
```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin  affiliation add  com.mederahealth
```
```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin  affiliation add  com.mederahealth.yiyuan
```
```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin  affiliation add  com.mederahealth.shaoyifu
```
现在的组织结构如下

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client -H ./fabric-ca-files/admin affiliation list
affiliation: com
   affiliation: com.mederahealth
      affiliation: com.mederahealth.shaoyifu
      affiliation: com.mederahealth.yiyuan
```
## 2.4. 注册节点和用户

注册mederahealth.com的管理员 admin@mederahealth.com 

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client register -H ./fabric-ca-files/admin --id.name admin@mederahealth.com --id.type client --id.affiliation "com.mederahealth" --id.attrs '"hf.Registrar.Roles=client,orderer,peer,user","hf.Registrar.DelegateRoles=client,orderer,peer,user",hf.Registrar.Attributes=*,hf.GenCRL=true,hf.Revoker=true,hf.AffiliationMgr=true,hf.IntermediateCA=true,role=admin:ecert' --id.secret=admin@mederahealth.com 

```

注册yiyuan.mederahealth.com的管理员 admin@yiyuan.mederahealth.com

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client register -H ./fabric-ca-files/admin --id.name admin@yiyuan.mederahealth.com --id.type client --id.affiliation "com.mederahealth.yiyuan" --id.attrs '"hf.Registrar.Roles=client,orderer,peer,user","hf.Registrar.DelegateRoles=client,orderer,peer,user",hf.Registrar.Attributes=*,hf.GenCRL=true,hf.Revoker=true,hf.AffiliationMgr=true,hf.IntermediateCA=true,role=admin:ecert' --id.secret=admin@yiyuan.mederahealth.com 
```

注册shaoyifu.mederahealth.com的管理员 admin@shaoyifu.mederahealth.com

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client register -H ./fabric-ca-files/admin --id.name admin@shaoyifu.mederahealth.com --id.type client --id.affiliation "com.mederahealth.shaoyifu" --id.attrs '"hf.Registrar.Roles=client,orderer,peer,user","hf.Registrar.DelegateRoles=client,orderer,peer,user",hf.Registrar.Attributes=*,hf.GenCRL=true,hf.Revoker=true,hf.AffiliationMgr=true,hf.IntermediateCA=true,role=admin:ecert' --id.secret=admin@shaoyifu.mederahealth.com 
```



# 3. orderer.mederahealth.com配置

以下操作均在orderer服务器上执行

拉取配置文件仓库

```
root@fabric-60-22:~/workspaces/# git clone http://github.com/Hadigan/fangtu.git
```

将orderer对应的配置文件烤出

```
root@fabric-60-22:~/workspaces/# cp -r fangtu/orderer.mederahealth.com ./
```


获取mederahealth.com组织的msp

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# fabric-ca-client getcacert -u http://172.19.60.21:7054 -M $(pwd)/crypto-config/mederahealth.com/msp
```
执行后会将mederahealth的msp获取到crypto-config目录下

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# tree crypto-config/
crypto-config/
└── mederahealth.com
    └── msp
        ├── cacerts
        │   └── 172-19-60-21-7054.pem
        ├── intermediatecerts
        │   └── 172-19-60-21-7054.pem
        ├── keystore
        └── signcerts
```

使用admin@mederahealth.com的账号密码获取证书

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# fabric-ca-client enroll -u http://admin@mederahealth.com:admin@mederahealth.com@172.19.60.21:7054 -H $(pwd)/crypto-config/mederahealth.com/users/admin@mederahealth.com
```

这时候发现crypto-config/mederahealth.com/users/admin 多了以下文件

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# ls crypto-config/mederahealth.com/users/admin@mederahealth.com/
fabric-ca-client-config.yaml  msp
```



## 3.1. 用admin@mederahealth.com账号创建其他账号

注册orderer.mederahealth.com

```
root@fabric-60-21:~/workspaces/ca-server# fabric-ca-client register -H $(pwd)/crypto-config/mederahealth.com/users/admin@mederahealth.com --id.name orderer.mederahealth.com --id.type orderer --id.affiliation "com.mederahealth" --id.attrs 'role=orderer:ecert' --id.secret=orderer.mederahealth.com 
```

获取orderer.mederahealth.com 的证书

```
fabric-ca-client enroll -u http://orderer.mederahealth.com:orderer.mederahealth.com@172.19.60.21:7054 -H $(pwd)/crypto-config/mederahealth.com/orderers/orderer.mederahealth.com/
```



## 3.2. 复制管理员的证书
==这一步争议 要不要复制 先不复制吧==


将admin@mederahealth.com用户的证书复制到```./crypto-config/mederahealth.com/msp/admincerts/```下，==如果不复制的话创建创世区块时将会报错==

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# mkdir -p crypto-config/mederahealth.com/msp/admincerts
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# cp crypto-config/mederahealth.com/users/admin@mederahealth.com/msp/signcerts/cert.pem crypto-config/mederahealth.com/msp/admincerts/
```

将admin@mederahealth.com用户的证书复制到admin@mederahealth.com用户的msp的admincerts下

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# mkdir -p crypto-config/mederahealth.com/users/admin@mederahealth.com/msp/admincerts

root@fabric-60-23:~/workspaces/orderer.mederahealth.com# cp crypto-config/mederahealth.com/users/admin@mederahealth.com/msp/signcerts/cert.pem crypto-config/mederahealth.com/users/admin@mederahealth.com/msp/admincerts/
```

将admin@mederahealth.com用户的证书复制到orderer.mederahealth.com账号的msp的admincerts下 ==如果不复制的话启动orderer容器将会报错==

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# mkdir -p crypto-config/mederahealth.com/orderers/orderer.mederahealth.com/msp/admincerts

root@fabric-60-23:~/workspaces/orderer.mederahealth.com# cp crypto-config/mederahealth.com/users/admin@mederahealth.com/msp/signcerts/cert.pem crypto-config/mederahealth.com/orderers/orderer.mederahealth.com/msp/admincerts/
```

## 3.3. 获取其他组织的msp

==这里有错误，以后再修改，用enroll获取的sig文件夹下的cert.pem文件不一样，因此需要采用复制的方法==

获取其他组织的msp以创建创世区块，这里只能获取每个组织的ca证书，没有密钥等文件，因为我们只有一个根ca服务器，因此通过getcacert获取的每个组织的ca证书都是一样。

获取yiyuan.mederahealth.com的msp

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# fabric-ca-client getcacert -u http://172.19.60.21:7054 -M $(pwd)/crypto-config/yiyuan.mederahealth.com/msp
```

获取admin@yiyuan.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/orderer.mederahealth.com# fabric-ca-client enroll -u http://admin@yiyuan.mederahealth.com:admin@yiyuan.mederahealth.com@172.19.60.21:7054 -H $(pwd)/crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com
```

复制证书

```
root@fabric-60-22:~/workspaces/orderer.mederahealth.com# mkdir -p crypto-config/yiyuan.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/orderer.mederahealth.com# cp crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com/msp/signcerts/cert.pem crypto-config/yiyuan.mederahealth.com/msp/admincerts/
```

获取shaoyifu.mederahealth.com的msp

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# fabric-ca-client getcacert -u http://172.19.60.21:7054 -M $(pwd)/crypto-config/shaoyifu.mederahealth.com/msp
```

获取admin@shaoyifu.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/orderer.mederahealth.com# fabric-ca-client enroll -u http://admin@shaoyifu.mederahealth.com:admin@shaoyifu.mederahealth.com@172.19.60.21:7054 -H $(pwd)/crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com
```

复制证书

```
root@fabric-60-22:~/workspaces/orderer.mederahealth.com# mkdir -p crypto-config/shaoyifu.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/orderer.mederahealth.com# cp crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/signcerts/cert.pem crypto-config/shaoyifu.mederahealth.com/msp/admincerts/
```

## 3.4 创建创世区块

设置配置文件的路径

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# export FABRIC_CFG_PATH=$PWD
```
生成创世区块

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```

## 3.5 启动orderer服务

```
root@fabric-60-23:~/workspaces/orderer.mederahealth.com# docker-compose -f ./docker-compose-Orderer.yaml up -d
```





# 4. yiyuan.mederahealth.com配置

以下操作均在yiyuan server上完成

拉取配置文件仓库

```
root@fabric-60-22:~/workspaces/# git clone http://github.com/Hadigan/fangtu.git
```

将yiyuan对应的配置文件烤出

```
root@fabric-60-22:~/workspaces/# cp -r fangtu/yiyuan.mederahealth.com ./
```

获取yiyuan.mederahealth.com的msp

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client getcacert -u http://172.19.60.21:7054 -M $(pwd)/crypto-config/yiyuan.mederahealth.com/msp
```

获取admin@yiyuan.mederahealth.com账号的凭证，该账户是由ca注册，密码为admin@yiyuan.mederahealth.com

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client enroll -u http://admin@yiyuan.mederahealth.com:admin@yiyuan.mederahealth.com@172.19.60.21:7054 -H $(pwd)/crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com
```

## 4.1. 使用admin@yiyuan.mederahealth.com账号创建其他账号

注册 peer0.yiyuan.mederahealth.com 

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client register -H $(pwd)/crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com --id.name peer0.yiyuan.mederahealth.com --id.type peer --id.affiliation "com.mederahealth.yiyuan" --id.attrs 'role=peer:ecert' --id.secret=peer0.yiyuan.mederahealth.com 
```

获取peer0.yiyuan.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client enroll -u http://peer0.yiyuan.mederahealth.com:peer0.yiyuan.mederahealth.com@172.19.60.21:7054 -H /root/workspaces/yiyuan.mederahealth.com/crypto-config/yiyuan.mederahealth.com/peers/peer0.yiyuan.mederahealth.com/
```

注册 peer1.yiyuan.mederahealth.com 

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client register -H $(pwd)/crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com --id.name peer1.yiyuan.mederahealth.com --id.type peer --id.affiliation "com.mederahealth.yiyuan" --id.attrs 'role=peer:ecert' --id.secret=peer1.yiyuan.mederahealth.com 
```

获取peer1.yiyuan.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client enroll -u http://peer1.yiyuan.mederahealth.com:peer1.yiyuan.mederahealth.com@172.19.60.21:7054 -H /root/workspaces/yiyuan.mederahealth.com/crypto-config/yiyuan.mederahealth.com/peers/peer1.yiyuan.mederahealth.com/
```

注册普通用户 user1.yiyuan.mederahealth.com

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client register -H $(pwd)/crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com --id.name user1@yiyuan.mederahealth.com --id.type client --id.affiliation "com.mederahealth.yiyuan" --id.attrs '"hf.Registrar.Roles=","hf.Registrar.DelegateRoles=",role=app:ecert' --id.secret=user1@yiyuan.mederahealth.com 
```

获取user1@yiyuan.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# fabric-ca-client enroll -u http://user1@yiyuan.mederahealth.com:user1@yiyuan.mederahealth.com@172.19.60.21:7054 -H /root/workspaces/yiyuan.mederahealth.com/crypto-config/yiyuan.mederahealth.com/users/user1@yiyuan.mederahealth.com/
```


## 4.2. 复制admin的凭证
==这一步争议待定,暂时不复制==

将admin@yiyuan.mederahealth.com的证书复制到crypto-config/yiyuan.mederahealth.com/msp/admincerts/ ==不复制的话无法生成channel.tx==

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# mkdir -p crypto-config/yiyuan.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# cp crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com/msp/signcerts/cert.pem crypto-config/yiyuan.mederahealth.com/msp/admincerts/
```

将admin@yiyuan.mederahealth.com的证书复制到crypto-config/yiyuan.mederahealth.com/users/admin/msp/admincerts/ ==必须复制，否则client使用admin权限将出错==

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# mkdir -p crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# cp crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com/msp/signcerts/cert.pem crypto-config/yiyuan.mederahealth.com/users/adminyiyuan.mederahealth.com/msp/admincerts/
```

将admin@yiyuan.mederahealth.com的证书复制到crypto-config/yiyuan.mederahealth.com/peers/peer0.yiyuan.mederahealth.com/msp/admincerts/ ==不复制的话无法启动peer节点==

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# mkdir -p crypto-config/yiyuan.mederahealth.com/peers/peer0.yiyuan.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# cp crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com/msp/signcerts/cert.pem crypto-config/yiyuan.mederahealth.com/peers/peer0.yiyuan.mederahealth.com/msp/admincerts/
```

将admin@yiyuan.mederahealth.com的证书复制到crypto-config/yiyuan.mederahealth.com/peers/peer1.yiyuan.mederahealth.com/msp/admincerts/ ==不复制的话无法启动peer节点==

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# mkdir -p crypto-config/yiyuan.mederahealth.com/peers/peer1.yiyuan.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# cp crypto-config/yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com/msp/signcerts/cert.pem crypto-config/yiyuan.mederahealth.com/peers/peer1.yiyuan.mederahealth.com/msp/admincerts/
```

将admin@yiyuan.mederahealth.com的证书复制到user1@yiyuan.mederahealth.com/msp/admincerts/

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# mkdir -p crypto-config/yiyuan.mederahealth.com/users/user1@yiyuan.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# cp yiyuan.mederahealth.com/users/admin@yiyuan.mederahealth.com/msp/signcerts/cert.pem crypto-config/yiyuan.mederahealth.com/users/user1@yiyuan.mederahealth.com/msp/admincerts/
```






# 5. shaoyifu.mederahealth.com配置

以下操作均在shaoyifu server上完成

拉取配置文件仓库

```
root@fabric-60-22:~/workspaces/# git clone http://github.com/Hadigan/fangtu.git
```

将shaoyifu对应的配置文件烤出

```
root@fabric-60-22:~/workspaces/# cp -r fangtu/shaoyifu.mederahealth.com ./
```

获取shaoyifu.mederahealth.com的msp

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client getcacert -u http://172.19.60.21:7054 -M $(pwd)/crypto-config/shaoyifu.mederahealth.com/msp
```

获取admin@shaoyifu.mederahealth.com账号的凭证，该账户是由ca注册，密码为admin@shaoyifu.mederahealth.com

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client enroll -u http://admin@shaoyifu.mederahealth.com:admin@shaoyifu.mederahealth.com@172.19.60.21:7054 -H $(pwd)/crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com
```

## 5.1. 使用admin@shaoyifu.mederahealth.com账号创建其他账号

注册 peer0.shaoyifu.mederahealth.com 

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client register -H $(pwd)/crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com --id.name peer0.shaoyifu.mederahealth.com --id.type peer --id.affiliation "com.mederahealth.shaoyifu" --id.attrs 'role=peer:ecert' --id.secret=peer0.shaoyifu.mederahealth.com 
```

获取peer0.shaoyifu.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client enroll -u http://peer0.shaoyifu.mederahealth.com:peer0.shaoyifu.mederahealth.com@172.19.60.21:7054 -H /root/workspaces/shaoyifu.mederahealth.com/crypto-config/shaoyifu.mederahealth.com/peers/peer0.shaoyifu.mederahealth.com/
```

注册 peer1.shaoyifu.mederahealth.com 

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client register -H $(pwd)/crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com --id.name peer1.shaoyifu.mederahealth.com --id.type peer --id.affiliation "com.mederahealth.shaoyifu" --id.attrs 'role=peer:ecert' --id.secret=peer1.shaoyifu.mederahealth.com 
```

获取peer1.shaoyifu.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client enroll -u http://peer1.shaoyifu.mederahealth.com:peer1.shaoyifu.mederahealth.com@172.19.60.21:7054 -H /root/workspaces/shaoyifu.mederahealth.com/crypto-config/shaoyifu.mederahealth.com/peers/peer1.shaoyifu.mederahealth.com/
```

注册普通用户 user1.shaoyifu.mederahealth.com

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client register -H $(pwd)/crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com --id.name user1@shaoyifu.mederahealth.com --id.type client --id.affiliation "com.mederahealth.shaoyifu" --id.attrs '"hf.Registrar.Roles=","hf.Registrar.DelegateRoles=",role=app:ecert' --id.secret=user1@shaoyifu.mederahealth.com 
```

获取user1@shaoyifu.mederahealth.com的证书

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# fabric-ca-client enroll -u http://user1@shaoyifu.mederahealth.com:user1@shaoyifu.mederahealth.com@172.19.60.21:7054 -H /root/workspaces/shaoyifu.mederahealth.com/crypto-config/shaoyifu.mederahealth.com/users/user1@shaoyifu.mederahealth.com/
```


## 5.2. 复制admin的凭证
==这一步争议待定,暂时不复制==

将admin@shaoyifu.mederahealth.com的证书复制到shaoyifu.mederahealth.com/msp/admincerts/ ==必须复制，提供给其他组织的msp必须包含admincerts文件夹==

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# mkdir -p crypto-config/shaoyifu.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# cp crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/signcerts/cert.pem crypto-config/shaoyifu.mederahealth.com/msp/admincerts/
```

将admin@shaoyifu.mederahealth.com的证书复制到shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/admincerts/ ==不复制的话无法peer将无法使用admin的权限==

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# mkdir -p crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# cp crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/signcerts/cert.pem crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/admincerts/
```

将admin@shaoyifu.mederahealth.com的证书复制到shaoyifu.mederahealth.com/peers/peer0.shaoyifu.mederahealth.com/msp/admincerts/ ==不复制的话无法启动peer节点==

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# mkdir -p crypto-config/shaoyifu.mederahealth.com/peers/peer0.shaoyifu.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# cp crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/signcerts/cert.pem crypto-config/shaoyifu.mederahealth.com/peers/peer0.shaoyifu.mederahealth.com/msp/admincerts/
```

将admin@shaoyifu.mederahealth.com的证书复制到shaoyifu.mederahealth.com/peers/peer1.shaoyifu.mederahealth.com/msp/admincerts/ ==不复制的话无法启动peer节点==

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# mkdir -p crypto-config/shaoyifu.mederahealth.com/peers/peer1.shaoyifu.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# cp crypto-config/shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/signcerts/cert.pem crypto-config/shaoyifu.mederahealth.com/peers/peer1.shaoyifu.mederahealth.com/msp/admincerts/
```

将admin@shaoyifu.mederahealth.com的证书复制到user1@shaoyifu.mederahealth.com/msp/admincerts/

```
root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# mkdir -p crypto-config/shaoyifu.mederahealth.com/users/user1@shaoyifu.mederahealth.com/msp/admincerts/

root@fabric-60-22:~/workspaces/shaoyifu.mederahealth.com# cp shaoyifu.mederahealth.com/users/admin@shaoyifu.mederahealth.com/msp/signcerts/cert.pem crypto-config/shaoyifu.mederahealth.com/users/user1@shaoyifu.mederahealth.com/msp/admincerts/
```

# 6. 通道配置

获得shaoyifu.mederahealth.com 的 完整msp，需要msp文件夹下有admincerts目录，且其中存放了admin@shaoyifu.mederahealth.com 的signcert.值只能让shaoyifu这个组织发给yiyuan组织。这个msp中没有私钥的信息，因此是可以提供给其他的组织的。

我们通过scp来发送，实际场景中应该通过其他的途径（因为shaoyifu不可能知道yiyuan的服务器的账户密码）。

发完之后，yiyuan的crypto-config/shaoyifu.mederahealth.com/msp下就有admincerts文件夹

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com/crypto-config/shaoyifu.mederahealth.com/msp# ls
admincerts  cacerts  intermediatecerts  keystore  signcerts
```

在yiyuan.mederahealth.com服务器中创建通道配置交易

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# export CHANNEL_NAME=mychannel export FABRIC_CFG_PATH=$PWD

root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
```

## 定义通道的anchorpeer

生成mychannel通道中yiyuan组织的anchor peer更新，该操作在yiyuan server上执行

```
export FABRIC_CFG_PATH=$PWD
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/YiyuanMSPanchors.tx -channelID mychannel -asOrg YiyuanMSP
```

生成mychannel通道中shaoyifu组织的anchor peer更新，该操作在shaoyifu server上执行

```
export FABRIC_CFG_PATH=$PWD
root@fabric-60-20:~/workspaces/shaoyifu.mederahealth.com# configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/ShaoyifuMSPanchors.tx -channelID mychannel -asOrg ShaoyifuMSP
```

# 7. 启动peer

启动yiyuan的peer

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# docker-compose -f docker-compose-yiyuan-peer.yaml up -d
```

启动shaoyifu的peer

```
root@fabric-60-20:~/workspaces/shaoyifu.mederahealth.com# docker-compose -f docker-compose-shaoyifu-peer.yaml up -d
```


# 8. 创建通道并加入

## 对peer0.yiyuan.mederahealth.com进行操作

进入yiyuan的peer0的cli

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# docker exec -it cli.peer0 bash
```
在cli.peer0中以admin@yiyuan.mederahealth.com的名义创建名为 mychannel的通道

```
root@537a30b2957e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel create -o orderer.mederahealth.com:37050 -c mychannel -f ./channel-artifacts/channel.tx
```

创建成功后当前目录下将多出一个mychannel.block的文件

```
root@537a30b2957e:/opt/gopath/src/github.com/hyperledger/fabric/peer# ls
channel-artifacts  crypto  mychannel.block
```

让peer0加入mychannel

```
root@537a30b2957e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel join -b mychannel.block
```

## 对peer1.yiyuan.mederahealth.com进行操作

进入peer1的cli容器

```
root@fabric-60-22:~/workspaces/yiyuan.mederahealth.com# docker exec -it cli.peer1 bash
```
获取mychannel的第一个block

```
root@3d1799b91d0e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel fetch oldest -c mychannel -o orderer.mederahealth.com:37050
```
此时当前目录下将会多出mychannel_oldest.block，这个block内容与上面peer0生成的mychannel.block是一样的，可以使用md5sum验证一下。

接下来让peer0加入mychannel
```
root@3d1799b91d0e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel join -b mychannel_oldest.block
```
## 对peer0.shaoyifu.mederahealth.com进行操作

进入peer0的cli容器

```
root@fabric-60-20:~/workspaces/shaoyifu.mederahealth.com# docker exec -it cli.peer0 bash
```
获取mychannel的第一个block

```
root@f39c5c8293e1:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel fetch oldest -c mychannel -o orderer.mederahealth.com:37050
```
此时当前目录下将会多出mychannel_oldest.block，这个block内容与上面生成的mychannel.block是一样的，可以使用md5sum验证一下。

接下来让peer0加入mychannel

```
root@f39c5c8293e1:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel join -b mychannel_oldest.block
```
## 对peer1.shaoyifu.mederahealth.com进行操作

进入peer1的cli容器

```
root@fabric-60-20:~/workspaces/shaoyifu.mederahealth.com# docker exec -it cli.peer1 bash
```
获取mychannel的第一个block

```
root@9982d4115973:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel fetch oldest -c mychannel -o orderer.mederahealth.com:37050
```
此时当前目录下将会多出mychannel_oldest.block，这个block内容与上面生成的mychannel.block是一样的，可以使用md5sum验证一下。

接下来让peer0加入mychannel

```
root@9982d4115973:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel join -b mychannel_oldest.block
```

# 更新mychannel的anchor peer

更新mychannel中yiyuan组织的anchor peer，我们只要有admin权限即可，以下操作在yiyuan的服务器的cli.peer0中进行

```
root@537a30b2957e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel update -o orderer.mederahealth.com:37050 -c mychannel -f ./channel-artifacts/YiyuanMSPanchors.tx
```

更新mychannel中shaoyifu组织的anchor peer，我们只要有admin权限即可，以下操作在shaoyifu服务器的cli.peer0中进行

```
root@f39c5c8293e1:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer channel update -o orderer.mederahealth.com:37050 -c mychannel -f ./channel-artifacts/ShaoyifuMSPanchors.tx
```

# 安装以及初始化 chaincode

安装了chaincode的节点才可以执行背书，我们在这里让所有的四个节点都安装chaincode。

我们这里将合约的名称取为 letter


首先是 peer0.yiyuan，下面操作在cli.peer0上执行

```
root@537a30b2957e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode install -n letter4 -v 1.0 -p github.com/chaincode/fangtu/
```

peer1.yiyuan, 下面操作在cli.peer1上执行

```
root@3d1799b91d0e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode install -n letter3 -v 1.0 -p github.com/chaincode/fangtu/
```

peer0.shaoyifu，下面操作在cli.peer0上执行

```
root@f39c5c8293e1:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode install -n letter -v 1.0 -p github.com/chaincode/fangtu/
```

peer1.shaoyifu，下面操作在cli.peer1上执行

```
root@9982d4115973:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode install -n letter -v 1.0 -p github.com/chaincode/fangtu/
```

==初始化合约==，下面操作在cli.peer0上执行，合约的初始化只需要做一次就行了，
```
root@e7c24a19f38e:/opt/gopath/src/github.com/hyperledger/fabric/peer# peer chaincode instantiate -o orderer.mederahealth.com:37050 -C mychannel -n letter4 -v 1.0 -c '{"Args":["init"]}' -P "AND ('YiyuanMSP.member', 'ShaoyifuMSP.member')"
```

```
peer chaincode invoke -o orderer.mederahealth.com:37050 -C mychannel -n letter --peerAddresses peer0.yiyuan.mederahealth.com:37051 --peerAddresses peer0.shaoyifu.mederahealth.com:37051 -c '{"Args":["update","11111111111111111111111111111115"]}'
```

```
peer chaincode invoke -o orderer.mederahealth.com:37050 -C mychannel -n letter4 -c '{"Args":["update","11111111111111111111111111111116"]}'
```

peer chaincode query -C mychannel -n letter4 -c '{"Args":["query", "11111111111111111111111111111116"]}'

```
peer chaincode upgrade -o orderer.mederahealth.com:37050 -C mychannel -P "OR ('YiyuanMSP.peer', 'ShaoyifuMSP.peer')" -n letter -v 2.0 -c '{"Args":["init"]}' 

peer chaincode package ccpack.out -n letter -p github.com/chaincode/fangtu/ -v 1.1 -s -S


```
