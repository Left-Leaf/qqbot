package example

import (
	"context"
	"errors"

	"github.com/tencent-connect/botgo/dto"
)

var ErrorCMD errorCMD

type errorCMD struct{}

func (c errorCMD) Handle(ctx context.Context, data *dto.WSATMessageData) error {
	return errors.New("指令错误示例")
}
