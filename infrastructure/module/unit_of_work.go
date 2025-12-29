package module

import (
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"gorm.io/gorm"
)

var _ repository.UnitOfWork = (*UnitOfWork)(nil)

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Begin() (repository.UnitOfWorkProcessor, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &UnitOfWorkProcessor{tx: tx}, nil
}

type UnitOfWorkProcessor struct {
	tx *gorm.DB
}

func (u *UnitOfWorkProcessor) ModuleWriter() repository.ModuleWriter {
	return NewModuleWriterRepository(u.tx)
}

func (u *UnitOfWorkProcessor) Rollback() error {
	return u.tx.Rollback().Error
}

func (u *UnitOfWorkProcessor) Commit() error {
	return u.tx.Commit().Error
}
