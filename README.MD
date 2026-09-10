
用于 GoFrame 框架

代码示例 `/internal/cmd/cmd.go`

```go
package cmd

import (
	"context"

	"github.com/geekor/wookoon-gf-cas/library/eventbus"
	"github.com/geekor/wookoon-gf-cas/library/wookooncas"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"wookoon-apps-go/internal/controller/storage"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func:  main,
	}
)

func main(ctx context.Context, parser *gcmd.Parser) (err error) {
	// ✨ 初始化 内存事件总线
	bus := eventbus.NewEventBus()

	// ------------------
	// ✨ 初始化 后台服务
	s := g.Server()
	cs := wookooncas.NewServer(bus)

	// 注册路由
	s.Group("/api", func(g1 *ghttp.RouterGroup) {
		g1.Middleware(ghttp.MiddlewareHandlerResponse)

		cs.RegisterRoutes(g1)

		g1.Group("/", func(authed *ghttp.RouterGroup) {
			authed.Middleware(wookooncas.CasAuthed(cs.Jwtcfg))

			authed.Bind(storage.NewV1())
		})
	})

	s.Run()
	return nil
}

```