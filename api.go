package adapter_onebot

import (
	Pichubot "github.com/yeswearebot/go-Pichubot"
)

// 发送群消息
func (a *OnebotAdapter) SendGroupMsg(groupID int64, message string) error {
	_, err := Pichubot.SendGroupMsg(message, groupID)
	return err
}

// 发送私聊消息
func (a *OnebotAdapter) SendPrivateMsg(userID int64, message string) error {
	_, err := Pichubot.SendPrivateMsg(message, userID)
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

func (a *OnebotAdapter) GetGroupInfo(groupID int64, noCache bool) (map[string]interface{}, error) {
	return Pichubot.GetGroupInfo(groupID, noCache)
}

func (a *OnebotAdapter) GetGroupMemberInfo(groupID, userID int64, noCache bool) (map[string]interface{}, error) {
	return Pichubot.GetGroupMemberInfo(groupID, userID, noCache)
}

func (a *OnebotAdapter) GetGroupMemberList(groupID int64) (map[string]interface{}, error) {
	return Pichubot.GetGroupMemberList(groupID)
}
