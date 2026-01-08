package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
)

type PatientUseCase struct {
	patientRepo domain.PatientRepository
}

func NewPatientUseCase(patientRepo domain.PatientRepository) *PatientUseCase {
	return &PatientUseCase{
		patientRepo: patientRepo,
	}
}

// GetPatient получает пациента по ID
func (u *PatientUseCase) GetPatient(id int) (*domain.Patient, error) {
	if id <= 0 {
		return nil, errors.New("invalid patient ID")
	}
	return u.patientRepo.GetByID(id)
}

// GetAllPatients получает всех пациентов
func (u *PatientUseCase) GetAllPatients() ([]*domain.Patient, error) {
	return u.patientRepo.GetAll()
}

// CreatePatient создает нового пациента
func (u *PatientUseCase) CreatePatient(patient *domain.Patient) error {
	if err := u.ValidatePatient(patient); err != nil {
		return err
	}

	// Проверяем, не существует ли уже пациент с таким ИИН
	if patient.IIN != "" {
		existingPatient, err := u.patientRepo.GetByIIN(patient.IIN)
		if err == nil && existingPatient != nil {
			return errors.New("пациент с таким ИИН уже существует")
		}
	}

	// Проверяем, не существует ли уже пациент с таким телефоном
	if patient.Phone != "" {
		existingPatient, err := u.patientRepo.GetByPhone(patient.Phone)
		if err == nil && existingPatient != nil {
			return errors.New("пациент с таким номером телефона уже существует")
		}
	}

	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()

	return u.patientRepo.Create(patient)
}

// UpdatePatient обновляет пациента
func (u *PatientUseCase) UpdatePatient(patient *domain.Patient) error {
	if err := u.ValidatePatient(patient); err != nil {
		return err
	}

	// Проверяем, не существует ли уже другой пациент с таким ИИН
	if patient.IIN != "" {
		existingPatient, err := u.patientRepo.GetByIIN(patient.IIN)
		if err == nil && existingPatient != nil && existingPatient.ID != patient.ID {
			return errors.New("пациент с таким ИИН уже существует")
		}
	}

	// Проверяем, не существует ли уже другой пациент с таким телефоном
	if patient.Phone != "" {
		existingPatient, err := u.patientRepo.GetByPhone(patient.Phone)
		if err == nil && existingPatient != nil && existingPatient.ID != patient.ID {
			return errors.New("пациент с таким номером телефона уже существует")
		}
	}

	patient.UpdatedAt = time.Now()

	return u.patientRepo.Update(patient)
}

// DeletePatient удаляет пациента
func (u *PatientUseCase) DeletePatient(id int) error {
	if id <= 0 {
		return errors.New("invalid patient ID")
	}
	return u.patientRepo.Delete(id)
}

// SearchPatients ищет пациентов по запросу
func (u *PatientUseCase) SearchPatients(query string) ([]*domain.Patient, error) {
	if strings.TrimSpace(query) == "" {
		return u.patientRepo.GetAll()
	}
	return u.patientRepo.Search(query)
}

// ValidatePatient валидирует данные пациента
func (u *PatientUseCase) ValidatePatient(patient *domain.Patient) error {
	if patient == nil {
		return errors.New("patient cannot be nil")
	}

	if strings.TrimSpace(patient.Name) == "" {
		return errors.New("patient name is required")
	}

	if len(patient.Name) > 100 {
		return errors.New("patient name is too long")
	}

	if patient.IIN != "" && len(patient.IIN) != 12 {
		return errors.New("ИИН должен содержать 12 символов")
	}

	if patient.Phone != "" && len(patient.Phone) > 20 {
		return errors.New("phone number is too long")
	}

	if patient.Email != "" {
		if len(patient.Email) > 100 {
			return errors.New("email is too long")
		}
		if !strings.Contains(patient.Email, "@") {
			return errors.New("invalid email format")
		}
	}

	if patient.Address != "" && len(patient.Address) > 200 {
		return errors.New("address is too long")
	}

	if patient.Notes != "" && len(patient.Notes) > 500 {
		return errors.New("notes are too long")
	}

	return nil
}
