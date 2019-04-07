package http


func stack() {}

func Recovery() Handler{
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {

			}

			ctx.WriteString(500, "InternalServerError")
		}()
		ctx.Next()
	}
}
