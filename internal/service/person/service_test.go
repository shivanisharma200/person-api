package person

import (
	"net/http"
	"person-api/internal/model"
	"reflect"
	"testing"

	"person-api/internal/store"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
)

var person = model.Person{
	ID:      "1",
	Name:    "Abc",
	Age:     34,
	Address: "Bangalore",
}

func Test_GetByID(t *testing.T) {
	var ctx *gofr.Context
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		id            string
		mockCall      *gomock.Call
		expectedError error
	}{
		{
			id:            "1",
			mockCall:      mockPerson.EXPECT().GetByID(ctx, "1").Return(&person, nil),
			expectedError: nil,
		},
		{
			id:            "-1",
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPerson)

	for _, testCase := range testCases {
		_, err := p.GetByID(ctx, testCase.id)
		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_GetAll(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		output        []model.Person
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			output:        []model.Person{person},
			mockCall:      mockPerson.EXPECT().Get(ctx).Return([]*model.Person{&person}, nil),
			expectedError: nil,
		},
		// Failure
		{
			mockCall:      mockPerson.EXPECT().Get(ctx).Return(nil, errors.Error("error")),
			expectedError: errors.Error("error"),
		},
	}

	p := New(mockPerson)
	for _, testCase := range testCases {
		_, err := p.Get(ctx)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Create(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		input         model.Person
		output        model.Person
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			input:         person,
			mockCall:      mockPerson.EXPECT().Create(ctx, &person).Return(&person, nil),
			expectedError: nil,
		},
		// Failure
		{
			input: person,
			mockCall: mockPerson.EXPECT().Create(ctx, &person).Return(nil, &errors.Response{
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
		// Invalid Id
		{
			input: model.Person{
				ID:      "-1",
				Name:    "Abc",
				Age:     34,
				Address: "Bangalore",
			},
			expectedError: &errors.Response{
				StatusCode: http.StatusBadRequest,
				Code:       http.StatusText(http.StatusBadRequest),
				Reason:     "Invalid fields provided",
			},
		},
		// Invalid Name
		{
			input: model.Person{
				ID:      "1",
				Name:    "",
				Age:     34,
				Address: "Bangalore",
			},
			expectedError: &errors.Response{
				StatusCode: http.StatusBadRequest,
				Code:       http.StatusText(http.StatusBadRequest),
				Reason:     "Invalid fields provided",
			},
		},
	}
	p := New(mockPerson)

	for _, testCase := range testCases {
		_, err := p.Create(ctx, &testCase.input)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Update(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		id            string
		output        model.Person
		mockCall      []*gomock.Call
		expectedError error
		status        int
	}{

		// Success
		{
			id: "1",
			mockCall: []*gomock.Call{mockPerson.EXPECT().Update(ctx, 1, &person).Return(&person, nil),
				mockPerson.EXPECT().GetByID(ctx, 1).Return(&person, nil),
			},
			expectedError: nil,
		},
		// Failure Invalid Id
		{
			id:            "-1",
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPerson)

	for _, testCase := range testCases {
		_, err := p.Update(ctx, testCase.id, &person)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}

func Test_Delete(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		id            string
		mockCall      []*gomock.Call
		expectedError error
		status        int
	}{
		// Success
		{
			id: "1",
			mockCall: []*gomock.Call{mockPerson.EXPECT().Delete(ctx, 1).Return(nil),
				mockPerson.EXPECT().GetByID(ctx, 1).Return(&person, nil),
			},
			expectedError: nil,
		},
		// Failure Invalid Id
		{
			id:            "-1",
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	p := New(mockPerson)

	for _, testCase := range testCases {
		err := p.Delete(ctx, testCase.id)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("Expected error: %v Got %v", testCase.expectedError, err)
		}
	}
}