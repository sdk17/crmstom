package domain

import (
	"time"
)

// Patient представляет пациента в доменной модели
type Patient struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	Address   string    `json:"address"`
	Notes     string    `json:"notes"`
	LastVisit time.Time `json:"last_visit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PatientRepository определяет интерфейс для работы с пациентами
type PatientRepository interface {
	GetByID(id int) (*Patient, error)
	GetAll() ([]*Patient, error)
	Create(patient *Patient) error
	Update(patient *Patient) error
	Delete(id int) error
	Search(query string) ([]*Patient, error)
	GetByPhone(phone string) (*Patient, error)
}

// PatientService определяет бизнес-логику для работы с пациентами
type PatientService interface {
	GetPatient(id int) (*Patient, error)
	GetAllPatients() ([]*Patient, error)
	CreatePatient(patient *Patient) error
	UpdatePatient(patient *Patient) error
	DeletePatient(id int) error
	SearchPatients(query string) ([]*Patient, error)
	ValidatePatient(patient *Patient) error
}
