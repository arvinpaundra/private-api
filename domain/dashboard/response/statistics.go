package response

type DashboardStatistics struct {
	TotalModules           int `json:"total_modules"`
	TotalSubjects          int `json:"total_subjects"`
	TotalGrades            int `json:"total_grades"`
	TotalSubmittedSubmissions int `json:"total_submitted_submissions"`
}
