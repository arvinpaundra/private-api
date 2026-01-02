package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/domain/dashboard/response"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type GetDashboardStatistics struct {
	authStorage   interfaces.AuthenticatedUser
	moduleACL     repository.ModuleACL
	subjectACL    repository.SubjectACL
	gradeACL      repository.GradeACL
	submissionACL repository.SubmissionACL
	userACL       repository.UserACL
}

func NewGetDashboardStatistics(
	authStorage interfaces.AuthenticatedUser,
	moduleACL repository.ModuleACL,
	subjectACL repository.SubjectACL,
	gradeACL repository.GradeACL,
	submissionACL repository.SubmissionACL,
	userACL repository.UserACL,
) *GetDashboardStatistics {
	return &GetDashboardStatistics{
		authStorage:   authStorage,
		moduleACL:     moduleACL,
		subjectACL:    subjectACL,
		gradeACL:      gradeACL,
		submissionACL: submissionACL,
		userACL:       userACL,
	}
}

func (s *GetDashboardStatistics) Execute(ctx context.Context) (*response.DashboardStatistics, error) {
	userID := s.authStorage.GetUserId()

	// Get user details
	user, err := s.userACL.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Aggregate counts from all domains via ACLs in parallel
	type result struct {
		modules     int
		subjects    int
		grades      int
		submissions int
		err         error
	}

	resultChan := make(chan result, 1)

	go func() {
		var res result

		// Count modules
		modules, err := s.moduleACL.CountModulesByUserID(ctx, userID)
		if err != nil {
			res.err = err
			resultChan <- res
			return
		}
		res.modules = modules

		// Count subjects
		subjects, err := s.subjectACL.CountSubjectsByUserID(ctx, userID)
		if err != nil {
			res.err = err
			resultChan <- res
			return
		}
		res.subjects = subjects

		// Count grades
		grades, err := s.gradeACL.CountGradesByUserID(ctx, userID)
		if err != nil {
			res.err = err
			resultChan <- res
			return
		}
		res.grades = grades

		// Count submitted submissions (not scoped to user, global count)
		submissions, err := s.submissionACL.CountSubmittedSubmissions(ctx)
		if err != nil {
			res.err = err
			resultChan <- res
			return
		}
		res.submissions = submissions

		resultChan <- res
	}()

	res := <-resultChan
	if res.err != nil {
		return nil, res.err
	}

	return &response.DashboardStatistics{
		UserFullname:              user.Fullname,
		TotalModules:              res.modules,
		TotalSubjects:             res.subjects,
		TotalGrades:               res.grades,
		TotalSubmittedSubmissions: res.submissions,
	}, nil
}
