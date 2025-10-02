package domain

import "time"

// Doctor представляет врача в системе
type Doctor struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Login     string    `json:"login"`
	Password  string    `json:"password,omitempty"` // omitempty для безопасности при отправке на фронт
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DoctorRepository определяет методы для работы с врачами
type DoctorRepository interface {
	Create(doctor *Doctor) error
	GetByID(id int) (*Doctor, error)
	GetAll() ([]*Doctor, error)
	Update(doctor *Doctor) error
	Delete(id int) error
	GetByLogin(login string) (*Doctor, error)
}
