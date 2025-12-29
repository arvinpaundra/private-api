package response

import (
	"github.com/arvinpaundra/private-api/domain/module/constant"
)

type Module struct {
	ID             string              `json:"id"`
	UserID         string              `json:"user_id"`
	SubjectID      string              `json:"subject_id"`
	GradeID        string              `json:"grade_id"`
	Title          string              `json:"title"`
	Slug           string              `json:"slug"`
	Description    *string             `json:"description"`
	Type           constant.ModuleType `json:"type"`
	IsPublished    bool                `json:"is_published"`
	QuestionsCount int                 `json:"questions_count"`
	Subject        *Subject            `json:"subject,omitempty"`
	Grade          *Grade              `json:"grade,omitempty"`
}

type Subject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Grade struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ModuleDetail struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Slug        string              `json:"slug"`
	Description *string             `json:"description"`
	Type        constant.ModuleType `json:"type"`
	IsPublished bool                `json:"is_published"`
	Subject     *Subject            `json:"subject"`
	Grade       *Grade              `json:"grade"`
	Questions   []*Question         `json:"questions"`
}

type Question struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	Slug    string    `json:"slug"`
	Choices []*Choice `json:"choices"`
}

type Choice struct {
	ID              string `json:"id"`
	Content         string `json:"content"`
	IsCorrectAnswer bool   `json:"is_correct_answer"`
}
