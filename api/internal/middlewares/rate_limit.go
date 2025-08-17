package middlewares

import (
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func Limit(rps float64) gin.HandlerFunc {
	limiter := tollbooth.NewLimiter(rps, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Minute,
	})
	return tollbooth_gin.LimitHandler(limiter)
}
