package person

import (
	"database/sql"
	"net/http"
	"strconv"

	"person-api/model"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Person struct {
}

func New() Person {
	return Person{}
}

func updateFunc(person *model.Person) (query string, values []interface{}) {
	if person.Name != "" {
		query += "name=$1,"

		values = append(values, person.Name)
	}

	if person.Address != "" {
		query += "address=$2,"

		values = append(values, person.Address)
	}

	if len(query) > 0 {
		query = query[:len(query)-1]
	}

	return query, values
}

func (p Person) GetByID(ctx *gofr.Context, id int) (*model.Person, error) {
	const q = "SELECT id, name, age, address FROM person WHERE id=$1"

	person := model.Person{}

	err := ctx.DB().QueryRowContext(ctx, q, id).Scan(&person.ID, &person.Name, &person.Age, &person.Address)

	if err == sql.ErrNoRows {
		idString := strconv.Itoa(id)
		return nil, errors.EntityNotFound{Entity: "Person", ID: idString}
	}

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannnot fetch row",
		}
	}
	return &person, nil
}

func (p Person) Get(ctx *gofr.Context) ([]*model.Person, error) {
	q := "SELECT id, name, age, address FROM person"
	rows, err := ctx.DB().QueryContext(ctx, q)

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "Person"}
	}

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot fetch rows",
		}
	}

	var persons []*model.Person

	defer rows.Close()

	for rows.Next() {
		var person model.Person
		_ = rows.Scan(&person.ID, &person.Name, &person.Age, &person.Address)
		persons = append(persons, &person)
	}

	return persons, nil
}

func (p Person) Create(ctx *gofr.Context, person *model.Person) (*model.Person, error) {
	var lastInsertedID int
	err := ctx.DB().QueryRowContext(ctx, "INSERT INTO person(name, age, address) VALUES($1, $2, $3) RETURNING id", person.Name, person.Age, person.Address).Scan(
		&lastInsertedID)

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot create new person",
		}
	}
	return p.GetByID(ctx, lastInsertedID)
}

func (p Person) Update(ctx *gofr.Context, id int, person *model.Person) (*model.Person, error) {
	query := "UPDATE person SET "

	resQuery, values := updateFunc(person)
	query += resQuery
	query += " WHERE id=$3"

	values = append(values, id)

	_, err := ctx.DB().ExecContext(ctx, query, values...)

	if err != nil {
		return nil, &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot update rows",
		}
	}

	return p.GetByID(ctx, id)
}

func (p Person) Delete(ctx *gofr.Context, id int) (err error) {
	_, err = ctx.DB().ExecContext(ctx, "DELETE FROM person WHERE id=$1", id)

	if err != nil {
		return &errors.Response{
			StatusCode: http.StatusInternalServerError,
			Code:       http.StatusText(http.StatusInternalServerError),
			Reason:     "cannot delete row",
		}
	}

	return nil
}
