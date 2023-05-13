package process

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

type Process struct {
	Api openapi.OpenAPI
}

var processor Process

func InitProcessor(api openapi.OpenAPI) {
	processor = Process{Api: api}
}

func GetProcessor() Process {
	return processor
}

func ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	//解析指令
	cmd := message.ParseCommand(input)

	switch cmd.Cmd {
	case "hi":
		hi_run(ctx, data)
	case "pin":
		pin_run()
	}

	return nil
}
