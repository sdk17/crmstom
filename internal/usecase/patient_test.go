package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/sdk17/crmstom/gen/mocks/repository"
	"github.com/sdk17/crmstom/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPatientUseCase_GetPatient(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockPatientRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "John Doe"}, nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "invalid patient ID",
		},
		{
			name:    "invalid id negative",
			id:      -1,
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "invalid patient ID",
		},
		{
			name: "patient not found",
			id:   999,
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByID(999).Return(nil, errors.New("patient not found"))
			},
			wantErr: true,
			errMsg:  "patient not found",
		},
		{
			name: "repository error",
			id:   1,
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByID(1).Return(nil, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockPatientRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewPatientUseCase(mockRepo)

			patient, err := uc.GetPatient(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, patient)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, patient)
				assert.Equal(t, tt.id, patient.ID)
			}
		})
	}
}

func TestPatientUseCase_GetAllPatients(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*repository.MockPatientRepository)
		want    int
		wantErr bool
	}{
		{
			name: "success with patients",
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetAll().Return([]*domain.Patient{
					{ID: 1, Name: "John"},
					{ID: 2, Name: "Jane"},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success empty list",
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetAll().Return([]*domain.Patient{}, nil)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockPatientRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewPatientUseCase(mockRepo)

			patients, err := uc.GetAllPatients()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, patients, tt.want)
			}
		})
	}
}

func TestPatientUseCase_CreatePatient(t *testing.T) {
	tests := []struct {
		name    string
		patient *domain.Patient
		setup   func(*repository.MockPatientRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			patient: &domain.Patient{
				Name:  "John Doe",
				Phone: "+1234567890",
				Email: "john@example.com",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByPhone("+1234567890").Return(nil, errors.New("not found"))
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success without phone",
			patient: &domain.Patient{
				Name: "John Doe",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "nil patient",
			patient: nil,
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "patient cannot be nil",
		},
		{
			name: "empty name",
			patient: &domain.Patient{
				Name: "",
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "patient name is required",
		},
		{
			name: "whitespace only name",
			patient: &domain.Patient{
				Name: "   ",
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "patient name is required",
		},
		{
			name: "name too long",
			patient: &domain.Patient{
				Name: string(make([]byte, 101)),
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "patient name is too long",
		},
		{
			name: "phone too long",
			patient: &domain.Patient{
				Name:  "John",
				Phone: string(make([]byte, 21)),
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "phone number is too long",
		},
		{
			name: "email too long",
			patient: &domain.Patient{
				Name:  "John",
				Email: string(make([]byte, 101)),
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "email is too long",
		},
		{
			name: "invalid email format",
			patient: &domain.Patient{
				Name:  "John",
				Email: "invalid-email",
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "invalid email format",
		},
		{
			name: "address too long",
			patient: &domain.Patient{
				Name:    "John",
				Address: string(make([]byte, 201)),
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "address is too long",
		},
		{
			name: "notes too long",
			patient: &domain.Patient{
				Name:  "John",
				Notes: string(make([]byte, 501)),
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "notes are too long",
		},
		{
			name: "iin too short",
			patient: &domain.Patient{
				Name: "John",
				IIN:  "12345678901",
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "ИИН должен содержать 12 символов",
		},
		{
			name: "iin too long",
			patient: &domain.Patient{
				Name: "John",
				IIN:  "1234567890123",
			},
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "ИИН должен содержать 12 символов",
		},
		{
			name: "success with valid iin",
			patient: &domain.Patient{
				Name: "John Doe",
				IIN:  "123456789012",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByIIN("123456789012").Return(nil, errors.New("not found"))
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "duplicate iin",
			patient: &domain.Patient{
				Name: "John Doe",
				IIN:  "123456789012",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByIIN("123456789012").Return(&domain.Patient{ID: 1, IIN: "123456789012"}, nil)
			},
			wantErr: true,
			errMsg:  "пациент с таким ИИН уже существует",
		},
		{
			name: "duplicate phone",
			patient: &domain.Patient{
				Name:  "John Doe",
				Phone: "+1234567890",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByPhone("+1234567890").Return(&domain.Patient{ID: 1, Phone: "+1234567890"}, nil)
			},
			wantErr: true,
			errMsg:  "пациент с таким номером телефона уже существует",
		},
		{
			name: "repository create error",
			patient: &domain.Patient{
				Name: "John Doe",
			},
			setup: func(m *repository.MockPatientRepository) {
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

			mockRepo := repository.NewMockPatientRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewPatientUseCase(mockRepo)

			err := uc.CreatePatient(tt.patient)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.False(t, tt.patient.CreatedAt.IsZero())
				assert.False(t, tt.patient.UpdatedAt.IsZero())
			}
		})
	}
}

func TestPatientUseCase_UpdatePatient(t *testing.T) {
	tests := []struct {
		name    string
		patient *domain.Patient
		setup   func(*repository.MockPatientRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			patient: &domain.Patient{
				ID:    1,
				Name:  "John Updated",
				Phone: "+1234567890",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByPhone("+1234567890").Return(&domain.Patient{ID: 1, Phone: "+1234567890"}, nil)
				m.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "change phone to unique",
			patient: &domain.Patient{
				ID:    1,
				Name:  "John",
				Phone: "+9999999999",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByPhone("+9999999999").Return(nil, errors.New("not found"))
				m.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "duplicate phone from another patient",
			patient: &domain.Patient{
				ID:    1,
				Name:  "John",
				Phone: "+9999999999",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByPhone("+9999999999").Return(&domain.Patient{ID: 2, Phone: "+9999999999"}, nil)
			},
			wantErr: true,
			errMsg:  "пациент с таким номером телефона уже существует",
		},
		{
			name: "update success with iin",
			patient: &domain.Patient{
				ID:   1,
				Name: "John Updated",
				IIN:  "123456789012",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByIIN("123456789012").Return(&domain.Patient{ID: 1, IIN: "123456789012"}, nil)
				m.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "change iin to unique",
			patient: &domain.Patient{
				ID:   1,
				Name: "John",
				IIN:  "999999999999",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByIIN("999999999999").Return(nil, errors.New("not found"))
				m.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "duplicate iin from another patient",
			patient: &domain.Patient{
				ID:   1,
				Name: "John",
				IIN:  "999999999999",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetByIIN("999999999999").Return(&domain.Patient{ID: 2, IIN: "999999999999"}, nil)
			},
			wantErr: true,
			errMsg:  "пациент с таким ИИН уже существует",
		},
		{
			name: "patient not found",
			patient: &domain.Patient{
				ID:   999,
				Name: "John",
			},
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().Update(gomock.Any()).Return(errors.New("patient not found"))
			},
			wantErr: true,
			errMsg:  "patient not found",
		},
		{
			name:    "nil patient",
			patient: nil,
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "patient cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockPatientRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewPatientUseCase(mockRepo)

			err := uc.UpdatePatient(tt.patient)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.False(t, tt.patient.UpdatedAt.IsZero())
			}
		})
	}
}

func TestPatientUseCase_DeletePatient(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockPatientRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().Delete(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "invalid patient ID",
		},
		{
			name:    "invalid id negative",
			id:      -5,
			setup:   func(m *repository.MockPatientRepository) {},
			wantErr: true,
			errMsg:  "invalid patient ID",
		},
		{
			name: "patient not found",
			id:   999,
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().Delete(999).Return(errors.New("patient not found"))
			},
			wantErr: true,
			errMsg:  "patient not found",
		},
		{
			name: "repository error",
			id:   1,
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().Delete(1).Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockPatientRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewPatientUseCase(mockRepo)

			err := uc.DeletePatient(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPatientUseCase_SearchPatients(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		setup   func(*repository.MockPatientRepository)
		want    int
		wantErr bool
	}{
		{
			name:  "search with query",
			query: "john",
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().Search("john").Return([]*domain.Patient{{ID: 1, Name: "John Doe"}}, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:  "empty query returns all",
			query: "",
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetAll().Return([]*domain.Patient{
					{ID: 1, Name: "John"},
					{ID: 2, Name: "Jane"},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name:  "whitespace query returns all",
			query: "   ",
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().GetAll().Return([]*domain.Patient{{ID: 1, Name: "John"}}, nil)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:  "search error",
			query: "test",
			setup: func(m *repository.MockPatientRepository) {
				m.EXPECT().Search("test").Return(nil, errors.New("search failed"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repository.NewMockPatientRepository(ctrl)
			tt.setup(mockRepo)
			uc := NewPatientUseCase(mockRepo)

			patients, err := uc.SearchPatients(tt.query)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, patients, tt.want)
			}
		})
	}
}

func TestPatientUseCase_ValidatePatient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockPatientRepository(ctrl)
	uc := NewPatientUseCase(mockRepo)

	tests := []struct {
		name    string
		patient *domain.Patient
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil patient",
			patient: nil,
			wantErr: true,
			errMsg:  "patient cannot be nil",
		},
		{
			name:    "valid patient minimal",
			patient: &domain.Patient{Name: "John"},
			wantErr: false,
		},
		{
			name: "valid patient full",
			patient: &domain.Patient{
				Name:      "John Doe",
				Phone:     "+1234567890",
				Email:     "john@example.com",
				Address:   "123 Main St",
				Notes:     "Some notes",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:    "empty name",
			patient: &domain.Patient{Name: ""},
			wantErr: true,
			errMsg:  "patient name is required",
		},
		{
			name:    "valid email",
			patient: &domain.Patient{Name: "John", Email: "test@test.com"},
			wantErr: false,
		},
		{
			name:    "valid iin 12 chars",
			patient: &domain.Patient{Name: "John", IIN: "123456789012"},
			wantErr: false,
		},
		{
			name:    "invalid iin 11 chars",
			patient: &domain.Patient{Name: "John", IIN: "12345678901"},
			wantErr: true,
			errMsg:  "ИИН должен содержать 12 символов",
		},
		{
			name:    "invalid iin 13 chars",
			patient: &domain.Patient{Name: "John", IIN: "1234567890123"},
			wantErr: true,
			errMsg:  "ИИН должен содержать 12 символов",
		},
		{
			name:    "empty iin is valid",
			patient: &domain.Patient{Name: "John", IIN: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidatePatient(tt.patient)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
