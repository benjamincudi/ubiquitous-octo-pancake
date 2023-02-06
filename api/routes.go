package api

import (
	"atillm/api/handlers"
	"atillm/datastore"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

func AttachApiHandlers(r gin.IRoutes) {
	ds := datastore.NewMemoryDS()
	routeMap := gin.H{
		"/api/withdraw": handlers.QuickWithdrawController{ds},
		"/api/reset":    handlers.ResetController{ds},
	}
	mapRoutesToHandlers(r, routeMap)
}

type getHandler interface{ HandleGet(c *gin.Context) }
type postHandler interface{ HandlePost(c *gin.Context) }
type putHandler interface{ HandlePut(c *gin.Context) }
type deleteHandler interface{ HandleDelete(c *gin.Context) }

func mapRoutesToHandlers(r gin.IRoutes, routeMap gin.H) {
	for path, controller := range routeMap {
		matchedAnyMethod := false
		if getter, ok := controller.(getHandler); ok {
			matchedAnyMethod = true
			r.GET(path, getter.HandleGet)
		}
		if poster, ok := controller.(postHandler); ok {
			matchedAnyMethod = true
			r.POST(path, poster.HandlePost)
		}
		if putter, ok := controller.(putHandler); ok {
			matchedAnyMethod = true
			r.PUT(path, putter.HandlePut)
		}
		if deleter, ok := controller.(deleteHandler); ok {
			matchedAnyMethod = true
			r.DELETE(path, deleter.HandleDelete)
		}
		if !matchedAnyMethod {
			panic(fmt.Sprintf("handler %s for %s did not match any mapped method", reflect.TypeOf(controller), path))
		}
	}
}
