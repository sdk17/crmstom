package usecase

import (
	"errors"

	"github.com/sdk17/crmstom/internal/domain"
)

type DoctorUseCase struct {
	doctorRepo domain.DoctorRepository
}

func NewDoctorUseCase(doctorRepo domain.DoctorRepository) *DoctorUseCase {
	return &DoctorUseCase{
		doctorRepo: doctorRepo,
	}
}

// CreateDoctor создает нового врача
func (u *DoctorUseCase) CreateDoctor(doctor *domain.Doctor) error {
	if err := u.ValidateDoctor(doctor); err != nil {
		return err
	}

	return u.doctorRepo.Create(doctor)
}

// GetDoctor получает врача по ID
func (u *DoctorUseCase) GetDoctor(id int) (*domain.Doctor, error) {
	if id <= 0 {
		return nil, errors.New("invalid doctor ID")
	}
	return u.doctorRepo.GetByID(id)
}

// GetAllDoctors получает всех врачей
func (u *DoctorUseCase) GetAllDoctors() ([]*domain.Doctor, error) {
	return u.doctorRepo.GetAll()
}

// UpdateDoctor обновляет врача
func (u *DoctorUseCase) UpdateDoctor(doctor *domain.Doctor) error {
	if err := u.ValidateDoctor(doctor); err != nil {
		return err
	}

	return u.doctorRepo.Update(doctor)
}

// DeleteDoctor удаляет врача
func (u *DoctorUseCase) DeleteDoctor(id int) error {
	if id <= 0 {
		return errors.New("invalid doctor ID")
	}
	return u.doctorRepo.Delete(id)
}

// AuthenticateDoctor аутентифицирует врача
func (u *DoctorUseCase) AuthenticateDoctor(login, password string) (*domain.Doctor, error) {
	if login == "" || password == "" {
		return nil, errors.New("login and password are required")
	}

	doctor, err := u.doctorRepo.GetByLogin(login)
	if err != nil {
		return nil, err
	}

	if doctor == nil {
		return nil, errors.New("invalid login or password")
	}

	// Простая проверка пароля (в реальном приложении используйте bcrypt)
	if doctor.Password != password {
		return nil, errors.New("invalid login or password")
	}

	// Не возвращаем пароль на фронт
	doctor.Password = ""

	return doctor, nil
}

// ValidateDoctor валидирует данные врача
func (u *DoctorUseCase) ValidateDoctor(doctor *domain.Doctor) error {
	if doctor.Name == "" {
		return errors.New("doctor name is required")
	}

	if len(doctor.Name) > 255 {
		return errors.New("doctor name is too long")
	}

	if doctor.Login == "" {
		return errors.New("doctor login is required")
	}

	if len(doctor.Login) > 100 {
		return errors.New("doctor login is too long")
	}

	if doctor.Password == "" {
		return errors.New("doctor password is required")
	}

	if len(doctor.Password) < 4 {
		return errors.New("doctor password is too short")
	}

	return nil
}
