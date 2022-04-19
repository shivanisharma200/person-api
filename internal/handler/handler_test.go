package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"person-api/internal/model"
	"reflect"
	"testing"

	"person-api/internal/service"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/golang/mock/gomock"
)

var person = model.Person{
	ID:      "1",
	Name:    "Abc",
	Age:     34,
	Address: "Bangalore",
}

func Test_GetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		id            int
		mockCall      *gomock.Call
		expectedError error
	}{
		// Success
		{
			id:            1,
			mockCall:      mockPerson.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&person, nil),
			expectedError: nil,
		},
		// Failure
		{
			id:            -1,
			mockCall:      mockPerson.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPerson)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:9000/persons/{id}", nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		_, err := p.GetByID(ctx)

		if !reflect.DeepEqual(testCase.expectedError, err) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		mockCall      *gomock.Call
		expectedError error
	}{
		// Success
		{
			mockCall:      mockPerson.EXPECT().Get(gomock.Any()).Return([]*model.Person{&person}, nil),
			expectedError: nil,
		},
		// Failure
		{
			mockCall:      mockPerson.EXPECT().Get(gomock.Any()).Return(nil, errors.EntityNotFound{Entity: "Person"}),
			expectedError: errors.EntityNotFound{Entity: "Person"},
		},
	}

	p := New(mockPerson)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9000/persons", nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		_, err := p.Get(ctx)

		if !reflect.DeepEqual(testCase.expectedError, err) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		body          []byte
		input         model.Person
		mockCall      *gomock.Call
		expectedError error
	}{
		// Success
		{
			body: []byte(`{
				"name": "Abc",
				"age": "34",
				"address": "Bangalore"
				}`),
			input:         person,
			mockCall:      mockPerson.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&person, nil),
			expectedError: nil,
		},
		// Failure
		{
			body: []byte(`{
				"name": "Abc",
				"age": "young"
				"address": "Bangalore"
				}`),
			input: person,
			mockCall: mockPerson.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, &errors.Response{
				StatusCode: http.StatusBadRequest,
				Code:       http.StatusText(http.StatusBadRequest),
				Reason:     "Invalid fields provided",
			}),
			expectedError: &errors.Response{
				StatusCode: http.StatusBadRequest,
				Code:       http.StatusText(http.StatusBadRequest),
				Reason:     "Invalid fields provided",
			},
		},
	}

	p := New(mockPerson)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9000/persons", bytes.NewReader(testCase.body))
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		_, err := p.Create(ctx)

		if !reflect.DeepEqual(testCase.expectedError, err) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Update(t *testing.T) {
	var pat = model.Person{
		Name:    "Abc",
		Address: "Pune",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		body          []byte
		mockCall      *gomock.Call
		expectedError error
	}{
		// Success
		{
			body: []byte(`{
				"name": "Abc",
				"address": "Pune"
				}`),
			mockCall:      mockPerson.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&pat, nil),
			expectedError: nil,
		},
		// Failure
		{
			body: []byte(`{
				"name": "Abc",
				"address": "Pune"
				}`),
			mockCall: mockPerson.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.EntityNotFound{Entity: "Person", ID: "id"}),
			expectedError: errors.EntityNotFound{Entity: "Person", ID: "id"},
		},
		// Failure
		{
			body: []byte(`{
				"name": 1,
				"description": "person description"
				}`),
			expectedError: errors.InvalidParam{},
		},
	}

	p := New(mockPerson)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9000/persons/{id}", bytes.NewReader(testCase.body))
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		_, err := p.Update(ctx)

		if !reflect.DeepEqual(testCase.expectedError, err) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		mockCall      *gomock.Call
		expectedError error
	}{
		// Success
		{
			mockCall:      mockPerson.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
			expectedError: nil,
		},
		// Failure
		{
			mockCall: mockPerson.EXPECT().
				Delete(gomock.Any(), gomock.Any()).Return(errors.Error("unexpected error occuered in deleting row")),
			expectedError: errors.Error("unexpected error occuered in deleting row"),
		},
	}

	p := New(mockPerson)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9000/persons/{id}", nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		_, err := p.Delete(ctx)

		if !reflect.DeepEqual(testCase.expectedError, err) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}
