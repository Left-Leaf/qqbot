package command

import (
	"context"

	"github.com/tencent-connect/botgo/dto"
)

/*
	指令规范(不一定要遵循)
	1.一个指令最多只能拥有一个指令ID
	2.指令必须实现Handle方法和GetID方法
	3.指令如果有参数,必须检查参数是否正确,若不正确必须返回错误信息
	4.对于多参数指令,错误信息必须指明错误参数所在位置
*/

type Command interface {
	Handle(ctx context.Context, data *dto.WSATMessageData) error //执行指令
	GetID() string                                               //获取指令ID,将id设置为私有更加安全(可能没有意义)
}
