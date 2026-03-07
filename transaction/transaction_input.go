package transaction

import "crowdfunding-api/user"

type GetCampaignTransactionsInput struct {
	ID int `uri:"id" binding:"required"`
	User user.User
}

