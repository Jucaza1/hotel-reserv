package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/db"
)

type RoomHandler struct {
	roomStore db.RoomStore
}

func NewRoomHandler(rs db.RoomStore) *RoomHandler {
	return &RoomHandler{
		roomStore: rs,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	rooms, err := h.roomStore.GetRooms(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
