package types

import "fmt"

const (
	minHotelName     = 3
	minHotelLocation = 4
	minHotelRating   = 0
)

type Hotel struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Location string   `bson:"location" json:"location"`
	Rooms    []string `bson:"rooms" json:"rooms"`
	Rating   int      `bson:"rating" json:"rating"`
}

type CreateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

func (p CreateHotelParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(p.Name) < minHotelName {
		errors["name"] = fmt.Sprintf("hotel name should be at least %d characters", minHotelName)
	}
	if len(p.Location) < minHotelLocation {
		errors["location"] = fmt.Sprintf("hotel location should be at least %d characters", minHotelLocation)
	}
	if p.Rating < minHotelRating {
		errors["rating"] = fmt.Sprintf("hotel rating should be greater than %d", minHotelRating)
	}
	return errors
}
func NewHotelFromParams(params CreateHotelParams) (*Hotel, error) {
	return &Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
		Rooms:    []string{},
	}, nil
}

type UpdateHotel struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func ValidateHotelUpdate(updateMap UpdateHotel) (*map[string]any, error) {
	validUpdate := map[string]any{}
	if updateMap.Name != "" {
		validUpdate["name"] = updateMap.Name
	}
	if updateMap.Location != "" {
		validUpdate["location"] = updateMap.Location
	}
	if updateMap.Name == "" && updateMap.Location == "" {
		return nil, fmt.Errorf("no valid update parameters for hotel")
	}
	return &validUpdate, nil
}
