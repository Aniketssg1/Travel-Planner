package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth        *handlers.AuthHandler
	Destination *handlers.DestinationHandler
	Review      *handlers.ReviewHandler
	Itinerary   *handlers.ItineraryHandler
}

func SetupRoutes(r *gin.Engine, h *Handlers) {
	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	destinations := api.Group("/destinations")
	{
		destinations.GET("", h.Destination.GetDestinations)
		destinations.GET("/:id", h.Destination.GetDestination)
		destinations.POST("", middleware.AuthRequired(), middleware.AdminRequired(), h.Destination.CreateDestination)
		destinations.PUT("/:id", middleware.AuthRequired(), middleware.AdminRequired(), h.Destination.UpdateDestination)
		destinations.DELETE("/:id", middleware.AuthRequired(), middleware.AdminRequired(), h.Destination.DeleteDestination)
		destinations.POST("/:id/reviews", middleware.AuthRequired(), h.Review.SubmitReview)
	}

	reviews := api.Group("/reviews", middleware.AuthRequired(), middleware.AdminRequired())
	{
		reviews.GET("", h.Review.GetAllReviews)
		reviews.DELETE("/:id", h.Review.DeleteReview)
	}

	itineraries := api.Group("/itineraries", middleware.AuthRequired())
	{
		itineraries.GET("", h.Itinerary.GetItineraries)
		itineraries.POST("", h.Itinerary.CreateItinerary)
		itineraries.PUT("/:id", h.Itinerary.UpdateItinerary)
		itineraries.DELETE("/:id", h.Itinerary.DeleteItinerary)
	}
}
