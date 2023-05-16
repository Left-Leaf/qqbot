package command

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

type Command interface {
	Handle(ctx context.Context, data *dto.WSATMessageData) error //执行指令
	Is(cmd string) bool                                          //指令匹配,请尽量简单化,遍历列表匹配会耗费较多时间
	GetID() string                                               //获取指令ID,目前好像没什么用，以后可能删除
}
