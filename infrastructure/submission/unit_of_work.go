package submission

import (
	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"gorm.io/gorm"
)

var _ repository.UnitOfWork = (*UnitOfWork)(nil)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{
		db: db,
	}
}

func (u *UnitOfWork) Begin() (repository.UnitOfWorkProcessor, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &UnitOfWorkProcessor{
		tx: tx,
	}, nil
}

var _ repository.UnitOfWorkProcessor = (*UnitOfWorkProcessor)(nil)

type UnitOfWorkProcessor struct {
	tx *gorm.DB
}

func (p *UnitOfWorkProcessor) SubmissionWriter() repository.SubmissionWriter {
	return NewSubmissionWriterRepository(p.tx)
}

func (p *UnitOfWorkProcessor) Commit() error {
	return p.tx.Commit().Error
}

func (p *UnitOfWorkProcessor) Rollback() error {
	return p.tx.Rollback().Error
}
