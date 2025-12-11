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

func TestDashboardUseCase_GetDashboardStats(t *testing.T) {
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	tests := []struct {
		name                  string
		setup                 func(*repository.MockPatientRepository, *repository.MockAppointmentRepository, *repository.MockServiceRepository)
		wantTodayAppointments int
		wantTodayRevenue      float64
		wantTotalPatients     int
		wantErr               bool
	}{
		{
			name: "success with data",
			setup: func(p *repository.MockPatientRepository, a *repository.MockAppointmentRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetAll().Return([]*domain.Patient{
					{ID: 1, Name: "John"},
					{ID: 2, Name: "Jane"},
					{ID: 3, Name: "Bob"},
				}, nil)
				a.EXPECT().GetAll().Return([]*domain.Appointment{
					{ID: 1, Date: today, Status: domain.StatusScheduled, Price: 1000},
					{ID: 2, Date: today, Status: domain.StatusCompleted, Price: 2000},
					{ID: 3, Date: today, Status: domain.StatusCompleted, Price: 3000},
					{ID: 4, Date: yesterday, Status: domain.StatusCompleted, Price: 5000},
				}, nil)
			},
			wantTodayAppointments: 3,
			wantTodayRevenue:      5000,
			wantTotalPatients:     3,
			wantErr:               false,
		},
		{
			name: "success with no today appointments",
			setup: func(p *repository.MockPatientRepository, a *repository.MockAppointmentRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetAll().Return([]*domain.Patient{
					{ID: 1, Name: "John"},
				}, nil)
				a.EXPECT().GetAll().Return([]*domain.Appointment{
					{ID: 1, Date: yesterday, Status: domain.StatusCompleted, Price: 5000},
				}, nil)
			},
			wantTodayAppointments: 0,
			wantTodayRevenue:      0,
			wantTotalPatients:     1,
			wantErr:               false,
		},
		{
			name: "success empty data",
			setup: func(p *repository.MockPatientRepository, a *repository.MockAppointmentRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetAll().Return([]*domain.Patient{}, nil)
				a.EXPECT().GetAll().Return([]*domain.Appointment{}, nil)
			},
			wantTodayAppointments: 0,
			wantTodayRevenue:      0,
			wantTotalPatients:     0,
			wantErr:               false,
		},
		{
			name: "only scheduled appointments today - no revenue",
			setup: func(p *repository.MockPatientRepository, a *repository.MockAppointmentRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetAll().Return([]*domain.Patient{{ID: 1}}, nil)
				a.EXPECT().GetAll().Return([]*domain.Appointment{
					{ID: 1, Date: today, Status: domain.StatusScheduled, Price: 1000},
					{ID: 2, Date: today, Status: domain.StatusCancelled, Price: 2000},
				}, nil)
			},
			wantTodayAppointments: 2,
			wantTodayRevenue:      0,
			wantTotalPatients:     1,
			wantErr:               false,
		},
		{
			name: "patient repository error",
			setup: func(p *repository.MockPatientRepository, a *repository.MockAppointmentRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "appointment repository error",
			setup: func(p *repository.MockPatientRepository, a *repository.MockAppointmentRepository, s *repository.MockServiceRepository) {
				p.EXPECT().GetAll().Return([]*domain.Patient{}, nil)
				a.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockPatientRepo, mockAppointmentRepo, mockServiceRepo)

			uc := NewDashboardUseCase(mockPatientRepo, mockAppointmentRepo, mockServiceRepo)
			stats, err := uc.GetDashboardStats()

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, stats)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, stats)
				assert.Equal(t, tt.wantTodayAppointments, stats.TodayAppointments)
				assert.Equal(t, tt.wantTodayRevenue, stats.TodayRevenue)
				assert.Equal(t, tt.wantTotalPatients, stats.TotalPatients)
			}
		})
	}
}

func TestDashboardUseCase_GetFinanceReport(t *testing.T) {
	date1 := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	date2 := time.Date(2025, 1, 16, 14, 0, 0, 0, time.UTC)
	date3 := time.Date(2025, 1, 20, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name            string
		setup           func(*repository.MockAppointmentRepository)
		wantTotalIncome float64
		wantDayCount    int
		wantWeekCount   int
		wantErr         bool
	}{
		{
			name: "success with completed appointments",
			setup: func(a *repository.MockAppointmentRepository) {
				a.EXPECT().GetAll().Return([]*domain.Appointment{
					{ID: 1, Date: date1, Status: domain.StatusCompleted, Price: 1000},
					{ID: 2, Date: date1, Status: domain.StatusCompleted, Price: 2000},
					{ID: 3, Date: date2, Status: domain.StatusCompleted, Price: 3000},
					{ID: 4, Date: date3, Status: domain.StatusCompleted, Price: 4000},
				}, nil)
			},
			wantTotalIncome: 10000,
			wantDayCount:    3,
			wantWeekCount:   2,
			wantErr:         false,
		},
		{
			name: "only completed appointments counted",
			setup: func(a *repository.MockAppointmentRepository) {
				a.EXPECT().GetAll().Return([]*domain.Appointment{
					{ID: 1, Date: date1, Status: domain.StatusCompleted, Price: 1000},
					{ID: 2, Date: date1, Status: domain.StatusScheduled, Price: 2000},
					{ID: 3, Date: date2, Status: domain.StatusCancelled, Price: 3000},
				}, nil)
			},
			wantTotalIncome: 1000,
			wantDayCount:    1,
			wantWeekCount:   1,
			wantErr:         false,
		},
		{
			name: "empty appointments",
			setup: func(a *repository.MockAppointmentRepository) {
				a.EXPECT().GetAll().Return([]*domain.Appointment{}, nil)
			},
			wantTotalIncome: 0,
			wantDayCount:    0,
			wantWeekCount:   0,
			wantErr:         false,
		},
		{
			name: "all appointments not completed",
			setup: func(a *repository.MockAppointmentRepository) {
				a.EXPECT().GetAll().Return([]*domain.Appointment{
					{ID: 1, Date: date1, Status: domain.StatusScheduled, Price: 1000},
					{ID: 2, Date: date2, Status: domain.StatusCancelled, Price: 2000},
				}, nil)
			},
			wantTotalIncome: 0,
			wantDayCount:    0,
			wantWeekCount:   0,
			wantErr:         false,
		},
		{
			name: "repository error",
			setup: func(a *repository.MockAppointmentRepository) {
				a.EXPECT().GetAll().Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPatientRepo := repository.NewMockPatientRepository(ctrl)
			mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
			mockServiceRepo := repository.NewMockServiceRepository(ctrl)
			tt.setup(mockAppointmentRepo)

			uc := NewDashboardUseCase(mockPatientRepo, mockAppointmentRepo, mockServiceRepo)
			report, err := uc.GetFinanceReport()

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, report)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, report)
				assert.Equal(t, tt.wantTotalIncome, report.TotalIncome)
				assert.Len(t, report.ByDay, tt.wantDayCount)
				assert.Len(t, report.ByWeek, tt.wantWeekCount)
			}
		})
	}
}

func TestDashboardUseCase_GetFinanceReport_IncomeAggregation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	date := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

	mockPatientRepo := repository.NewMockPatientRepository(ctrl)
	mockAppointmentRepo := repository.NewMockAppointmentRepository(ctrl)
	mockServiceRepo := repository.NewMockServiceRepository(ctrl)

	mockAppointmentRepo.EXPECT().GetAll().Return([]*domain.Appointment{
		{ID: 1, Date: date, Status: domain.StatusCompleted, Price: 1000},
		{ID: 2, Date: date, Status: domain.StatusCompleted, Price: 2000},
		{ID: 3, Date: date, Status: domain.StatusCompleted, Price: 3000},
	}, nil)

	uc := NewDashboardUseCase(mockPatientRepo, mockAppointmentRepo, mockServiceRepo)
	report, err := uc.GetFinanceReport()

	require.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 6000.0, report.TotalIncome)
	assert.Len(t, report.ByDay, 1)
	assert.Equal(t, 6000.0, report.ByDay[0].Income)
	assert.Equal(t, "2025-01-15", report.ByDay[0].Date)
}

func TestFormatWeekKey(t *testing.T) {
	tests := []struct {
		year int
		week int
		want string
	}{
		{year: 2025, week: 1, want: "2025-W01"},
		{year: 2025, week: 10, want: "2025-W10"},
		{year: 2025, week: 52, want: "2025-W52"},
		{year: 2024, week: 5, want: "2024-W05"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			result := formatWeekKey(tt.year, tt.week)
			assert.Equal(t, tt.want, result)
		})
	}
}

