package infrastructure

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/sdk17/crm_ar/internal/domain"
)

// MemoryPatientRepository реализует PatientRepository в памяти
type MemoryPatientRepository struct {
	patients []*domain.Patient
	nextID   int
	mutex    sync.RWMutex
}

// NewMemoryPatientRepository создает новый экземпляр MemoryPatientRepository
func NewMemoryPatientRepository() *MemoryPatientRepository {
	repo := &MemoryPatientRepository{
		patients: make([]*domain.Patient, 0),
		nextID:   1,
	}

	// Добавляем тестовые данные
	repo.seedData()
	return repo
}

// seedData добавляет тестовые данные
func (r *MemoryPatientRepository) seedData() {
	now := time.Now()
	testPatients := []*domain.Patient{
		{
			ID:        1,
			Name:      "Иванов Иван Иванович",
			Phone:     "+7 (777) 123-45-67",
			Email:     "ivanov@example.com",
			BirthDate: time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC),
			Address:   "г. Алматы, ул. Абая, 150",
			Notes:     "Аллергия на пенициллин",
			LastVisit: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Name:      "Петрова Анна Сергеевна",
			Phone:     "+7 (777) 234-56-78",
			Email:     "petrova@example.com",
			BirthDate: time.Date(1990, 7, 22, 0, 0, 0, 0, time.UTC),
			Address:   "г. Алматы, ул. Достык, 45",
			Notes:     "Беременность - 2 триместр",
			LastVisit: time.Date(2025, 1, 18, 0, 0, 0, 0, time.UTC),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        3,
			Name:      "Сидоров Петр Александрович",
			Phone:     "+7 (777) 345-67-89",
			Email:     "sidorov@example.com",
			BirthDate: time.Date(1978, 11, 8, 0, 0, 0, 0, time.UTC),
			Address:   "г. Алматы, ул. Сатпаева, 78",
			Notes:     "Диабет 2 типа",
			LastVisit: time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	r.patients = append(r.patients, testPatients...)
	r.nextID = 4
}

// GetByID получает пациента по ID
func (r *MemoryPatientRepository) GetByID(id int) (*domain.Patient, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, patient := range r.patients {
		if patient.ID == id {
			return patient, nil
		}
	}
	return nil, errors.New("patient not found")
}

// GetAll получает всех пациентов
func (r *MemoryPatientRepository) GetAll() ([]*domain.Patient, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*domain.Patient, len(r.patients))
	copy(result, r.patients)
	return result, nil
}

// Create создает нового пациента
func (r *MemoryPatientRepository) Create(patient *domain.Patient) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	patient.ID = r.nextID
	r.nextID++
	r.patients = append(r.patients, patient)
	return nil
}

// Update обновляет пациента
func (r *MemoryPatientRepository) Update(patient *domain.Patient) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, p := range r.patients {
		if p.ID == patient.ID {
			r.patients[i] = patient
			return nil
		}
	}
	return errors.New("patient not found")
}

// Delete удаляет пациента
func (r *MemoryPatientRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, patient := range r.patients {
		if patient.ID == id {
			r.patients = append(r.patients[:i], r.patients[i+1:]...)
			return nil
		}
	}
	return errors.New("patient not found")
}

// Search ищет пациентов по запросу
func (r *MemoryPatientRepository) Search(query string) ([]*domain.Patient, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	query = strings.ToLower(query)
	var result []*domain.Patient

	for _, patient := range r.patients {
		if strings.Contains(strings.ToLower(patient.Name), query) ||
			strings.Contains(patient.Phone, query) ||
			strings.Contains(strings.ToLower(patient.Email), query) {
			result = append(result, patient)
		}
	}

	return result, nil
}

// GetByPhone получает пациента по номеру телефона
func (r *MemoryPatientRepository) GetByPhone(phone string) (*domain.Patient, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, patient := range r.patients {
		if patient.Phone == phone {
			return patient, nil
		}
	}
	return nil, errors.New("patient not found")
}
