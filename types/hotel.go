package types

type Hotel struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Location string   `bson:"location" json:"location"`
	Rooms    []string `bson:"rooms" json:"rooms"`
	Rating   int      `bson:"rating" json:"rating"`
}

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
