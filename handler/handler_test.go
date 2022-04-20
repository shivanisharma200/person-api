package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"person-api/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"

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

func httpReq(httpMethod, pathUrl string, body []byte) *gofr.Context {
	req := httptest.NewRequest(httpMethod, "http://localhost:9000/"+pathUrl, bytes.NewReader(body))
	ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
	return ctx
}
func Test_GetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		id            int
		out           interface{}
		mockCall      *gomock.Call
		expectedError error
	}{
		{
			desc:          "success test case",
			out:           &person,
			id:            1,
			mockCall:      mockPerson.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&person, nil),
			expectedError: nil,
		},
		{
			desc:          "failure test case",
			out:           nil,
			id:            -1,
			mockCall:      mockPerson.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPerson)

	for i, testCase := range testCases {
		ctx := httpReq(http.MethodGet, "persons/{id}", nil)
		out, err := p.GetByID(ctx)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_Get(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		out           interface{}
		mockCall      *gomock.Call
		expectedError error
	}{
		{
			desc:          "success test case",
			out:           []*model.Person{&person},
			mockCall:      mockPerson.EXPECT().Get(gomock.Any()).Return([]*model.Person{&person}, nil),
			expectedError: nil,
		},
		{
			desc:          "failure test case",
			out:           nil,
			mockCall:      mockPerson.EXPECT().Get(gomock.Any()).Return(nil, errors.EntityNotFound{Entity: "Person"}),
			expectedError: errors.EntityNotFound{Entity: "Person"},
		},
	}

	p := New(mockPerson)

	for i, testCase := range testCases {
		ctx := httpReq(http.MethodGet, "persons", nil)
		out, err := p.Get(ctx)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		body          []byte
		input         model.Person
		out           interface{}
		mockCall      *gomock.Call
		expectedError error
	}{
		{
			desc: "success test case",
			body: []byte(`{
				"name": "Abc",
				"age": 34,
				"address": "Bangalore"
				}`),
			input:         person,
			out:           &person,
			mockCall:      mockPerson.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&person, nil),
			expectedError: nil,
		},
		{
			desc: "failure test case",
			body: []byte(`{
				"name": "Abc",
				"age": "young"
				"address": "Bangalore"
				}`),
			input:         person,
			out:           nil,
			expectedError: errors.InvalidParam{},
		},
	}

	p := New(mockPerson)

	for i, testCase := range testCases {
		ctx := httpReq(http.MethodPost, "persons", testCase.body)
		out, err := p.Create(ctx)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		body          []byte
		out           interface{}
		mockCall      *gomock.Call
		expectedError error
	}{
		{
			desc: "success test case",
			body: []byte(`{
				"name": "Abc",
				"address": "Pune"
				}`),
			out:           &person,
			mockCall:      mockPerson.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&person, nil),
			expectedError: nil,
		},
		{
			desc: "failure test case",
			body: []byte(`{
				"name": "Abc",
				"address": "Pune"
				}`),
			out: nil,
			mockCall: mockPerson.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.EntityNotFound{Entity: "Person", ID: "id"}),
			expectedError: errors.EntityNotFound{Entity: "Person", ID: "id"},
		},
		{
			desc: "Invalid Name Provided",
			body: []byte(`{
				"name": 1,
				"description": "person description"
				}`),
			out:           nil,
			expectedError: errors.InvalidParam{},
		},
	}

	p := New(mockPerson)

	for i, testCase := range testCases {
		ctx := httpReq(http.MethodPut, "persons/{id}", testCase.body)
		out, err := p.Update(ctx)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := service.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		mockCall      *gomock.Call
		expectedError error
	}{
		{
			desc:          "success test case",
			mockCall:      mockPerson.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
			expectedError: nil,
		},
		{
			desc: "failure test case",
			mockCall: mockPerson.EXPECT().
				Delete(gomock.Any(), gomock.Any()).Return(errors.Error("unexpected error occuered in deleting row")),
			expectedError: errors.Error("unexpected error occuered in deleting row"),
		},
	}

	p := New(mockPerson)

	for i, testCase := range testCases {
		ctx := httpReq(http.MethodDelete, "persons/{id}", nil)
		_, err := p.Delete(ctx)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}
