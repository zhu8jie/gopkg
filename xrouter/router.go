package router

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/zhu8jie/gopkg/xutils"
)

var webStartTime = time.Now().Format(xutils.YYYYMMDDHHIISS)

// AppPid 健康检查
func AppPid(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"pid":        os.Getpid(),
		"start_time": webStartTime,
		"code":       http.StatusOK,
	})
	return
}

func GetRouter(model string, writer io.Writer) *gin.Engine {
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(writer, os.Stdout)
	r := gin.Default()
	r.MaxMultipartMemory = 10 << 20 // 10M
	if model == gin.DebugMode {
		ginpprof.Wrap(r)
	}
	gin.SetMode(model)
	r.GET("/pid", AppPid)
	return r
}
