package infrastructure

import (
	"github.com/sdk17/crmstom/internal/domain"
)

// Интерфейсы репозиториев - используем доменные интерфейсы

type PatientRepository = domain.PatientRepository
type ServiceRepository = domain.ServiceRepository
type AppointmentRepository = domain.AppointmentRepository
