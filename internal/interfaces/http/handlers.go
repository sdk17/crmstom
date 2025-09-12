package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/sdk17/crmstom/internal/domain"
	"github.com/sdk17/crmstom/internal/usecase"
)

// Handler содержит все HTTP обработчики
type Handler struct {
	patientUseCase     *usecase.PatientUseCase
	appointmentUseCase *usecase.AppointmentUseCase
	serviceUseCase     *usecase.ServiceUseCase
	dashboardUseCase   *usecase.DashboardUseCase
}

// NewHandler создает новый экземпляр Handler
func NewHandler(
	patientUseCase *usecase.PatientUseCase,
	appointmentUseCase *usecase.AppointmentUseCase,
	serviceUseCase *usecase.ServiceUseCase,
	dashboardUseCase *usecase.DashboardUseCase,
) *Handler {
	return &Handler{
		patientUseCase:     patientUseCase,
		appointmentUseCase: appointmentUseCase,
		serviceUseCase:     serviceUseCase,
		dashboardUseCase:   dashboardUseCase,
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
	var patient domain.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.patientUseCase.CreatePatient(&patient); err != nil {
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

// SetupRoutes настраивает маршруты
func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	// API маршруты
	mux.HandleFunc("/api/patients", h.PatientsHandler)
	mux.HandleFunc("/api/patients/", h.PatientHandler)
	// TODO: Добавить остальные маршруты
}
