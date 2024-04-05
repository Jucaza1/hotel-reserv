package types

import (
	"fmt"
	"time"
)

type Booking struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string    `bson:"userID,omitempty" json:"userID,omitempty"`
	HotelID     string    `bson:"hotelID,omitempty" json:"hotelID,omitempty"`
	RoomID      string    `bson:"roomID,omitempty" json:"roomID,omitempty"`
	FromDate    time.Time `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	ToDate      time.Time `bson:"toDate,omitempty" json:"toDate,omitempty"`
	CreatedDate time.Time `bson:"createDate,omitempty" json:"CreatedDate,omitempty"`
	CancelledAt time.Time `bson:"cancelledAt" json:"cancelledAt"`
	Cancelled   bool      `bson:"cancelled" json:"cancelled"`
}
type CreateBookingParams struct {
	FromDate time.Time `json:"fromDate,omitempty"`
	ToDate   time.Time `json:"toDate,omitempty"`
}

func (p CreateBookingParams) Validate() error {
	if p.FromDate.After(p.ToDate) && time.Now().After(p.FromDate) {
		return fmt.Errorf("invalid date")
	}
	return nil
}
func NewBookingFromParams(params CreateBookingParams, userID string, hotelID string, roomID string) (*Booking, error) {
	return &Booking{
		UserID:      userID,
		RoomID:      roomID,
		HotelID:     hotelID,
		FromDate:    params.FromDate,
		ToDate:      params.ToDate,
		CreatedDate: time.Now(),
	}, nil
}
