package router

import (
	"voting/rest/handler"

	"github.com/gin-gonic/gin"
)

// Init returns a new router
func Init(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/vote/create", h.CreateVote)
	router.DELETE("/vote", h.DeleteVote)
	router.POST("/vote/update", h.UpdateVote)
	router.GET("/vote/:voteID", h.GetVote)
	router.GET("/vote", h.ListVote)
	router.GET("/vote/:voteID/result", h.GetResult)

	router.POST("/ballot/add", h.AddBallot)
	router.POST("/ballot/cast", h.CastBallot)
	router.POST("/ballot/count", h.CountBallot)
	router.POST("/ballot", h.GetBallot)

	return router
}
