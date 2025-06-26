-- Включаем расширение для работы с UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица users: id передаётся извне
CREATE TABLE users (
    id uuid PRIMARY KEY,  -- значение всегда передаётся из приложения
    email varchar(255) UNIQUE NOT NULL,
    first_name varchar(100),
    date_of_birth date,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    last_login timestamp,
    is_active boolean DEFAULT true
);

-- Таблица user_profiles
CREATE TABLE user_profiles (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid UNIQUE REFERENCES users(id),
    avg_cycle_length integer,
    avg_period_length integer,
    usage_goals varchar[],
    language varchar(10) DEFAULT 'en',
    theme varchar(20) DEFAULT 'system',
    notifications_enabled boolean DEFAULT true,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

-- Остальные таблицы из вашего init.sql...
-- (копируем всё остальное из deploy/dev/db-test/init.sql)

-- Таблица menstrual_cycles
CREATE TABLE menstrual_cycles (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid REFERENCES users(id),
    start_date date NOT NULL,
    end_date date,
    period_end_date date,
    notes text,
    is_predicted boolean DEFAULT false,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_menstrual_cycles_start_date ON menstrual_cycles(start_date);
CREATE INDEX idx_menstrual_cycles_user_start ON menstrual_cycles(user_id, start_date);

-- Таблица symptoms
CREATE TABLE symptoms (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    category varchar(50),
    name varchar(100) NOT NULL,
    icon varchar(50),
    display_order integer,
    is_active boolean DEFAULT true
);

CREATE INDEX idx_symptoms_category ON symptoms(category);

-- Таблица user_symptoms
CREATE TABLE user_symptoms (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid REFERENCES users(id),
    date date NOT NULL,
    symptom_id uuid REFERENCES symptoms(id),
    intensity integer CHECK (intensity BETWEEN 1 AND 5),
    notes text,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date, symptom_id)
);

-- Таблица moods
CREATE TABLE moods (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(100) NOT NULL,
    icon varchar(50),
    is_positive boolean,
    is_active boolean DEFAULT true
);

-- Таблица user_moods
CREATE TABLE user_moods (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid REFERENCES users(id),
    date date NOT NULL,
    mood_id uuid REFERENCES moods(id),
    intensity integer CHECK (intensity BETWEEN 1 AND 5),
    notes text,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date, mood_id)
);

-- Таблица medications
CREATE TABLE medications (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid REFERENCES users(id),
    name varchar(100) NOT NULL,
    type varchar(50),
    dosage varchar(100),
    start_date date NOT NULL,
    end_date date,
    reminder_enabled boolean DEFAULT false,
    reminder_time time,
    notes text,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

-- Таблица notes
CREATE TABLE notes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid REFERENCES users(id),
    date date NOT NULL,
    text text,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notes_user_date ON notes(user_id, date);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at
    BEFORE UPDATE ON user_profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menstrual_cycles_updated_at
    BEFORE UPDATE ON menstrual_cycles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_symptoms_updated_at
    BEFORE UPDATE ON user_symptoms
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_moods_updated_at
    BEFORE UPDATE ON user_moods
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_medications_updated_at
    BEFORE UPDATE ON medications
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notes_updated_at
    BEFORE UPDATE ON notes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();