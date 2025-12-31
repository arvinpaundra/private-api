package repository

type UnitOfWork interface {
	Begin() (UnitOfWorkProcessor, error)
}

type UnitOfWorkProcessor interface {
	SubmissionWriter() SubmissionWriter

	Commit() error
	Rollback() error
}
