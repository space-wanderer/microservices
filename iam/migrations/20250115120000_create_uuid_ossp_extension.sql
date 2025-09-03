-- +goose Up
-- Создаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose Down
-- Удаляем расширение
DROP EXTENSION IF EXISTS "uuid-ossp";
