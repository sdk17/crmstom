-- Инициализация базы данных для CRM Стоматология

-- Создание таблицы пациентов
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    birth_date DATE,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы услуг (упрощенная схема)
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы врачей
CREATE TABLE IF NOT EXISTS doctors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    login VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы записей
CREATE TABLE IF NOT EXISTS appointments (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER REFERENCES patients(id) ON DELETE CASCADE,
    service_id INTEGER REFERENCES services(id) ON DELETE CASCADE,
    doctor_id INTEGER REFERENCES doctors(id) ON DELETE SET NULL,
    appointment_date TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled',
    price DECIMAL(10,2),
    duration_minutes INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_appointments_date ON appointments(appointment_date);
CREATE INDEX IF NOT EXISTS idx_appointments_patient ON appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_service ON appointments(service_id);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor ON appointments(doctor_id);
CREATE INDEX IF NOT EXISTS idx_doctors_login ON doctors(login);

-- Вставка тестовых данных
INSERT INTO patients (name, phone, email, birth_date, address) VALUES
('Иванов Иван Иванович', '+7 (777) 123-45-67', 'ivanov@example.com', '1985-03-15', 'ул. Абая 150, кв. 25'),
('Петрова Анна Сергеевна', '+7 (777) 234-56-78', 'petrova@example.com', '1990-07-22', 'пр. Назарбаева 45, кв. 12'),
('Сидоров Петр Александрович', '+7 (777) 345-67-89', 'sidorov@example.com', '1978-11-08', 'ул. Сатпаева 78, кв. 8')
ON CONFLICT DO NOTHING;

-- Вставка тестовых врачей
INSERT INTO doctors (name, email, login, password, is_admin) VALUES
('Др. Смит', 'smith@clinic.com', 'dr_smith', 'password123', false),
('Др. Джонс', 'jones@clinic.com', 'dr_jones', 'password123', false),
('Др. Уилсон', 'wilson@clinic.com', 'dr_wilson', 'password123', true)
ON CONFLICT (login) DO NOTHING;

INSERT INTO services (name, type, notes) VALUES
-- Консультации
('Первичная консультация', 'Консультация', 'Осмотр и составление плана лечения'),
('Повторная консультация', 'Консультация', 'Контрольный осмотр после лечения'),

-- Лечение кариеса
('Лечение кариеса (поверхностный)', 'Лечение кариеса', 'Пломбирование поверхностного кариеса'),
('Лечение кариеса (средний)', 'Лечение кариеса', 'Пломбирование среднего кариеса'),
('Лечение кариеса (глубокий)', 'Лечение кариеса', 'Пломбирование глубокого кариеса'),

-- Лечение пульпита
('Лечение пульпита (одноканальный)', 'Лечение пульпита', 'Эндодонтическое лечение одноканального зуба'),
('Лечение пульпита (двухканальный)', 'Лечение пульпита', 'Эндодонтическое лечение двухканального зуба'),
('Лечение пульпита (трехканальный)', 'Лечение пульпита', 'Эндодонтическое лечение трехканального зуба'),

-- Удаление зубов
('Удаление зуба (простое)', 'Удаление зубов', 'Простое удаление зуба'),
('Удаление зуба (сложное)', 'Удаление зубов', 'Сложное удаление зуба с разрезом'),
('Удаление зуба мудрости', 'Удаление зубов', 'Удаление зуба мудрости'),

-- Гигиена
('Профессиональная чистка', 'Гигиена', 'Ультразвуковая чистка и полировка'),
('Чистка Air Flow', 'Гигиена', 'Чистка методом Air Flow'),

-- Протезирование
('Коронка металлокерамическая', 'Протезирование', 'Установка металлокерамической коронки'),
('Коронка керамическая', 'Протезирование', 'Установка керамической коронки'),
('Мост (3 зуба)', 'Протезирование', 'Установка мостовидного протеза на 3 зуба'),

-- Имплантация
('Имплант (установка)', 'Имплантация', 'Установка зубного импланта'),
('Имплант (с коронкой)', 'Имплантация', 'Имплант с установкой коронки'),

-- Детская стоматология
('Лечение молочного зуба', 'Детская стоматология', 'Лечение кариеса молочного зуба'),
('Герметизация фиссур', 'Детская стоматология', 'Профилактическая герметизация')
ON CONFLICT DO NOTHING;

INSERT INTO appointments (patient_id, service_id, appointment_date, status, price, duration_minutes, notes) VALUES
(1, 1, '2024-09-20 10:00:00', 'scheduled', 5000.00, 30, 'Первичная консультация'),
(2, 12, '2024-09-20 14:00:00', 'scheduled', 15000.00, 60, 'Профессиональная чистка'),
(1, 4, '2024-09-25 11:00:00', 'scheduled', 15000.00, 60, 'Лечение кариеса на 6-м зубе'),
(3, 1, '2024-09-21 09:00:00', 'scheduled', 5000.00, 30, 'Консультация по протезированию'),
(1, 7, '2024-09-26 15:00:00', 'scheduled', 35000.00, 120, 'Лечение пульпита 7-го зуба'),
(2, 10, '2024-09-22 11:00:00', 'scheduled', 20000.00, 90, 'Удаление зуба мудрости'),
(3, 15, '2024-09-23 10:00:00', 'scheduled', 120000.00, 120, 'Установка коронки')
ON CONFLICT DO NOTHING;
