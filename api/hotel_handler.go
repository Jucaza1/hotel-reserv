package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
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
	if len(id) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(hotel)
}
func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	err := h.hotelStore.DeleteHotel(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Next()
}
func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var params types.CreateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return types.ErrInvalidParams(err)
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(http.StatusBadRequest).JSON(errors)
	}
	hotel := types.NewHotelFromParams(params)
	insertedHotel, err := h.hotelStore.InsertHotel(c.Context(), hotel)
	if err != nil {
		return err
	}
	return c.JSON(insertedHotel)
}
func (h *HotelHandler) HandlePatchHotel(c *fiber.Ctx) error {
	var (
		hotelID   = c.Params("id")
		updateMap types.UpdateHotel
	)
	if len(hotelID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	if err := c.BodyParser(&updateMap); err != nil {
		return types.ErrInvalidParams(err)
	}
	validUpdate, err := types.ValidateHotelUpdate(updateMap)
	if err != nil {
		return types.ErrInvalidParams(err)
	}
	h.hotelStore.UpdateHotel(c.Context(), hotelID, *validUpdate)
	return c.JSON(types.MsgUpdated{Updated: hotelID})
}
