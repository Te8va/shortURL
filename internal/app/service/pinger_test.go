package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Te8va/shortURL/internal/app/service"
	"github.com/Te8va/shortURL/internal/app/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestURLService_PingPg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPinger := mocks.NewMockPingerServ(ctrl)

	svc := service.NewURLService(nil, nil, mockPinger, nil)

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "ping successful",
			mockSetup: func() {
				mockPinger.EXPECT().
					PingPg(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ping failed",
			mockSetup: func() {
				mockPinger.EXPECT().
					PingPg(gomock.Any()).
					Return(errors.New("connection error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := svc.PingPg(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
