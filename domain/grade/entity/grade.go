package entity

import (
	"github.com/arvinpaundra/private-api/core/trait"
	"github.com/arvinpaundra/private-api/core/util"
)

type Grade struct {
	trait.Createable
	trait.Updateable
	trait.Removeable

	ID          string
	UserID      string
	Name        string
	Description *string
}

func NewGrade(userID, name string, description *string) *Grade {
	grade := &Grade{
		ID:          util.GenerateUUID(),
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	grade.MarkCreate()

	return grade
}

func (g *Grade) Update(name string, description *string) {
	g.Name = name
	g.Description = description

	g.MarkUpdate()
}
