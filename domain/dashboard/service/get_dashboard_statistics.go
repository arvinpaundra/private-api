package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/domain/dashboard/response"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"golang.org/x/sync/errgroup"
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

	// Aggregate counts from all domains using errgroup
	var (
		modulesCount     int
		subjectsCount    int
		gradesCount      int
		submissionsCount int
	)

	g, ctx := errgroup.WithContext(ctx)

	// Count modules in parallel
	g.Go(func() error {
		count, err := s.countModules(ctx, userID)
		if err != nil {
			return err
		}
		modulesCount = count
		return nil
	})

	// Count subjects in parallel
	g.Go(func() error {
		count, err := s.countSubjects(ctx, userID)
		if err != nil {
			return err
		}
		subjectsCount = count
		return nil
	})

	// Count grades in parallel
	g.Go(func() error {
		count, err := s.countGrades(ctx, userID)
		if err != nil {
			return err
		}
		gradesCount = count
		return nil
	})

	// Count submitted submissions in parallel
	g.Go(func() error {
		count, err := s.countSubmissions(ctx)
		if err != nil {
			return err
		}
		submissionsCount = count
		return nil
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &response.DashboardStatistics{
		UserFullname:              user.Fullname,
		TotalModules:              modulesCount,
		TotalSubjects:             subjectsCount,
		TotalGrades:               gradesCount,
		TotalSubmittedSubmissions: submissionsCount,
	}, nil
}

func (s *GetDashboardStatistics) countModules(ctx context.Context, userID string) (int, error) {
	return s.moduleACL.CountModulesByUserID(ctx, userID)
}

func (s *GetDashboardStatistics) countSubjects(ctx context.Context, userID string) (int, error) {
	return s.subjectACL.CountSubjectsByUserID(ctx, userID)
}

func (s *GetDashboardStatistics) countGrades(ctx context.Context, userID string) (int, error) {
	return s.gradeACL.CountGradesByUserID(ctx, userID)
}

func (s *GetDashboardStatistics) countSubmissions(ctx context.Context) (int, error) {
	return s.submissionACL.CountSubmittedSubmissions(ctx)
}
