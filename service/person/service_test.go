package person

import (
	"net/http"
	"person-api/model"
	"testing"

	"github.com/stretchr/testify/assert"

	"person-api/store"

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
		desc          string
		id            string
		out           *model.Person
		mockCall      *gomock.Call
		expectedError error
	}{
		{
			desc:          "success test case",
			id:            "1",
			out:           &person,
			mockCall:      mockPerson.EXPECT().GetByID(ctx, 1).Return(&person, nil),
			expectedError: nil,
		},
		{
			desc:          "failure test case",
			id:            "-1",
			out:           nil,
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPerson)

	for i, testCase := range testCases {
		out, err := p.GetByID(ctx, testCase.id)
		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_GetAll(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		out           []*model.Person
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		{
			desc:          "success test case",
			out:           []*model.Person{&person},
			mockCall:      mockPerson.EXPECT().Get(ctx).Return([]*model.Person{&person}, nil),
			expectedError: nil,
		},
		{
			desc:          "failure test case",
			out:           nil,
			mockCall:      mockPerson.EXPECT().Get(ctx).Return(nil, errors.Error("error")),
			expectedError: errors.Error("error"),
		},
	}

	p := New(mockPerson)
	for i, testCase := range testCases {
		out, err := p.Get(ctx)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_Create(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		input         model.Person
		out           *model.Person
		mockCall      *gomock.Call
		expectedError error
		status        int
	}{
		{
			desc:          "success test case",
			input:         person,
			out:           &person,
			mockCall:      mockPerson.EXPECT().Create(ctx, &person).Return(&person, nil),
			expectedError: nil,
		},
		{
			desc:  "failure test case",
			input: person,
			out:   nil,
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
		{
			desc: "Invalid name",
			input: model.Person{
				ID:      "1",
				Name:    "",
				Age:     34,
				Address: "Bangalore",
			},
			out: nil,
			expectedError: &errors.Response{
				StatusCode: http.StatusBadRequest,
				Code:       http.StatusText(http.StatusBadRequest),
				Reason:     "Invalid fields provided",
			},
		},
	}
	p := New(mockPerson)

	for i, testCase := range testCases {
		out, err := p.Create(ctx, &testCase.input)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_Update(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		id            string
		out           *model.Person
		mockCall      []*gomock.Call
		expectedError error
		status        int
	}{
		{
			desc: "success test case",
			id:   "1",
			out:  &person,
			mockCall: []*gomock.Call{mockPerson.EXPECT().Update(ctx, 1, &person).Return(&person, nil),
				mockPerson.EXPECT().GetByID(ctx, 1).Return(&person, nil),
			},
			expectedError: nil,
		},
		{
			desc:          "failure test case",
			id:            "-1",
			out:           nil,
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	p := New(mockPerson)

	for i, testCase := range testCases {
		out, err := p.Update(ctx, testCase.id, &person)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
		assert.Equal(t, testCase.out, out, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}

func Test_Delete(t *testing.T) {
	var ctx *gofr.Context

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockPerson := store.NewMockPerson(mockCtrl)
	testCases := []struct {
		desc          string
		id            string
		mockCall      []*gomock.Call
		expectedError error
		status        int
	}{
		{
			desc: "success test case",
			id:   "1",
			mockCall: []*gomock.Call{mockPerson.EXPECT().Delete(ctx, 1).Return(nil),
				mockPerson.EXPECT().GetByID(ctx, 1).Return(&person, nil),
			},
			expectedError: nil,
		},
		{
			desc:          "failure test case",
			id:            "-1",
			expectedError: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	p := New(mockPerson)

	for i, testCase := range testCases {
		err := p.Delete(ctx, testCase.id)

		assert.Equal(t, testCase.expectedError, err, "TEST[%d], failed.\n%s", i, testCase.desc)
	}
}
