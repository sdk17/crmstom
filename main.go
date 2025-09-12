package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Transaction представляет финансовую транзакцию
type Transaction struct {
	ID     int     `json:"id"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
	Type   string  `json:"type"` // "income" или "expense"
}

// Patient представляет пациента
type Patient struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	BirthDate string `json:"birth_date"`
	Address   string `json:"address"`
	Notes     string `json:"notes"`
	LastVisit string `json:"last_visit"`
}

// Appointment представляет запись на приём
type Appointment struct {
	ID          int    `json:"id"`
	PatientID   int    `json:"patient_id"`
	PatientName string `json:"patient_name"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Service     string `json:"service"`
	Doctor      string `json:"doctor"`
	Status      string `json:"status"`
	Cost        int    `json:"cost"`
	Notes       string `json:"notes"`
}

// FinanceReport представляет отчёт по финансам
type FinanceReport struct {
	TotalIncome float64      `json:"total_income"`
	ByDay       []DayIncome  `json:"by_day"`
	ByWeek      []WeekIncome `json:"by_week"`
}

// DayIncome представляет доход по дням
type DayIncome struct {
	Date   string  `json:"date"`
	Income float64 `json:"income"`
}

// WeekIncome представляет доход по неделям
type WeekIncome struct {
	Week   string  `json:"week"`
	Income float64 `json:"income"`
}

// Временное хранилище данных (в реальном проекте это была бы база данных)
var transactions = []Transaction{
	{ID: 1, Amount: 20000, Date: "2025-01-15", Type: "income"},
	{ID: 2, Amount: 15000, Date: "2025-01-15", Type: "income"},
	{ID: 3, Amount: 25000, Date: "2025-01-16", Type: "income"},
	{ID: 4, Amount: 18000, Date: "2025-01-16", Type: "income"},
	{ID: 5, Amount: 30000, Date: "2025-01-17", Type: "income"},
	{ID: 6, Amount: 12000, Date: "2025-01-18", Type: "income"},
	{ID: 7, Amount: 22000, Date: "2025-01-19", Type: "income"},
	{ID: 8, Amount: 19000, Date: "2025-01-20", Type: "income"},
}

var patients = []Patient{
	{ID: 1, Name: "Иванов Иван Иванович", Phone: "+7 (777) 123-45-67", Email: "ivanov@example.com", BirthDate: "1985-03-15", Address: "г. Алматы, ул. Абая, 150", Notes: "Аллергия на пенициллин", LastVisit: "2025-01-15"},
	{ID: 2, Name: "Петрова Анна Сергеевна", Phone: "+7 (777) 234-56-78", Email: "petrova@example.com", BirthDate: "1990-07-22", Address: "г. Алматы, ул. Достык, 45", Notes: "Беременность - 2 триместр", LastVisit: "2025-01-18"},
	{ID: 3, Name: "Сидоров Петр Александрович", Phone: "+7 (777) 345-67-89", Email: "sidorov@example.com", BirthDate: "1978-11-08", Address: "г. Алматы, ул. Сатпаева, 78", Notes: "Диабет 2 типа", LastVisit: "2025-01-20"},
}

var appointments = []Appointment{
	{ID: 1, PatientID: 1, PatientName: "Иванов Иван Иванович", Date: "2025-01-25", Time: "10:00", Service: "consultation", Doctor: "Др. Смит", Status: "scheduled", Cost: 5000, Notes: "Первый визит"},
	{ID: 2, PatientID: 2, PatientName: "Петрова Анна Сергеевна", Date: "2025-01-25", Time: "14:30", Service: "cavity", Doctor: "Др. Джонс", Status: "scheduled", Cost: 15000, Notes: "Беременная пациентка"},
	{ID: 3, PatientID: 3, PatientName: "Сидоров Петр Александрович", Date: "2025-01-26", Time: "09:00", Service: "prosthetics", Doctor: "Др. Уилсон", Status: "completed", Cost: 50000, Notes: "Установка коронки"},
}

// Service представляет услугу
type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Duration    int    `json:"duration"`
	Notes       string `json:"notes"`
}

var services = []Service{
	{ID: 1, Name: "Консультация", Category: "diagnosis", Description: "Первичный осмотр и консультация", Price: 5000, Duration: 30, Notes: "Включает рентген"},
	{ID: 2, Name: "Лечение кариеса", Category: "treatment", Description: "Пломбирование зуба", Price: 15000, Duration: 60, Notes: "Световая пломба"},
	{ID: 3, Name: "Протезирование", Category: "prosthetics", Description: "Изготовление и установка коронки", Price: 50000, Duration: 120, Notes: "Металлокерамика"},
	{ID: 4, Name: "Имплантация", Category: "surgery", Description: "Установка зубного импланта", Price: 100000, Duration: 180, Notes: "Титан"},
	{ID: 5, Name: "Отбеливание", Category: "cosmetic", Description: "Профессиональное отбеливание зубов", Price: 25000, Duration: 90, Notes: "Без вреда для эмали"},
	{ID: 6, Name: "Удаление зуба", Category: "surgery", Description: "Хирургическое удаление зуба", Price: 8000, Duration: 45, Notes: "Простое удаление"},
	{ID: 7, Name: "Лечение каналов", Category: "treatment", Description: "Эндодонтическое лечение", Price: 20000, Duration: 90, Notes: "Многоканальный зуб"},
	{ID: 8, Name: "Чистка зубов", Category: "prevention", Description: "Профессиональная гигиена", Price: 12000, Duration: 60, Notes: "Ультразвук + Air Flow"},
}

func main() {
	// Настройка маршрутов - API endpoints должны быть ПЕРЕД статическими файлами
	http.HandleFunc("/reports/finance", financeReportHandler)
	http.HandleFunc("/api/test", testHandler)
	http.HandleFunc("/api/patients", patientsHandler)
	http.HandleFunc("/api/patients/", patientHandler)
	http.HandleFunc("/api/appointments", appointmentsHandler)
	http.HandleFunc("/api/appointments/", appointmentHandler)
	http.HandleFunc("/api/services", servicesHandler)
	http.HandleFunc("/api/services/", serviceHandler)
	http.HandleFunc("/api/dashboard", dashboardHandler)

	// HTML страницы
	http.HandleFunc("/reports.html", serveReports)
	http.HandleFunc("/patients.html", servePatients)
	http.HandleFunc("/appointments.html", serveAppointments)
	http.HandleFunc("/patients-appointments.html", servePatientsAppointments)
	http.HandleFunc("/services.html", serveServices)

	// Статические файлы - только для папки static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Главная страница - должна быть в самом конце
	http.HandleFunc("/", serveIndex)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// financeReportHandler обрабатывает запросы к /reports/finance
func financeReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовки CORS
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Группируем доходы от записей по дням
	dayIncome := make(map[string]float64)
	weekIncome := make(map[string]float64)
	totalIncome := 0.0

	// Считаем доходы только от завершенных записей
	for _, appointment := range appointments {
		if appointment.Status == "completed" {
			income := float64(appointment.Cost)
			totalIncome += income
			dayIncome[appointment.Date] += income

			// Получаем номер недели
			date, err := time.Parse("2006-01-02", appointment.Date)
			if err == nil {
				year, week := date.ISOWeek()
				weekKey := fmt.Sprintf("%d-W%02d", year, week)
				weekIncome[weekKey] += income
			}
		}
	}

	// Формируем отчёт по дням
	var byDay []DayIncome
	for date, income := range dayIncome {
		byDay = append(byDay, DayIncome{Date: date, Income: income})
	}

	// Формируем отчёт по неделям
	var byWeek []WeekIncome
	for week, income := range weekIncome {
		byWeek = append(byWeek, WeekIncome{Week: week, Income: income})
	}

	report := FinanceReport{
		TotalIncome: totalIncome,
		ByDay:       byDay,
		ByWeek:      byWeek,
	}

	// Отправляем JSON ответ
	json.NewEncoder(w).Encode(report)
}

// serveIndex отдаёт главную страницу
func serveIndex(w http.ResponseWriter, r *http.Request) {
	// Только для точного пути "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

// testHandler для тестирования API
func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("TEST API: %s %s\n", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "API работает!", "path": "` + r.URL.Path + `"}`))
}

// serveReports отдаёт страницу отчётов
func serveReports(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/reports.html")
}

// servePatients отдаёт страницу пациентов
func servePatients(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/patients.html")
}

// serveAppointments отдаёт страницу записей
func serveAppointments(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/appointments.html")
}

// serveServices отдаёт страницу услуг
func serveServices(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/services.html")
}

func servePatientsAppointments(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/patients-appointments.html")
}

// patientsHandler обрабатывает запросы к /api/patients
func patientsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("API пациенты: %s %s\n", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		query := r.URL.Query().Get("query")
		var filteredPatients []Patient

		if query == "" {
			filteredPatients = patients
		} else {
			for _, patient := range patients {
				if contains(patient.Name, query) || contains(patient.Phone, query) {
					filteredPatients = append(filteredPatients, patient)
				}
			}
		}

		json.NewEncoder(w).Encode(filteredPatients)
		return
	}

	if r.Method == http.MethodPost {
		var newPatient Patient
		if err := json.NewDecoder(r.Body).Decode(&newPatient); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		// Проверяем дубликаты по номеру телефона (если номер указан)
		if newPatient.Phone != "" {
			for _, existingPatient := range patients {
				if existingPatient.Phone == newPatient.Phone {
					http.Error(w, "Пациент с таким номером телефона уже существует", http.StatusConflict)
					return
				}
			}
		}

		// Генерируем новый ID
		newID := 1
		if len(patients) > 0 {
			newID = patients[len(patients)-1].ID + 1
		}
		newPatient.ID = newID
		newPatient.LastVisit = time.Now().Format("2006-01-02")

		patients = append(patients, newPatient)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Пациент добавлен",
			"patient": newPatient,
		})
		return
	}

	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// patientHandler обрабатывает запросы к /api/patients/{id}
func patientHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Извлекаем ID из URL
	path := r.URL.Path
	idStr := path[len("/api/patients/"):]

	if idStr == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updatePatient(w, r, id)
	case http.MethodDelete:
		deletePatient(w, r, id)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func updatePatient(w http.ResponseWriter, r *http.Request, id int) {
	var updatedPatient Patient
	if err := json.NewDecoder(r.Body).Decode(&updatedPatient); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Проверяем дубликаты по номеру телефона (если номер указан)
	if updatedPatient.Phone != "" {
		for _, existingPatient := range patients {
			if existingPatient.Phone == updatedPatient.Phone && existingPatient.ID != id {
				http.Error(w, "Пациент с таким номером телефона уже существует", http.StatusConflict)
				return
			}
		}
	}

	for i, patient := range patients {
		if patient.ID == id {
			updatedPatient.ID = id // Сохраняем оригинальный ID
			patients[i] = updatedPatient
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "success",
				"message": "Пациент обновлён",
				"patient": updatedPatient,
			})
			return
		}
	}

	http.Error(w, "Пациент не найден", http.StatusNotFound)
}

func deletePatient(w http.ResponseWriter, r *http.Request, id int) {
	for i, patient := range patients {
		if patient.ID == id {
			patients = append(patients[:i], patients[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Пациент удалён",
			})
			return
		}
	}

	http.Error(w, "Пациент не найден", http.StatusNotFound)
}

// appointmentsHandler обрабатывает запросы к /api/appointments
func appointmentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		doctor := r.URL.Query().Get("doctor")
		status := r.URL.Query().Get("status")
		var filteredAppointments []Appointment

		filteredAppointments = appointments
		if doctor != "" {
			var temp []Appointment
			for _, apt := range filteredAppointments {
				if apt.Doctor == doctor {
					temp = append(temp, apt)
				}
			}
			filteredAppointments = temp
		}
		if status != "" {
			var temp []Appointment
			for _, apt := range filteredAppointments {
				if apt.Status == status {
					temp = append(temp, apt)
				}
			}
			filteredAppointments = temp
		}

		json.NewEncoder(w).Encode(filteredAppointments)
		return
	}

	if r.Method == http.MethodPost {
		var newAppointment Appointment
		if err := json.NewDecoder(r.Body).Decode(&newAppointment); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		// Проверяем конфликт времени
		for _, apt := range appointments {
			if apt.Date == newAppointment.Date && apt.Time == newAppointment.Time {
				http.Error(w, "Это время уже занято", http.StatusConflict)
				return
			}
		}

		// Генерируем новый ID
		newID := 1
		if len(appointments) > 0 {
			newID = appointments[len(appointments)-1].ID + 1
		}
		newAppointment.ID = newID
		newAppointment.Status = "scheduled"

		appointments = append(appointments, newAppointment)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "success",
			"message":     "Запись добавлена",
			"appointment": newAppointment,
		})
		return
	}

	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// appointmentHandler обрабатывает запросы к /api/appointments/{id}
func appointmentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Извлекаем ID из URL
	path := r.URL.Path
	idStr := path[len("/api/appointments/"):]

	if idStr == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updateAppointment(w, r, id)
	case http.MethodDelete:
		deleteAppointment(w, r, id)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func updateAppointment(w http.ResponseWriter, r *http.Request, id int) {
	var updatedAppointment Appointment
	if err := json.NewDecoder(r.Body).Decode(&updatedAppointment); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Проверяем конфликт времени (исключая текущую запись)
	for _, apt := range appointments {
		if apt.ID != id && apt.Date == updatedAppointment.Date && apt.Time == updatedAppointment.Time {
			http.Error(w, "Это время уже занято", http.StatusConflict)
			return
		}
	}

	for i, appointment := range appointments {
		if appointment.ID == id {
			updatedAppointment.ID = id // Сохраняем оригинальный ID
			appointments[i] = updatedAppointment
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":      "success",
				"message":     "Запись обновлена",
				"appointment": updatedAppointment,
			})
			return
		}
	}

	http.Error(w, "Запись не найдена", http.StatusNotFound)
}

func deleteAppointment(w http.ResponseWriter, r *http.Request, id int) {
	for i, appointment := range appointments {
		if appointment.ID == id {
			appointments = append(appointments[:i], appointments[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Запись удалена",
			})
			return
		}
	}

	http.Error(w, "Запись не найдена", http.StatusNotFound)
}

// contains проверяет, содержит ли строка подстроку
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// servicesHandler обрабатывает запросы к /api/services
func servicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(services)
		return
	}

	if r.Method == http.MethodPost {
		var newService Service
		if err := json.NewDecoder(r.Body).Decode(&newService); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		// Генерируем новый ID
		newID := 1
		if len(services) > 0 {
			newID = services[len(services)-1].ID + 1
		}
		newService.ID = newID

		services = append(services, newService)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Услуга добавлена",
			"service": newService,
		})
		return
	}

	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// serviceHandler обрабатывает запросы к /api/services/{id}
func serviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Извлекаем ID из URL
	path := r.URL.Path
	idStr := path[len("/api/services/"):]

	if idStr == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updateService(w, r, id)
	case http.MethodDelete:
		deleteService(w, r, id)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func updateService(w http.ResponseWriter, r *http.Request, id int) {
	var updatedService Service
	if err := json.NewDecoder(r.Body).Decode(&updatedService); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	for i, service := range services {
		if service.ID == id {
			updatedService.ID = id // Сохраняем оригинальный ID
			services[i] = updatedService
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "success",
				"message": "Услуга обновлена",
				"service": updatedService,
			})
			return
		}
	}

	http.Error(w, "Услуга не найдена", http.StatusNotFound)
}

func deleteService(w http.ResponseWriter, r *http.Request, id int) {
	for i, service := range services {
		if service.ID == id {
			services = append(services[:i], services[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "success",
				"message": "Услуга удалена",
			})
			return
		}
	}

	http.Error(w, "Услуга не найдена", http.StatusNotFound)
}

// DashboardStats представляет статистику дашборда
type DashboardStats struct {
	TotalPatients         int     `json:"total_patients"`
	TotalAppointments     int     `json:"total_appointments"`
	CompletedAppointments int     `json:"completed_appointments"`
	TotalRevenue          float64 `json:"total_revenue"`
	TodayAppointments     int     `json:"today_appointments"`
	PendingAppointments   int     `json:"pending_appointments"`
}

// dashboardHandler обрабатывает запросы к /api/dashboard
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем текущую дату
	today := time.Now().Format("2006-01-02")

	// Подсчитываем статистику
	totalPatients := len(patients)
	totalAppointments := len(appointments)
	completedAppointments := 0
	totalRevenue := 0.0
	todayAppointments := 0
	pendingAppointments := 0

	for _, appointment := range appointments {
		if appointment.Status == "completed" {
			completedAppointments++
			totalRevenue += float64(appointment.Cost)
		}
		if appointment.Date == today {
			todayAppointments++
		}
		if appointment.Status == "scheduled" {
			pendingAppointments++
		}
	}

	stats := DashboardStats{
		TotalPatients:         totalPatients,
		TotalAppointments:     totalAppointments,
		CompletedAppointments: completedAppointments,
		TotalRevenue:          totalRevenue,
		TodayAppointments:     todayAppointments,
		PendingAppointments:   pendingAppointments,
	}

	json.NewEncoder(w).Encode(stats)
}
