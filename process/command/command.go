package command

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

/*
	指令规范(不一定要遵循)
	1.一个指令有且只有一个对应的指令ID(降低查询难度)
	2.指令必须实现Handle方法
	3.指令如果有参数,必须检查参数是否正确,若不正确必须返回错误信息
	4.对于多参数指令,错误信息必须指明错误参数所在位置
*/

type Command interface {
	Handle(ctx context.Context, data *dto.WSATMessageData) error //执行指令
}
