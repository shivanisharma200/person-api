package service

//go:generate mockgen -destination=interface_mock.go -package=store person-api/internal/service Person

import (
	"person-api/internal/model"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person interface {
	GetByID(ctx *gofr.Context, id string) (*model.Person, error)
	Create(ctx *gofr.Context, person *model.Person) (*model.Person, error)
	Get(ctx *gofr.Context) ([]*model.Person, error)
	Update(ctx *gofr.Context, id string, patient *model.Person) (*model.Person, error)
	Delete(ctx *gofr.Context, id string) error
}
