package notification

import "github.com/nikoksr/notify/service/slack"

func (r sNotification) AddSlack() {
	if r.config.Slack == nil {
		return
	}
	s := slack.New(r.config.Slack.APIToken)
	s.AddReceivers(r.config.Slack.ChannelIds...)
	r.notify.UseServices(s)
}
