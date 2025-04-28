package usecase

import (
	"context"
	"testing"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_UserUsecase_Register(t *testing.T) {
	type testCase struct {
		name             string
		wantError        bool
		request          models.RegisterUserRequest
		onRegister       func(mock *mocks.MockUserRepository)
		onGetUserByEmail func(mock *mocks.MockUserRepository)
	}

	var testTable []testCase
	testTable = append(testTable, testCase{
		name:      "success",
		wantError: false,
		request: models.RegisterUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		},
		onGetUserByEmail: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByEmail(gomock.Any(), "john@example.com").Return(models.User{}, nil)
		},
		onRegister: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().Register(gomock.Any(), models.RegisterUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			}).Return(int64(1), nil)
		},
	})

	testTable = append(testTable, testCase{
		name:      "email already registered",
		wantError: true,
		request: models.RegisterUserRequest{
			Name:     "Jane Doe",
			Email:    "jane@example.com",
			Password: "password123",
		},
		onGetUserByEmail: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByEmail(gomock.Any(), "jane@example.com").Return(models.User{UserID: 2}, nil)
		},
	})

	testTable = append(testTable, testCase{
		name:      "error checking user by email",
		wantError: true,
		request: models.RegisterUserRequest{
			Name:     "Error User",
			Email:    "error@example.com",
			Password: "password123",
		},
		onGetUserByEmail: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByEmail(gomock.Any(), "error@example.com").Return(models.User{}, serror.New("database error"))
		},
	})

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			userRepo := mocks.NewMockUserRepository(mockCtrl)

			if tc.onGetUserByEmail != nil {
				tc.onGetUserByEmail(userRepo)
			}

			if tc.onRegister != nil {
				tc.onRegister(userRepo)
			}

			usecase := &UserUsecase{userRepo: userRepo}

			err := usecase.Register(context.Background(), tc.request)

			if tc.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
