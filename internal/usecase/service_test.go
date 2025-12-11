package usecase

import (
	"errors"
	"testing"

	"github.com/sdk17/crmstom/gen/mocks/repository"
	"github.com/sdk17/crmstom/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestServiceUseCase_GetService(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockServiceRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetByID(1).Return(&domain.Service{ID: 1, Name: "Консультация", Type: "consultation"}, nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "invalid service ID",
		},
		{
			name:    "invalid id negative",
			id:      -1,
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "invalid service ID",
		},
		{
			name: "service not found",
			id:   999,
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetByID(999).Return(nil, errors.New("service not found"))
			},
			wantErr: true,
			errMsg:  "service not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewServiceUseCase(mockRepo)

			service, err := uc.GetService(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, service)
				assert.Equal(t, tt.id, service.ID)
			}
		})
	}
}

func TestServiceUseCase_GetAllServices(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*repository.MockServiceRepository)
		want    int
		wantErr bool
	}{
		{
			name: "success with services",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetAll().Return([]*domain.Service{
					{ID: 1, Name: "Консультация", Type: "consultation"},
					{ID: 2, Name: "Лечение", Type: "treatment"},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success empty list",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetAll().Return([]*domain.Service{}, nil)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewServiceUseCase(mockRepo)

			services, err := uc.GetAllServices()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, services, tt.want)
			}
		})
	}
}

func TestServiceUseCase_CreateService(t *testing.T) {
	tests := []struct {
		name    string
		service *domain.Service
		setup   func(*repository.MockServiceRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			service: &domain.Service{
				Name:  "Консультация",
				Type:  "consultation",
				Notes: "Первичный осмотр",
			},
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "nil service",
			service: nil,
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service cannot be nil",
		},
		{
			name: "empty name",
			service: &domain.Service{
				Name: "",
				Type: "consultation",
			},
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service name is required",
		},
		{
			name: "whitespace only name",
			service: &domain.Service{
				Name: "   ",
				Type: "consultation",
			},
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service name is required",
		},
		{
			name: "name too long",
			service: &domain.Service{
				Name: string(make([]byte, 101)),
				Type: "consultation",
			},
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service name is too long",
		},
		{
			name: "empty type",
			service: &domain.Service{
				Name: "Консультация",
				Type: "",
			},
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service type is required",
		},
		{
			name: "type too long",
			service: &domain.Service{
				Name: "Консультация",
				Type: string(make([]byte, 51)),
			},
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service type is too long",
		},
		{
			name: "notes too long",
			service: &domain.Service{
				Name:  "Консультация",
				Type:  "consultation",
				Notes: string(make([]byte, 501)),
			},
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service notes are too long",
		},
		{
			name: "repository error",
			service: &domain.Service{
				Name: "Консультация",
				Type: "consultation",
			},
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Create(gomock.Any()).Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewServiceUseCase(mockRepo)

			err := uc.CreateService(tt.service)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.False(t, tt.service.CreatedAt.IsZero())
				assert.False(t, tt.service.UpdatedAt.IsZero())
			}
		})
	}
}

func TestServiceUseCase_UpdateService(t *testing.T) {
	tests := []struct {
		name    string
		service *domain.Service
		setup   func(*repository.MockServiceRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			service: &domain.Service{
				ID:   1,
				Name: "Консультация Updated",
				Type: "consultation",
			},
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "nil service",
			service: nil,
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "service cannot be nil",
		},
		{
			name: "service not found",
			service: &domain.Service{
				ID:   999,
				Name: "Test",
				Type: "test",
			},
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Update(gomock.Any()).Return(errors.New("service not found"))
			},
			wantErr: true,
			errMsg:  "service not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewServiceUseCase(mockRepo)

			err := uc.UpdateService(tt.service)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.False(t, tt.service.UpdatedAt.IsZero())
			}
		})
	}
}

func TestServiceUseCase_DeleteService(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockServiceRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Delete(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "invalid service ID",
		},
		{
			name:    "invalid id negative",
			id:      -1,
			setup:   func(m *repository.MockServiceRepository) {},
			wantErr: true,
			errMsg:  "invalid service ID",
		},
		{
			name: "not found",
			id:   999,
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Delete(999).Return(errors.New("service not found"))
			},
			wantErr: true,
			errMsg:  "service not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewServiceUseCase(mockRepo)

			err := uc.DeleteService(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServiceUseCase_GetServicesByCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		setup    func(*repository.MockServiceRepository)
		want     int
		wantErr  bool
	}{
		{
			name:     "success with category",
			category: "consultation",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetByCategory("consultation").Return([]*domain.Service{
					{ID: 1, Name: "Консультация", Type: "consultation"},
				}, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:     "empty category returns all",
			category: "",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetAll().Return([]*domain.Service{
					{ID: 1, Name: "Service 1"},
					{ID: 2, Name: "Service 2"},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name:     "whitespace category returns all",
			category: "   ",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetAll().Return([]*domain.Service{
					{ID: 1, Name: "Service 1"},
				}, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:     "repository error",
			category: "consultation",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetByCategory("consultation").Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewServiceUseCase(mockRepo)

			services, err := uc.GetServicesByCategory(tt.category)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, services, tt.want)
			}
		})
	}
}

func TestServiceUseCase_SearchServices(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		setup   func(*repository.MockServiceRepository)
		want    int
		wantErr bool
	}{
		{
			name:  "search with query",
			query: "консульт",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Search("консульт").Return([]*domain.Service{
					{ID: 1, Name: "Консультация"},
				}, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:  "empty query returns all",
			query: "",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetAll().Return([]*domain.Service{
					{ID: 1, Name: "Service 1"},
					{ID: 2, Name: "Service 2"},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name:  "whitespace query returns all",
			query: "   ",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().GetAll().Return([]*domain.Service{
					{ID: 1, Name: "Service 1"},
				}, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:  "search error",
			query: "test",
			setup: func(m *repository.MockServiceRepository) {
				m.EXPECT().Search("test").Return(nil, errors.New("search failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewServiceUseCase(mockRepo)

			services, err := uc.SearchServices(tt.query)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, services, tt.want)
			}
		})
	}
}

func TestServiceUseCase_ValidateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockServiceRepository(ctrl)
	uc := NewServiceUseCase(mockRepo)

	tests := []struct {
		name    string
		service *domain.Service
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil service",
			service: nil,
			wantErr: true,
			errMsg:  "service cannot be nil",
		},
		{
			name: "valid service minimal",
			service: &domain.Service{
				Name: "Консультация",
				Type: "consultation",
			},
			wantErr: false,
		},
		{
			name: "valid service full",
			service: &domain.Service{
				Name:  "Консультация",
				Type:  "consultation",
				Notes: "Первичный осмотр",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			service: &domain.Service{
				Name: "",
				Type: "consultation",
			},
			wantErr: true,
			errMsg:  "service name is required",
		},
		{
			name: "empty type",
			service: &domain.Service{
				Name: "Консультация",
				Type: "",
			},
			wantErr: true,
			errMsg:  "service type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateService(tt.service)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
