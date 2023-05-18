package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"qqbot/process/command"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

type Process struct {
	Api openapi.OpenAPI
	// CmdList []command.Command
	CmdMap map[string]command.Command //使用map集合拥有更快的查找速度
}

// 定义一个消息处理器
var processor Process

// 初始化消息处理器
func InitProcessor(api openapi.OpenAPI) {
	processor = Process{
		Api:    api,
		CmdMap: make(map[string]command.Command),
	}
}

// 注册指令
func RegisterCmd(c command.Command) {
	// processor.CmdList = append(processor.CmdList, c)
	processor.CmdMap[c.GetID()] = c
}

// 获取消息处理器，目前好像没用，以后可能会去掉
func GetProcessor() Process {
	return processor
}

// 处理消息
func ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	//解析指令
	cmd := message.ParseCommand(input)
	// //遍历指令列表(已废弃)
	// for _, c := range processor.CmdList {
	// 	if c.Is(cmd.Cmd) {
	// 		if err := c.Handle(ctx, data); err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// }
	c := processor.CmdMap[cmd.Cmd]
	err := c.Handle(ctx, data)
	if err != nil {
		toCreate := BuildRMessage(err.Error(), data.ID)
		SendReply(ctx, data.ChannelID, toCreate)
	}
	return nil
}

// 打印消息
func PrintMessage(data *dto.Message) error {
	ctx := context.Background()
	content := data.Content
	if data.Mentions != nil {
		userList := data.Mentions
		input := message.ETLInput(content)
		head := ""
		for _, user := range userList {
			head += ("@" + user.Username)
		}
		content = head + " " + input
	}
	if data.Attachments != nil {
		attachments := data.Attachments
		end := "(附件: "
		for _, a := range attachments {
			end += ("{" + a.URL + "}")
		}
		end += ")"
		content += end
	}
	if data.DirectMessage {
		fmt.Printf("[message] %s [私信消息] %s(%s) -> %s\n", data.Timestamp, data.Author.Username, data.Author.ID, content)
	} else {
		guild, err := processor.Api.Guild(ctx, data.GuildID)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("[message] %s [%s] %s(%s) -> %s\n", data.Timestamp, guild.Name, data.Author.Username, data.Author.ID, content)
	}
	return nil
}

func MemberChange(eventType dto.EventType, data *dto.WSGuildMemberData) error {
	date := time.Now().Format("2006-01-02T15:04:05+08:00")
	ctx := context.Background()
	guild, err := processor.Api.Guild(ctx, data.GuildID)
	if err != nil {
		log.Println(err)
	}
	username := data.User.Username
	userID := data.User.ID
	output := fmt.Sprintf("[change] %s [%s] System -> %s(%s)", date, guild.Name, username, userID)
	if eventType == "GUILD_MEMBER_REMOVE" {
		fmt.Printf("%s离开频道\n", output)
	} else if eventType == "GUILD_MEMBER_ADD" {
		fmt.Printf("%s加入频道\n", output)
	} else if eventType == "GUILD_MEMBER_UPDATE" {
		fmt.Printf("%s频道属性发生变化\n", output)
	}
	return nil
}

func GuildChange(eventType dto.EventType, data *dto.WSGuildData) error {
	date := time.Now().Format("2006-01-02T15:04:05+08:00")
	output := fmt.Sprintf("[change] %s [%s] System -> ", date, data.Name)
	if eventType == "GUILD_CREATE" {
		fmt.Printf("%sbot加入频道\n", output)
	} else if eventType == "GUILD_UPDATE" {
		fmt.Printf("%s频道信息变更\n", output)
	} else if eventType == "GUILD_DELETE" {
		fmt.Printf("%sbot离开频道\n", output)
	}
	return nil
}

func ChannelChange(eventType dto.EventType, data *dto.WSChannelData) error {
	date := time.Now().Format("2006-01-02T15:04:05+08:00")
	ctx := context.Background()
	guild, err := processor.Api.Guild(ctx, data.GuildID)
	if err != nil {
		log.Println(err)
	}
	output := fmt.Sprintf("[change] %s [%s] System -> 子频道%s", date, guild.Name, data.Name)
	if eventType == "CHANNEL_CREATE" {
		fmt.Printf("%s被创建\n", output)
	} else if eventType == "GUILD_UPDATE" {
		fmt.Printf("%s信息变更\n", output)
	} else if eventType == "GUILD_DELETE" {
		fmt.Printf("%s被删除\n", output)
	}
	return nil
}

// ProcessInlineSearch is a function to process inline search
func ProcessInlineSearch(interaction *dto.WSInteractionData) error {
	if interaction.Data.Type != dto.InteractionDataTypeChatSearch {
		return fmt.Errorf("interaction data type not chat search")
	}
	search := &dto.SearchInputResolved{}
	if err := json.Unmarshal(interaction.Data.Resolved, search); err != nil {
		log.Println(err)
		return err
	}
	if search.Keyword != "test" {
		return fmt.Errorf("resolved search key not allowed")
	}
	searchRsp := &dto.SearchRsp{
		Layouts: []dto.SearchLayout{
			{
				LayoutType: 0,
				ActionType: 0,
				Title:      "内联搜索",
				Records: []dto.SearchRecord{
					{
						Cover: "https://pub.idqqimg.com/pc/misc/files/20211208/311cfc87ce394c62b7c9f0508658cf25.png",
						Title: "内联搜索标题",
						Tips:  "内联搜索 tips",
						URL:   "https://www.qq.com",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(searchRsp)
	if err := processor.Api.PutInteraction(context.Background(), interaction.ID, string(body)); err != nil {
		log.Println("api call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}
