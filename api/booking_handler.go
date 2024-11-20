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
	if len(roomID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	user, ok := c.Context().UserValue("user").(types.User)
	if !ok {
		return types.ErrUnauthorized(fmt.Errorf("user not found"))
	}
	if user.IsAdmin {
		bookings, err := h.bookStore.GetBookingsByRoom(c.Context(), roomID)
		if err != nil {
			return err
		}
		return c.JSON(bookings)
	} else {
		bookings, err := h.bookStore.GetBookingsByUserAndRoom(c.Context(), user.ID, roomID)
		if err != nil {
			return err
		}
		return c.JSON(bookings)
	}
}

func (h *BookingHandler) HandlePostBooking(c *fiber.Ctx) error {
	roomID := c.Params("id")
	if len(roomID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	room, err := h.roomStore.GetRoom(c.Context(), roomID)
	if err != nil {
		return err
	}
	var params types.CreateBookingParams
	if err := c.BodyParser(&params); err != nil {
		return types.ErrInvalidParams(err)
	}
	if err := params.Validate(); err != nil {
		return types.ErrInvalidParams(err)
	}
	userID := c.Context().UserValue("user").(types.User).ID
	booking, _ := types.NewBookingFromParams(params, userID, room.HotelID, roomID)
	InsertedBooking, err := h.bookStore.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(InsertedBooking)
}

func (h *BookingHandler) HandleGetBookingsByHotel(c *fiber.Ctx) error {
	hotelID := c.Params("hid")
	if len(hotelID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	user, ok := c.Context().UserValue("user").(types.User)
	if !ok {
		return types.ErrUnauthorized(fmt.Errorf("user not found"))
	}
	if user.IsAdmin {
		bookings, err := h.bookStore.GetBookingsByHotel(c.Context(), hotelID)
		if err != nil {
			return err
		}
		return c.JSON(bookings)
	} else {
		bookings, err := h.bookStore.GetBookingsByUserAndHotel(c.Context(), user.ID, hotelID)
		if err != nil {
			return err
		}
		return c.JSON(bookings)
	}

}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(types.User)
	if !ok {
		return types.ErrUnauthorized(fmt.Errorf("user not found"))
	}
	if user.IsAdmin {
		bookings, err := h.bookStore.GetBookings(c.Context())
		if err != nil {
			return err
		}
		return c.JSON(bookings)

	} else {
		bookings, err := h.bookStore.GetBookingsByUser(c.Context(), user.ID)
		if err != nil {
			return err
		}
		return c.JSON(bookings)
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	if len(bookingID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	user, ok := c.Context().UserValue("user").(types.User)
	if !ok {
		return types.ErrUnauthorized(fmt.Errorf("user not found"))
	}
	if user.IsAdmin {
		if err := h.bookStore.CancelBooking(c.Context(), bookingID); err != nil {
			return err
		}
	} else {
		booking, err := h.bookStore.GetBookingByID(c.Context(), bookingID)
		if err != nil {
			return err
		}
		if booking.UserID == user.ID {
			if err := h.bookStore.CancelBooking(c.Context(), bookingID); err != nil {
				return err
			}
		} else {
			return types.ErrUnauthorized(fmt.Errorf("unauthorized cancel on different user"))
		}
	}
	return c.JSON(types.MsgCancelled{Cancelled: bookingID})
}

func (h *BookingHandler) HandleDeleteBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	if len(bookingID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	if err := h.bookStore.DeleteBooking(c.Context(), bookingID); err != nil {
		return err
	}
	return c.JSON(types.MsgDeleted{Deleted: bookingID})
}
