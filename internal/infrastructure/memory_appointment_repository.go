package infrastructure

import (
	"errors"
	"sync"
	"time"

	"github.com/sdk17/crm_ar/internal/domain"
)

// MemoryAppointmentRepository реализует AppointmentRepository в памяти
type MemoryAppointmentRepository struct {
	appointments []*domain.Appointment
	nextID       int
	mutex        sync.RWMutex
}

// NewMemoryAppointmentRepository создает новый экземпляр MemoryAppointmentRepository
func NewMemoryAppointmentRepository() *MemoryAppointmentRepository {
	repo := &MemoryAppointmentRepository{
		appointments: make([]*domain.Appointment, 0),
		nextID:       1,
	}

	// Добавляем тестовые данные
	repo.seedData()
	return repo
}

// seedData добавляет тестовые данные
func (r *MemoryAppointmentRepository) seedData() {
	now := time.Now()
	testAppointments := []*domain.Appointment{
		{
			ID:          1,
			PatientID:   1,
			PatientName: "Иванов Иван Иванович",
			Date:        time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC),
			Time:        "10:00",
			Service:     "diagnosis",
			Doctor:      "Др. Смит",
			Status:      domain.StatusScheduled,
			Price:       5000,
			Duration:    30,
			Notes:       "Первый визит",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			PatientID:   2,
			PatientName: "Петрова Анна Сергеевна",
			Date:        time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC),
			Time:        "14:30",
			Service:     "treatment",
			Doctor:      "Др. Джонс",
			Status:      domain.StatusScheduled,
			Price:       15000,
			Duration:    60,
			Notes:       "Беременная пациентка",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          3,
			PatientID:   3,
			PatientName: "Сидоров Петр Александрович",
			Date:        time.Date(2025, 1, 26, 0, 0, 0, 0, time.UTC),
			Time:        "09:00",
			Service:     "prosthetics",
			Doctor:      "Др. Уилсон",
			Status:      domain.StatusCompleted,
			Price:       50000,
			Duration:    120,
			Notes:       "Установка коронки",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	r.appointments = append(r.appointments, testAppointments...)
	r.nextID = 4
}

// GetByID получает запись по ID
func (r *MemoryAppointmentRepository) GetByID(id int) (*domain.Appointment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, appointment := range r.appointments {
		if appointment.ID == id {
			return appointment, nil
		}
	}
	return nil, errors.New("appointment not found")
}

// GetAll получает все записи
func (r *MemoryAppointmentRepository) GetAll() ([]*domain.Appointment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*domain.Appointment, len(r.appointments))
	copy(result, r.appointments)
	return result, nil
}

// Create создает новую запись
func (r *MemoryAppointmentRepository) Create(appointment *domain.Appointment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	appointment.ID = r.nextID
	r.nextID++
	r.appointments = append(r.appointments, appointment)
	return nil
}

// Update обновляет запись
func (r *MemoryAppointmentRepository) Update(appointment *domain.Appointment) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, a := range r.appointments {
		if a.ID == appointment.ID {
			r.appointments[i] = appointment
			return nil
		}
	}
	return errors.New("appointment not found")
}

// Delete удаляет запись
func (r *MemoryAppointmentRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, appointment := range r.appointments {
		if appointment.ID == id {
			r.appointments = append(r.appointments[:i], r.appointments[i+1:]...)
			return nil
		}
	}
	return errors.New("appointment not found")
}

// GetByPatientID получает записи по ID пациента
func (r *MemoryAppointmentRepository) GetByPatientID(patientID int) ([]*domain.Appointment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*domain.Appointment
	for _, appointment := range r.appointments {
		if appointment.PatientID == patientID {
			result = append(result, appointment)
		}
	}
	return result, nil
}

// GetByDate получает записи по дате
func (r *MemoryAppointmentRepository) GetByDate(date time.Time) ([]*domain.Appointment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	dateStr := date.Format("2006-01-02")
	var result []*domain.Appointment
	for _, appointment := range r.appointments {
		appointmentDateStr := appointment.Date.Format("2006-01-02")
		if appointmentDateStr == dateStr {
			result = append(result, appointment)
		}
	}
	return result, nil
}

// GetByDateRange получает записи в диапазоне дат
func (r *MemoryAppointmentRepository) GetByDateRange(start, end time.Time) ([]*domain.Appointment, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*domain.Appointment
	for _, appointment := range r.appointments {
		if appointment.Date.After(start) && appointment.Date.Before(end) {
			result = append(result, appointment)
		}
	}
	return result, nil
}

// CheckTimeConflict проверяет конфликт времени
func (r *MemoryAppointmentRepository) CheckTimeConflict(date time.Time, time string, excludeID int) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	dateStr := date.Format("2006-01-02")
	for _, appointment := range r.appointments {
		if appointment.ID == excludeID {
			continue
		}
		appointmentDateStr := appointment.Date.Format("2006-01-02")
		if appointmentDateStr == dateStr && appointment.Time == time {
			return true, nil
		}
	}
	return false, nil
}
