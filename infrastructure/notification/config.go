package notification

type (
	SConfig struct {
		Discord *sDiscord
		Slack   *sSlack
	}
	sDiscord struct {
		BotToken   string   `validate:"required"`
		ChannelIds []string `validate:"required"`
	}
	sSlack struct {
		APIToken   string   `validate:"required"`
		ChannelIds []string `validate:"required"`
	}
)
