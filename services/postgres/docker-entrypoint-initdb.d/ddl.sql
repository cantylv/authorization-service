-------- EXTENSIONS --------
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-------- DDL --------
-- Эта таблица содержит данные о пользователях
CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT,
    password TEXT,
    first_name TEXT,
    last_name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Эта таблица содержит названия микросервисов
CREATE TABLE microservice (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Эта таблица содержит данные о группах
CREATE TABLE "group" (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_id UUID REFERENCES "user"(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Эта таблица содержит названия процессов
CREATE TABLE privelege (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    microservice_id INT REFERENCES microservice(id) ON DELETE CASCADE,
    name TEXT, /*request_url.method*/
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Эта таблица содержит доступы пользователей к процессам
CREATE TABLE role (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID REFERENCES "user" (id) ON DELETE CASCADE,
    privelege_id INT REFERENCES privelege(id) ON DELETE CASCADE
);

-- Эта таблица содержит доступы пользователей к процессам
CREATE TABLE participation (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
    group_id INT REFERENCES "group"(id) ON DELETE CASCADE
);

-- Эта таблица содержит доступы групп к микросервисам 
CREATE TABLE workteam (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    group_id INT REFERENCES "group"(id) ON DELETE CASCADE,
    microservice_id INT REFERENCES microservice(id) ON DELETE CASCADE
);

-------- INDEXES --------
CREATE INDEX user_privelege ON role (user_id, privelege_id);

-------- TABLE CONSTRAINTS --------
-- table 'user'
ALTER TABLE "user"
ADD CONSTRAINT user_unique_email UNIQUE (email),
ADD CONSTRAINT user_email_length CHECK (LENGTH(email) <= 50 AND LENGTH(email) >= 6),
ADD CONSTRAINT user_password_length CHECK (LENGTH(password) = 145),
ADD CONSTRAINT user_first_name_length CHECK (LENGTH(first_name) <= 50 AND LENGTH(first_name) >= 2),
ADD CONSTRAINT user_last_name_length CHECK (LENGTH(last_name) <= 50 AND LENGTH(last_name) >= 2);

ALTER TABLE "user"
ALTER COLUMN email SET NOT NULL,
ALTER COLUMN password SET NOT NULL,
ALTER COLUMN first_name SET NOT NULL,
ALTER COLUMN last_name SET NOT NULL,
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

-- table 'privelege'
ALTER TABLE privelege
ADD CONSTRAINT privelege_name_length CHECK (LENGTH(name) <= 200 AND LENGTH(name) >= 10);

ALTER TABLE privelege
ALTER COLUMN name SET NOT NULL,
ALTER COLUMN created_at SET NOT NULL;

-- table 'group'
ALTER TABLE "group"
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

-- table 'microservice'
ALTER TABLE microservice
ADD CONSTRAINT microservice_name_unique UNIQUE(name),
ADD CONSTRAINT microservice_name_length CHECK (LENGTH(name) <= 50 AND LENGTH(name) >= 2);

ALTER TABLE microservice
ALTER COLUMN name SET NOT NULL,
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

-------- FUNCTIONS AND TRIGGERS --------
-- table 'user'
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_group_updated_at
BEFORE UPDATE ON "group"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_microservice_updated_at
BEFORE UPDATE ON microservice
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();