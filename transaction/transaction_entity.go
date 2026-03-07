package transaction

import (
	"time"
	"crowdfunding-api/user"
)

type Transaction struct {
	ID int
	CampaignID int
	UserID int
	Amount int
	Status string
	Code string
	User user.User
	CreatedAt time.Time
	UpdatedAt time.Time
}