package microcore

import (
	"github.com/fasthttp/router"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func InitRouter(r *router.Router) {
	if r != nil {
		r.ServeFilesCustom("/docs/{filepath:*}", &fasthttp.FS{
			Root: "./generations/swagger/",
		})
		r.GET("/api/{filepath:*}", fasthttpadaptor.NewFastHTTPHandlerFunc(httpSwagger.Handler(func(config *httpSwagger.Config) {
			config.URL = "/docs/swagger.json"
			config.DeepLinking = false
			config.DocExpansion = "list"
			config.DomID = "#swagger-ui"
			config.PersistAuthorization = true
		})))
	}
}
