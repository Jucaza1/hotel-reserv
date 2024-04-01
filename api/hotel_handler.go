package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/db"
)

type HotelHandler struct {
	hotelStore db.HotelStore
}

func NewHotelHandler(hs db.HotelStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(hotel)
}
