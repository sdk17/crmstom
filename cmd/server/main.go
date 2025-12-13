package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	httphandler "github.com/sdk17/crmstom/internal/interfaces/http"
	"github.com/sdk17/crmstom/internal/repository"
	"github.com/sdk17/crmstom/internal/usecase"
)

func main() {
	// Подключение к PostgreSQL
	if os.Getenv("DB_HOST") == "" {
		log.Fatal("DB_HOST environment variable is required. Please configure PostgreSQL connection.")
	}

	fmt.Println("Подключение к PostgreSQL...")
	config := repository.NewDatabaseConfig()
	db, err := repository.ConnectToDatabase(config)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()
	fmt.Println("Подключение к PostgreSQL успешно")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	// Инициализация репозиториев
	patientRepo := repository.NewPatientRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)

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
	mux.HandleFunc("/doctors.html", serveDoctors)
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

func serveDoctors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/doctors.html")
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

func runMigrations(db *sql.DB) error {
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	fmt.Printf("Применение миграций из %s...\n", migrationsPath)
	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	fmt.Println("Миграции применены успешно")

	return nil
}
