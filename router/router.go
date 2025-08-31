package router

import (
	"do_aug25/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	apiV1Router := r.Group("/api/v1")

	catRouter := apiV1Router.Group("/cat")
	{
		catRouter.GET("/", handlers.GetCats)
		catRouter.GET("/:id", handlers.GetCat)
		catRouter.POST("/", handlers.CreateCat)
		catRouter.PUT("/:id", handlers.UpdateCatsSalary)
		catRouter.PUT("/assignMission/:cat_id/:mission_id", handlers.AssignMissionToCat)
		catRouter.DELETE("/:id", handlers.DeleteCat)
	}

	missionRouter := apiV1Router.Group("/mission")
	{
		missionRouter.GET("/", handlers.GetMissions)
		missionRouter.GET("/:id", handlers.GetMission)
		missionRouter.POST("/", handlers.CreateMission)
		missionRouter.POST("/:mission_id/createTarget", handlers.CreateTarget)
		missionRouter.DELETE("/deleteTarget/:mission_id/:target_id", handlers.DeleteTarget)
		missionRouter.DELETE("/:id", handlers.DeleteMission)
		missionRouter.PUT("/:id/markCompleted", handlers.MarkMissionCompleted)
	}

	targetRouter := apiV1Router.Group("/target")
	{
		targetRouter.POST("/:target_id/createNote", handlers.CreateTargetNote)
	}

	// targetNoteRouter := apiV1Router.Group("/targetnote")
	{
		// catRouter.GET("/", handlers.GetCats)
		// targetNoteRouter.POST("/", handlers.CreateTargetNote)
		// catRouter.GET("/:id", handlers.GetCat)
	}
}
