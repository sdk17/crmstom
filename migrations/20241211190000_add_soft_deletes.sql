-- +goose Up
-- Add soft delete support to all tables

ALTER TABLE patients ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
ALTER TABLE services ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
ALTER TABLE doctors ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
ALTER TABLE appointments ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;

-- Create indexes for soft delete queries
CREATE INDEX IF NOT EXISTS idx_patients_deleted_at ON patients(deleted_at);
CREATE INDEX IF NOT EXISTS idx_services_deleted_at ON services(deleted_at);
CREATE INDEX IF NOT EXISTS idx_doctors_deleted_at ON doctors(deleted_at);
CREATE INDEX IF NOT EXISTS idx_appointments_deleted_at ON appointments(deleted_at);
