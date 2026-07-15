package adapter_onebot

import (
	"fmt"

	Pichubot "github.com/yeswearebot/go-Pichubot"
	"github.com/yeswearebot/plugin_config"
	"github.com/yeswearebot/yes-core/core"
)

// OnebotAdapter 适配器实例
type OnebotAdapter struct {
	bot  *Pichubot.Bot
	ctx  *core.SystemContext
	conf OnebotConfig
}

func init() {
	core.Register(func() core.Plugin { return &OnebotAdapter{} })
}

func (a *OnebotAdapter) Name() string        { return "adapter-onebot" }
func (a *OnebotAdapter) DependsOn() []string { return []string{"config"} } // 必须依赖配置中心

func (a *OnebotAdapter) Init(ctx *core.SystemContext) error {
	a.ctx = ctx // 保存上下文，供事件桥接使用

	if err := plugin_config.Unmarshal(ctx, "onebot", &a.conf); err != nil {
		return fmt.Errorf("[Adapter-Onebot] 读取配置失败: %w", err)
	}

	a.bot = Pichubot.NewBot()
	a.bot.Config = Pichubot.Config{
		Loglvl:   Pichubot.LOGGER_LEVEL_INFO,
		Host:     a.conf.Host,
		Path:     a.conf.Path,
		MsgAwait: a.conf.MsgAwait,
	}

	// 注册事件桥接
	a.registerEventListeners()

	fmt.Println("[Adapter-Onebot] 初始化完成，准备连接 WebSocket...")
	return nil
}

func (a *OnebotAdapter) Start(ctx *core.SystemContext) error {
	// 在后台启动连接，防止阻塞内核主循环
	go func() {
		fmt.Printf("[Adapter-Onebot] 正在连接到 %s%s...\n", a.conf.Host, a.conf.Path)
		a.bot.Run()
	}()
	return nil
}

func (a *OnebotAdapter) Stop(ctx *core.SystemContext) error {
	// PichuBot 目前没有显式的关闭连接 API
	fmt.Println("[Adapter-Onebot] 已停止")
	return nil
}
