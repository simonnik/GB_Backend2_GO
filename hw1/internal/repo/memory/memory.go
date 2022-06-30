package memory

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/simonnik/GB_Backend2_GO/internal/entity"
)

// InMemory struct works with memory
type InMemory struct {
	users      []entity.User
	groups     []entity.Group
	memberShip map[uuid.UUID][]uuid.UUID // [groupId][]userId
}

func (i InMemory) CreateUser(ctx context.Context, u entity.User) error {
	i.users = append(i.users, u)

	return nil
}

func (i InMemory) CreateGroup(ctx context.Context, g entity.Group) error {
	i.groups = append(i.groups, g)

	return nil
}

func (i InMemory) AddToGroup(ctx context.Context, uid, gid uuid.UUID) error {
	i.memberShip[gid] = append(i.memberShip[gid], uid)

	return nil
}

func (i InMemory) RemoveFromGroup(ctx context.Context, uid, gid uuid.UUID) error {
	if len(i.memberShip[gid]) < 1 {
		return fmt.Errorf("user [%s] is not a member of a group [%s]", uid, gid)
	}
	var memberShip []uuid.UUID
	for _, id := range i.memberShip[gid] {
		if id != uid {
			memberShip = append(memberShip, uid)
		}
	}
	i.memberShip[gid] = memberShip

	return nil
}

func (i InMemory) SearchUser(ctx context.Context, name string, gids []uuid.UUID) ([]entity.User, error) {
	var users []entity.User
	for _, u := range i.users {
		if u.Name == name {
			if len(gids) > 0 {
				for _, gid := range gids {
					if len(i.memberShip[gid]) > 0 {
						isMember := false
						for _, id := range i.memberShip[gid] {
							if id == u.ID {
								isMember = true
								break
							}
						}

						if !isMember {
							return nil, fmt.Errorf("user[%s] is not a member of a group[%s]", u.ID, gid)
						}
					}
				}
			}
			users = append(users, u)
		}
	}

	return users, nil
}

func (i InMemory) SearchGroup(ctx context.Context, name string, uids []uuid.UUID) ([]entity.Group, error) {
	var groups []entity.Group
	for _, g := range i.groups {
		if g.Name == name {
			if len(uids) > 0 {
				for _, uid := range uids {
					if len(i.memberShip[g.ID]) > 0 {
						isMember := false
						for _, id := range i.memberShip[g.ID] {
							if id == uid {
								isMember = true
								break
							}
						}

						if !isMember {
							return nil, fmt.Errorf("user[%s] is not a member of a group[%s]", uid, g.ID)
						}
					}
				}
			}
			groups = append(groups, g)
		}
	}

	return groups, nil
}

// NewMemory create struct
func NewMemory() *InMemory {
	return &InMemory{}
}
