# bilicoin 自动投币工具

## 说明
1. 在B站，每日会自动投入补全 50 经验，帮助你快速升级。   
2. 主要筛选热门鬼畜区的视频进行投币。  
3. 不会影响你手动投币  
   由于程序会自动检查当天有没有拿到 50 经验，并在只会在每天的最后时刻补全，也就是说如果当天已经手动投入超过5个币，程序就不会帮你投了。  
4. 支持[方糖](http://sc.ftqq.com/ "ftqq")进行微信通知。  
5. 支持 `QRCode` 登录，无需手动粘贴 `Cookie`  
6. 支持多用户批量处理  
7. 支持通过 `API` 进行控制  
8. [Demo 查看](http://r3inbowari.top:9090/version "Demo")

## 获取工具
项目请自行编译或者从 Release 中下载  
完整的项目包括以下两个文件: `bilicoin_os_arch`, `bili.json[自动生成]`  
1. 编译  
    ```
    git clone https://github.com/r3inbowari/bilicoin.git
    cd bilicoin
    ./build.bat or ./build_upx.bat
    ```

2. 下载  
   从 [Release](https://github.com/r3inbowari/bilicoin/releases "Releases Download") 中下载  

## 基本使用  

1. 命令行输入下面内容，会弹出 `QRCode` 使用B站客户端扫码，添加用户  
    ```
    ./bilicoin_linux_amd64 -n
    ```
    <img src="qrcode.png" style="height:300px" />

2. 登录成功后使用命令行输入下面内容即可开启投币服务  
    ```
    ./bilicoin_linux_amd64 -s
    ```

## 通过 `API` 使用(需要自己开发界面)  
1. 命令行输入下面内容，进入服务器模式  
    ```
    ./bilicoin_linux_amd64 -a
    ```
2. 基本请求  
   详细的请求和响应格式可以看[这里](https://docs.apipost.cn/view/8ab6ae6778a3b405 "API DOC")  
   
    ```
    获得所有用户
    GET /users
    
    添加用户请求
    POST /user
    Response oauthData(二维码地址)
    
    轮询是否登陆成功
    POST /user?oauth=3835a3c053dcda56c0c0136110f69ec9 
    
    试图删除一个UID
    DETETE/user?uid=3077202
    
    试图修改UID的Cron表达式
    GET /{id}/cron?spec=cron表达式
    
    试图修改UID的FTQQ的key或者是开关
    GET /{id}/ft?key=方糖key&sw=开关
    ```

## 其他命令  

1. 查询当前配置文件中所有的 `UID`:  
      
    ```
    ./bilicoin_linux_amd64 -l
    ```
2. 从配置文件中删除指定的 `UID`:  
    
    ```
    ./bilicoin_linux_amd64 [UID] -d
    // example
    // 1. 尝试删除 UID 为 30772 的登录信息
    ./bilicoin_linux_amd64 30722 -d
    ```
3. 配置方糖微信通知[可选]  
   
    ```
    ./bilicoin_linux_amd64 -f [用户ID UID] [方糖 SecretKey]
    // example: 
    // 1. 添加方糖key
    ./bilicoin_linux_amd64 -f 933330 SCUxxxxxTe034cxxxxx732b1xxxxx23f7exxxxxd05eaxxxxxxxxxx

    // 2. 清除方糖key
    ./bilicoin_linux_amd64 -f 933330
    ```
 
 4. 修改Cron表达式(默认是30 50 23 * * ?)  
   
    ```
    ./bilicoin_linux_amd64 -c [用户ID UID] [Cron Spec]
    // example: 
    // 1. 修改cron
    ./bilicoin_linux_amd64 -c 933330 0 10 20 * * ?
    ```

## 使用 Docker 构建  

你也可以使用 `docker` 进行部署，通过使用api进行控制。  
1. 构建镜像 
   
    ```
    // build image
    docker build -t r3inbowari/bilicoin:v1.0.3 .

    // prune dangling image: builder
    docker image prune --filter label=stage=builder
    ```

2. 如果不想构建的话可以直接拉取已经构建好的镜像[linux/amd64](https://hub.docker.com/repository/docker/r3inbowari/bilicoin "DockerHub Page")  
   
    ```
    docker pull r3inbowari/bilicoin
    ```

3. 直接运行即可  
   
    ```
    // run
    docker run \
    --name bilicoin \
    -p 9090:9090 \
    -itd --restart=always \
    r3inbowari/bilicoin:v1.0.3
    ```
    
4. 浏览器打开地址验证是否开启  
   ```
   GET http://localhost:9090/version
   ```
   
## 其他问题  
1. `bili.json` 中的 `canvas_finger` 的值可以选择修改一下，不过不影响使用。  
2. 多用户投币重复使用二维码方式登录即可。  
3. 重复登录同一个账号时，该账号的上一次登录信息将会被覆盖。  
