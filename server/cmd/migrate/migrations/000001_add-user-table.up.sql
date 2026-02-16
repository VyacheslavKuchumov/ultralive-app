DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = 'public' AND table_name = 'users'
  ) THEN
    CREATE TABLE users (
      id BIGSERIAL PRIMARY KEY,
      first_name VARCHAR(255),
      last_name VARCHAR(255),
      name VARCHAR(255),
      email VARCHAR(255),
      password VARCHAR(255),
      role VARCHAR(50) DEFAULT 'user',
      created_at TIMESTAMPTZ DEFAULT NOW()
    );
  END IF;
END $$;

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS first_name VARCHAR(255),
  ADD COLUMN IF NOT EXISTS last_name VARCHAR(255),
  ADD COLUMN IF NOT EXISTS name VARCHAR(255),
  ADD COLUMN IF NOT EXISTS email VARCHAR(255),
  ADD COLUMN IF NOT EXISTS password VARCHAR(255),
  ADD COLUMN IF NOT EXISTS role VARCHAR(50),
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'users' AND column_name = 'user_uid'
  ) AND EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = 'public' AND table_name = 'auths'
  ) THEN
    UPDATE users u
    SET
      email = COALESCE(NULLIF(u.email, ''), LOWER(a.email)),
      password = COALESCE(NULLIF(u.password, ''), a.password)
    FROM auths a
    WHERE a.auth_uid = u.user_uid;
  END IF;
END $$;

UPDATE users
SET first_name = COALESCE(
  NULLIF(first_name, ''),
  COALESCE(NULLIF(split_part(name, ' ', 1), ''), 'Legacy')
)
WHERE first_name IS NULL OR first_name = '';

UPDATE users
SET last_name = COALESCE(
  NULLIF(last_name, ''),
  CASE
    WHEN POSITION(' ' IN COALESCE(name, '')) > 0 THEN
      COALESCE(NULLIF(TRIM(SUBSTRING(name FROM POSITION(' ' IN name) + 1)), ''), 'User')
    ELSE 'User'
  END
)
WHERE last_name IS NULL OR last_name = '';

UPDATE users
SET name = TRIM(first_name || ' ' || last_name)
WHERE name IS NULL OR name = '';

UPDATE users
SET role = COALESCE(NULLIF(role, ''), 'user')
WHERE role IS NULL OR role = '';

UPDATE users
SET email = CONCAT('legacy_user_', id, '@ultralive.local')
WHERE email IS NULL OR email = '';

UPDATE users
SET password = '!'
WHERE password IS NULL OR password = '';

UPDATE users
SET created_at = NOW()
WHERE created_at IS NULL;

ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN last_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN name SET NOT NULL;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ALTER COLUMN password SET NOT NULL;
ALTER TABLE users ALTER COLUMN role SET NOT NULL;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'user';
ALTER TABLE users ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE users ALTER COLUMN created_at SET DEFAULT NOW();

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'users_email_key'
  ) THEN
    ALTER TABLE users
      ADD CONSTRAINT users_email_key UNIQUE (email);
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_name ON users (name);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
