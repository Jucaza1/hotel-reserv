package types

type RoomSize int

const (
	Small RoomSize = iota + 1
	Normal
	Large
	Extra
)

type Room struct {
	ID      string   `bson:"_id,omitempty" json:"id,omitempty"`
	Size    RoomSize `bson:"type" json:"type"`
	Price   float64  `bson:"price" json:"price"`
	HotelID string   `bson:"hotelID" json:"hotelID"`
}
type CreateRoomParams struct {
	Size  RoomSize `json:"type"`
	Price float64  `json:"price"`
}

func NewRoomFromParams(params CreateRoomParams) *Room {
	return &Room{
		Size:  params.Size,
		Price: params.Price,
	}
}
