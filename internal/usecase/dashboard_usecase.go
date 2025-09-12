package usecase

import (
	"fmt"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
)

// DashboardUseCase реализует бизнес-логику для дашборда
type DashboardUseCase struct {
	patientRepo     domain.PatientRepository
	appointmentRepo domain.AppointmentRepository
}

// NewDashboardUseCase создает новый экземпляр DashboardUseCase
func NewDashboardUseCase(
	patientRepo domain.PatientRepository,
	appointmentRepo domain.AppointmentRepository,
) *DashboardUseCase {
	return &DashboardUseCase{
		patientRepo:     patientRepo,
		appointmentRepo: appointmentRepo,
	}
}

// GetDashboardStats получает статистику дашборда
func (u *DashboardUseCase) GetDashboardStats() (*domain.DashboardStats, error) {
	// Получаем всех пациентов
	patients, err := u.patientRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Получаем все записи
	appointments, err := u.appointmentRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Подсчитываем статистику
	today := time.Now().Truncate(24 * time.Hour)
	totalPatients := len(patients)
	totalAppointments := len(appointments)
	completedAppointments := 0
	totalRevenue := 0.0
	todayAppointments := 0
	pendingAppointments := 0

	for _, appointment := range appointments {
		if appointment.Status == domain.StatusCompleted {
			completedAppointments++
			totalRevenue += float64(appointment.Cost)
		}
		
		appointmentDate := appointment.Date.Truncate(24 * time.Hour)
		if appointmentDate.Equal(today) {
			todayAppointments++
		}
		
		if appointment.Status == domain.StatusScheduled {
			pendingAppointments++
		}
	}

	return &domain.DashboardStats{
		TotalPatients:         totalPatients,
		TotalAppointments:     totalAppointments,
		CompletedAppointments: completedAppointments,
		TotalRevenue:          totalRevenue,
		TodayAppointments:     todayAppointments,
		PendingAppointments:   pendingAppointments,
	}, nil
}

// GetFinanceReport получает финансовый отчет
func (u *DashboardUseCase) GetFinanceReport() (*domain.FinanceReport, error) {
	appointments, err := u.appointmentRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Группируем доходы по дням и неделям
	dayIncome := make(map[string]float64)
	weekIncome := make(map[string]float64)
	totalIncome := 0.0

	for _, appointment := range appointments {
		if appointment.Status == domain.StatusCompleted {
			income := float64(appointment.Cost)
			totalIncome += income
			
			// Доход по дням
			dateStr := appointment.Date.Format("2006-01-02")
			dayIncome[dateStr] += income

			// Доход по неделям
			year, week := appointment.Date.ISOWeek()
			weekKey := formatWeekKey(year, week)
			weekIncome[weekKey] += income
		}
	}

	// Формируем отчет по дням
	var byDay []domain.DayIncome
	for date, income := range dayIncome {
		byDay = append(byDay, domain.DayIncome{
			Date:   date,
			Income: income,
		})
	}

	// Формируем отчет по неделям
	var byWeek []domain.WeekIncome
	for week, income := range weekIncome {
		byWeek = append(byWeek, domain.WeekIncome{
			Week:   week,
			Income: income,
		})
	}

	return &domain.FinanceReport{
		TotalIncome: totalIncome,
		ByDay:       byDay,
		ByWeek:      byWeek,
	}, nil
}

// formatWeekKey форматирует ключ недели
func formatWeekKey(year, week int) string {
	return fmt.Sprintf("%d-W%02d", year, week)
}
