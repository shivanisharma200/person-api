package person

import (
	"net/http"
	"strconv"

	"person-api/internal/model"
	"person-api/internal/store"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person struct {
	PersonStoreHandler store.Person
}

func New(str store.Person) *Person {
	return &Person{PersonStoreHandler: str}
}

func (p *Person) GetByID(ctx *gofr.Context, idString string) (*model.Person, error) {
	id, _ := strconv.Atoi(idString)
	if !IsIDValid(id) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return p.PersonStoreHandler.GetByID(ctx, id)
}

func (p *Person) Get(ctx *gofr.Context) ([]*model.Person, error) {
	return p.PersonStoreHandler.Get(ctx)
}

func (p *Person) Create(ctx *gofr.Context, person *model.Person) (*model.Person, error) {
	err := person.Validate()

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusBadRequest,
			Code:       http.StatusText(http.StatusBadRequest),
			Reason:     "Invalid fields provided",
		}
	}

	return p.PersonStoreHandler.Create(ctx, person)
}

func (p *Person) Update(ctx *gofr.Context, idString string, person *model.Person) (*model.Person, error) {
	id, _ := strconv.Atoi(idString)
	if !IsIDValid(id) {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	_, err := p.GetByID(ctx, idString)

	if err != nil {
		return nil, errors.EntityNotFound{Entity: "Person", ID: idString}
	}

	return p.PersonStoreHandler.Update(ctx, id, person)
}

func (p *Person) Delete(ctx *gofr.Context, idString string) error {
	id, _ := strconv.Atoi(idString)
	if !IsIDValid(id) {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	_, err := p.GetByID(ctx, idString)

	if err != nil {
		return errors.EntityNotFound{Entity: "Person", ID: idString}
	}

	return p.PersonStoreHandler.Delete(ctx, id)
}
