package gen

//go:generate mockgen -destination=mocks/repository/patient_repository_mock.go -package=repository github.com/sdk17/crmstom/internal/domain PatientRepository
//go:generate mockgen -destination=mocks/repository/appointment_repository_mock.go -package=repository github.com/sdk17/crmstom/internal/domain AppointmentRepository
//go:generate mockgen -destination=mocks/repository/service_repository_mock.go -package=repository github.com/sdk17/crmstom/internal/domain ServiceRepository
//go:generate mockgen -destination=mocks/repository/doctor_repository_mock.go -package=repository github.com/sdk17/crmstom/internal/domain DoctorRepository
