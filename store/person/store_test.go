package person

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"person-api/model"
	"testing"

	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
)

var person *model.Person = &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore"}

func NewMock() (db *sql.DB, mock sqlmock.Sqlmock, store Person, ctx *gofr.Context) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}

	store = New()
	ctx = gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return
}

func Test_GetByID(t *testing.T) {
	db, mock, store, ctx := NewMock()
	q := "SELECT id, name, age, address FROM person WHERE id=$1"

	defer db.Close()

	testCases := []struct {
		desc          string
		id            int
		out           *model.Person
		mockQuery     interface{}
		expectedError error
	}{
		{
			desc: "success test case",
			id:   1,
			out:  person,
			mockQuery: mock.ExpectQuery(q).WithArgs(1).WillReturnRows(mock.NewRows([]string{"id", "name", "age", "address"}).
				AddRow(1, "Abc", 34, "Bangalore")),
			expectedError: nil,
		},

		{
			desc:          "failure test case",
			id:            1,
			out:           nil,
			mockQuery:     mock.ExpectQuery(q).WithArgs(1).WillReturnError(sql.ErrNoRows),
			expectedError: errors.EntityNotFound{Entity: "Person", ID: "1"},
		},
	}

	for i, testCase := range testCases {
		testCase := testCase

		t.Run("desc", func(t *testing.T) {
			out, err := store.GetByID(ctx, testCase.id)

			assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
			assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
		})
	}
}

func Test_Get(t *testing.T) {
	db, mock, store, ctx := NewMock()

	defer db.Close()

	q := "SELECT id, name, age, address FROM person"

	testCases := []struct {
		desc          string
		out           []*model.Person
		mockQuery     interface{}
		expectedError error
	}{
		{
			desc: "success test case",
			out: []*model.Person{
				&model.Person{"1", "Abc", 34, "Bangalore"},
				&model.Person{"2", "Xyz", 29, "Pune"},
			},
			mockQuery: mock.
				ExpectQuery(q).
				WillReturnRows(mock.NewRows([]string{"id", "name", "age", "address"}).
					AddRow(1, "Abc", 34, "Bangalore").
					AddRow(2, "Xyz", 29, "Pune")),
			expectedError: nil,
		},
		{
			desc: "failure test case",
			out:  nil,
			mockQuery: mock.
				ExpectQuery(q).
				WillReturnError(sql.ErrNoRows),
			expectedError: errors.EntityNotFound{Entity: "Person"},
		},
		{
			desc: "failure test case",
			out:  nil,
			mockQuery: mock.
				ExpectQuery(q).
				WillReturnError(&errors.Response{
					StatusCode: http.StatusInternalServerError,
					Code:       http.StatusText(http.StatusInternalServerError),
					Reason:     "cannot fetch rows",
				}),
			expectedError: &errors.Response{
				StatusCode: http.StatusInternalServerError,
				Code:       http.StatusText(http.StatusInternalServerError),
				Reason:     "cannot fetch rows",
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			out, err := store.Get(ctx)

			assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
			assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
		})
	}
}

func Test_Create(t *testing.T) {
	db, mock, store, ctx := NewMock()

	const q = "INSERT INTO person(name, age, address) VALUES($1, $2, $3) RETURNING id"
	defer db.Close()

	testCases := []struct {
		desc          string
		input         *model.Person
		out           *model.Person
		mockQuery     []interface{}
		expectedError error
	}{
		// {
		// 	desc:   "success test case",
		// 	input:  &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore"},
		// 	out: person,
		// 	mockQuery: []interface{}{mock.ExpectQuery(q).
		// 		WithArgs("Abc", 34, "Bangalore").
		// 		WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1)).WillReturnError(nil),
		// 		mock.
		// 			ExpectQuery(q).
		// 			WithArgs(1).
		// 			WillReturnRows(mock.NewRows([]string{"id", "name", "address", "age"}).
		// 				AddRow(1, "Abc", 34, "Bangalore")).WillReturnError(nil),
		// 	},
		// 	expectedError: nil,
		// },
		{
			desc:  "failure test case",
			input: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore"},
			out:   nil,
			mockQuery: []interface{}{mock.ExpectQuery(q).
				WithArgs("Abc", 34, "Bangalore").
				WillReturnError(errors.Error("Failed to create person")),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnError(&errors.Response{
						StatusCode: http.StatusInternalServerError,
						Code:       http.StatusText(http.StatusInternalServerError),
						Reason:     "cannot create new person",
					}),
			},
			expectedError: &errors.Response{
				StatusCode: http.StatusInternalServerError,
				Code:       http.StatusText(http.StatusInternalServerError),
				Reason:     "cannot create new person",
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			out, err := store.Create(ctx, testCase.input)

			assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
			assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
		})
	}
}

func Test_Update(t *testing.T) {
	db, mock, store, ctx := NewMock()

	q := "SELECT id, name, age, address FROM person"

	defer db.Close()
	testCases := []struct {
		desc          string
		id            int
		input         *model.Person
		out           *model.Person
		mockQuery     []interface{}
		expectedError error
	}{
		// {
		// 	desc:   "success test case",
		// 	id:     1,
		// 	input:  &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore"},
		// 	out: person,
		// 	mockQuery: []interface{}{
		// 		mock.ExpectExec("UPDATE person SET name=$1,address=$2 WHERE id=$3").
		// 			WithArgs("Abc", 34, "Bangalore").
		// 			WillReturnResult(sqlmock.NewResult(1, 1)),
		// 		mock.
		// 			ExpectQuery(q).
		// 			WithArgs(1).
		// 			WillReturnRows(mock.NewRows([]string{"id", "name", "address", "age"}).
		// 				AddRow(1, "Abc", 34, "Bangalore")),
		// 	},
		// 	expectedError: nil,
		// },
		{
			desc:  "failure test case",
			id:    1,
			input: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore"},
			out:   nil,
			mockQuery: []interface{}{
				mock.ExpectExec("UPDATE person SET name=$1,address=$2 WHERE id=$3").
					WithArgs("Abc", 34, "Bangalore").
					WillReturnError(&errors.Response{
						StatusCode: http.StatusInternalServerError,
						Code:       http.StatusText(http.StatusInternalServerError),
						Reason:     "cannot update rows",
					}),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnError(&errors.Response{
						StatusCode: http.StatusInternalServerError,
						Code:       http.StatusText(http.StatusInternalServerError),
						Reason:     "cannot update rows",
					}),
			},
			expectedError: &errors.Response{
				StatusCode: http.StatusInternalServerError,
				Code:       http.StatusText(http.StatusInternalServerError),
				Reason:     "cannot update rows",
			},
		},
	}
	for i, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			out, err := store.Update(ctx, testCase.id, testCase.input)

			assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
			assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
		})
	}
}

func Test_Delete(t *testing.T) {
	db, mock, store, ctx := NewMock()

	defer db.Close()

	testCases := []struct {
		desc          string
		id            int
		mockQuery     interface{}
		expectedError error
	}{
		// Success

		{
			desc: "success test case",
			id:   1,
			mockQuery: mock.ExpectExec("DELETE FROM person WHERE id=$1").
				WithArgs(1).
				WillReturnResult(sqlmock.NewResult(1, 1)),
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			id:   1,
			mockQuery: mock.ExpectExec("DELETE FROM person WHERE id=$1").
				WithArgs(1).
				WillReturnError(&errors.Response{
					StatusCode: http.StatusInternalServerError,
					Code:       http.StatusText(http.StatusInternalServerError),
					Reason:     "cannot delete row",
				}),
			expectedError: &errors.Response{
				StatusCode: http.StatusInternalServerError,
				Code:       http.StatusText(http.StatusInternalServerError),
				Reason:     "cannot delete row",
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			err := store.Delete(ctx, testCase.id)

			assert.Equal(t, err, testCase.expectedError)
		})
	}
}
