package entity

type Question struct {
	ID               string
	Content          string
	Slug             string
	NextQuestionSlug *string
	Choices          []*Choice
}

func (q *Question) GetChoiceByID(choiceID string) (*Choice, bool) {
	for _, choice := range q.Choices {
		if choice.ID == choiceID {
			return choice, true
		}
	}
	return nil, false
}

type Choice struct {
	ID              string
	Content         string
	IsCorrectAnswer bool
}
