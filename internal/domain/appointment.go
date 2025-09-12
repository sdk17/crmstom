package domain

import (
	"time"
)

// AppointmentStatus представляет статус записи
type AppointmentStatus string

const (
	StatusScheduled AppointmentStatus = "scheduled"
	StatusCompleted AppointmentStatus = "completed"
	StatusCancelled AppointmentStatus = "cancelled"
)

// Appointment представляет запись в доменной модели
type Appointment struct {
	ID          int               `json:"id"`
	PatientID   int               `json:"patient_id"`
	PatientName string            `json:"patient_name"`
	Date        time.Time         `json:"date"`
	Time        string            `json:"time"`
	Service     string            `json:"service"`
	Doctor      string            `json:"doctor"`
	Status      AppointmentStatus `json:"status"`
	Cost        int               `json:"cost"`
	Notes       string            `json:"notes"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AppointmentRepository определяет интерфейс для работы с записями
type AppointmentRepository interface {
	GetByID(id int) (*Appointment, error)
	GetAll() ([]*Appointment, error)
	Create(appointment *Appointment) error
	Update(appointment *Appointment) error
	Delete(id int) error
	GetByPatientID(patientID int) ([]*Appointment, error)
	GetByDate(date time.Time) ([]*Appointment, error)
	GetByDateRange(start, end time.Time) ([]*Appointment, error)
	CheckTimeConflict(date time.Time, time string, excludeID int) (bool, error)
}

// AppointmentService определяет бизнес-логику для работы с записями
type AppointmentService interface {
	GetAppointment(id int) (*Appointment, error)
	GetAllAppointments() ([]*Appointment, error)
	CreateAppointment(appointment *Appointment) error
	UpdateAppointment(appointment *Appointment) error
	DeleteAppointment(id int) error
	GetAppointmentsByPatient(patientID int) ([]*Appointment, error)
	GetAppointmentsByDate(date time.Time) ([]*Appointment, error)
	CompleteAppointment(id int) error
	CancelAppointment(id int) error
	ValidateAppointment(appointment *Appointment) error
}
