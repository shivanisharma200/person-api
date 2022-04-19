package handler

import (
	"person-api/internal/model"
	"person-api/internal/service"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
)

type API struct {
	PersonService service.Person 
}

func New(personService service.Person) *API {
	return &API{PersonService: personService}
}

func (p *API) GetByID(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	person, err := p.PersonService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return types.Response{
		Data: person,
	}, nil
}

func (p *API) Get(ctx *gofr.Context) (interface{}, error) {
	persons, err := p.PersonService.Get(ctx)

	if err != nil {
		return nil, err
	}

	return types.Response{
		Data: persons,
	}, nil
}

func (p *API) Create(ctx *gofr.Context) (interface{}, error) {
	var person model.Person

	err := ctx.Bind(&person)
	if err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{}
	}

	personCreated, err := p.PersonService.Create(ctx, &person)

	if err != nil {
		return nil, err
	}

	return types.Response{
		Data: personCreated,
	}, nil
}

func (p *API) Update(ctx *gofr.Context) (interface{}, error) {
	idString := ctx.PathParam("id")

	var patient *model.Person

	err := ctx.Bind(&patient)
	if err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{}
	}

	person, err := p.PersonService.Update(ctx, idString, patient)

	if err != nil {
		return nil, err
	}

	return types.Response{
		Data: person,
	}, nil
}

func (p *API) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := ctx.PathParam("id")
	err := p.PersonService.Delete(ctx, idString)

	if err != nil {
		return nil, err
	}

	return "", nil
}