package infrastructure

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/sdk17/crm_ar/internal/domain"
)

// MemoryServiceRepository реализует ServiceRepository в памяти
type MemoryServiceRepository struct {
	services []*domain.Service
	nextID   int
	mutex    sync.RWMutex
}

// NewMemoryServiceRepository создает новый экземпляр MemoryServiceRepository
func NewMemoryServiceRepository() *MemoryServiceRepository {
	repo := &MemoryServiceRepository{
		services: make([]*domain.Service, 0),
		nextID:   1,
	}

	// Добавляем тестовые данные
	repo.seedData()
	return repo
}

// seedData добавляет тестовые данные
func (r *MemoryServiceRepository) seedData() {
	now := time.Now()
	testServices := []*domain.Service{
		{
			ID:        1,
			Name:      "Консультация",
			Type:      "Консультация",
			Notes:     "Первичный осмотр и консультация",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Name:      "Лечение кариеса",
			Type:      "Лечение кариеса",
			Notes:     "Пломбирование зуба",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        3,
			Name:      "Протезирование",
			Type:      "Протезирование",
			Notes:     "Изготовление и установка коронки",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        4,
			Name:      "Имплантация",
			Type:      "Имплантация",
			Notes:     "Установка зубного импланта",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        5,
			Name:      "Отбеливание",
			Type:      "Гигиена",
			Notes:     "Профессиональное отбеливание зубов",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        6,
			Name:      "Удаление зуба",
			Type:      "Удаление зубов",
			Notes:     "Хирургическое удаление зуба",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        7,
			Name:      "Лечение каналов",
			Type:      "Лечение пульпита",
			Notes:     "Эндодонтическое лечение",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        8,
			Name:      "Чистка зубов",
			Type:      "Гигиена",
			Notes:     "Профессиональная гигиена",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	r.services = append(r.services, testServices...)
	r.nextID = 9
}

// GetByID получает услугу по ID
func (r *MemoryServiceRepository) GetByID(id int) (*domain.Service, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, service := range r.services {
		if service.ID == id {
			return service, nil
		}
	}
	return nil, errors.New("service not found")
}

// GetAll получает все услуги
func (r *MemoryServiceRepository) GetAll() ([]*domain.Service, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*domain.Service, len(r.services))
	copy(result, r.services)
	return result, nil
}

// Create создает новую услугу
func (r *MemoryServiceRepository) Create(service *domain.Service) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	service.ID = r.nextID
	r.nextID++
	r.services = append(r.services, service)
	return nil
}

// Update обновляет услугу
func (r *MemoryServiceRepository) Update(service *domain.Service) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, s := range r.services {
		if s.ID == service.ID {
			r.services[i] = service
			return nil
		}
	}
	return errors.New("service not found")
}

// Delete удаляет услугу
func (r *MemoryServiceRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, service := range r.services {
		if service.ID == id {
			r.services = append(r.services[:i], r.services[i+1:]...)
			return nil
		}
	}
	return errors.New("service not found")
}

// GetByCategory получает услуги по категории
func (r *MemoryServiceRepository) GetByCategory(category string) ([]*domain.Service, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*domain.Service
	for _, service := range r.services {
		if service.Type == category {
			result = append(result, service)
		}
	}
	return result, nil
}

// Search ищет услуги по запросу
func (r *MemoryServiceRepository) Search(query string) ([]*domain.Service, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	query = strings.ToLower(query)
	var result []*domain.Service

	for _, service := range r.services {
		if strings.Contains(strings.ToLower(service.Name), query) ||
			strings.Contains(strings.ToLower(service.Type), query) ||
			strings.Contains(strings.ToLower(service.Notes), query) {
			result = append(result, service)
		}
	}

	return result, nil
}
