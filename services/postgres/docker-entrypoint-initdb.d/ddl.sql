-------- EXTENSIONS --------
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-------- DDL --------
CREATE TABLE "user" (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
);

CREATE TABLE session (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES "user" (id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL,
    fingerprint TEXT NOT NULL,
    user_ip_address TEXT NOT NULL,  /*IPv4 or IPv6*/
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-------- TABLE CONSTRAINTS --------
-- table 'user'
ALTER TABLE user
ADD CONSTRAINT field_max_length CHECK (
    LENGTH(email) <= 50 AND 
    LENGTH(username) <= 50 AND 
    LENGTH(first_name) <= 50
);

-- table 'session'
ALTER TABLE session
ADD CONSTRAINT field_max_length CHECK (
    LENGTH(refresh_token) <= 60 AND 
    LENGTH(fingerprint) <= 150 AND 
    LENGTH(user_ip_address) <= 39 AND LENGTH(user_ip_address) >= 7
);

-------- FUNCTIONS AND TRIGGERS --------
-- table 'user'
CREATE OR REPLACE FUNCTION update_updated_at_user_column()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON user
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_user_column();