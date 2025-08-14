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

func TestURLService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mocks.NewMockURLGetterServ(ctrl)
	svc := service.NewURLService(nil, mockGetter, nil, nil)

	testCases := []struct {
		name           string
		id             string
		mockSetup      func()
		expectedURL    string
		expectedExists bool
		expectedDel    bool
	}{
		{
			name: "found and not deleted",
			id:   "http://localhost/abc123",
			mockSetup: func() {
				mockGetter.
					EXPECT().
					Get(gomock.Any(), "http://localhost/abc123").
					Return("https://example.com", true, false)
			},
			expectedURL:    "https://example.com",
			expectedExists: true,
			expectedDel:    false,
		},
		{
			name: "not found",
			id:   "http://localhost/404",
			mockSetup: func() {
				mockGetter.
					EXPECT().
					Get(gomock.Any(), "http://localhost/404").
					Return("", false, false)
			},
			expectedURL:    "",
			expectedExists: false,
			expectedDel:    false,
		},
		{
			name: "found but deleted",
			id:   "http://localhost/deleted",
			mockSetup: func() {
				mockGetter.
					EXPECT().
					Get(gomock.Any(), "http://localhost/deleted").
					Return("", true, true)
			},
			expectedURL:    "",
			expectedExists: true,
			expectedDel:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			url, exists, deleted := svc.Get(context.Background(), tc.id)
			assert.Equal(t, tc.expectedURL, url)
			assert.Equal(t, tc.expectedExists, exists)
			assert.Equal(t, tc.expectedDel, deleted)
		})
	}
}

func TestURLService_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mocks.NewMockURLGetterServ(ctrl)
	svc := service.NewURLService(nil, mockGetter, nil, nil)

	t.Run("success", func(t *testing.T) {
		expected := []map[string]string{
			{"short": "abc123", "original": "https://google.com"},
		}
		mockGetter.
			EXPECT().
			GetUserURLs(gomock.Any(), 1).
			Return(expected, nil)

		result, err := svc.GetUserURLs(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("error from getter", func(t *testing.T) {
		mockGetter.
			EXPECT().
			GetUserURLs(gomock.Any(), 2).
			Return(nil, errors.New("db error"))

		result, err := svc.GetUserURLs(context.Background(), 2)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
