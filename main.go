package main

import (
	"context"
	"log"
	"os"
	"qqbot/database"
	"qqbot/mylog"
	"qqbot/process"
	"qqbot/process/command/example"
	"time"

	"github.com/spf13/viper"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

// 指令队列
var MessageChan chan *dto.WSATMessageData

func main() {
	log.Println("启动程序")
	ctx := context.Background()

	InitConfig()
	//初始化数据库
	database.InitDB()
	defer database.CloseDB()

	// 初始化新的文件 logger，并使用相对路径来作为日志存放位置，设置最小日志界别为 DebugLevel
	logger, err := mylog.New("./logs", mylog.DebugLevel)
	if err != nil {
		log.Fatalln("error log new", err)
	}
	log.Print("日志初始化成功")
	botgo.SetLogger(logger)

	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	botToken.AppID = viper.GetUint64("appid")
	botToken.AccessToken = viper.GetString("token")

	// 初始化 openapi，正式环境
	// api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 沙箱环境
	api := botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)

	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	//初始化消息处理器
	process.InitProcessor(api)
	log.Println("消息处理器初始化成功")

	//监听消息
	MessageChan = make(chan *dto.WSATMessageData, 10)
	go func() {
		for {
			message := <-MessageChan
			go process.ProcessMessage(message)
		}
	}()

	//注册消息
	process.RegisterCmd("hi", example.Hello)
	process.RegisterCmd("help", example.Help)
	process.RegisterCmd("error", example.ErrorCMD)
	log.Println("指令注册完成")

	// websocket.RegisterResumeSignal(syscall.SIGUSR1)
	// 根据不同的回调，生成 intents
	intent := websocket.RegisterHandlers(
		// at 机器人事件
		ATMessageEventHandler(),
		// 私信，目前只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		DirectMessageHandler(),
		// 频道消息，只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		CreateMessageHandler(),
		// 频道事件
		GuildEventHandler(),
		// 成员事件
		MemberEventHandler(),
		// 子频道事件
		ChannelEventHandler(),
		// 互动事件
		InteractionHandler(),
		// 发帖事件
		ThreadEventHandler(),
	)

	log.Println("bot启动")
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}

}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler() event.ATMessageEventHandler {
	return func(_ *dto.WSPayload, data *dto.WSATMessageData) error {
		MessageChan <- data
		return nil
	}
}

// CreateMessageHandler 处理消息事件
func CreateMessageHandler() event.MessageEventHandler {
	return func(_ *dto.WSPayload, data *dto.WSMessageData) error {
		return process.MessageChange((*dto.Message)(data))
	}
}

// DirectMessageHandler 处理私信事件
func DirectMessageHandler() event.DirectMessageEventHandler {
	return func(_ *dto.WSPayload, data *dto.WSDirectMessageData) error {
		return process.MessageChange((*dto.Message)(data))
	}
}

// GuildEventHandler 处理频道事件
func GuildEventHandler() event.GuildEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildData) error {
		return process.GuildChange(event.Type, data)
	}
}

// ChannelEventHandler 处理子频道事件(待办)
func ChannelEventHandler() event.ChannelEventHandler {
	return func(event *dto.WSPayload, data *dto.WSChannelData) error {
		return process.ChannelChange(event.Type, data)
	}
}

// MemberEventHandler 处理成员变更事件
func MemberEventHandler() event.GuildMemberEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		return process.MemberChange(event.Type, data)
	}
}

// InteractionHandler 处理内联交互事件(待办)
func InteractionHandler() event.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		return process.ProcessInlineSearch(data)
	}
}

// ThreadEventHandler 处理发帖事件
func ThreadEventHandler() event.ThreadEventHandler {
	return func(event *dto.WSPayload, data *dto.WSThreadData) error {
		return nil
	}
}

// 初始化配置系统
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
