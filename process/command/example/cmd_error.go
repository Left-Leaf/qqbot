package example

import (
	"context"
	"errors"

	"github.com/tencent-connect/botgo/dto"
)

type ErrorCMD struct {
	id string
}

func (c ErrorCMD) Handle(ctx context.Context, data *dto.WSATMessageData) error {
	return errors.New("指令错误示例")
}

func (c ErrorCMD) GetID() string {
	return c.id
}
