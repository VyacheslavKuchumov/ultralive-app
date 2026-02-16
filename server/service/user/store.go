package user

import (
	"VyacheslavKuchumov/test-backend/types"
	"database/sql"
	"fmt"
)

type Store struct {
	db      *sql.DB
	initErr error
}

func NewStore(db *sql.DB) *Store {
	store := &Store{db: db}
	if db != nil {
		store.initErr = store.ensureUserSchema()
	}
	return store
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	if err := s.ensureReady(); err != nil {
		return nil, err
	}

	row := s.db.QueryRow(
		"SELECT id, first_name, last_name, name, email, password, role, created_at FROM users WHERE email = $1",
		email,
	)
	u, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	if err := s.ensureReady(); err != nil {
		return nil, err
	}

	row := s.db.QueryRow(
		"SELECT id, first_name, last_name, name, email, password, role, created_at FROM users WHERE id = $1",
		id,
	)
	u, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	if err := s.ensureReady(); err != nil {
		return err
	}

	_, err := s.db.Exec(
		"INSERT INTO users (first_name, last_name, name, email, password, role) VALUES ($1, $2, $3, $4, $5, $6)",
		user.FirstName, user.LastName, user.Name, user.Email, user.Password, user.Role,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateUserProfile(userID int, payload types.UpdateProfilePayload) (*types.User, error) {
	if err := s.ensureReady(); err != nil {
		return nil, err
	}

	row := s.db.QueryRow(
		`UPDATE users
		 SET first_name = CAST($1 AS VARCHAR(255)),
		     last_name = CAST($2 AS VARCHAR(255)),
		     name = TRIM(CAST($1 AS VARCHAR(255)) || ' ' || CAST($2 AS VARCHAR(255))),
		     email = COALESCE(NULLIF(CAST($3 AS VARCHAR(255)), ''), email)
		 WHERE id = $4
		 RETURNING id, first_name, last_name, name, email, password, role, created_at`,
		payload.FirstName,
		payload.LastName,
		payload.Email,
		userID,
	)

	u, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) UpdateUserPassword(userID int, hashedPassword string) error {
	if err := s.ensureReady(); err != nil {
		return err
	}

	result, err := s.db.Exec(
		`UPDATE users
		 SET password = $1
		 WHERE id = $2`,
		hashedPassword,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (s *Store) ListUsers() ([]*types.UserLookup, error) {
	if err := s.ensureReady(); err != nil {
		return nil, err
	}

	rows, err := s.db.Query(
		`SELECT id, name
		 FROM users
		 ORDER BY name, id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*types.UserLookup, 0)
	for rows.Next() {
		user := new(types.UserLookup)
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (s *Store) GetUserByName(name string) (*types.User, error) {
	if err := s.ensureReady(); err != nil {
		return nil, err
	}

	row := s.db.QueryRow(
		"SELECT id, first_name, last_name, name, email, password, role, created_at FROM users WHERE name = $1",
		name,
	)
	u, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRowIntoUser(row rowScanner) (*types.User, error) {
	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) ensureReady() error {
	if s.initErr != nil {
		return s.initErr
	}
	if s.db == nil {
		return fmt.Errorf("database is not configured")
	}
	return nil
}

func (s *Store) ensureUserSchema() error {
	statements := []string{
		`ALTER TABLE users
			ADD COLUMN IF NOT EXISTS first_name VARCHAR(255),
			ADD COLUMN IF NOT EXISTS last_name VARCHAR(255),
			ADD COLUMN IF NOT EXISTS name VARCHAR(255),
			ADD COLUMN IF NOT EXISTS role VARCHAR(50),
			ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ`,
		`UPDATE users
		 SET first_name = COALESCE(
			 NULLIF(first_name, ''),
			 COALESCE(NULLIF(split_part(COALESCE(name, ''), ' ', 1), ''), 'User')
		 )
		 WHERE first_name IS NULL OR first_name = ''`,
		`UPDATE users
		 SET last_name = COALESCE(
			 NULLIF(last_name, ''),
			 CASE
				 WHEN POSITION(' ' IN COALESCE(name, '')) > 0
					 THEN COALESCE(NULLIF(TRIM(SUBSTRING(name FROM POSITION(' ' IN name) + 1)), ''), 'User')
				 ELSE 'User'
			 END
		 )
		 WHERE last_name IS NULL OR last_name = ''`,
		`UPDATE users
		 SET name = TRIM(COALESCE(first_name, '') || ' ' || COALESCE(last_name, ''))
		 WHERE name IS NULL OR name = ''`,
		`UPDATE users
		 SET role = 'user'
		 WHERE role IS NULL OR role = ''`,
		`UPDATE users
		 SET created_at = NOW()
		 WHERE created_at IS NULL`,
	}

	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}
