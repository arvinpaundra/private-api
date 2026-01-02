package response

import (
	"github.com/arvinpaundra/private-api/domain/submission/constant"
)

type Submission struct {
	ID          string                    `json:"id"`
	ModuleID    string                    `json:"module_id"`
	Code        string                    `json:"code"`
	StudentName string                    `json:"student_name"`
	Status      constant.SubmissionStatus `json:"status"`
	Module      *Module                   `json:"module,omitempty"`
	Answers     []*SubmissionAnswer       `json:"answers,omitempty"`
}

type Module struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type SubmissionAnswer struct {
	ID           string `json:"id"`
	SubmissionID string `json:"submission_id"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	IsCorrect    bool   `json:"is_correct"`
}

type SubmissionDetail struct {
	ID          string                    `json:"id"`
	Code        string                    `json:"code"`
	StudentName string                    `json:"student_name"`
	Status      constant.SubmissionStatus `json:"status"`
	Module      *Module                   `json:"module"`
	Answers     []*SubmissionAnswer       `json:"answers"`
}

type StartSubmissionResponse struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

type SubmitAnswerResponse struct {
	IsCorrect            bool    `json:"is_correct"`
	CorrectChoiceID      string  `json:"correct_choice_id"`
	CorrectChoiceContent string  `json:"correct_choice_content"`
	NextQuestionSlug     *string `json:"next_question_slug"`
}

type FinalizeSubmissionResponse struct {
	StudentName string `json:"student_name"`
	Score       int    `json:"score"`
	Total       int    `json:"total"`
	Status      string `json:"status"`
}

type ModuleSubmissionGroup struct {
	Module           *ModuleWithRelations `json:"module"`
	TotalSubmissions int                  `json:"total_submissions"`
	Submissions      []*SubmissionSummary `json:"submissions"`
}

type ModuleWithRelations struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Slug    string   `json:"slug"`
	Grade   *Grade   `json:"grade"`
	Subject *Subject `json:"subject"`
}

type Grade struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Subject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SubmissionSummary struct {
	StudentName    string `json:"student_name"`
	TotalCorrect   int    `json:"total_correct"`
	TotalQuestions int    `json:"total_questions"`
	SubmittedAt    string `json:"submitted_at"`
}
