-------- EXTENSIONS --------
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-------- DDL --------
-- Эта таблица содержит данные о пользователях
CREATE TABLE group (
    id INT GENERATED ALWAYS AS IDENTITY,
    owner_id UUID REFERENCES "user"(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Эта таблица содержит доступы пользователей к процессам
CREATE TABLE role (
    id INT GENERATED ALWAYS AS IDENTITY,
    user_id UUID REFERENCES "user" (id) ON DELETE CASCADE,
    privelege_id REFERENCES privelege(id) ON DELETE CASCADE
);

-- Эта таблица содержит названия процессов
CREATE TABLE privelege (
    id INT GENERATED ALWAYS AS IDENTITY,
    microservice_id REFERENCES microsevice(id) ON DELETE CASCADE,
    name TEXT, /*request_url.method*/
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Эта таблица содержит данные о группах
CREATE TABLE group (
    id INT GENERATED ALWAYS AS IDENTITY,
    owner_id UUID REFERENCES "user"(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Эта таблица содержит доступы пользователей к процессам
CREATE TABLE participation (
    id INT GENERATED ALWAYS AS IDENTITY,
    user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
    group_id INT REFERENCES group(id) ON DELETE CASCADE
)

-- Эта таблица содержит доступы групп к микросервисам 
CREATE TABLE scope (
    id INT GENERATED ALWAYS AS IDENTITY,
    group_id
)


-------- INDEXES --------
CREATE INDEX user_privelege ON role (user_id, privelege_id);

-------- TABLE CONSTRAINTS --------
-- table 'user'
ALTER TABLE "user"
ADD CONSTRAINT fields_length CHECK (
    LENGTH(email) <= 50 AND LENGTH(email) >= 6 AND 
    LENGTH(password) = 145 AND 
    LENGTH(first_name) <= 50 AND LENGTH(first_name) >= 2 AND
    LENGTH(last_name) <= 50 AND LENGTH(last_name) >= 2
);

-- table roles
ALTER TABLE privelege
ADD CONSTRAINT name_privelege_max_length CHECK (
    LENGTH(privelege) <= 200 AND LENGTH(privelege) >= 10
);

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