package types

import (
	"encoding/json"
	"fmt"
)

type RoomSize int

const (
	Small RoomSize = iota + 1
	Normal
	Large
	Extra
)

type Room struct {
	ID      string   `bson:"_id,omitempty" json:"id,omitempty"`
	Size    RoomSize `bson:"size" json:"size"`
	Price   float64  `bson:"price" json:"price"`
	HotelID string   `bson:"hotelID" json:"hotelID"`
}
type CreateRoomParams struct {
	Size  RoomSize `json:"size"`
	Price float64  `json:"price"`
}

func NewRoomFromParams(params CreateRoomParams) *Room {
	return &Room{
		Size:  params.Size,
		Price: params.Price,
	}
}
func (rs RoomSize) String() string {
	switch rs {
	case Small:
		return "Small"
	case Normal:
		return "Normal"
	case Large:
		return "Large"
	case Extra:
		return "Extra"
	default:
		return "Unknown"
	}
}

// Implementing json.Marshaler
func (rs RoomSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.String())
}

// Implementing json.Unmarshaler
func (rs *RoomSize) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Small":
		*rs = Small
	case "Normal":
		*rs = Normal
	case "Large":
		*rs = Large
	case "Extra":
		*rs = Extra
	default:
		return fmt.Errorf("invalid RoomSize: %s", s)
	}
	return nil
}
