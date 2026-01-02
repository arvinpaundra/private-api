package response

type DashboardStatistics struct {
	UserFullname              string `json:"user_fullname"`
	TotalModules              int    `json:"total_modules"`
	TotalSubjects             int    `json:"total_subjects"`
	TotalGrades               int    `json:"total_grades"`
	TotalSubmittedSubmissions int    `json:"total_submitted_submissions"`
}
