package command

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

type Command interface {
	Handle(ctx context.Context, data *dto.WSATMessageData) error
	GetID() string
	Is(cmd string) bool
}
