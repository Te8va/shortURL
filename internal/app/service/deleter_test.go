package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Te8va/shortURL/internal/app/service"
	"github.com/Te8va/shortURL/internal/app/service/mocks"
)

func TestURLService_DeleteUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDeleter := mocks.NewMockURLDeleteServ(ctrl)

	testCases := []struct {
		name        string
		userID      int
		ids         []string
		mockSetup   func()
		expectedErr error
	}{
		{
			name:   "success",
			userID: 1,
			ids:    []string{"abc123", "xyz789"},
			mockSetup: func() {
				mockDeleter.
					EXPECT().
					DeleteUserURLs(gomock.Any(), []string{"abc123", "xyz789"}, 1).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "deleter error",
			userID: 2,
			ids:    []string{"badid"},
			mockSetup: func() {
				mockDeleter.
					EXPECT().
					DeleteUserURLs(gomock.Any(), []string{"badid"}, 2).
					Return(errors.New("delete failed"))
			},
			expectedErr: errors.New("delete failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			svc := service.NewURLService(nil, nil, nil, mockDeleter)
			err := svc.DeleteUserURLs(context.Background(), tc.ids, tc.userID)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
