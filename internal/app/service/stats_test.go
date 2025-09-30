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

func TestURLService_GetStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStats := mocks.NewMockURLStatsServ(ctrl)

	svc := service.NewURLService(nil, nil, nil, nil, mockStats)

	tests := []struct {
		name        string
		mockSetup   func()
		wantURLs    int
		wantUsers   int
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful stats retrieval",
			mockSetup: func() {
				mockStats.EXPECT().
					GetStats(gomock.Any()).
					Return(150, 25, nil)
			},
			wantURLs:  150,
			wantUsers: 25,
			wantErr:   false,
		},
		{
			name: "empty stats - no URLs or users",
			mockSetup: func() {
				mockStats.EXPECT().
					GetStats(gomock.Any()).
					Return(0, 0, nil)
			},
			wantURLs:  0,
			wantUsers: 0,
			wantErr:   false,
		},
		{
			name: "stats service returns error",
			mockSetup: func() {
				mockStats.EXPECT().
					GetStats(gomock.Any()).
					Return(0, 0, errors.New("database connection failed"))
			},
			wantURLs:    0,
			wantUsers:   0,
			wantErr:     true,
			expectedErr: errors.New("database connection failed"),
		},
		{
			name: "stats with large numbers",
			mockSetup: func() {
				mockStats.EXPECT().
					GetStats(gomock.Any()).
					Return(10000, 500, nil)
			},
			wantURLs:  10000,
			wantUsers: 500,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			urlsCount, usersCount, err := svc.GetStats(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantURLs, urlsCount)
				assert.Equal(t, tt.wantUsers, usersCount)
			}
		})
	}
}
