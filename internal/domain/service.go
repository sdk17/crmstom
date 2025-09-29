package domain

import "time"

// Service представляет услугу в доменной модели (упрощенная схема)
type Service struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServiceRepository определяет интерфейс для работы с услугами
type ServiceRepository interface {
	GetByID(id int) (*Service, error)
	GetAll() ([]*Service, error)
	Create(service *Service) error
	Update(service *Service) error
	Delete(id int) error
	GetByCategory(category string) ([]*Service, error)
	Search(query string) ([]*Service, error)
}

// ServiceService определяет бизнес-логику для работы с услугами
type ServiceService interface {
	GetService(id int) (*Service, error)
	GetAllServices() ([]*Service, error)
	CreateService(service *Service) error
	UpdateService(service *Service) error
	DeleteService(id int) error
	GetServicesByCategory(category string) ([]*Service, error)
	SearchServices(query string) ([]*Service, error)
	ValidateService(service *Service) error
}
