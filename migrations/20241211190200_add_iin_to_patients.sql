-- +goose Up
-- Add IIN (Individual Identification Number) field to patients table
ALTER TABLE patients ADD COLUMN iin VARCHAR(12) UNIQUE;

-- Create index for IIN
CREATE INDEX IF NOT EXISTS idx_patients_iin ON patients(iin);

-- +goose Down
DROP INDEX IF EXISTS idx_patients_iin;
ALTER TABLE patients DROP COLUMN IF EXISTS iin;
