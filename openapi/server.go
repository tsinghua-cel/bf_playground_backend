package openapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/bf_playground_backend/config"
	"time"
)

type OpenAPI struct {
	conf *config.Config
}

func NewOpenAPI(conf *config.Config) *OpenAPI {
	return &OpenAPI{conf: conf}
}

func (s *OpenAPI) Run() error {
	url := fmt.Sprintf("%s:%d", s.conf.Server.Host, s.conf.Server.Port)
	return s.startHttp(url)
}

func (s *OpenAPI) startHttp(address string) error {
	router := gin.Default()
	router.Use(cors())
	router.Use(ginLogrus())
	handler := apiHandler{conf: s.conf}
	v1 := router.Group("/bfapi/v1")
	{
		v1.GET("/project-list", handler.GetProjectList)
		v1.GET("/top-strategies", handler.GetTopStrategies)

		// project page
		v1.GET("/project/:id", handler.GetProjectDetail)
		// download strategy
		v1.GET("/download/:id", handler.DownloadProject)
	}

	ch := make(chan error)
	go func() {
		err := router.Run(address)
		ch <- err
	}()
	time.Sleep(100 * time.Millisecond)
	select {
	case v := <-ch:
		return v
	default:
		return nil
	}
}

// gin use logrus
func ginLogrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"ip":     c.ClientIP(),
		}).Info("request")
		c.Next()
	}
}

// enable cors
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
