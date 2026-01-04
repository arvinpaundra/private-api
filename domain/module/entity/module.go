package entity

import (
	"github.com/arvinpaundra/private-api/core/trait"
	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/module/constant"
)

type Module struct {
	trait.Createable
	trait.Updateable
	trait.Removeable

	ID          string
	UserID      string
	SubjectID   string
	GradeID     string
	Title       string
	Slug        string
	Description *string
	Type        constant.ModuleType
	IsPublished bool

	Questions []*Question
}

func NewModule(userID, subjectID, gradeID, title string, description *string) (*Module, error) {
	module := &Module{
		ID:          util.GenerateUUID(),
		UserID:      userID,
		SubjectID:   subjectID,
		GradeID:     gradeID,
		Title:       title,
		Description: description,
		Type:        constant.MultipleChoice,
		IsPublished: false,
	}

	err := module.GenSlug()
	if err != nil {
		return nil, err
	}

	module.MarkCreate()

	return module, nil
}

func (m *Module) GenSlug() error {
	slug, err := util.RandomAlphanumeric(12)
	if err != nil {
		return err
	}

	m.Slug = slug

	return nil
}

func (m *Module) Publish() {
	m.IsPublished = true
	m.MarkUpdate()
}

func (m *Module) Unpublish() {
	m.IsPublished = false
	m.MarkUpdate()
}

func (m *Module) AddQuestion(question *Question) {
	m.Questions = append(m.Questions, question)
	m.MarkUpdate()
}
