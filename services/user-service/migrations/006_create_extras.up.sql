CREATE UNIQUE INDEX uq_users_email_lower
ON users (LOWER(email));

CREATE UNIQUE INDEX uq_user_default_address
ON addresses(user_id)
WHERE is_default = TRUE;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();