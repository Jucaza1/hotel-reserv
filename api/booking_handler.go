package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
)

type BookingHandler struct {
	bookStore db.BookingStore
	roomStore db.RoomStore
}

func NewBookingHandler(bs db.BookingStore, rs db.RoomStore) *BookingHandler {
	return &BookingHandler{
		bookStore: bs,
		roomStore: rs,
	}
}

func (h *BookingHandler) HandleGetBookingsByRoom(c *fiber.Ctx) error {
	roomID := c.Params("id")
	bookings, err := h.bookStore.GetBookings(c.Context(), roomID)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}
func (h *BookingHandler) HandleGetBookingsByHotel(c *fiber.Ctx) error {
	roomID := c.Params("idh")
	bookings, err := h.bookStore.GetBookingsByHotel(c.Context(), roomID)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandlePostBooking(c *fiber.Ctx) error {
	roomID := c.Params("id")
	hotelID := c.Params("idh")
	room, err := h.roomStore.GetRoom(c.Context(), roomID)
	if err != nil {
		return err
	}
	if hotelID != room.HotelID {
		return fmt.Errorf("invalid request")
	}
	var params types.CreateBookingParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}
	userID := c.Get("userID")
	booking, err := types.NewBookingFromParams(params, userID, hotelID, roomID)
	if err != nil {
		return err
	}
	InsertedBooking, err := h.bookStore.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(InsertedBooking)
}
