CREATE TABLE IF NOT EXISTS schema_migrations (
    version bigint PRIMARY KEY,
    dirty boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE app_user (
    id uuid PRIMARY KEY,
    name text,
    username text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX app_user_username_idx
ON app_user (username);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_app_user_updated_at
BEFORE UPDATE ON app_user
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
