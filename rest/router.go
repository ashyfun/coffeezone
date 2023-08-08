package main

import (
	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	*gin.RouterGroup
}

type Router struct {
	RouterGroup
	*gin.Engine
}

func NewRouter() *Router {
	router := gin.Default()
	return &Router{RouterGroup{&router.RouterGroup}, router}
}

func (r *RouterGroup) WrapGroup(path string) *RouterGroup {
	return &RouterGroup{r.Group(path)}
}

type HandlerFunc func(*Context)

func (r *RouterGroup) View(method string, path string, handlers ...HandlerFunc) {
	var predefHandlers []gin.HandlerFunc
	for _, h := range handlers {
		predefHandlers = append(predefHandlers, func(ginCtx *gin.Context) {
			ctx := &Context{ginCtx}
			h(ctx)
		})
	}

	r.Handle(method, path, predefHandlers...)
}
