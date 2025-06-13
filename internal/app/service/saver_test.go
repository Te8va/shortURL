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

func TestURLService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := mocks.NewMockURLSaverServ(ctrl)

	svc := service.NewURLService(mockSaver, nil, nil, nil)

	tests := []struct {
		name      string
		userID    int
		url       string
		mockSetup func()
		wantID    string
		wantErr   bool
	}{
		{
			name:   "successful save",
			userID: 1,
			url:    "https://example.com",
			mockSetup: func() {
				mockSaver.EXPECT().
					Save(gomock.Any(), 1, "https://example.com").
					Return("short1234", nil)
			},
			wantID:  "short1234",
			wantErr: false,
		},
		{
			name:   "save error",
			userID: 2,
			url:    "https://fail.com",
			mockSetup: func() {
				mockSaver.EXPECT().
					Save(gomock.Any(), 2, "https://fail.com").
					Return("", errors.New("save failed"))
			},
			wantID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			gotID, err := svc.Save(context.Background(), tt.userID, tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, gotID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, gotID)
			}
		})
	}
}

func TestURLService_SaveBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaver := mocks.NewMockURLSaverServ(ctrl)

	svc := service.NewURLService(mockSaver, nil, nil, nil)

	batchInput := map[string]string{
		"corr1": "https://example1.com",
		"corr2": "https://example2.com",
	}

	batchOutput := map[string]string{
		"corr1": "short1",
		"corr2": "short2",
	}

	tests := []struct {
		name      string
		userID    int
		urls      map[string]string
		mockSetup func()
		wantMap   map[string]string
		wantErr   bool
	}{
		{
			name:   "successful batch save",
			userID: 1,
			urls:   batchInput,
			mockSetup: func() {
				mockSaver.EXPECT().
					SaveBatch(gomock.Any(), 1, batchInput).
					Return(batchOutput, nil)
			},
			wantMap: batchOutput,
			wantErr: false,
		},
		{
			name:   "batch save error",
			userID: 1,
			urls:   batchInput,
			mockSetup: func() {
				mockSaver.EXPECT().
					SaveBatch(gomock.Any(), 1, batchInput).
					Return(nil, errors.New("batch save failed"))
			},
			wantMap: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			gotMap, err := svc.SaveBatch(context.Background(), tt.userID, tt.urls)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, gotMap)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMap, gotMap)
			}
		})
	}
}
