package example

import (
	"context"
	"log"

	"github.com/tencent-connect/botgo/dto"
)

// 单指令有参
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

func (c Help) Is(cmd string) bool {
	return c.id == cmd
}

func (c Help) GetID() string {
	return c.id
}
