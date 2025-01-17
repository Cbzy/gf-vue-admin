package boot

import (
	"github.com/flipped-aurora/gf-vue-admin/app/router/middleware"
	"github.com/flipped-aurora/gf-vue-admin/app/router/system"
	"github.com/gogf/gf/frame/g"
)

var Routers = new(_router)

type _router struct{}

// Initialize 路由初始化
// Author [SliverHorn](https://github.com/SliverHorn)
func (r *_router) Initialize() {
	public := g.Server().Group("")
	{
		system.NewCaptchaGroup(public).Public()
		system.NewUserRouter(public).Public()
	} // 无需鉴权中间件
	private := g.Server().Group("")
	private.Middleware(middleware.JwtAuth, middleware.Casbin)
	{
		system.NewApiRouter(private).Private()
		system.NewUserRouter(private).Private()
		system.NewMenuRouter(private).Private()
		system.NewDictionaryRouter(private).Private()
		system.NewAuthorityMenuRouter(private).Private()
		system.NewDictionaryDetailRouter(private).Private()
	} // 需要Jwt鉴权, casbin鉴权
}
