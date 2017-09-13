package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port            int    `default:"1203"`
	VerifyToken     string `required:"true"`
	PageAccessToken string `required:"true"`
	WitToken        string `required:"true"`
	TextFile        string `required:"true"`

	// Time in minutes that a silenced bot will wait from last echo message.
	// Bot will exit Silence state if there is a message from user and last echo message is more than <SilenceWaitTime> minutes old
	SilenceWaitTime string `split_words:"true" default:"5"`

	// Time in minutes for dialog to be expired if there's no message from user
	// Dialog will be restart if the last message is more than <DialogExpiredTime> minutes old
	DialogExpiredTime string `split_words:"true" default:"360"`

	OrderApiUrl   string `split_words:"true" required:"true"`
	OrderApiKey   string `split_words:"true" required:"true"`
	SearchApiUrl  string `split_words:"true" required:"true"`
	ProductApiUrl string `split_words:"true" required:"true"`
	BotLogLevel   string `split_words:"true" default:"error"`
	AppLogLevel   string `split_words:"true" default:"info"`

	DbHost     string `split_words:"true" default:"localhost"`
	DbPort     string `split_words:"true" default:"3306"`
	DbUsername string `split_words:"true"`
	DbPassword string `split_words:"true"`
	DbName     string `split_words:"true"`

	ApiPort          int    `split_words:"true" default:"1204"`
	ApiSecurityToken string `split_words:"true"`

	ResizedImgUrl    string `split_words:"true" default:"https://cdn.fptshop.com.vn/Uploads/Originals/%s"`
	ColorImgUrl      string `split_words:"true"`
	InStockThreshold int    `split_words:"true" default:"1"`
}

var Env Config

func init() {
	if err := envconfig.Process("fptshop", &Env); err != nil {
		log.Fatal(" LỖI TO Failed to read configurations: ", err)
	}
}

func GetStaffId() string {
	// Disable for FPTShop since it uses Pancake v2
	return ""
}
