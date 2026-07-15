package adapter_onebot

// OnebotConfig 适配器配置，只关心网络连接参数
type OnebotConfig struct {
	Host     string `mapstructure:"host"`      // WebSocket 地址，如 127.0.0.1:6700
	Path     string `mapstructure:"path"`      // WebSocket 路径，如 /
	Token    string `mapstructure:"token"`     // 鉴权 Token (如果有的话)
	MsgAwait bool   `mapstructure:"msg_await"` // 是否开启消息随机延迟
}
