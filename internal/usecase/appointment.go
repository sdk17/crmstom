package usecase

import (
	"errors"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
)

type AppointmentUseCase struct {
	appointmentRepo domain.AppointmentRepository
	patientRepo     domain.PatientRepository
	serviceRepo     domain.ServiceRepository
}

func NewAppointmentUseCase(
	appointmentRepo domain.AppointmentRepository,
	patientRepo domain.PatientRepository,
	serviceRepo domain.ServiceRepository,
) *AppointmentUseCase {
	return &AppointmentUseCase{
		appointmentRepo: appointmentRepo,
		patientRepo:     patientRepo,
		serviceRepo:     serviceRepo,
	}
}

// GetAppointment получает запись по ID
func (u *AppointmentUseCase) GetAppointment(id int) (*domain.Appointment, error) {
	if id <= 0 {
		return nil, errors.New("invalid appointment ID")
	}
	return u.appointmentRepo.GetByID(id)
}

// GetAllAppointments получает все записи
func (u *AppointmentUseCase) GetAllAppointments() ([]*domain.Appointment, error) {
	return u.appointmentRepo.GetAll()
}

// CreateAppointment создает новую запись
func (u *AppointmentUseCase) CreateAppointment(appointment *domain.Appointment) error {
	if err := u.ValidateAppointment(appointment); err != nil {
		return err
	}

	// Проверяем, что пациент существует
	patient, err := u.patientRepo.GetByID(appointment.PatientID)
	if err != nil {
		return errors.New("patient not found")
	}
	appointment.PatientName = patient.Name

	appointment.Status = domain.StatusScheduled
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()

	return u.appointmentRepo.Create(appointment)
}

// UpdateAppointment обновляет запись
func (u *AppointmentUseCase) UpdateAppointment(appointment *domain.Appointment) error {
	if err := u.ValidateAppointment(appointment); err != nil {
		return err
	}

	// Проверяем, что пациент существует
	patient, err := u.patientRepo.GetByID(appointment.PatientID)
	if err != nil {
		return errors.New("patient not found")
	}
	appointment.PatientName = patient.Name

	// Проверяем конфликт времени (исключая текущую запись)
	hasConflict, err := u.appointmentRepo.CheckTimeConflict(appointment.Date, appointment.Time, appointment.ID)
	if err != nil {
		return err
	}
	if hasConflict {
		return errors.New("time slot is already occupied")
	}

	appointment.UpdatedAt = time.Now()

	return u.appointmentRepo.Update(appointment)
}

// DeleteAppointment удаляет запись
func (u *AppointmentUseCase) DeleteAppointment(id int) error {
	if id <= 0 {
		return errors.New("invalid appointment ID")
	}
	return u.appointmentRepo.Delete(id)
}

// GetAppointmentsByPatient получает записи по пациенту
func (u *AppointmentUseCase) GetAppointmentsByPatient(patientID int) ([]*domain.Appointment, error) {
	if patientID <= 0 {
		return nil, errors.New("invalid patient ID")
	}
	return u.appointmentRepo.GetByPatientID(patientID)
}

// GetAppointmentsByDate получает записи по дате
func (u *AppointmentUseCase) GetAppointmentsByDate(date time.Time) ([]*domain.Appointment, error) {
	return u.appointmentRepo.GetByDate(date)
}

// CompleteAppointment завершает запись
func (u *AppointmentUseCase) CompleteAppointment(id int) error {
	appointment, err := u.appointmentRepo.GetByID(id)
	if err != nil {
		return err
	}

	appointment.Status = domain.StatusCompleted
	appointment.UpdatedAt = time.Now()

	return u.appointmentRepo.Update(appointment)
}

// CancelAppointment отменяет запись
func (u *AppointmentUseCase) CancelAppointment(id int) error {
	appointment, err := u.appointmentRepo.GetByID(id)
	if err != nil {
		return err
	}

	appointment.Status = domain.StatusCancelled
	appointment.UpdatedAt = time.Now()

	return u.appointmentRepo.Update(appointment)
}

// ValidateAppointment валидирует данные записи
func (u *AppointmentUseCase) ValidateAppointment(appointment *domain.Appointment) error {
	if appointment == nil {
		return errors.New("appointment cannot be nil")
	}

	if appointment.PatientID <= 0 {
		return errors.New("patient ID is required")
	}

	if appointment.Date.IsZero() {
		return errors.New("date is required")
	}

	if appointment.Service == "" {
		return errors.New("service is required")
	}

	return nil
}
