package utils

type TgConfig struct {
	Enable bool   `toml:"enable"`
	Token  string `toml:"token"`
	ChatID string `toml:"chat_id"`
}

type PdConfig struct {
	Enable      bool      `toml:"enable"`
	RoutingKey  string    `toml:"routing_key"`
	EventAction string    `toml:"event_action"`
	Payload     PdpConfig `toml:"pdp"`
}

type PdpConfig struct {
	Severity string `toml:"severity"`
	Source   string `toml:"source"`
}

type SlConfig struct {
	Enable     bool   `toml:"enable"`
	WebhookURL string `toml:"webhookURL"`
}
