package filters

import "github.com/astaxie/beego/context"

func BaseFilter(ctx * context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", ctx.Input.Domain())
	ctx.Output.Header("Access-Control-Allow-Methods", "*")
}

