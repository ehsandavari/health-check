package notification

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/nikoksr/notify/service/discord"
)

func (r sNotification) AddDiscord() {
	if r.config.Discord == nil {
		return
	}
	d := discord.New()
	if err := d.AuthenticateWithBotToken(r.config.Discord.BotToken); err != nil {
		r.logger.WithError(err).Fatal(contextplus.Background(), "error in Authenticate discord")
	}
	d.AddReceivers(r.config.Discord.ChannelIds...)
	r.notify.UseServices(d)
}
