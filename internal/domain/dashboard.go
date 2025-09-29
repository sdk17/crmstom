package domain

// DashboardStats представляет статистику дашборда
type DashboardStats struct {
	TodayAppointments int     `json:"today_appointments"`
	TodayRevenue      float64 `json:"today_revenue"`
	TotalPatients     int     `json:"total_patients"`
}

// FinanceReport представляет финансовый отчет
type FinanceReport struct {
	TotalIncome float64      `json:"total_income"`
	ByDay       []DayIncome  `json:"by_day"`
	ByWeek      []WeekIncome `json:"by_week"`
}

// DayIncome представляет доход за день
type DayIncome struct {
	Date   string  `json:"date"`
	Income float64 `json:"income"`
}

// WeekIncome представляет доход за неделю
type WeekIncome struct {
	Week   string  `json:"week"`
	Income float64 `json:"income"`
}

// DashboardService определяет бизнес-логику для дашборда
type DashboardService interface {
	GetDashboardStats() (*DashboardStats, error)
	GetFinanceReport() (*FinanceReport, error)
}
