package process

import (
	"context"
	"log"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
)

func sendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	if _, err := processor.api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println(err)
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
