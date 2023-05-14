package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"qqbot/process/command"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

type Process struct {
	Api     openapi.OpenAPI
	CmdList []command.Command
}

// 定义一个消息处理器
var processor Process

// 初始化消息处理器
func InitProcessor(api openapi.OpenAPI) {
	processor = Process{
		Api: api,
	}
}

// 注册指令
func RegisterCmd(c command.Command) {
	processor.CmdList = append(processor.CmdList, c)
}

// 获取消息处理器
func GetProcessor() Process {
	return processor
}

// 处理消息
func ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	//解析指令
	cmd := message.ParseCommand(input)
	//遍历指令列表
	for _, c := range processor.CmdList {
		if c.Is(cmd.Cmd) {
			if err := c.Handle(ctx, data); err != nil {
				log.Println(err)
			}
		}
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
