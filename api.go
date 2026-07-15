package adapter_onebot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	adapter_manager "github.com/yeswearebot/adapter-manager"
	Pichubot "github.com/yeswearebot/go-Pichubot"
)

func (a *OnebotAdapter) ParseMessage(rawEvent any) []adapter_manager.MessageSegment {
	// 根据事件类型断言，提取原始消息字符串
	var rawMsg string
	switch event := rawEvent.(type) {
	case Pichubot.MessageGroup:
		rawMsg = event.RawMessage
	case Pichubot.MessagePrivate:
		rawMsg = event.RawMessage
	default:
		return []adapter_manager.MessageSegment{adapter_manager.Text("[不支持的事件类型]")}
	}
	return parseCQCode(rawMsg)
}

// SendGroupMsg 发送群消息 (将 Segment 转为 CQ 码发送)
func (a *OnebotAdapter) SendGroupMsg(groupID int64, message []adapter_manager.MessageSegment) error {
	cqString := segmentsToCQCode(message)
	_, err := Pichubot.SendGroupMsg(cqString, groupID)
	return err
}

// SendPrivateMsg 发送私聊消息
func (a *OnebotAdapter) SendPrivateMsg(userID int64, message []adapter_manager.MessageSegment) error {
	cqString := segmentsToCQCode(message)
	_, err := Pichubot.SendPrivateMsg(cqString, userID)
	return err
}

// 撤回消息
func (a *OnebotAdapter) DeleteMsg(messageID int32) error {
	return Pichubot.DeleteMsg(messageID)
}

// 群组单人禁言
func (a *OnebotAdapter) SetGroupBan(groupID, userID, duration int64) error {
	return Pichubot.SetGroupBan(groupID, userID, duration)
}

// SendMsg 发送消息（非 OneBot 原生，但通用）
func (a *OnebotAdapter) SendMsg(msgType, message string, toID int64) (map[string]interface{}, error) {
	return Pichubot.SendMsg(msgType, message, toID)
}

// GetMsg 获取消息
func (a *OnebotAdapter) GetMsg(messageID int32) (map[string]interface{}, error) {
	return Pichubot.GetMsg(messageID)
}

// GetForwardMsg 获取合并转发消息
func (a *OnebotAdapter) GetForwardMsg(id string) (map[string]interface{}, error) {
	return Pichubot.GetForwardMsg(id)
}

// SendLike 发送好友赞
func (a *OnebotAdapter) SendLike(userID int64, times int64) error {
	return Pichubot.SendLike(userID, times)
}

// SetGroupKick 群组踢人
func (a *OnebotAdapter) SetGroupKick(groupID, userID int64, rejectAddRequest bool) error {
	return Pichubot.SetGroupKick(groupID, userID, rejectAddRequest)
}

// SetGroupAnonymousBan 群组匿名用户禁言
func (a *OnebotAdapter) SetGroupAnonymousBan(groupID int64, anonymousFlag string, duration int64) error {
	return Pichubot.SetGroupAnonymousBan(groupID, anonymousFlag, duration)
}

// SetGroupWholeBan 群全员禁言
func (a *OnebotAdapter) SetGroupWholeBan(groupID int64, enable bool) error {
	return Pichubot.SetGroupWholeBan(groupID, enable)
}

// SetGroupAdmin 群组设置管理员
func (a *OnebotAdapter) SetGroupAdmin(groupID, userID int64, enable bool) error {
	return Pichubot.SetGroupAdmin(groupID, userID, enable)
}

// SetGroupAnonymous 群组匿名设置
func (a *OnebotAdapter) SetGroupAnonymous(groupID int64, enable bool) error {
	return Pichubot.SetGroupAnonymous(groupID, enable)
}

// SetGroupCard 设置群名片
func (a *OnebotAdapter) SetGroupCard(groupID, userID int64, card string) error {
	return Pichubot.SetGroupCard(groupID, userID, card)
}

// SetGroupName 设置群名
func (a *OnebotAdapter) SetGroupName(groupID int64, groupName string) error {
	return Pichubot.SetGroupName(groupID, groupName)
}

// SetGroupLeave 退群
func (a *OnebotAdapter) SetGroupLeave(groupID int64, isDismiss bool) error {
	return Pichubot.SetGroupLeave(groupID, isDismiss)
}

// SetGroupSpecialTitle 设置群组专属头衔
func (a *OnebotAdapter) SetGroupSpecialTitle(groupID, userID int64, specialTitle string) error {
	return Pichubot.SetGroupSpecialTitle(groupID, userID, specialTitle)
}

// SetFriendAddRequest 处理加好友请求
func (a *OnebotAdapter) SetFriendAddRequest(flag string, approve bool) error {
	return Pichubot.SetFriendAddRequest(flag, approve)
}

// SetGroupAddRequest 处理加群请求
func (a *OnebotAdapter) SetGroupAddRequest(flag string, approve bool, reason string) error {
	return Pichubot.SetGroupAddRequest(flag, approve, reason)
}

// SetGroupInviteRequest 处理加群邀请
func (a *OnebotAdapter) SetGroupInviteRequest(flag string, approve bool, reason string) error {
	return Pichubot.SetGroupInviteRequest(flag, approve, reason)
}

// GetLoginInfo 获取登录号信息
func (a *OnebotAdapter) GetLoginInfo() (map[string]interface{}, error) {
	return Pichubot.GetLoginInfo()
}

// GetImage 获取图片信息
func (a *OnebotAdapter) GetImage(file string) (map[string]interface{}, error) {
	return Pichubot.GetImage(file)
}

// OCRImage 图片OCR
func (a *OnebotAdapter) OCRImage(imageFile string) (map[string]interface{}, error) {
	return Pichubot.OCRImage(imageFile)
}

func (a *OnebotAdapter) GetGroupInfo(groupID int64, noCache bool) (any, error) {
	return Pichubot.GetGroupInfo(groupID, noCache)
}

func (a *OnebotAdapter) GetGroupMemberInfo(groupID, userID int64, noCache bool) (any, error) {
	return Pichubot.GetGroupMemberInfo(groupID, userID, noCache)
}

func (a *OnebotAdapter) GetGroupMemberList(groupID int64) (any, error) {
	return Pichubot.GetGroupMemberList(groupID)
}

// segmentsToCQCode 将标准消息段数组转换为 Onebot 的 CQ 码字符串
func segmentsToCQCode(segments []adapter_manager.MessageSegment) string {
	var sb strings.Builder
	for _, seg := range segments {
		switch seg.Type {
		case adapter_manager.SegText:
			sb.WriteString(seg.GetString("text"))
		case adapter_manager.SegImage:
			// 兼容 url 或 file 字段
			url := seg.GetString("url")
			if url == "" {
				url = seg.GetString("file")
			}
			sb.WriteString(fmt.Sprintf("[CQ:image,file=%s]", url))
		case adapter_manager.SegAt:
			sb.WriteString(fmt.Sprintf("[CQ:at,qq=%d]", seg.GetInt64("user_id")))
		case adapter_manager.SegReply:
			sb.WriteString(fmt.Sprintf("[CQ:reply,id=%d]", seg.GetInt64("message_id")))
		case adapter_manager.SegFace:
			sb.WriteString(fmt.Sprintf("[CQ:face,id=%s]", seg.GetString("id")))
		default:
			// 未知类型原样输出，或者直接丢弃，这里选择保留类型名以防万一
			sb.WriteString(fmt.Sprintf("[CQ:%s]", seg.Type))
		}
	}
	return sb.String()
}

// parseCQCode 将 Onebot 的 CQ 码字符串解析为标准消息段数组
var cqReg = regexp.MustCompile(`\[CQ:([^,\]]+)(?:,([^\]]*))?\]`)

func parseCQCode(rawMsg string) []adapter_manager.MessageSegment {
	segments := make([]adapter_manager.MessageSegment, 0)
	lastIndex := 0

	for _, match := range cqReg.FindAllStringSubmatchIndex(rawMsg, -1) {
		start := match[0]
		end := match[1]

		// 提取 CQ 码前面的纯文本
		if start > lastIndex {
			text := rawMsg[lastIndex:start]
			segments = append(segments, adapter_manager.Text(text))
		}

		// 提取 CQ 码类型和参数
		cqType := rawMsg[match[2]:match[3]]
		paramsStr := ""
		if match[4] != -1 {
			paramsStr = rawMsg[match[4]:match[5]]
		}

		data := make(map[string]any)
		if paramsStr != "" {
			for _, param := range strings.Split(paramsStr, ",") {
				kv := strings.SplitN(param, "=", 2)
				if len(kv) == 2 {
					data[kv[0]] = kv[1]
				}
			}
		}

		// 映射到标准 Segment
		switch cqType {
		case "text":
			// 一般 CQ 码不会有 text，但保险起见
			segments = append(segments, adapter_manager.Text(data["text"].(string)))
		case "image":
			// 把 file 字段映射为标准 url 字段 (如果存在)
			if file, ok := data["file"]; ok {
				data["url"] = file
			}
			segments = append(segments, adapter_manager.MessageSegment{Type: adapter_manager.SegImage, Data: data})
		case "at":
			if qqStr, ok := data["qq"]; ok {
				if id, err := strconv.ParseInt(qqStr.(string), 10, 64); err == nil {
					data["user_id"] = id
				}
			}
			segments = append(segments, adapter_manager.MessageSegment{Type: adapter_manager.SegAt, Data: data})
		case "reply":
			if idStr, ok := data["id"]; ok {
				if id, err := strconv.ParseInt(idStr.(string), 10, 64); err == nil {
					data["message_id"] = id
				}
			}
			segments = append(segments, adapter_manager.MessageSegment{Type: adapter_manager.SegReply, Data: data})
		case "face":
			segments = append(segments, adapter_manager.MessageSegment{Type: adapter_manager.SegFace, Data: data})
		default:
			// 其他未知类型直接保存
			segments = append(segments, adapter_manager.MessageSegment{Type: adapter_manager.SegmentType(cqType), Data: data})
		}

		lastIndex = end
	}

	// 提取末尾的纯文本
	if lastIndex < len(rawMsg) {
		text := rawMsg[lastIndex:]
		segments = append(segments, adapter_manager.Text(text))
	}

	return segments
}
