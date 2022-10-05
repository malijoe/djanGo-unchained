package django

import "github.com/gin-gonic/gin"

type view interface {
	asView(router *gin.RouterGroup)
}

func Path(eng *gin.Engine, path string, view view) {

}
