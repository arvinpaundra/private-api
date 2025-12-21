package entity

import (
	"github.com/arvinpaundra/private-api/core/trait"
	"github.com/arvinpaundra/private-api/core/util"
)

type Subject struct {
	trait.Createable
	trait.Updateable
	trait.Removeable

	ID          string
	UserID      string
	Name        string
	Description *string
}

func NewSubject(userID, name string, description *string) *Subject {
	subject := &Subject{
		ID:          util.GenerateUUID(),
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	subject.MarkCreate()

	return subject
}

func (s *Subject) Update(name string, description *string) {
	s.Name = name
	s.Description = description

	s.MarkUpdate()
}
