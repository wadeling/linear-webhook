package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wadeling/linear-webhook/api/linear"
	"github.com/wadeling/linear-webhook/pkg/config"
	"github.com/wadeling/linear-webhook/pkg/notify"
	_ "github.com/wadeling/linear-webhook/pkg/notify/all"
	"github.com/wadeling/linear-webhook/pkg/staff"
	"github.com/wadeling/linear-webhook/pkg/webhook"
	"net/http"
)

func init() {
	flag.StringVar(&config.CfgFullPath, "config", "", "config-file-path")
}

func ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info().Msg("recv ping")
		c.JSON(200, gin.H{
			"message": "pong",
			"code":    0,
		})
	}
}

func webhookHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get header
		delivery := c.GetHeader("Linear-Delivery")
		event := c.GetHeader("Linear-Event")
		userAgent := c.GetHeader("User-Agent")

		log.Info().Str("delivery", delivery).Str("event", event).Str("userAgent", userAgent).Msg("recv webhook")

		body := webhook.PayLoad{}
		err := c.BindJSON(&body)
		if err != nil {
			log.Err(err).Msg("bad request")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		log.Debug().Interface("body", body).Msg("webhook detail")

		go func() {
			// send notify
			_ = notify.DeliverAll(notify.Config{
				Url:       config.Config.Feishu.WebhookUrl,
				LinConfig: linear.Config{Host: config.Config.Linear.ApiAddr, ApiKeys: config.Config.Linear.ApiKeys},
			}, body)
		}()

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"code":    0,
		})
	}
}

func main() {
	flag.Parse()

	// default log level-debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// read config file
	err := config.InitConfig(config.CfgFullPath)
	if err != nil {
		log.Err(err).Msg("failed to read config file")
		return
	}

	// show config param
	log.Debug().Interface("config", config.Config).Msg("config content")

	// load staff config
	if len(config.Config.Feishu.StaffFile) > 0 {
		if err := staff.Instance().LoadStaffInfoFromFile(config.Config.Feishu.StaffFile); err != nil {
			log.Err(err).Msg("load staff config err")
			return
		}
	}

	// set router
	r := gin.Default()
	r.GET("/ping", ping())
	r.POST("/webhook", webhookHandler())

	if len(config.Config.Server.Https.CertFile) > 0 && len(config.Config.Server.Https.KeyFile) > 0 {
		err := r.RunTLS(config.Config.Server.ListenAddr, config.Config.Server.Https.CertFile, config.Config.Server.Https.KeyFile)
		if err != nil {
			log.Err(err).Msg("https server exit")
		}
	} else {
		err := r.Run(config.Config.Server.ListenAddr)
		if err != nil {
			log.Err(err).Msg("http server exit")
		}
	}

}
