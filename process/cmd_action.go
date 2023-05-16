package process

import (
	"context"
	"regexp"

	"github.com/tencent-connect/botgo/dto"
)

// 发送消息
func SendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) error {
	if _, err := processor.Api.PostMessage(ctx, channelID, toCreate); err != nil {
		regexp, _ := regexp.Compile("^code:[0-9]+")
		code := regexp.FindString(err.Error())[5:]
		switch code {
		case "202":
			return nil
		default:
			return err
		}
	}
	return nil
}
