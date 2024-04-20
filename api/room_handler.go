package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jucaza1/hotel-reserv/db"
	"github.com/jucaza1/hotel-reserv/types"
)

type RoomHandler struct {
	roomStore  db.RoomStore
	hotelStore db.HotelStore
}

func NewRoomHandler(rs db.RoomStore, hs db.HotelStore) *RoomHandler {
	return &RoomHandler{
		roomStore:  rs,
		hotelStore: hs,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	hid := c.Params("hid")
	if len(hid) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	rooms, err := h.roomStore.GetRooms(c.Context(), hid)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleGetRoomByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	rooms, err := h.roomStore.GetRoom(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleDeleteRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	err := h.roomStore.DeleteRoom(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(types.MsgDeleted{Deleted: id})
}

func (h *RoomHandler) HandleDeleteRoomsByHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	err := h.roomStore.DeleteRoomsByHotel(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(types.MsgDeleted{Deleted: id})
}

func (h *RoomHandler) HandlePostRoom(c *fiber.Ctx) error {
	hotelID := c.Params("hid")
	if len(hotelID) == 0 {
		return types.ErrInvalidID(fmt.Errorf("missing params in path"))
	}
	if _, err := h.hotelStore.GetHotelByID(c.Context(), hotelID); err != nil {
		return err
	}
	var params types.CreateRoomParams
	if err := c.BodyParser(&params); err != nil {
		return types.ErrInvalidParams(err)
	}
	room := types.NewRoomFromParams(params)
	room.HotelID = hotelID
	insertedRoom, err := h.roomStore.InsertRoom(c.Context(), room)
	if err != nil {
		return err
	}
	return c.JSON(insertedRoom)
}
