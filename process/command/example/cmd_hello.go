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
	toCreate := process.BuildRMessage("默认回复"+message.Emoji(307), data.ID)
	return process.SendReply(ctx, data.ChannelID, toCreate)
}

func (c Hello) GetID() string {
	return c.id
}
