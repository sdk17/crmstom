package usecase

import (
	"errors"
	"strings"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
)

// ServiceUseCase реализует бизнес-логику для услуг
type ServiceUseCase struct {
	serviceRepo domain.ServiceRepository
}

// NewServiceUseCase создает новый экземпляр ServiceUseCase
func NewServiceUseCase(serviceRepo domain.ServiceRepository) *ServiceUseCase {
	return &ServiceUseCase{
		serviceRepo: serviceRepo,
	}
}

// GetService получает услугу по ID
func (u *ServiceUseCase) GetService(id int) (*domain.Service, error) {
	if id <= 0 {
		return nil, errors.New("invalid service ID")
	}
	return u.serviceRepo.GetByID(id)
}

// GetAllServices получает все услуги
func (u *ServiceUseCase) GetAllServices() ([]*domain.Service, error) {
	return u.serviceRepo.GetAll()
}

// CreateService создает новую услугу
func (u *ServiceUseCase) CreateService(service *domain.Service) error {
	if err := u.ValidateService(service); err != nil {
		return err
	}

	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	return u.serviceRepo.Create(service)
}

// UpdateService обновляет услугу
func (u *ServiceUseCase) UpdateService(service *domain.Service) error {
	if err := u.ValidateService(service); err != nil {
		return err
	}

	service.UpdatedAt = time.Now()

	return u.serviceRepo.Update(service)
}

// DeleteService удаляет услугу
func (u *ServiceUseCase) DeleteService(id int) error {
	if id <= 0 {
		return errors.New("invalid service ID")
	}
	return u.serviceRepo.Delete(id)
}

// GetServicesByCategory получает услуги по категории
func (u *ServiceUseCase) GetServicesByCategory(category string) ([]*domain.Service, error) {
	if strings.TrimSpace(category) == "" {
		return u.serviceRepo.GetAll()
	}
	return u.serviceRepo.GetByCategory(category)
}

// SearchServices ищет услуги по запросу
func (u *ServiceUseCase) SearchServices(query string) ([]*domain.Service, error) {
	if strings.TrimSpace(query) == "" {
		return u.serviceRepo.GetAll()
	}
	return u.serviceRepo.Search(query)
}

// ValidateService валидирует данные услуги
func (u *ServiceUseCase) ValidateService(service *domain.Service) error {
	if service == nil {
		return errors.New("service cannot be nil")
	}

	if strings.TrimSpace(service.Name) == "" {
		return errors.New("service name is required")
	}

	if len(service.Name) > 100 {
		return errors.New("service name is too long")
	}

	if strings.TrimSpace(service.Category) == "" {
		return errors.New("service category is required")
	}

	if len(service.Category) > 50 {
		return errors.New("service category is too long")
	}

	if service.Price < 0 {
		return errors.New("service price cannot be negative")
	}

	if service.Duration < 0 {
		return errors.New("service duration cannot be negative")
	}

	if service.Duration > 480 { // 8 часов
		return errors.New("service duration is too long")
	}

	if service.Description != "" && len(service.Description) > 500 {
		return errors.New("service description is too long")
	}

	if service.Notes != "" && len(service.Notes) > 500 {
		return errors.New("service notes are too long")
	}

	return nil
}
