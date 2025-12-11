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

func TestDoctorUseCase_GetDoctor(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockDoctorRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetByID(1).Return(&domain.Doctor{
					ID:    1,
					Name:  "Dr. Smith",
					Login: "drsmith",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "invalid doctor ID",
		},
		{
			name:    "invalid id negative",
			id:      -1,
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "invalid doctor ID",
		},
		{
			name: "doctor not found",
			id:   999,
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetByID(999).Return(nil, errors.New("doctor not found"))
			},
			wantErr: true,
			errMsg:  "doctor not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockDoctorRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewDoctorUseCase(mockRepo)

			doctor, err := uc.GetDoctor(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, doctor)
				assert.Equal(t, tt.id, doctor.ID)
			}
		})
	}
}

func TestDoctorUseCase_GetAllDoctors(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*repository.MockDoctorRepository)
		want    int
		wantErr bool
	}{
		{
			name: "success with doctors",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetAll().Return([]*domain.Doctor{
					{ID: 1, Name: "Dr. Smith"},
					{ID: 2, Name: "Dr. Jones"},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success empty list",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetAll().Return([]*domain.Doctor{}, nil)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockDoctorRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewDoctorUseCase(mockRepo)

			doctors, err := uc.GetAllDoctors()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, doctors, tt.want)
			}
		})
	}
}

func TestDoctorUseCase_CreateDoctor(t *testing.T) {
	tests := []struct {
		name    string
		doctor  *domain.Doctor
		setup   func(*repository.MockDoctorRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Email:    "drsmith@example.com",
				Login:    "drsmith",
				Password: "password123",
			},
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "empty name",
			doctor: &domain.Doctor{
				Name:     "",
				Login:    "drsmith",
				Password: "password123",
			},
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "doctor name is required",
		},
		{
			name: "name too long",
			doctor: &domain.Doctor{
				Name:     string(make([]byte, 256)),
				Login:    "drsmith",
				Password: "password123",
			},
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "doctor name is too long",
		},
		{
			name: "empty login",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "",
				Password: "password123",
			},
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "doctor login is required",
		},
		{
			name: "login too long",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    string(make([]byte, 101)),
				Password: "password123",
			},
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "doctor login is too long",
		},
		{
			name: "empty password",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "drsmith",
				Password: "",
			},
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "doctor password is required",
		},
		{
			name: "password too short",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "drsmith",
				Password: "123",
			},
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "doctor password is too short",
		},
		{
			name: "repository error",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "drsmith",
				Password: "password123",
			},
			setup: func(m *repository.MockDoctorRepository) {
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

			mockRepo := repository.NewMockDoctorRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewDoctorUseCase(mockRepo)

			err := uc.CreateDoctor(tt.doctor)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDoctorUseCase_UpdateDoctor(t *testing.T) {
	tests := []struct {
		name    string
		doctor  *domain.Doctor
		setup   func(*repository.MockDoctorRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			doctor: &domain.Doctor{
				ID:       1,
				Name:     "Dr. Smith Updated",
				Login:    "drsmith",
				Password: "newpassword",
			},
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "validation error",
			doctor: &domain.Doctor{
				ID:       1,
				Name:     "",
				Login:    "drsmith",
				Password: "password",
			},
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "doctor name is required",
		},
		{
			name: "not found",
			doctor: &domain.Doctor{
				ID:       999,
				Name:     "Dr. Unknown",
				Login:    "unknown",
				Password: "password",
			},
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().Update(gomock.Any()).Return(errors.New("doctor not found"))
			},
			wantErr: true,
			errMsg:  "doctor not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockDoctorRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewDoctorUseCase(mockRepo)

			err := uc.UpdateDoctor(tt.doctor)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDoctorUseCase_DeleteDoctor(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockDoctorRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().Delete(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "invalid doctor ID",
		},
		{
			name:    "invalid id negative",
			id:      -1,
			setup:   func(m *repository.MockDoctorRepository) {},
			wantErr: true,
			errMsg:  "invalid doctor ID",
		},
		{
			name: "not found",
			id:   999,
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().Delete(999).Return(errors.New("doctor not found"))
			},
			wantErr: true,
			errMsg:  "doctor not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockDoctorRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewDoctorUseCase(mockRepo)

			err := uc.DeleteDoctor(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDoctorUseCase_AuthenticateDoctor(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		setup    func(*repository.MockDoctorRepository)
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "success",
			login:    "drsmith",
			password: "password123",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetByLogin("drsmith").Return(&domain.Doctor{
					ID:       1,
					Name:     "Dr. Smith",
					Login:    "drsmith",
					Password: "password123",
					IsAdmin:  false,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "success admin",
			login:    "admin",
			password: "adminpass",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetByLogin("admin").Return(&domain.Doctor{
					ID:       1,
					Name:     "Admin",
					Login:    "admin",
					Password: "adminpass",
					IsAdmin:  true,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "empty login",
			login:    "",
			password: "password",
			setup:    func(m *repository.MockDoctorRepository) {},
			wantErr:  true,
			errMsg:   "login and password are required",
		},
		{
			name:     "empty password",
			login:    "drsmith",
			password: "",
			setup:    func(m *repository.MockDoctorRepository) {},
			wantErr:  true,
			errMsg:   "login and password are required",
		},
		{
			name:     "both empty",
			login:    "",
			password: "",
			setup:    func(m *repository.MockDoctorRepository) {},
			wantErr:  true,
			errMsg:   "login and password are required",
		},
		{
			name:     "user not found",
			login:    "unknown",
			password: "password",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetByLogin("unknown").Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name:     "user nil returned",
			login:    "unknown",
			password: "password",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetByLogin("unknown").Return(nil, nil)
			},
			wantErr: true,
			errMsg:  "invalid login or password",
		},
		{
			name:     "wrong password",
			login:    "drsmith",
			password: "wrongpassword",
			setup: func(m *repository.MockDoctorRepository) {
				m.EXPECT().GetByLogin("drsmith").Return(&domain.Doctor{
					ID:       1,
					Name:     "Dr. Smith",
					Login:    "drsmith",
					Password: "correctpassword",
				}, nil)
			},
			wantErr: true,
			errMsg:  "invalid login or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockDoctorRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewDoctorUseCase(mockRepo)

			doctor, err := uc.AuthenticateDoctor(tt.login, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, doctor)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, doctor)
				assert.Empty(t, doctor.Password)
			}
		})
	}
}

func TestDoctorUseCase_ValidateDoctor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockDoctorRepository(ctrl)
	uc := NewDoctorUseCase(mockRepo)

	tests := []struct {
		name    string
		doctor  *domain.Doctor
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid doctor",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "drsmith",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "valid doctor with email",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Email:    "drsmith@example.com",
				Login:    "drsmith",
				Password: "password123",
				IsAdmin:  true,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			doctor: &domain.Doctor{
				Name:     "",
				Login:    "drsmith",
				Password: "password",
			},
			wantErr: true,
			errMsg:  "doctor name is required",
		},
		{
			name: "empty login",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "",
				Password: "password",
			},
			wantErr: true,
			errMsg:  "doctor login is required",
		},
		{
			name: "empty password",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "drsmith",
				Password: "",
			},
			wantErr: true,
			errMsg:  "doctor password is required",
		},
		{
			name: "password too short",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "drsmith",
				Password: "abc",
			},
			wantErr: true,
			errMsg:  "doctor password is too short",
		},
		{
			name: "password exactly 4 chars",
			doctor: &domain.Doctor{
				Name:     "Dr. Smith",
				Login:    "drsmith",
				Password: "abcd",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateDoctor(tt.doctor)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

