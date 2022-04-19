package person

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"person-api/internal/model"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
)

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
	q := "SELECT id, name, age, address FROM person WHERE id=?"

	defer db.Close()

	testCases := []struct {
		desc          string
		id            int
		mockQuery     interface{}
		expectedError error
	}{
		{
			desc: "success test case",
			id:   1,
			mockQuery: mock.ExpectQuery(q).WithArgs(1).WillReturnRows(mock.NewRows([]string{"id", "name", "age", "address"}).
				AddRow(1, "Abc", 34, "Bangalore")),
			expectedError: nil,
		},

		{
			desc:          "failure test case",
			id:            1,
			mockQuery:     mock.ExpectQuery(q).WithArgs(1).WillReturnError(sql.ErrNoRows),
			expectedError: errors.EntityNotFound{Entity: "Person", ID: "1"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run("desc", func(t *testing.T) {
			_, err := store.GetByID(ctx, testCase.id)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error: %v, got: %v", testCase.expectedError, err)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	db, mock, store, ctx := NewMock()

	defer db.Close()

	q := "SELECT id, name, age, address FROM person"

	testCases := []struct {
		desc          string
		mockQuery     interface{}
		expectedError error
	}{
		// Success

		{
			desc: "success test case",

			mockQuery: mock.
				ExpectQuery(q).
				WillReturnRows(mock.NewRows([]string{"id", "name", "age", "address"}).
					AddRow(1, "Abc", "34", "Bangalore").
					AddRow(2, "Xyz", "29", "Pune")),
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			mockQuery: mock.
				ExpectQuery(q).
				WillReturnError(sql.ErrNoRows),
			expectedError: errors.EntityNotFound{Entity: "Person"},
		},
		// Failure
		{
			desc: "failure test case",
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
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := store.Get(ctx)

			if err != nil && err.Error() != testCase.expectedError.Error() {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
		})
	}
}

func Test_Create(t *testing.T) {
	db, mock, store, ctx := NewMock()

	const q = "INSERT INTO person(name, age, address) VALUES($, $, $)"

	defer db.Close()

	testCases := []struct {
		desc          string
		input         *model.Person
		output        *model.Person
		mockQuery     []interface{}
		expectedError error
	}{
		// Success

		{
			desc:  "success test case",
			input: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			output: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			mockQuery: []interface{}{mock.ExpectExec(q).
				WithArgs("Abc", 34, "Bangalore").
				WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnRows(mock.NewRows([]string{"id", "name", "address", "age"}).
						AddRow(1, "Abc", 34, "Bangalore")).WillReturnError(nil),
			},
			expectedError: nil,
		},
		// Failure
		{
			desc:  "failure test case",
			input: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			output: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			mockQuery: []interface{}{mock.ExpectExec(q).
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
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := store.Create(ctx, testCase.input)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
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
		output        *model.Person
		mockQuery     []interface{}
		expectedError error
	}{
		// Success
		{
			desc: "success test case",
			id:   1,
			input: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			output: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			mockQuery: []interface{}{
				mock.ExpectExec("UPDATE person SET name=?,address=? WHERE id=?").
					WithArgs("Abc", 34, "Bangalore").
					WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.
					ExpectQuery(q).
					WithArgs(1).
					WillReturnRows(mock.NewRows([]string{"id", "name", "address", "age"}).
						AddRow(1, "Abc", 34, "Bangalore")),
			},
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			id:   1,
			input: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			output: &model.Person{ID: "1", Name: "Abc", Age: 34, Address: "Bangalore" },
			mockQuery: []interface{}{
				mock.ExpectExec("UPDATE person SET name=?,address=? WHERE id=?").
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
	for _, testCase := range testCases {
		testCase := testCase

		t.Run("", func(t *testing.T) {
			_, err := store.Update(ctx, testCase.id, testCase.input)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
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
			mockQuery: mock.ExpectExec( "DELETE FROM person WHERE id=$").
				WithArgs(1).
				WillReturnResult(sqlmock.NewResult(1, 1)),
			expectedError: nil,
		},
		// Failure
		{
			desc: "failure test case",
			id:   1,
			mockQuery: mock.ExpectExec("DELETE FROM person WHERE id=$").
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

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected error:%v, got:%v", testCase.expectedError, err)
			}
		})
	}
}