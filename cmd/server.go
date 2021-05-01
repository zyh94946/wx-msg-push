package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/zyh94946/wx-msg-push-tencent/api"
	"github.com/zyh94946/wx-msg-push-tencent/config"
	"github.com/zyh94946/wx-msg-push/conf"
	"log"
	"net/http"
	"os"
)

func Server() *cobra.Command {
	var (
		cfgPath string
	)
	cmdServer := &cobra.Command{
		Use:   "server",
		Short: "Start Run wx-msg-push server",
		Run: func(cmd *cobra.Command, args []string) {
			if len(cfgPath) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
			conf.Init(cfgPath)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			err := run()
			if err != nil {
				log.Println("router.Run error", err)
			}
			return nil
		},
	}
	cmdServer.Flags().StringVarP(&cfgPath, "conf", "c", "", "server config [toml]")
	return cmdServer
}

func run() error {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/:SECRET", GoSendMsg)
	router.POST("/:SECRET", GoSendMsg)

	PrintVersion()
	cfg := conf.GetConfig()
	log.Println("server is running at", cfg.Server.Addr)

	s := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.MaxHTTPTime.Duration,
		WriteTimeout: cfg.Server.MaxHTTPTime.Duration,
	}
	_ = s.ListenAndServe()

	return nil
}

type msgData struct {
	SECRET  string `uri:"SECRET"`
	Title   string `form:"title"`
	Content string `form:"content"`
	Type    string `form:"type"`
	ToUser  string `form:"touser"`
	ToParty string `form:"toparty"`
	ToTag   string `form:"totag"`
}

func GoSendMsg(c *gin.Context) {

	var msg msgData

	if err := c.ShouldBindUri(&msg); err != nil {
		log.Println("ShouldBindUri fail", err)
		c.JSON(http.StatusBadRequest, gin.H{"errorCode": -2, "errorMessage": "fail"})
		return
	}

	cfg := conf.GetConfig()
	if _, isExist := cfg.WeChatConf[msg.SECRET]; !isExist {
		c.JSON(http.StatusBadRequest, gin.H{"errorCode": -1, "errorMessage": "request check fail"})
		return
	}

	if err := c.ShouldBind(&msg); err != nil {
		log.Println(err)
	}

	if msg.ToUser == "" && msg.ToParty == "" && msg.ToTag == "" {
		msg.ToUser = "@all"
	}

	weChatConf := cfg.WeChatConf[msg.SECRET]

	at := &api.AccessToken{
		CorpId:     weChatConf.CorpId,
		CorpSecret: weChatConf.CorpSecret,
	}

	var appMsg api.AppMsg
	opts := &api.MsgOpts{
		ToUser:                 msg.ToUser,
		ToParty:                msg.ToParty,
		ToTag:                  msg.ToTag,
		Title:                  msg.Title,
		Content:                msg.Content,
		AgentId:                weChatConf.AgentId,
		MediaId:                weChatConf.MediaId,
		EnableDuplicateCheck:   weChatConf.EnableDuplicateCheck,
		DuplicateCheckInterval: weChatConf.DuplicateCheckInterval,
	}

	switch msg.Type {
	case config.MsgTypeMpNews:
		appMsg = api.NewMpNews(opts)

	case config.MsgTypeText:
		appMsg = api.NewText(opts)

	default:
		appMsg = api.NewMpNews(opts)
	}

	if err := api.Send(appMsg, at); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadGateway, gin.H{"errorCode": -3, "errorMessage": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errorCode": 0, "errorMessage": ""})
	return
}
