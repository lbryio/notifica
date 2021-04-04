package config

import (
	"github.com/johntdyer/slackrus"
	"github.com/lbryio/notifica/app/env"
	"github.com/sirupsen/logrus"
)

func InitLogging(conf *env.Config) {
	IsDebugMode = true // conf.IsDebug
	if IsDebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if conf.SlackHookURL != "" {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        conf.SlackHookURL,
			AcceptedLevels: slackrus.LevelThreshold(logrus.InfoLevel),
			Channel:        conf.SlackChannel,
			IconEmoji:      ":bar_chart:",
			Username:       "Rick",
		})
	}

}
