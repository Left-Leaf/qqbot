package example

import (
	"context"
	"log"

	"github.com/tencent-connect/botgo/dto"
)

var Help help

// 有参指令
type help struct{}

func (c help) Handle(ctx context.Context, data *dto.WSATMessageData) error {
	log.Println("执行help命令")
	return nil
}
