package repository

type UnitOfWork interface {
	Begin() (UnitOfWorkProcessor, error)
}

type UnitOfWorkProcessor interface {
	UserWriter() UserWriter

	Commit() error
	Rollback() error
}
