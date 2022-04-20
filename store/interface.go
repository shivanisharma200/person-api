package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"person-api/model"
)

type Person interface {
	GetByID(ctx *gofr.Context, id int) (*model.Person, error)
	Get(ctx *gofr.Context) ([]*model.Person, error)
	Create(ctx *gofr.Context, person *model.Person) (*model.Person, error)
	Update(ctx *gofr.Context, id int, patient *model.Person) (*model.Person, error)
	Delete(ctx *gofr.Context, id int) error
}
