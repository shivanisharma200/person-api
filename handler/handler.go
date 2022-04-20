package handler

import (
	"person-api/model"
	"person-api/service"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type personHandler struct {
	PersonService service.Person
}

func New(personService service.Person) *personHandler {
	return &personHandler{PersonService: personService}
}

func (p *personHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	person, err := p.PersonService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p *personHandler) Get(ctx *gofr.Context) (interface{}, error) {

	persons, err := p.PersonService.Get(ctx)

	if err != nil {
		return nil, err
	}


	return persons, nil
}

func (p *personHandler) Create(ctx *gofr.Context) (interface{}, error) {
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

	return personCreated, nil
}

func (p *personHandler) Update(ctx *gofr.Context) (interface{}, error) {
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

	return person, nil
}

func (p *personHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := ctx.PathParam("id")
	err := p.PersonService.Delete(ctx, idString)

	if err != nil {
		return nil, err
	}

	return "", nil
}
