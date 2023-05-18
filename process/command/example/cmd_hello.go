package example

import (
	"context"
	"qqbot/process"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
)

// 无参指令
type Hello struct {
	id string
}

func NewHello() *Hello {
	return &Hello{
		id: "hi",
	}
}

func (c Hello) Handle(ctx context.Context, data *dto.WSATMessageData) error {
	toCreate := &dto.MessageToCreate{
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}
	return process.SendReply(ctx, data.ChannelID, toCreate)
}

func (c Hello) GetID() string {
	return c.id
}
