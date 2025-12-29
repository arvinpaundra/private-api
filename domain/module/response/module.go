package response

import "github.com/arvinpaundra/private-api/domain/module/constant"

type Module struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	SubjectID   string              `json:"subject_id"`
	GradeID     string              `json:"grade_id"`
	Title       string              `json:"title"`
	Slug        string              `json:"slug"`
	Description *string             `json:"description"`
	Type        constant.ModuleType `json:"type"`
	IsPublished bool                `json:"is_published"`
}
