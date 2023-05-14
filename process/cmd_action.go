package process

import (
	"context"
	"log"
	"regexp"

	"github.com/tencent-connect/botgo/dto"
)

// 发送消息
func SendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	if _, err := processor.Api.PostMessage(ctx, channelID, toCreate); err != nil {
		regexp, _ := regexp.Compile("^code:[0-9]+")
		code := regexp.FindString(err.Error())[5:]
		switch code {
		case "202":
			return
		default:
			log.Println(err)
		}
	}
}
