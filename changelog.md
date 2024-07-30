<a name="unreleased"></a>
## [Unreleased]


<a name="v1.2.17"></a>
## [v1.2.17] - 2024-07-30
### Chore
- 规范化文件命名和日志清理
- update CHANGELOG
- **deps:** bump google.golang.org/grpc

### Feat
- **http:** 增加模式设置功能和请求头常量定义
- **logger:** 增强日志功能，支持从环境变量注入附加字段
- **trace:** 添加详细文档注释和新功能函数

### Fix
- **logger:** 修正 init 函数中 Caller 函数的调用深度

### Refactor
- **errors:** 移除 errors 包并使用本地定义的错误变量
- **exit:** 适配 exit 包的重构
- **exit:** 重构退出事件处理机制并添加单元测试
- **logger:** 代码重构和优化测试用例
- **logger:** 重构日志模块，简化日志选项和初始化流程

### Style
- **trace:** golangci-lint fix

### Pull Requests
- Merge pull request [#24](https://github.com/alioth-center/infrastructure/issues/24) from Jecosine/main
- Merge pull request [#23](https://github.com/alioth-center/infrastructure/issues/23) from sunist-c/main
- Merge pull request [#22](https://github.com/alioth-center/infrastructure/issues/22) from sunist-c/main
- Merge pull request [#21](https://github.com/alioth-center/infrastructure/issues/21) from alioth-center/dependabot/go_modules/go_modules-00bc3d2421


<a name="v1.2.16"></a>
## [v1.2.16] - 2024-06-28
### Chore
- golangci-lint fix
- 更新readme，添加贡献者列表
- update CHANGELOG

### Conf
- 配置 codecov 忽略无法测试的驱动与三方实现
- 使ChangLog工作流可以手动触发

### Feat
- 添加了 openai.client 的 embedding 接口与相关测试
- 添加了 openai.client 的测试
- 为 openai.client 添加了 CalculateToken 方法，可以计算文本所占用的 token 数量
- 添加三目运算符的函数实现与对应测试方法
- 添加彩云天气的第三方集成
- 为http.client添加一个接收url.Values的方法
- 添加新的Logger实现，它不会进行任何日志，并且适配测试文件

### Fix
- 修复 golangci-lint 错误
- 修复GitHub Actions中的报错

### Update
- 升级 golang 版本和 go.mod 依赖
- 更新 openai 的接口封装，适配最新的版本

### Pull Requests
- Merge pull request [#19](https://github.com/alioth-center/infrastructure/issues/19) from sunist-c/main
- Merge pull request [#17](https://github.com/alioth-center/infrastructure/issues/17) from sunist-c/main


<a name="v1.2.15"></a>
## [v1.2.15] - 2024-06-04
### Chore
- update CHANGELOG

### Conf
- 移除codefactor集成，将changelog更改为release时触发
- 配置codefactor和changelog的设置
- 更新 git-chglog 工作流配置，调整为在 pr 时触发
- 更新 git-chglog 工作流配置，添加 debug 输出
- 更新 git-chglog 工作流配置

### Fix
- 修复 HTTP 模块中的部分硬编码用例
- 修复依赖，执行golangci-lint

### Init
- 初始化 README 文档

### Optimize
- 优化了cli的实现结构，添加了config的写入功能

### Update
- 添加了http模块的部分单元测试用例
- 为 http 模块新增单元测试内容
- 为cli添加 go-prompt 依赖并简单定义了接口

### Pull Requests
- Merge pull request [#16](https://github.com/alioth-center/infrastructure/issues/16) from sunist-c/main
- Merge pull request [#14](https://github.com/alioth-center/infrastructure/issues/14) from sunist-c/cli


<a name="v1.2.14"></a>
## [v1.2.14] - 2024-05-12
### Feat
- 新增通过配置和注册的函数自动构建路由，启动服务的功能
- 为http请求和客户端添加自动附加tid的功能
- utils数组增加工具函数，新增获取指针函数
- 为http服务框架添加默认处理函数和跟踪中间件，拓展了部分功能与中间件

### Fix
- 修复 997e87db 为上级节点添加 gin 中间件会影响子节点的问题

### Optimize
- 项目执行golang-ci fix

### Update
- 更新golangci超时配置

### Pull Requests
- Merge pull request [#12](https://github.com/alioth-center/infrastructure/issues/12) from sunist-c/main


<a name="v1.2.13"></a>
## [v1.2.13] - 2024-04-26
### Delete
- 去除eto相关实现

### Feat
- 更新了日志器的枚举声明和注释，引入了CLS日志器实现
- 添加API端点组与端点构建器的实现
- 添加了货币种类的ISO定义枚举
- 添加了时间戳的rfc标准定义
- 为utils包添加了方法
- 添加了eto拓展，部分实现
- 添加地区和语言的枚举定义

### Fix
- 修复了 73e3dc1e 中不当的测试端口复用

### Optimize
- 优化了concurrency中的代码实现
- 清理了废弃函数
- 清理了database冗余函数

### Update
- 更新go依赖
- 为时间戳标准添加注释
- 更新go.mod

### Pull Requests
- Merge pull request [#11](https://github.com/alioth-center/infrastructure/issues/11) from sunist-c/main


<a name="v1.2.12"></a>
## [v1.2.12] - 2024-03-22
### Feat
- 重构了http服务端实现，优化了预处理方式，修复了重复响应的错误与500覆写的错误
- 重构了StringTemplate的实现，引入了反射处理，适配了单元测试
- 优化thirdparty/openai的错误处理机制
- 重构了trace.Context的实现，优化了方法签名
- 重构了http客户端实现，改为了请求构建器加客户端的方式，并支持mock
- 为values包添加了字符串模板构建，并完成了values包的单元测试

### Update
- 更新项目依赖

### Pull Requests
- Merge pull request [#10](https://github.com/alioth-center/infrastructure/issues/10) from sunist-c/main


<a name="v1.2.11"></a>
## [v1.2.11] - 2024-02-05
### Feat
- 重构了时区包，优化了设计逻辑，添加了UTC时区，拓展了工具方法

### Update
- 更新了trace工具包，废弃了部分函数，新增了部分函数，并进行适配
- 更新了日志器，暴露了部分方法，优化了无效方法

### Pull Requests
- Merge pull request [#9](https://github.com/alioth-center/infrastructure/issues/9) from sunist-c/main


<a name="v1.2.10"></a>
## [v1.2.10] - 2024-02-04
### Conf
- 更新gitignure以忽略本地调试日志，更新codecov配置文件降低CI要求

### Feat
- 为数据库接口新增事务处理功能
- 为数据库orm拓展新增事务接口
- 为数据库模块新增了一个SQL模板构建器，可以自动进行模板替换并适配注入防御
- 为logger包新增了一个构造函数，用于自定义日志实现
- 优化了工具包的values和concurrency模块的代码，添加了注释和使用示例
- 为cache模块新增了一个标识自身的DriverNam常量和对应方法
- 为http请求客户端接口新增了两个方法，用于从现有请求复制一个请求进行操作
- 为concurrency包添加了一个并发安全的slice实现和两个并发安全的map实现
- 为encrypt包添加了rsa加解密的实现

### Fix
- 修复了可能导致问题的unsafe使用，更新了hashmap的哈希实现

### Optimize
- 清理了rpc包中部分无效代码

### Update
- 更新了部分依赖项
- 为utils包的部分函数添加了一些样例和注释

### Pull Requests
- Merge pull request [#8](https://github.com/alioth-center/infrastructure/issues/8) from sunist-c/main


<a name="v1.2.9"></a>
## [v1.2.9] - 2024-01-25
### Conf
- 更新codecov的配置文件，对patch进行配置
- 更新codecov的配置文件，降低CI通过要求

### Feat
- 为http协议封装了框架，在gin的基础上自动解析参数与封装响应，并实现了部分注入接口
- 添加了字符串、数组和数字类型的工具函数，主要为类型转换的封装

### Update
- 添加了http模块的部分单元测试用例
- 更新rpc的部分函数签名，与http统一
- 将go-gin引入依赖项
- 升级部分依赖项以解决 CVE-2023-48795 的影响

### Pull Requests
- Merge pull request [#6](https://github.com/alioth-center/infrastructure/issues/6) from sunist-c/main


<a name="v1.2.8"></a>
## [v1.2.8] - 2024-01-22
### Conf
- 更新错误的工作流配置文件
- 更新 codecov 的配置文件，在PullRequest时开启CI检查
- 更新 codecov 的配置文件，更新了版本和 pr comment
- 更新了codecov工作流与.gitignore配置

### Fix
- 修复redis实现的计数器接口测试时的依赖项
- 修复 10538a75 更新的redis计数器在key不存在时可能返回Failed的错误
- 修复 10538a75 更新的内存计数器接口可能修改非计数器的错误
- 修复 10538a75 更新的内存缓存模块的update方法可能导致的并发错误

### Optimize
- 优化了缓存模块，添加了单元测试，重构了计数器接口，适配了影响内容

### Update
- 添加了redis实现计数器的单元测试用例，并优化了部分redis实现的缓存接口的测试用例
- 添加了计数器的单元测试用例，并为内存缓存的单元测试增加了部分测试样例
- 修复了部分错误的redis样例

### Pull Requests
- Merge pull request [#5](https://github.com/alioth-center/infrastructure/issues/5) from sunist-c/main


<a name="v1.2.7"></a>
## [v1.2.7] - 2024-01-15
### Conf
- 更新 codecov 的配置文件，先进行测试再lint，关闭了部分linter
- 配置自动化工作流程并适配部分测试文件
- 添加codecov的配置文件

### Feat
- 为数据库模块新增拓展机制，并将orm功能移动至拓展内
- 为数据库模块添加了五个方法与其默认实现，并做出了适配

### Fix
- 修复 logger 会关闭 stdout/stderr 的问题，可能导致其余功能异常

### Pull Requests
- Merge pull request [#3](https://github.com/alioth-center/infrastructure/issues/3) from sunist-c/main


<a name="v1.2.6"></a>
## [v1.2.6] - 2023-12-05
### Feat
- 添加了部分工具函数
- 为trace包添加了一个fork函数，该函数可以从已有的context中复制trace_id供并行使用

### Fix
- 修复 a94a4d40 修改导致的单元测试过时问题，修复创建 redis 实例中的错误

### Update
- 更新了事件处理机制，并进行了适配
- 更新了飞书的事件订阅验证逻辑

### Pull Requests
- Merge pull request [#2](https://github.com/alioth-center/infrastructure/issues/2) from sunist-c/main


<a name="v1.2.5"></a>
## [v1.2.5] - 2023-11-27
### Feat
- 添加了部分飞书事件的集成
- 添加了飞书事件的集成，提供了事件相关的请求封装与转换函数，提供了签名验证方法，提供了默认实现

### Fix
- 修复 c5ad296f 中不恰当的阻塞方法，该方法会导致未在 main 方法中调用阻塞函数时程序无法从 ctrl+c 退出

### Update
- 修复了飞书集成中发送文本信息的适配错误，添加了下载语音文件的接口，将音频长度计算方法公有
- 更改了飞书消息接收方的部分字段名称


<a name="v1.2.4"></a>
## [v1.2.4] - 2023-11-24
### Feat
- 为 openai 集成添加自定义 UserAgent 的功能，并适配相关逻辑
- 更新了model object的定义
- 修了model中的小typo，补充了client中ListFiles的query构建
- 实现了Fine-tuning和File API(retrieve content除外)

### Fix
- 修复 36859bf3 中的部分无效分配警告，该警告所述内容与程序预期不符
- 修复 5bd655af 中封装的聊天完成接口的json字段错误

### Pull Requests
- Merge pull request [#1](https://github.com/alioth-center/infrastructure/issues/1) from Jecosine/main: 追加了Fine tuning 和 File API的实现


<a name="v1.2.3"></a>
## [v1.2.3] - 2023-11-23
### Feat
- 添加了openai的客户端实现，实现了大部分非预览接口
- 添加了openai的配置模型封装
- 添加了openai的请求与响应模型封装
- 添加了部分常用http status code封装

### Fix
- 修复 cff8fdb6 提交中的multipart文件处理错误


<a name="v1.2.2"></a>
## [v1.2.2] - 2023-11-21
### Feat
- 为飞书集成添加图片和音频消息发送接口，进行了部分接口测试
- 添加飞书集成，可以发送文本信息和上传文件
- 添加一个http的客户端实现，封装了请求和响应
- 将utils包细分功能模块

### Fix
- 解决CVE-2023-39325

### License
- 将许可证设置为MIT协议
- 将许可证设置为MIT协议

### Update
- 适配utils包变更逻辑


<a name="v1.2.1"></a>
## [v1.2.1] - 2023-11-07
### Feat
- 为rpc模块添加了一个请求限制器
- 更新了缓存模块，实现了基于内存map的计数器

### Update
- 改进了rpc请求限制器的实现


<a name="v1.2.0"></a>
## [v1.2.0] - 2023-11-06
### Feat
- 添加了redis依赖项
- 修改了cache/memory的部分配置字段
- 基于内建map实现了kv缓存接口，添加了主动淘汰策略，进行了单元测试
- 基于redis实现了kv缓存接口
- 添加了kv缓存的接口定义
- 为sqlite/mysql/postgres的数据库驱动都添加了优雅退出事件
- 新增了邮件发送功能，使用smtp进行邮件发送
- 更新了rpc服务封装，现在可以自动进行大部份rpc相关的操作了
- 添加了追踪包的堆栈追踪函数，用于打印堆栈
- 在工具包中添加了await/async的封装，默认nil值生成的函数
- 更新了日志器的部分不一致配置字段，适配了日志器的优雅退出逻辑
- 更新了gorm封装，添加了分页查询接口，适配了优雅退出逻辑
- 更新了exit功能，优化了主函数退出逻辑，添加了信息打印

### Fix
- 添加了工具包中缺失的部分错误定义
- 移除了无意义的注册退出事件时的堆栈跟踪


<a name="v1.1.0"></a>
## [v1.1.0] - 2023-11-02
### Feat
- 新增rpc相关的支持，使用类似gin的模式


<a name="v1.0.1"></a>
## [v1.0.1] - 2023-11-01
### Feat
- 将时区编码进代码中，并添加单元测试
- 将text包更名为utils，并新增部分编码与随机功能函数
- 在封装的数据库执行中默认使用事务，并更新mysql单元测试文件避免产生误解


<a name="v1.0.0"></a>
## v1.0.0 - 2023-10-26
### Feat
- 初始化基础设施建设，添加若干个模块

### Fix
- 删除无效文件


[Unreleased]: https://github.com/alioth-center/infrastructure/compare/v1.2.17...HEAD
[v1.2.17]: https://github.com/alioth-center/infrastructure/compare/v1.2.16...v1.2.17
[v1.2.16]: https://github.com/alioth-center/infrastructure/compare/v1.2.15...v1.2.16
[v1.2.15]: https://github.com/alioth-center/infrastructure/compare/v1.2.14...v1.2.15
[v1.2.14]: https://github.com/alioth-center/infrastructure/compare/v1.2.13...v1.2.14
[v1.2.13]: https://github.com/alioth-center/infrastructure/compare/v1.2.12...v1.2.13
[v1.2.12]: https://github.com/alioth-center/infrastructure/compare/v1.2.11...v1.2.12
[v1.2.11]: https://github.com/alioth-center/infrastructure/compare/v1.2.10...v1.2.11
[v1.2.10]: https://github.com/alioth-center/infrastructure/compare/v1.2.9...v1.2.10
[v1.2.9]: https://github.com/alioth-center/infrastructure/compare/v1.2.8...v1.2.9
[v1.2.8]: https://github.com/alioth-center/infrastructure/compare/v1.2.7...v1.2.8
[v1.2.7]: https://github.com/alioth-center/infrastructure/compare/v1.2.6...v1.2.7
[v1.2.6]: https://github.com/alioth-center/infrastructure/compare/v1.2.5...v1.2.6
[v1.2.5]: https://github.com/alioth-center/infrastructure/compare/v1.2.4...v1.2.5
[v1.2.4]: https://github.com/alioth-center/infrastructure/compare/v1.2.3...v1.2.4
[v1.2.3]: https://github.com/alioth-center/infrastructure/compare/v1.2.2...v1.2.3
[v1.2.2]: https://github.com/alioth-center/infrastructure/compare/v1.2.1...v1.2.2
[v1.2.1]: https://github.com/alioth-center/infrastructure/compare/v1.2.0...v1.2.1
[v1.2.0]: https://github.com/alioth-center/infrastructure/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/alioth-center/infrastructure/compare/v1.0.1...v1.1.0
[v1.0.1]: https://github.com/alioth-center/infrastructure/compare/v1.0.0...v1.0.1
