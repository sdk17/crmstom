package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
	"github.com/sdk17/crmstom/internal/usecase"
)

// Handler содержит все HTTP обработчики
type Handler struct {
	patientUseCase     *usecase.PatientUseCase
	appointmentUseCase *usecase.AppointmentUseCase
	serviceUseCase     *usecase.ServiceUseCase
	dashboardUseCase   *usecase.DashboardUseCase
	doctorUseCase      *usecase.DoctorUseCase
}

// NewHandler создает новый экземпляр Handler
func NewHandler(
	patientUseCase *usecase.PatientUseCase,
	appointmentUseCase *usecase.AppointmentUseCase,
	serviceUseCase *usecase.ServiceUseCase,
	dashboardUseCase *usecase.DashboardUseCase,
	doctorUseCase *usecase.DoctorUseCase,
) *Handler {
	return &Handler{
		patientUseCase:     patientUseCase,
		appointmentUseCase: appointmentUseCase,
		serviceUseCase:     serviceUseCase,
		dashboardUseCase:   dashboardUseCase,
		doctorUseCase:      doctorUseCase,
	}
}

// setCORSHeaders устанавливает CORS заголовки
func (h *Handler) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// writeJSONResponse записывает JSON ответ
func (h *Handler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse записывает JSON ответ с ошибкой
func (h *Handler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSONResponse(w, statusCode, map[string]string{"error": message})
}

// writeSuccessResponse записывает JSON ответ с успехом
func (h *Handler) writeSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := map[string]interface{}{
		"status":  "success",
		"message": message,
	}
	if data != nil {
		response["data"] = data
	}
	h.writeJSONResponse(w, http.StatusOK, response)
}

// PatientsHandler обрабатывает запросы к пациентам
func (h *Handler) PatientsHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetPatients(w, r)
	case http.MethodPost:
		h.handleCreatePatient(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleGetPatients обрабатывает GET запросы для пациентов
func (h *Handler) handleGetPatients(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	var patients []*domain.Patient
	var err error

	if query != "" {
		patients, err = h.patientUseCase.SearchPatients(query)
	} else {
		patients, err = h.patientUseCase.GetAllPatients()
	}

	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusOK, patients)
}

// handleCreatePatient обрабатывает POST запросы для создания пациента
func (h *Handler) handleCreatePatient(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		BirthDate string `json:"birth_date"`
		Address   string `json:"address"`
		Notes     string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Создаем пациента
	patient := &domain.Patient{
		Name:    request.Name,
		Phone:   request.Phone,
		Email:   request.Email,
		Address: request.Address,
		Notes:   request.Notes,
	}

	// Парсим дату рождения, если она указана
	if request.BirthDate != "" {
		if birthDate, err := time.Parse("2006-01-02", request.BirthDate); err == nil {
			patient.BirthDate = birthDate
		}
	}

	if err := h.patientUseCase.CreatePatient(patient); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "уже существует") {
			statusCode = http.StatusConflict
		}
		h.writeErrorResponse(w, statusCode, err.Error())
		return
	}

	h.writeSuccessResponse(w, "Patient created successfully", patient)
}

// PatientHandler обрабатывает запросы к конкретному пациенту
func (h *Handler) PatientHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Извлекаем ID из URL
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/patients/")

	if idStr == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Patient ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid patient ID")
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.handleUpdatePatient(w, r, id)
	case http.MethodDelete:
		h.handleDeletePatient(w, r, id)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleUpdatePatient обрабатывает PUT запросы для обновления пациента
func (h *Handler) handleUpdatePatient(w http.ResponseWriter, r *http.Request, id int) {
	var patient domain.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	patient.ID = id

	if err := h.patientUseCase.UpdatePatient(&patient); err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "уже существует") {
			statusCode = http.StatusConflict
		}
		h.writeErrorResponse(w, statusCode, err.Error())
		return
	}

	h.writeSuccessResponse(w, "Patient updated successfully", patient)
}

// handleDeletePatient обрабатывает DELETE запросы для удаления пациента
func (h *Handler) handleDeletePatient(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.patientUseCase.DeletePatient(id); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeSuccessResponse(w, "Patient deleted successfully", nil)
}

// ServicesHandler обрабатывает запросы к /api/services
func (h *Handler) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetServices(w, r)
	case http.MethodPost:
		h.handleCreateService(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleGetServices получает список услуг
func (h *Handler) handleGetServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.serviceUseCase.GetAllServices()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get services")
		return
	}

	h.writeSuccessResponse(w, "Services retrieved successfully", services)
}

// handleCreateService создает новую услугу
func (h *Handler) handleCreateService(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Notes string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Создаем услугу
	service := &domain.Service{
		Name:  request.Name,
		Type:  request.Type,
		Notes: request.Notes,
	}

	if err := h.serviceUseCase.CreateService(service); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create service")
		return
	}

	h.writeSuccessResponse(w, "Service created successfully", service)
}

// ServiceHandler обрабатывает запросы к /api/services/{id}
func (h *Handler) ServiceHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/services/")
	id, err := strconv.Atoi(path)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid service ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetService(w, r, id)
	case http.MethodPut:
		h.handleUpdateService(w, r, id)
	case http.MethodDelete:
		h.handleDeleteService(w, r, id)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleGetService получает услугу по ID
func (h *Handler) handleGetService(w http.ResponseWriter, r *http.Request, id int) {
	service, err := h.serviceUseCase.GetService(id)
	if err != nil {
		h.writeErrorResponse(w, http.StatusNotFound, "Service not found")
		return
	}

	h.writeSuccessResponse(w, "Service retrieved successfully", service)
}

// handleUpdateService обновляет услугу
func (h *Handler) handleUpdateService(w http.ResponseWriter, r *http.Request, id int) {
	var service domain.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	service.ID = id
	if err := h.serviceUseCase.UpdateService(&service); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to update service")
		return
	}

	h.writeSuccessResponse(w, "Service updated successfully", service)
}

// handleDeleteService удаляет услугу
func (h *Handler) handleDeleteService(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.serviceUseCase.DeleteService(id); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete service")
		return
	}

	h.writeSuccessResponse(w, "Service deleted successfully", nil)
}

// AppointmentsHandler обрабатывает запросы к /api/appointments
func (h *Handler) AppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetAppointments(w, r)
	case http.MethodPost:
		h.handleCreateAppointment(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleGetAppointments получает список записей
func (h *Handler) handleGetAppointments(w http.ResponseWriter, r *http.Request) {
	appointments, err := h.appointmentUseCase.GetAllAppointments()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get appointments")
		return
	}

	h.writeSuccessResponse(w, "Appointments retrieved successfully", appointments)
}

// handleCreateAppointment создает новую запись
func (h *Handler) handleCreateAppointment(w http.ResponseWriter, r *http.Request) {
	var request struct {
		PatientID int     `json:"patient_id"`
		Service   string  `json:"service"`
		Date      string  `json:"date"`
		Time      string  `json:"time"`
		Doctor    string  `json:"doctor"`
		Status    string  `json:"status"`
		Price     float64 `json:"price"`
		Duration  int     `json:"duration"`
		Notes     string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Создаем запись
	appointment := &domain.Appointment{
		PatientID: request.PatientID,
		Service:   request.Service,
		Time:      request.Time,
		Doctor:    request.Doctor,
		Price:     request.Price,
		Duration:  request.Duration,
		Notes:     request.Notes,
	}

	// Парсим дату, если она указана
	if request.Date != "" {
		if date, err := time.Parse("2006-01-02T15:04:05Z", request.Date); err == nil {
			appointment.Date = date
		} else if date, err := time.Parse("2006-01-02", request.Date); err == nil {
			appointment.Date = date
		}
	}

	// Устанавливаем статус
	if request.Status != "" {
		appointment.Status = domain.AppointmentStatus(request.Status)
	} else {
		appointment.Status = domain.StatusScheduled
	}

	if err := h.appointmentUseCase.CreateAppointment(appointment); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create appointment")
		return
	}

	h.writeSuccessResponse(w, "Appointment created successfully", appointment)
}

// AppointmentHandler обрабатывает запросы к /api/appointments/{id}
func (h *Handler) AppointmentHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/appointments/")
	id, err := strconv.Atoi(path)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetAppointment(w, r, id)
	case http.MethodPut:
		h.handleUpdateAppointment(w, r, id)
	case http.MethodDelete:
		h.handleDeleteAppointment(w, r, id)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleGetAppointment получает запись по ID
func (h *Handler) handleGetAppointment(w http.ResponseWriter, r *http.Request, id int) {
	appointment, err := h.appointmentUseCase.GetAppointment(id)
	if err != nil {
		h.writeErrorResponse(w, http.StatusNotFound, "Appointment not found")
		return
	}

	h.writeSuccessResponse(w, "Appointment retrieved successfully", appointment)
}

// handleUpdateAppointment обновляет запись
func (h *Handler) handleUpdateAppointment(w http.ResponseWriter, r *http.Request, id int) {
	var appointment domain.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	appointment.ID = id
	if err := h.appointmentUseCase.UpdateAppointment(&appointment); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to update appointment")
		return
	}

	h.writeSuccessResponse(w, "Appointment updated successfully", appointment)
}

// handleDeleteAppointment удаляет запись
func (h *Handler) handleDeleteAppointment(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.appointmentUseCase.DeleteAppointment(id); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete appointment")
		return
	}

	h.writeSuccessResponse(w, "Appointment deleted successfully", nil)
}

// DashboardHandler обрабатывает запросы к /api/dashboard
func (h *Handler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	stats, err := h.dashboardUseCase.GetDashboardStats()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get dashboard stats")
		return
	}

	h.writeSuccessResponse(w, "Dashboard stats retrieved successfully", stats)
}

// ReportsHandler обрабатывает запросы к /api/reports
func (h *Handler) ReportsHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	report, err := h.dashboardUseCase.GetFinanceReport()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get finance report")
		return
	}

	h.writeSuccessResponse(w, "Finance report retrieved successfully", report)
}

// DoctorsHandler обрабатывает запросы к /api/doctors
func (h *Handler) DoctorsHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetDoctors(w, r)
	case http.MethodPost:
		h.handleCreateDoctor(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// DoctorHandler обрабатывает запросы к /api/doctors/{id}
func (h *Handler) DoctorHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Извлекаем ID из URL
	path := strings.TrimPrefix(r.URL.Path, "/api/doctors/")
	id, err := strconv.Atoi(path)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetDoctor(w, r, id)
	case http.MethodPut:
		h.handleUpdateDoctor(w, r, id)
	case http.MethodDelete:
		h.handleDeleteDoctor(w, r, id)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleGetDoctors получает всех врачей
func (h *Handler) handleGetDoctors(w http.ResponseWriter, r *http.Request) {
	doctors, err := h.doctorUseCase.GetAllDoctors()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get doctors")
		return
	}

	// Не отправляем пароли на фронт
	for _, doctor := range doctors {
		doctor.Password = ""
	}

	h.writeSuccessResponse(w, "Doctors retrieved successfully", doctors)
}

// handleGetDoctor получает врача по ID
func (h *Handler) handleGetDoctor(w http.ResponseWriter, r *http.Request, id int) {
	doctor, err := h.doctorUseCase.GetDoctor(id)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to get doctor")
		return
	}

	if doctor == nil {
		h.writeErrorResponse(w, http.StatusNotFound, "Doctor not found")
		return
	}

	// Не отправляем пароль на фронт
	doctor.Password = ""

	h.writeSuccessResponse(w, "Doctor retrieved successfully", doctor)
}

// handleCreateDoctor создает нового врача
func (h *Handler) handleCreateDoctor(w http.ResponseWriter, r *http.Request) {
	var doctor domain.Doctor
	if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.doctorUseCase.CreateDoctor(&doctor); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Не отправляем пароль на фронт
	doctor.Password = ""

	h.writeSuccessResponse(w, "Doctor created successfully", doctor)
}

// handleUpdateDoctor обновляет врача
func (h *Handler) handleUpdateDoctor(w http.ResponseWriter, r *http.Request, id int) {
	var doctor domain.Doctor
	if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	doctor.ID = id

	if err := h.doctorUseCase.UpdateDoctor(&doctor); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Не отправляем пароль на фронт
	doctor.Password = ""

	h.writeSuccessResponse(w, "Doctor updated successfully", doctor)
}

// handleDeleteDoctor удаляет врача
func (h *Handler) handleDeleteDoctor(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.doctorUseCase.DeleteDoctor(id); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete doctor")
		return
	}

	h.writeSuccessResponse(w, "Doctor deleted successfully", nil)
}

// AuthHandler обрабатывает запросы к /api/auth
func (h *Handler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var authRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	doctor, err := h.doctorUseCase.AuthenticateDoctor(authRequest.Login, authRequest.Password)
	if err != nil {
		h.writeErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.writeSuccessResponse(w, "Authentication successful", doctor)
}

// SetupRoutes настраивает маршруты
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	// API маршруты для пациентов
	mux.HandleFunc("/api/patients", h.PatientsHandler)
	mux.HandleFunc("/api/patients/", h.PatientHandler)

	// API маршруты для услуг
	mux.HandleFunc("/api/services", h.ServicesHandler)
	mux.HandleFunc("/api/services/", h.ServiceHandler)

	// API маршруты для записей
	mux.HandleFunc("/api/appointments", h.AppointmentsHandler)
	mux.HandleFunc("/api/appointments/", h.AppointmentHandler)

	// API маршруты для дашборда
	mux.HandleFunc("/api/dashboard", h.DashboardHandler)

	// API маршрут для финансовых отчетов
	mux.HandleFunc("/api/reports", h.ReportsHandler)

	// API маршруты для врачей
	mux.HandleFunc("/api/doctors", h.DoctorsHandler)
	mux.HandleFunc("/api/doctors/", h.DoctorHandler)

	// API маршрут для авторизации
	mux.HandleFunc("/api/auth", h.AuthHandler)
}
