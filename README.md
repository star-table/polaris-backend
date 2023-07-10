# polaris-backend
极星后台，基于 gin 框架之上封装的业务框架。`polaris` 是极星的意思。   

## 常用命令行
所有的业务层服务可以在 service 目录下找到。由于是多个服务同时启动，因此关于多个服务的启动和关闭，需要一些脚本进行辅助

在运行这些脚本 shell 脚本之前，记得给对应文件**可执行**权限：`chmod +x xxx.sh `
 
### 重新生成 graphql
* `./app/gen.sh`

### 重新生成 facade 层
* `app/bin/build.sh`  

### app 目录下生成 restful 接口的 swagger 文档
* 可以尝试 `swag init --parseDependency --parseInternal`

### 启动服务
启动服务可以使用 `./bin/local/start.sh`，具体启动了哪些服务，可以参考 shell 脚本的具体实现

### 关闭服务
由于各个平台的方式不一样，这里以 Mac 系统为例。关闭服务可以使用 `ps -A | grep 'main --env=local' | awk '{print $1}' | xargs kill -9` 达到关闭多个 go 应用的目的

## 规范
从 graphql 的 schema 声明到 facade 的生成。根据 schema 生成的结构体位于 vo 层。

* 从 facade 中调用 service 中的接口，相关参数声明在 vo 层，如：`user_vo`。
* service 层调用 domain 层，domain 层使用到的参数结构体在 bo 层声明。（business object）
* domain 层可以调用 dao 层，dao 层的参数结构体在 po 层声明。（persistent object）

## 业务
### 动态
* 比如当新增一个任务后，需要生成一条动态，动态会记录到数据库中，其中需要存入前端展示的必要信息。
* 动态信息的入库，可以参考 service/platform/trendssvc/domain/trends_issue_domain.go 的 `assemblyTrendsBos` 函数
* 其他业务需要调用生成动态，则需要调用 http 接口：`/api/trendssvc/addIssueTrends`，调用方式可以参考 service/platform/projectsvc/domain/issue_work_hours_domain.go 中的 `PushIssueTrends` 函数

## 融合极星到无码系统
为避免对现有极星的影响，遂迁出一个新分支 feature/lesscode-polaris 作为融合版本的极星。使用 application.common.fuse_k8s.yaml 作为测试环境的配置文件，也就是说测试环境的 POL_ENV=fuse_k8s

它处于 k8s 网络环境中，一些配置需要注意：
* 日志路径，测试的 k8s 环境的路径需要是特定的，以 app 为例，其 default channel 的日志路径配置值为：`/data/logs/polaris-app/run.log`。
    * 同理，projectsvc 服务的日志文件配置值为：`/data/logs/polaris-projectsvc/run.log`

## 部署
### 私有化部署
目前只有移动这一个客户是私有化部署版本。私有化版本，需要更改配置文件中的 Application.RunMode 值为 3/4。且对应组织的 payLevel 改为 4，且 payEndTime 设置为一个久远的未来时间，
以保证在使用时不会过期。

### 注意事项
projectsvc 服务中有一些下载 excel 的服务，而这个服务对外暴露的容器未在本仓库中记录。其名称是 frontend-resourcesvc，对外提供文件下载服务。关于这个服务的详细信息可以找徐高人。

## reference 
- 暂无 
