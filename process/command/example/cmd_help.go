package example

import (
	"context"
	"log"

	"github.com/tencent-connect/botgo/dto"
)

// 有参指令
type Help struct {
	id string
}

func NewHelp() *Help {
	return &Help{
		id: "help",
	}
}

func (c Help) Handle(ctx context.Context, data *dto.WSATMessageData) error {
	log.Println("执行help命令")
	return nil
}
func (c Help) GetID() string {
	return c.id
}
