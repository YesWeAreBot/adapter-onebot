# adapter-onebot
OneBot 协议适配器，用于将 [PichuBot](https://github.com/0ojixueseno0/go-Pichubot) 的 OneBot 事件桥接到 [yes-core](https://github.com/YesWeAreBot/yes-core) 微内核的事件总线，并提供统一的 API 调用接口。

## 安装
在 `go.mod` 中添加依赖：
```bash
go get github.com/yeswearebot/adapter-onebot
```
然后，在你的主程序中匿名导入：
```go
import _ "github.com/yeswearebot/adapter-onebot"
```
> **注意**：需要同时导入配置中心插件 `_ "github.com/yeswearebot/plugin_config"`。
## 配置
适配器通过 `plugin_config` 读取 `onebot` 键下的配置。支持 `yaml`/`json`/`toml`，示例（`config.yaml`）：
```yaml
onebot:
  host: "127.0.0.1:6700"      # WebSocket 服务地址，必填
  path: "/ws"                 # WebSocket 路径，默认 /ws
  token: ""                   # 鉴权 Token（如有）
  msg_await: false            # 是否启用消息随机延迟（防风控）
```
| 字段      | 类型    | 说明                               |
| --------- | ------- | ---------------------------------- |
| host      | string  | OneBot WebSocket 服务地址，如 `127.0.0.1:6700` |
| path      | string  | WebSocket 路径，通常为 `/ws`       |
| token     | string  | 鉴权 Token（可选）                 |
| msg_await | bool    | 是否启用消息随机延迟                |
## 使用示例
### 1. 基本集成（main.go）
```go
package main
import (
    "github.com/yeswearebot/yes-core/core"
    
    _ "github.com/yeswearebot/plugin_config"           // 配置中心
    _ "github.com/yeswearebot/adapter-onebot"           // OneBot 适配器
    _ "github.com/your-repo/my-plugin"                 // 你的业务插件
)
func main() {
    app := core.NewApp()
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```
### 2. 编写业务插件：订阅事件
```go
package my_plugin
import (
    "fmt"
    "github.com/yeswearebot/yes-core/core"
    adapter_manager "github.com/yeswearebot/adapter-manager" // 引入消息段标准
)
type MyPlugin struct{}
func init() {
    core.Register(func() core.Plugin { return &MyPlugin{} })
}
func (p *MyPlugin) Name() string        { return "my-plugin" }
func (p *MyPlugin) DependsOn() []string { return []string{"adapter-onebot"} } // 确保适配器先启动
func (p *MyPlugin) Init(ctx *core.SystemContext) error {
    // 订阅消息事件
    ctx.Events.Subscribe("adapter.raw.message", func(payload any) {
        data, ok := payload.(map[string]any)
        if !ok { return }
        eventType, _ := data["event_type"].(string)
        switch eventType {
        case "group_msg":
            groupID := data["scene_id"].(int64)
            userID := data["user_id"].(int64)
            // 消息字段可能是字符串（CQ码）或已解析的消息段，根据需要处理
            message := data["message"]
            fmt.Printf("[群消息] 群 %d 用户 %d 说: %v\n", groupID, userID, message)
            
            // 获取原始事件进行更细致的操作
            if rawEvent, ok := data["raw_event"]; ok {
                // 根据 rawEvent 的具体类型进行断言和操作
                // 例如: if event, ok := rawEvent.(Pichubot.MessageGroup); ok { ... }
            }
        case "private_msg":
            userID := data["user_id"].(int64)
            message := data["message"]
            fmt.Printf("[私聊] 用户 %d 说: %v\n", userID, message)
        }
    })
    return nil
}
func (p *MyPlugin) Start(ctx *core.SystemContext) error { return nil }
func (p *MyPlugin) Stop(ctx *core.SystemContext) error  { return nil }
```
### 3. 主动调用适配器 API：发送结构化消息
```go
// 获取适配器实例
raw, exists := ctx.Registry.Get("adapter-onebot")
if !exists { return }
adapter := raw.(*adapter_onebot.OnebotAdapter)
// 构建结构化消息段（推荐）
message := []adapter_manager.MessageSegment{
    adapter_manager.Text("Hello, "),
    adapter_manager.At(userID), // 使用 At 消息段 @用户
    adapter_manager.Text("! 这是来自 yes-core 的图片消息:"),
    adapter_manager.Image("https://example.com/image.jpg"), // 支持 URL 或文件路径
}
// 发送群消息（会自动转换为 CQ 码）
err := adapter.SendGroupMsg(123456, message)
if err != nil {
    fmt.Println("发送失败:", err)
}
// 发送私聊消息
err = adapter.SendPrivateMsg(userID, message)
```
### 4. CQ 码与消息段转换
适配器内置了 CQ 码和标准消息段的转换逻辑，你可以在需要时直接使用：
```go
// 将 CQ 码字符串解析为消息段
rawMsg := "[CQ:at,qq=123456] 你好！[CQ:image,file=https://example.com/img.jpg]"
segments := adapter.ParseMessage(rawEvent) // 在事件处理中，rawEvent 是原始事件结构体
// segments 现在是一个包含 Text, At, Image 类型消息段的数组
// 将消息段转换为 CQ 码字符串
message := []adapter_manager.MessageSegment{
    adapter_manager.Text("测试 "),
    adapter_manager.Face("1"), // 发送一个表情
}
cqCode := adapter.segmentsToCQCode(message) // 注意：此方法未导出，仅供内部使用
// 发送时会自动转换，通常无需手动调用
```
## 事件说明
所有事件均以 `adapter.raw.message` 为主题发布，载荷为 `map[string]any`，包含以下通用字段：
| 字段名       | 类型               | 说明                           |
| ------------ | ------------------ | ------------------------------ |
| platform     | string             | 固定为 `"adapter-onebot"`      |
| event_type   | string             | 具体事件类型（见下表）         |
| scene        | string             | `"group"` 或 `"private"`       |
| scene_id     | int64              | 群号或对方 QQ 号               |
| user_id      | int64              | 发送者 QQ 号                   |
| self_id      | int64              | 机器人自身 QQ 号               |
| raw_event    | 原始结构体         | PichuBot 原始事件对象          |
| …            | 其他字段（依事件而定） | 参见下表各事件详细字段        |
**事件类型列表（`event_type`）**：
| event_type          | 对应 PichuBot 事件     | 场景      | 说明                     | 事件载荷额外字段                                                                                              |
| ------------------- | ---------------------- | --------- | ------------------------ | ------------------------------------------------------------------------------------------------------------- |
| `group_msg`         | `MessageGroup`         | group     | 群消息                   | `message` (消息内容), `message_id` (消息ID), `raw_message` (原始CQ码字符串), `nickname` (发送者昵称)              |
| `private_msg`       | `MessagePrivate`       | private   | 私聊消息                 | `message`, `message_id`, `raw_message`, `nickname`                                                             |
| `group_upload`      | `GroupUpload`          | group     | 群文件上传               | `file` (文件对象: `Id`, `Name`, `Size`, `Busid`)                                                               |
| `group_admin`       | `GroupAdmin`           | group     | 群管理员变动             | `sub_type` (`set`/`unset`)                                                                                   |
| `group_decrease`    | `GroupDecrease`        | group     | 群成员减少               | `sub_type` (`leave`/`kick`/`kick_me`), `operator_id` (操作者QQ)                                              |
| `group_increase`    | `GroupIncrease`        | group     | 群成员增加               | `sub_type` (`approve`/`invite`), `operator_id`                                                               |
| `group_ban`         | `GroupBan`             | group     | 群聊禁言                 | `sub_type` (`ban`/`lift_ban`), `duration` (禁言时长秒), `operator_id`                                        |
| `friend_add`        | `FriendAdd`            | private   | 好友添加                 | -                                                                                                             |
| `group_recall`      | `GroupRecall`          | group     | 群消息撤回               | `message_id` (被撤回的消息ID), `operator_id`                                                                  |
| `friend_recall`     | `FriendRecall`         | private   | 好友消息撤回             | `message_id`                                                                                                  |
| `notify`            | `Notify`               | group     | 戳一戳/运气王/荣誉等     | `sub_type` (`poke`/`lucky_king`/`honor`), `target_id` (被戳者/运气王QQ, 仅部分子类型), `honor_type` (荣誉类型) |
| `friend_request`    | `FriendRequest`        | private   | 加好友请求               | `flag` (请求标识), `comment` (验证消息)                                                                      |
| `group_request`     | `GroupRequest`         | group     | 加群请求/邀请            | `sub_type` (`add`/`invite`), `flag`, `comment`                                                               |
| `meta_lifecycle`    | `MetaLifecycle`        | -         | 生命周期事件             | `sub_type` (`enable`/`disable`/`connect`)                                                                    |
| `meta_heartbeat`    | `MetaHeartbeat`        | -         | 心跳包                   | `interval` (间隔毫秒), `status` (状态信息)                                                                   |
## API 方法（动作接口）
适配器提供以下公开方法，可在获取实例后调用。
### 消息发送
| 方法                                                                   | 说明                       |
| ---------------------------------------------------------------------- | -------------------------- |
| `SendGroupMsg(groupID int64, message []MessageSegment) error`           | 发送群聊消息（消息段自动转CQ码） |
| `SendPrivateMsg(userID int64, message []MessageSegment) error`          | 发送私聊消息（消息段自动转CQ码） |
| `SendMsg(msgType, message string, toID int64) (map[string]interface{}, error)` | 通用发送（非原生，直接发送CQ码字符串） |
### 消息管理
| 方法                                                                   | 说明                       |
| ---------------------------------------------------------------------- | -------------------------- |
| `DeleteMsg(messageID int32) error`                                     | 撤回消息                   |
| `GetMsg(messageID int32) (map[string]interface{}, error)`              | 获取消息详情               |
| `GetForwardMsg(id string) (map[string]interface{}, error)`             | 获取合并转发消息内容       |
### 群组管理
| 方法                                                                                                      | 说明                       |
| --------------------------------------------------------------------------------------------------------- | -------------------------- |
| `SetGroupKick(groupID, userID int64, rejectAddRequest bool) error`                                        | 踢人                       |
| `SetGroupBan(groupID, userID, duration int64) error`                                                      | 禁言/解禁（单位秒）        |
| `SetGroupAnonymousBan(groupID int64, anonymousFlag string, duration int64) error`                         | 匿名用户禁言               |
| `SetGroupWholeBan(groupID int64, enable bool) error`                                                      | 全员禁言                   |
| `SetGroupAdmin(groupID, userID int64, enable bool) error`                                                 | 设置/取消管理员            |
| `SetGroupAnonymous(groupID int64, enable bool) error`                                                     | 允许/禁止匿名              |
| `SetGroupCard(groupID, userID int64, card string) error`                                                  | 设置群名片                 |
| `SetGroupName(groupID int64, groupName string) error`                                                     | 修改群名                   |
| `SetGroupLeave(groupID int64, isDismiss bool) error`                                                      | 退群/解散群               |
| `SetGroupSpecialTitle(groupID, userID int64, specialTitle string) error`                                  | 设置专属头衔               |
| `GetGroupInfo(groupID int64, noCache bool) (any, error)`                                                 | 获取群信息                 |
| `GetGroupMemberInfo(groupID, userID int64, noCache bool) (any, error)`                                    | 获取群成员信息             |
| `GetGroupMemberList(groupID int64) (any, error)`                                                         | 获取群成员列表             |
### 请求处理
| method                                                                     | 说明                       |
| -------------------------------------------------------------------------- | -------------------------- |
| `SetFriendAddRequest(flag string, approve bool) error`                     | 处理加好友请求             |
| `SetGroupAddRequest(flag string, approve bool, reason string) error`       | 处理加群请求               |
| `SetGroupInviteRequest(flag string, approve bool, reason string) error`    | 处理加群邀请               |
### 其他
| method                                                | 说明                       |
| ----------------------------------------------------- | -------------------------- |
| `SendLike(userID int64, times int64) error`           | 发送好友赞                 |
| `GetLoginInfo() (map[string]interface{}, error)`       | 获取机器人登录信息         |
| `GetImage(file string) (map[string]interface{}, error)` | 获取图片信息              |
| `OCRImage(imageFile string) (map[string]interface{}, error)` | 图片文字识别（OCR） |
> 💡 所有方法的参数与 PichuBot 原生保持一致，详细说明请参考 PichuBot 文档。
## 注意事项
- 适配器依赖 `plugin_config` 插件，请确保已在 main 中导入。
- 配置文件中 `host` 和 `path` 必须正确，否则连接失败。
- WebSocket 连接在 `Start` 阶段异步建立，不影响框架启动。
- 事件载荷中的 `message` 字段在消息事件中是**已解析的消息内容**，可能不再是原始 CQ 码字符串。如需原始 CQ 码字符串，请访问 `raw_message` 字段。
- 请求类事件（好友/群请求）需要调用对应的处理 API（如 `SetFriendAddRequest`）才能完成响应，否则请求会超时。
## 贡献
欢迎提交 Issue 和 Pull Request。如有任何问题，请在仓库中提出。