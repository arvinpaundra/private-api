package constant

type SubmissionStatus string

func (ss SubmissionStatus) String() string {
	return string(ss)
}

const (
	InProgress SubmissionStatus = "inprogress"
	Submitted  SubmissionStatus = "submitted"
	Canceled   SubmissionStatus = "canceled"
)
