package infrastructure

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
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
			ID:          1,
			Name:        "Консультация",
			Category:    "diagnosis",
			Description: "Первичный осмотр и консультация",
			Price:       5000,
			Duration:    30,
			Notes:       "Включает рентген",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          2,
			Name:        "Лечение кариеса",
			Category:    "treatment",
			Description: "Пломбирование зуба",
			Price:       15000,
			Duration:    60,
			Notes:       "Световая пломба",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          3,
			Name:        "Протезирование",
			Category:    "prosthetics",
			Description: "Изготовление и установка коронки",
			Price:       50000,
			Duration:    120,
			Notes:       "Металлокерамика",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          4,
			Name:        "Имплантация",
			Category:    "surgery",
			Description: "Установка зубного импланта",
			Price:       100000,
			Duration:    180,
			Notes:       "Титан",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          5,
			Name:        "Отбеливание",
			Category:    "cosmetic",
			Description: "Профессиональное отбеливание зубов",
			Price:       25000,
			Duration:    90,
			Notes:       "Без вреда для эмали",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          6,
			Name:        "Удаление зуба",
			Category:    "surgery",
			Description: "Хирургическое удаление зуба",
			Price:       8000,
			Duration:    45,
			Notes:       "Простое удаление",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          7,
			Name:        "Лечение каналов",
			Category:    "treatment",
			Description: "Эндодонтическое лечение",
			Price:       20000,
			Duration:    90,
			Notes:       "Многоканальный зуб",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          8,
			Name:        "Чистка зубов",
			Category:    "prevention",
			Description: "Профессиональная гигиена",
			Price:       12000,
			Duration:    60,
			Notes:       "Ультразвук + Air Flow",
			CreatedAt:   now,
			UpdatedAt:   now,
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
		if service.Category == category {
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
			strings.Contains(strings.ToLower(service.Category), query) ||
			strings.Contains(strings.ToLower(service.Description), query) {
			result = append(result, service)
		}
	}

	return result, nil
}
