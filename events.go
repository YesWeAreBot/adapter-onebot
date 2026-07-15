package adapter_onebot

import (
	Pichubot "github.com/yeswearebot/go-Pichubot"
)

// 向底层驱动注册监听器
func (a *OnebotAdapter) registerEventListeners() {
	// 消息事件
	Pichubot.Listeners.OnGroupMsg = append(Pichubot.Listeners.OnGroupMsg, a.handleGroupMsg)
	Pichubot.Listeners.OnPrivateMsg = append(Pichubot.Listeners.OnPrivateMsg, a.handlePrivateMsg)

	// 提醒事件
	Pichubot.Listeners.OnGroupUpload = append(Pichubot.Listeners.OnGroupUpload, a.handleGroupUpload)
	Pichubot.Listeners.OnGroupAdmin = append(Pichubot.Listeners.OnGroupAdmin, a.handleGroupAdmin)
	Pichubot.Listeners.OnGroupDecrease = append(Pichubot.Listeners.OnGroupDecrease, a.handleGroupDecrease)
	Pichubot.Listeners.OnGroupIncrease = append(Pichubot.Listeners.OnGroupIncrease, a.handleGroupIncrease)
	Pichubot.Listeners.OnGroupBan = append(Pichubot.Listeners.OnGroupBan, a.handleGroupBan)
	Pichubot.Listeners.OnFriendAdd = append(Pichubot.Listeners.OnFriendAdd, a.handleFriendAdd)
	Pichubot.Listeners.OnGroupRecall = append(Pichubot.Listeners.OnGroupRecall, a.handleGroupRecall)
	Pichubot.Listeners.OnFriendRecall = append(Pichubot.Listeners.OnFriendRecall, a.handleFriendRecall)
	Pichubot.Listeners.OnNotify = append(Pichubot.Listeners.OnNotify, a.handleNotify)

	// 请求事件
	Pichubot.Listeners.OnFriendRequest = append(Pichubot.Listeners.OnFriendRequest, a.handleFriendRequest)
	Pichubot.Listeners.OnGroupRequest = append(Pichubot.Listeners.OnGroupRequest, a.handleGroupRequest)

	// 元事件
	Pichubot.Listeners.OnMetaLifecycle = append(Pichubot.Listeners.OnMetaLifecycle, a.handleMetaLifecycle)
	Pichubot.Listeners.OnMetaHeartbeat = append(Pichubot.Listeners.OnMetaHeartbeat, a.handleMetaHeartbeat)
}

func (a *OnebotAdapter) handleGroupMsg(event Pichubot.MessageGroup) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "group_msg",
		"scene":      "group",
		"scene_id":   event.GroupID,
		"user_id":    event.UserID,
		"nickname":   event.Sender.Nickname,
		"message":    event.Message,
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

// 处理私聊消息
func (a *OnebotAdapter) handlePrivateMsg(event Pichubot.MessagePrivate) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "private_msg",
		"scene":      "private",
		"scene_id":   event.UserID,
		"user_id":    event.UserID,
		"nickname":   event.Sender.Nickname,
		"message":    event.Message,
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleGroupUpload(event Pichubot.GroupUpload) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "group_upload",
		"scene":      "group",
		"group_id":   event.GroupId,
		"user_id":    event.UserId,
		"self_id":    event.SelfId,
		"file":       event.File, // 包含 Id, Name, Size, Busid
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleGroupAdmin(event Pichubot.GroupAdmin) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "group_admin",
		"scene":      "group",
		"group_id":   event.GroupId,
		"user_id":    event.UserId,
		"self_id":    event.SelfId,
		"sub_type":   event.SubType, // set / unset
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleGroupDecrease(event Pichubot.GroupDecrease) {
	payload := map[string]any{
		"platform":    "onebot",
		"event_type":  "group_decrease",
		"scene":       "group",
		"group_id":    event.GroupId,
		"user_id":     event.UserId,
		"operator_id": event.OperatorId,
		"self_id":     event.SelfId,
		"sub_type":    event.SubType, // leave / kick / kick_me
		"raw_event":   event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleGroupIncrease(event Pichubot.GroupIncrease) {
	payload := map[string]any{
		"platform":    "onebot",
		"event_type":  "group_increase",
		"scene":       "group",
		"group_id":    event.GroupId,
		"user_id":     event.UserId,
		"operator_id": event.OperatorId,
		"self_id":     event.SelfId,
		"sub_type":    event.SubType, // approve / invite
		"raw_event":   event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleGroupBan(event Pichubot.GroupBan) {
	payload := map[string]any{
		"platform":    "onebot",
		"event_type":  "group_ban",
		"scene":       "group",
		"group_id":    event.GroupId,
		"user_id":     event.UserId,
		"operator_id": event.OperatorId,
		"self_id":     event.SelfId,
		"sub_type":    event.SubType, // ban / lift_ban
		"duration":    event.Duration,
		"raw_event":   event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleFriendAdd(event Pichubot.FriendAdd) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "friend_add",
		"scene":      "private",
		"user_id":    event.UserId,
		"self_id":    event.SelfId,
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleGroupRecall(event Pichubot.GroupRecall) {
	payload := map[string]any{
		"platform":    "onebot",
		"event_type":  "group_recall",
		"scene":       "group",
		"group_id":    event.GroupId,
		"user_id":     event.UserId,
		"operator_id": event.OperatorId,
		"self_id":     event.SelfId,
		"message_id":  event.MessageId,
		"raw_event":   event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleFriendRecall(event Pichubot.FriendRecall) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "friend_recall",
		"scene":      "private",
		"user_id":    event.UserId,
		"self_id":    event.SelfId,
		"message_id": event.MessageId,
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleNotify(event Pichubot.Notify) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "notify",
		"scene":      "group", // Notify 均为群内事件
		"group_id":   event.GroupId,
		"user_id":    event.UserId,
		"self_id":    event.SelfId,
		"sub_type":   event.SubType,    // poke / lucky_king / honor
		"target_id":  event.TargetId,   // 仅部分子类型有值
		"honor_type": event.Honor_type, // 仅荣誉事件有值
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleFriendRequest(event Pichubot.FriendRequest) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "friend_request",
		"scene":      "private",
		"user_id":    event.UserId,
		"self_id":    event.SelfId,
		"comment":    event.Comment,
		"flag":       event.Flag,
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleGroupRequest(event Pichubot.GroupRequest) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "group_request",
		"scene":      "group",
		"group_id":   event.GroupId,
		"user_id":    event.UserId,
		"self_id":    event.SelfId,
		"sub_type":   event.SubType, // add / invite
		"comment":    event.Comment,
		"flag":       event.Flag,
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleMetaLifecycle(event Pichubot.MetaLifecycle) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "meta_lifecycle",
		"self_id":    event.SelfId,
		"sub_type":   event.SubType, // enable / disable / connect
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}

func (a *OnebotAdapter) handleMetaHeartbeat(event Pichubot.MetaHeartbeat) {
	payload := map[string]any{
		"platform":   "onebot",
		"event_type": "meta_heartbeat",
		"self_id":    event.SelfId,
		"interval":   event.Interval,
		"status":     event.Status,
		"raw_event":  event,
	}
	a.ctx.Events.Publish("adapter.message", payload)
}
