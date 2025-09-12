package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sdk17/crmstom/internal/infrastructure"
	httphandler "github.com/sdk17/crmstom/internal/interfaces/http"
	"github.com/sdk17/crmstom/internal/usecase"
)

func main() {
	// Инициализация репозиториев
	patientRepo := infrastructure.NewMemoryPatientRepository()
	appointmentRepo := infrastructure.NewMemoryAppointmentRepository()
	serviceRepo := infrastructure.NewMemoryServiceRepository()

	// Инициализация use cases
	patientUseCase := usecase.NewPatientUseCase(patientRepo)
	appointmentUseCase := usecase.NewAppointmentUseCase(appointmentRepo, patientRepo, serviceRepo)
	serviceUseCase := usecase.NewServiceUseCase(serviceRepo)
	dashboardUseCase := usecase.NewDashboardUseCase(patientRepo, appointmentRepo)

	// Инициализация HTTP handlers
	handler := httphandler.NewHandler(patientUseCase, appointmentUseCase, serviceUseCase, dashboardUseCase)

	// Настройка маршрутов
	mux := http.NewServeMux()

	// API маршруты
	handler.SetupRoutes(mux)

	// Статические файлы
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static/"))))

	// HTML страницы
	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/patients.html", servePatients)
	mux.HandleFunc("/appointments.html", serveAppointments)
	mux.HandleFunc("/patients-appointments.html", servePatientsAppointments)
	mux.HandleFunc("/services.html", serveServices)
	mux.HandleFunc("/reports.html", serveReports)

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	fmt.Println("📊 Clean Architecture + SOLID принципы")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Обработчики для статических файлов
func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "../../static/index.html")
}

func servePatients(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../static/patients.html")
}

func serveAppointments(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../static/appointments.html")
}

func servePatientsAppointments(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../static/patients-appointments.html")
}

func serveServices(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../static/services.html")
}

func serveReports(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../static/reports.html")
}
