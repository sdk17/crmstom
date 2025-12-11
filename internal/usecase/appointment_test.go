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

func TestAppointmentUseCase_GetAppointment(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockAppointmentRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByID(1).Return(&domain.Appointment{
					ID:        1,
					PatientID: 1,
					Service:   "Консультация",
					Status:    domain.StatusScheduled,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockAppointmentRepository) {},
			wantErr: true,
			errMsg:  "invalid appointment ID",
		},
		{
			name:    "invalid id negative",
			id:      -1,
			setup:   func(m *repository.MockAppointmentRepository) {},
			wantErr: true,
			errMsg:  "invalid appointment ID",
		},
		{
			name: "appointment not found",
			id:   999,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByID(999).Return(nil, errors.New("appointment not found"))
			},
			wantErr: true,
			errMsg:  "appointment not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			appointment, err := uc.GetAppointment(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, appointment)
				assert.Equal(t, tt.id, appointment.ID)
			}
		})
	}
}

func TestAppointmentUseCase_GetAllAppointments(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*repository.MockAppointmentRepository)
		want    int
		wantErr bool
	}{
		{
			name: "success with appointments",
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetAll().Return([]*domain.Appointment{
					{ID: 1, PatientID: 1, Service: "Консультация"},
					{ID: 2, PatientID: 2, Service: "Лечение"},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "success empty list",
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetAll().Return([]*domain.Appointment{}, nil)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "repository error",
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			appointments, err := uc.GetAllAppointments()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, appointments, tt.want)
			}
		})
	}
}

func TestAppointmentUseCase_CreateAppointment(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		appointment *domain.Appointment
		setup       func(*repository.MockAppointmentRepository, *repository.MockPatientRepository, *repository.MockServiceRepository)
		wantErr     bool
		errMsg      string
	}{
		{
			name: "success",
			appointment: &domain.Appointment{
				PatientID: 1,
				Date:      futureDate,
				Time:      "10:00",
				Service:   "Консультация",
				Doctor:    "Dr. Smith",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "John Doe"}, nil)
				a.EXPECT().Create(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "nil appointment",
			appointment: nil,
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
			},
			wantErr: true,
			errMsg:  "appointment cannot be nil",
		},
		{
			name: "missing patient id",
			appointment: &domain.Appointment{
				PatientID: 0,
				Date:      futureDate,
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
			},
			wantErr: true,
			errMsg:  "patient ID is required",
		},
		{
			name: "missing date",
			appointment: &domain.Appointment{
				PatientID: 1,
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
			},
			wantErr: true,
			errMsg:  "date is required",
		},
		{
			name: "missing service",
			appointment: &domain.Appointment{
				PatientID: 1,
				Date:      futureDate,
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
			},
			wantErr: true,
			errMsg:  "service is required",
		},
		{
			name: "patient not found",
			appointment: &domain.Appointment{
				PatientID: 999,
				Date:      futureDate,
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(999).Return(nil, errors.New("patient not found"))
			},
			wantErr: true,
			errMsg:  "patient not found",
		},
		{
			name: "repository create error",
			appointment: &domain.Appointment{
				PatientID: 1,
				Date:      futureDate,
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "John"}, nil)
				a.EXPECT().Create(gomock.Any()).Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			err := uc.CreateAppointment(tt.appointment)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, domain.StatusScheduled, tt.appointment.Status)
				assert.Equal(t, "John Doe", tt.appointment.PatientName)
				assert.False(t, tt.appointment.CreatedAt.IsZero())
			}
		})
	}
}

func TestAppointmentUseCase_UpdateAppointment(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		appointment *domain.Appointment
		setup       func(*repository.MockAppointmentRepository, *repository.MockPatientRepository, *repository.MockServiceRepository)
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "nil appointment validation error",
			appointment: nil,
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
			},
			wantErr: true,
			errMsg:  "appointment cannot be nil",
		},
		{
			name: "validation error missing service",
			appointment: &domain.Appointment{
				ID:        1,
				PatientID: 1,
				Date:      futureDate,
				Service:   "",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
			},
			wantErr: true,
			errMsg:  "service is required",
		},
		{
			name: "success",
			appointment: &domain.Appointment{
				ID:        1,
				PatientID: 1,
				Date:      futureDate,
				Time:      "10:00",
				Service:   "Консультация",
				Status:    domain.StatusScheduled,
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "John Doe"}, nil)
				a.EXPECT().CheckTimeConflict(futureDate, "10:00", 1).Return(false, nil)
				a.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success with price and duration",
			appointment: &domain.Appointment{
				ID:        1,
				PatientID: 1,
				Date:      futureDate,
				Time:      "14:00",
				Service:   "Лечение",
				Status:    domain.StatusScheduled,
				Price:     15000.00,
				Duration:  60,
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "Jane Doe"}, nil)
				a.EXPECT().CheckTimeConflict(futureDate, "14:00", 1).Return(false, nil)
				a.EXPECT().Update(gomock.Any()).DoAndReturn(func(apt *domain.Appointment) error {
					assert.Equal(t, 15000.00, apt.Price)
					assert.Equal(t, 60, apt.Duration)
					assert.Equal(t, "Jane Doe", apt.PatientName)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "time conflict",
			appointment: &domain.Appointment{
				ID:        1,
				PatientID: 1,
				Date:      futureDate,
				Time:      "10:00",
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "John"}, nil)
				a.EXPECT().CheckTimeConflict(futureDate, "10:00", 1).Return(true, nil)
			},
			wantErr: true,
			errMsg:  "time slot is already occupied",
		},
		{
			name: "patient not found on update",
			appointment: &domain.Appointment{
				ID:        1,
				PatientID: 999,
				Date:      futureDate,
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(999).Return(nil, errors.New("patient not found"))
			},
			wantErr: true,
			errMsg:  "patient not found",
		},
		{
			name: "check time conflict error",
			appointment: &domain.Appointment{
				ID:        1,
				PatientID: 1,
				Date:      futureDate,
				Time:      "10:00",
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "John"}, nil)
				a.EXPECT().CheckTimeConflict(futureDate, "10:00", 1).Return(false, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name: "update repository error",
			appointment: &domain.Appointment{
				ID:        1,
				PatientID: 1,
				Date:      futureDate,
				Time:      "10:00",
				Service:   "Консультация",
			},
			setup: func(a *repository.MockAppointmentRepository, p *repository.MockPatientRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetByID(1).Return(&domain.Patient{ID: 1, Name: "John"}, nil)
				a.EXPECT().CheckTimeConflict(futureDate, "10:00", 1).Return(false, nil)
				a.EXPECT().Update(gomock.Any()).Return(errors.New("update failed"))
			},
			wantErr: true,
			errMsg:  "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			err := uc.UpdateAppointment(tt.appointment)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.False(t, tt.appointment.UpdatedAt.IsZero())
			}
		})
	}
}

func TestAppointmentUseCase_DeleteAppointment(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockAppointmentRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().Delete(1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid id zero",
			id:      0,
			setup:   func(m *repository.MockAppointmentRepository) {},
			wantErr: true,
			errMsg:  "invalid appointment ID",
		},
		{
			name:    "invalid id negative",
			id:      -1,
			setup:   func(m *repository.MockAppointmentRepository) {},
			wantErr: true,
			errMsg:  "invalid appointment ID",
		},
		{
			name: "not found",
			id:   999,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().Delete(999).Return(errors.New("appointment not found"))
			},
			wantErr: true,
			errMsg:  "appointment not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			err := uc.DeleteAppointment(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAppointmentUseCase_GetAppointmentsByPatient(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockAppointmentRepository)
		want    int
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByPatientID(1).Return([]*domain.Appointment{
					{ID: 1, PatientID: 1},
					{ID: 2, PatientID: 1},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name:    "invalid patient id",
			id:      0,
			setup:   func(m *repository.MockAppointmentRepository) {},
			wantErr: true,
			errMsg:  "invalid patient ID",
		},
		{
			name: "empty result",
			id:   1,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByPatientID(1).Return([]*domain.Appointment{}, nil)
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			appointments, err := uc.GetAppointmentsByPatient(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Len(t, appointments, tt.want)
			}
		})
	}
}

func TestAppointmentUseCase_GetAppointmentsByDate(t *testing.T) {
	testDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		date    time.Time
		setup   func(*repository.MockAppointmentRepository)
		want    int
		wantErr bool
	}{
		{
			name: "success with appointments",
			date: testDate,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByDate(testDate).Return([]*domain.Appointment{
					{ID: 1, Date: testDate},
					{ID: 2, Date: testDate},
				}, nil)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "empty result",
			date: testDate,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByDate(testDate).Return([]*domain.Appointment{}, nil)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "repository error",
			date: testDate,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByDate(testDate).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			appointments, err := uc.GetAppointmentsByDate(tt.date)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, appointments, tt.want)
			}
		})
	}
}

func TestAppointmentUseCase_CompleteAppointment(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockAppointmentRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByID(1).Return(&domain.Appointment{
					ID:     1,
					Status: domain.StatusScheduled,
				}, nil)
				m.EXPECT().Update(gomock.Any()).DoAndReturn(func(apt *domain.Appointment) error {
					assert.Equal(t, domain.StatusCompleted, apt.Status)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "appointment not found",
			id:   999,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByID(999).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "update error",
			id:   1,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByID(1).Return(&domain.Appointment{ID: 1, Status: domain.StatusScheduled}, nil)
				m.EXPECT().Update(gomock.Any()).Return(errors.New("update failed"))
			},
			wantErr: true,
			errMsg:  "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			err := uc.CompleteAppointment(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAppointmentUseCase_CancelAppointment(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		setup   func(*repository.MockAppointmentRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success",
			id:   1,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByID(1).Return(&domain.Appointment{
					ID:     1,
					Status: domain.StatusScheduled,
				}, nil)
				m.EXPECT().Update(gomock.Any()).DoAndReturn(func(apt *domain.Appointment) error {
					assert.Equal(t, domain.StatusCancelled, apt.Status)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name: "appointment not found",
			id:   999,
			setup: func(m *repository.MockAppointmentRepository) {
				m.EXPECT().GetByID(999).Return(nil, errors.New("not found"))
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)
			err := uc.CancelAppointment(tt.id)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAppointmentUseCase_ValidateAppointment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
	mockPatientRepo := repository.NewMockPatientRepository(ctrl)
	mockServiceRepo := repository.NewMockServiceRepository(ctrl)
	uc := NewAppointmentUseCase(mockAppointmentRepo, mockPatientRepo, mockServiceRepo)

	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		appointment *domain.Appointment
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "nil appointment",
			appointment: nil,
			wantErr:     true,
			errMsg:      "appointment cannot be nil",
		},
		{
			name: "valid appointment",
			appointment: &domain.Appointment{
				PatientID: 1,
				Date:      futureDate,
				Service:   "Консультация",
			},
			wantErr: false,
		},
		{
			name: "missing patient id",
			appointment: &domain.Appointment{
				Date:    futureDate,
				Service: "Консультация",
			},
			wantErr: true,
			errMsg:  "patient ID is required",
		},
		{
			name: "missing date",
			appointment: &domain.Appointment{
				PatientID: 1,
				Service:   "Консультация",
			},
			wantErr: true,
			errMsg:  "date is required",
		},
		{
			name: "missing service",
			appointment: &domain.Appointment{
				PatientID: 1,
				Date:      futureDate,
			},
			wantErr: true,
			errMsg:  "service is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateAppointment(tt.appointment)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
