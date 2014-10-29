# 测试版 ubuntu 部署手册

- - -

## RVM 安装 ruby

参考 http://rvm.io/

1. `\curl -sSL https://get.rvm.io | bash -s stable`
2. `source /etc/profile.d/rvm.sh` # 根据终端提示
3. `rvm install 2.1.3` #首次安装需要较长一段时间
4. `rvm use 2.1.3 --default`

## NVM 安装 nodejs

参考 https://www.digitalocean.com/community/tutorials/how-to-install-node-js-with-nvm-node-version-manager-on-a-vps

1. `curl https://raw.githubusercontent.com/creationix/nvm/v0.11.1/install.sh | bash`
2. `nvm install 0.10.13`

nodejs 安装通过后，尝试

1. `npm install pm2 -g`
2. `npm install hiredis -g` # 可选

安装两个后面需要的全局的包

## 安装 golang 环境

参考 http://stackoverflow.com/questions/17480044/how-to-install-the-current-version-of-go-in-ubuntu

1. `sudo apt-get install python-software-properties`
2. `sudo add-apt-repository ppa:duh/golang`
3. `sudo apt-get update`
4. `sudo apt-get install golang`

这不是最好的安装方法，但是是最简单（gvm在国内被墙）的。

## 安装 redis 以及 mysql

参考 http://redis.io/topics/quickstart

1. `wget http://download.redis.io/redis-stable.tar.gz`
2. `tar xvzf redis-stable.tar.gz`
3. `cd redis-stable`
4. `make`
5. `make install`
6. `redis-server &` # 启动 redis-server

mysql 的安装这里就不多做描述了，ubuntu 下面直接 `apt-get install mysql-server`

## 下载并修改配置文件

1. `apt-get install git` 
2. `git clone https://github.com/jianxinio/postman.git` # 下载项目文件
3. `cd postman`
4. `cd middleware/ && npm install && cd ..` # 安装 middleware 依赖
5. `cd api/ && npm install && cd ..` # 安装 api 依赖
6. `sudo apt-get install libmysqlclient-dev` # 为了 ruby 的 mysql 安装成功/ centos 下需要安装对应的 `mysql-devel`
7. `cd website/ && bundle install && cd ..` # 安装 website 依赖
8. `cd postman/ && sh bootstrap.sh && cd ..` # 安装 postman 依赖

#### 修改关键配置文件

1. 进入 config 文件夹
2. `cp domain.example.json domain.json && cp database.example.json database.json`
3. 修改 domain.json 文件，如果是本地的调试的话，可以把域名都设置为 localhost，端口换为非 80 端口
4. 修改 database.json，主要修改 mysql 的 username 和 password

#### 初始化数据库

在数据库中建立对应的 database，例如

  1. mysql -u root -p
  2. CREATE SCHEMA `jianxin` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
  
    
1. 进入 api 文件夹
2. `node install.js`

#### 生成 tls 通信需要的证书 

参考 https://www.openssl.org/docs/HOWTO/certificates.txt

1. 进入 config 下面的 pems/middleware
2. `openssl genrsa -out privkey.pem 2048`
3. `openssl req -new -x509 -key privkey.pem -out cert.pem -days 1095`

- - - 
到这里所有前期准备就完成了，下面开始准备启动


## 启动项目

1. 启动 middleware，进入 middleware，执行 `pm2 start processes.json`
1. 启动 api，进入 api，执行 `pm2 start processes.json`
2. 启动 website，进入 website，执行 `padrino start -h 0.0.0.0 -e production &`

*如果在线上部署，一定请使用 iptables 进行处理！*

## 部署发送端

1. 登录 website，新建 sender，输入对应的 ip（如果本地测试请输入 127.0.0.1，正式环境需要同一台机器也要输入外网ip） 以及 hostname
2. 按照要求设置 dns，并验证通过
3. 点击 download config 下载 sender 的配置文件。
4. 进入 postman，运行 `source dev.sh`
5. 继续进入 src/postman/main/postman，并执行 `go install`
6. 将 postman/bin 中的 postman 可执行文件和 下载的 config.json 放在一起后，./postman & 启动发送器。(不要在同一个终端中运行，否则会因为开发环境变量导致运行出错)

## 使用说明
