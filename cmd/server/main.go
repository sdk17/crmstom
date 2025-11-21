package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sdk17/crm_ar/internal/infrastructure"
	httphandler "github.com/sdk17/crm_ar/internal/interfaces/http"
	"github.com/sdk17/crm_ar/internal/usecase"
)

func main() {
	// Инициализация репозиториев
	var patientRepo infrastructure.PatientRepository
	var appointmentRepo infrastructure.AppointmentRepository
	var serviceRepo infrastructure.ServiceRepository
	var doctorRepo *infrastructure.PostgresDoctorRepository

	// Проверяем, есть ли переменная окружения для базы данных
	if os.Getenv("DB_HOST") != "" {
		// Используем PostgreSQL
		fmt.Println("🗄️ Подключение к PostgreSQL...")
		config := infrastructure.NewDatabaseConfig()
		db, err := infrastructure.ConnectToDatabase(config)
		if err != nil {
			log.Fatalf("Ошибка подключения к базе данных: %v", err)
		}
		defer db.Close()
		fmt.Println("✅ Подключение к PostgreSQL успешно")

		patientRepo = infrastructure.NewPostgresPatientRepository(db)
		appointmentRepo = infrastructure.NewPostgresAppointmentRepository(db)
		serviceRepo = infrastructure.NewPostgresServiceRepository(db)
		doctorRepo = infrastructure.NewPostgresDoctorRepository(db)
	} else {
		// Используем память (для разработки)
		fmt.Println("💾 Использование in-memory хранилища...")
		patientRepo = infrastructure.NewMemoryPatientRepository()
		appointmentRepo = infrastructure.NewMemoryAppointmentRepository()
		serviceRepo = infrastructure.NewMemoryServiceRepository()
		log.Fatal("⚠️ Doctor repository requires PostgreSQL database. Please set DB_HOST environment variable.")
	}

	// Инициализация use cases
	patientUseCase := usecase.NewPatientUseCase(patientRepo)
	appointmentUseCase := usecase.NewAppointmentUseCase(appointmentRepo, patientRepo, serviceRepo)
	serviceUseCase := usecase.NewServiceUseCase(serviceRepo)
	dashboardUseCase := usecase.NewDashboardUseCase(patientRepo, appointmentRepo, serviceRepo)
	doctorUseCase := usecase.NewDoctorUseCase(doctorRepo)

	// Инициализация HTTP handlers
	handler := httphandler.NewHandler(patientUseCase, appointmentUseCase, serviceUseCase, dashboardUseCase, doctorUseCase)

	// Настройка маршрутов
	mux := http.NewServeMux()

	// API маршруты
	handler.SetupRoutes(mux)

	// Статические файлы
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// HTML страницы
	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/login.html", serveLogin)
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
	// Disable caching for HTML to always fetch latest UI
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/index.html")
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
	// Disable caching for HTML to always fetch latest UI
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/login.html")
}

func servePatients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/patients.html")
}

func serveAppointments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/appointments.html")
}

func servePatientsAppointments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/patients-appointments.html")
}

func serveServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/services.html")
}

func serveReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/reports.html")
}
