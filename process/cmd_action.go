package process

import (
	"context"
	"log"
	"regexp"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
)

func sendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
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

func hi_run(ctx context.Context, data *dto.WSATMessageData) {

	toCreate := &dto.MessageToCreate{
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}
	sendReply(ctx, data.ChannelID, toCreate)
}

func pin_run() {

}
