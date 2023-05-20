package process

import (
	"context"
	"regexp"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/log"
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
			log.Error(err)
			return err
		}
	}
	return nil
}

// 构建消息
func BuildMessage(content string, image string, dataID string) *dto.MessageToCreate {
	toCreate := &dto.MessageToCreate{
		Content: content,
		Image:   image,
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             dataID,
			IgnoreGetMessageError: true,
		},
	}
	return toCreate
}
