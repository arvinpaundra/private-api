package repository

type UnitOfWork interface {
	Begin() (UnitOfWorkProcessor, error)
}

type UnitOfWorkProcessor interface {
	ModuleWriter() ModuleWriter

	Commit() error
	Rollback() error
}
