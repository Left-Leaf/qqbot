package example

import (
	"context"
	"qqbot/process"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
)

var Hello hello

// 无参指令
type hello struct{}

func (c hello) Handle(ctx context.Context, data *dto.WSATMessageData) error {
	toCreate := process.BuildRMessage("默认回复"+message.Emoji(307), data.ID)
	return process.SendReply(ctx, data.ChannelID, toCreate)
}
