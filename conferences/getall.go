package conferences

import (
	"context"
	"time"
)

//GetAllParams ...
type GetAllParams struct {
}

// GetAllResponse ...
type GetAllResponse struct {
	Conferences []Conference
}

// GetAll retrieves all conferences
// encore:api public
func GetAll(ctx context.Context, params *GetAllParams) (*GetAllResponse, error) {
	return &GetAllResponse{
		Conferences: []Conference{
			{
				ID:   1,
				Name: "Gophercon",
				Slug: "gc",
				Events: []Event{
					{
						ID:        1,
						Name:      "GopherCon 2020",
						Slug:      "gc-2020",
						StartDate: time.Date(2020, time.November, 9, 17, 00, 00, 0, time.UTC),
						EndDate:   time.Date(2020, time.November, 13, 23, 45, 00, 0, time.UTC),
						Location:  "Online",
						Slots: []EventSlot{
							{
								ID:          1,
								Name:        "Pre-Conference Workshop: Getting a Jumpstart in Go",
								Description: "Description goes here",
								Cost:        400,
								StartDate:   time.Date(2020, time.November, 9, 17, 00, 00, 0, time.UTC),
								EndDate:     time.Date(2020, time.November, 9, 21, 00, 00, 0, time.UTC),
								// DependsOn:         nil,
								PurchaseableFrom:  time.Date(2020, time.May, 9, 17, 00, 00, 0, time.UTC),
								PurchaseableUntil: time.Date(2020, time.November, 4, 17, 00, 00, 0, time.UTC),
								AvailableToPublic: true,
							},
						},
					},
				},
			},
		},
	}, nil

}
