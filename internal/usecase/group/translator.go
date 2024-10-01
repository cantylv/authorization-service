package group

import (
	"github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
)

func newBidFromExistingGroup(g *entity.Group) *dto.Bid {
	return &dto.Bid{
		ID:        g.ID,
		GroupName: g.Name,
		UserId:    g.OwnerID,
		Status:    "approved",
	}
}
