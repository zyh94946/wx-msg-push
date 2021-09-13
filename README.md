# 企业微信应用消息推送服务 独立部署版

消息效果和企业微信注册相关见云函数版 [wx-msg-push-tencent](https://github.com/zyh94946/wx-msg-push-tencent)

## 特性

1. 支持多企业微信应用
2. 配置文件修改保存即生效，不需要重启服务

## 消息限制

目前支持发送的应用消息类型为：

1. 图文消息(mpnews)，消息内容支持html标签，不超过666K个字节，会自动生成摘要替换br标签为换行符，过滤html标签。
2. 文本消息(text)，消息内容最长不超过2048个字节，超过将截断。
3. markdown消息(markdown)，markdown内容，最长不超过2048个字节，必须是utf8编码。([支持的markdown语法](https://work.weixin.qq.com/api/doc/90000/90135/90236#%E6%94%AF%E6%8C%81%E7%9A%84markdown%E8%AF%AD%E6%B3%95))

每企业消息发送不可超过帐号上限数*30人次/天（注：若调用api一次发给1000人，算1000人次；若企业帐号上限是500人，则每天可发送15000人次的消息）

每应用对同一个成员不可超过30条/分，超过部分会被丢弃不下发

默认未启用重复消息推送检查，如果修改请参考配置文件。

## 部署运行

1. 推荐直接从 releases 中下载预编译的包。或 clone 本项目后执行编译 `make build`，编译其它平台参考 `Makefile`。或者 `go get -u github.com/zyh94946/wx-msg-push` 直接安装可执行文件至 `$GOPATH/bin`。

2. 修改配置文件 `config.toml`, 中的运行端口、企业微信配置。

3. 运行 `wx-msg-push server -c config.toml`

请自行处理进程守护，包括不限于以下方式

- `&`
- `nohup`
- `tmux`
- `supervisor`
- `crontab`
- `service`

## 使用方法

使用方法和返回结果和云函数版一致。请自行处理域名绑定和转发。

消息类型值：`text` 代表文本消息，`mpnews` 代表图文消息，`markdown` 代表 markdown 消息。为兼容旧版本，不传默认为图文消息。

支持推送消息至指定的 `touser`, `toparty`, `totag`。不传默认设置 `touser=@all`

GET方式

`http://ip:port/CORP_SECRET?title=消息标题&content=消息内容&type=消息类型`

POST方式

```bash
$ curl --location --request POST 'http://ip:port/CORP_SECRET' \
--header 'Content-Type: application/json;charset=utf-8' \
--data-raw '{"title":"消息标题","content":"消息内容","type":"消息类型"}'
```

发送成功状态码返回200，`"Content-Type":"application/json"` body `{"errorCode":0,"errorMessage":""}` 。

## 其它

发送失败问题排查：

- 请检查配置是否正确
- 请检查发送url中CORP_SECRET是否正确
- 根据返回结果、日志找原因


## TODO

- docker 镜像
- 日志处理
- 后台执行
- 重启优化
-
