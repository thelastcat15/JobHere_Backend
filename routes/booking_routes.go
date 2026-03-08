package routes

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/handlers"
)

// RegisterBookingRoutes registers Booking-related endpoints
func RegisterBookingRoutes(r fiber.Router) {
	r.Get("/", handlers.GetBookings)
	r.Get("/:id", handlers.GetBooking)
	r.Post("/", handlers.CreateBooking)
	r.Delete("/:id", handlers.DeleteBooking)

	r.Post("/:id/checkin", handlers.ConfirmUserAtParkingForSuccessBooking)
	r.Post("/:id/complete", handlers.CompleteBookingPayment)
	r.Post("/:id/cancel", handlers.CancelBooking)

	r.Get("/history", handlers.GetUserBookingHistory)
}
