注:此仓库为原项目的fork,原项目已被删除
# annie-public

### 原神游戏服务器的重新实现，Golang版本，由于开发过程中整合了部分非公开的系统框架，因此目前暂时仅将游戏服务器相关代码开源，但在Release发布中提供支持游戏服务器运行的必要核心组件的二进制文件

### 如何使用

##### 组件介绍

| 组件名          | 描述     |
|--------------|--------|
| air          | 注册中心   |
| annie_user   | 账号服务   |
| game_genshin | 游戏服务器  |
| gate_genshin | 网关服务   |
| gm_genshin   | GM后台服务 |

##### 目前依赖的第三方中间件有MySQL、MongoDB、Nats等，同时每个组件启动时将会读取名为application.toml的配置文件，其模板在源码的各个组件的cmd包下，请自行修改，其中游戏服务器还会读取配置表，即GameDataConfigTable目录
