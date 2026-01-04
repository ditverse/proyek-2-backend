package repositories

import (
	"backend-sarpras/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (nama, email, password_hash, role, organisasi_kode)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING kode_user, created_at
	`
	var orgKode interface{}
	if user.OrganisasiKode != nil {
		orgKode = *user.OrganisasiKode
	}
	err := r.DB.QueryRow(
		query,
		user.Nama,
		user.Email,
		user.PasswordHash,
		user.Role,
		orgKode,
	).Scan(&user.KodeUser, &user.CreatedAt)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT kode_user, nama, email, password_hash, role, organisasi_kode, created_at
		FROM users
		WHERE email = $1
	`
	var orgKode sql.NullString
	err := r.DB.QueryRow(query, email).Scan(
		&user.KodeUser,
		&user.Nama,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&orgKode,
		&user.CreatedAt,
	)
	if orgKode.Valid {
		user.OrganisasiKode = &orgKode.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByID(kodeUser string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT kode_user, nama, email, password_hash, role, organisasi_kode, created_at
		FROM users
		WHERE kode_user = $1
	`
	var orgKode sql.NullString
	err := r.DB.QueryRow(query, kodeUser).Scan(
		&user.KodeUser,
		&user.Nama,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&orgKode,
		&user.CreatedAt,
	)
	if orgKode.Valid {
		user.OrganisasiKode = &orgKode.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByRole(role models.RoleEnum) ([]models.User, error) {
	query := `
		SELECT kode_user, nama, email, role, organisasi_kode, created_at
		FROM users
		WHERE role = $1
		ORDER BY nama
	`
	rows, err := r.DB.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var orgKode sql.NullString
		err := rows.Scan(
			&u.KodeUser,
			&u.Nama,
			&u.Email,
			&u.Role,
			&orgKode,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if orgKode.Valid {
			u.OrganisasiKode = &orgKode.String
		}
		users = append(users, u)
	}
	return users, nil
}

// GetByIDs returns a map of users keyed by kode_user for batch loading
func (r *UserRepository) GetByIDs(kodes []string) map[string]*models.User {
	result := make(map[string]*models.User)
	if len(kodes) == 0 {
		return result
	}

	query := `
		SELECT kode_user, nama, email, role, organisasi_kode, created_at
		FROM users
		WHERE kode_user = ANY($1)
	`
	rows, err := r.DB.Query(query, kodes)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		u := &models.User{}
		var orgKode sql.NullString
		err := rows.Scan(
			&u.KodeUser,
			&u.Nama,
			&u.Email,
			&u.Role,
			&orgKode,
			&u.CreatedAt,
		)
		if err == nil {
			if orgKode.Valid {
				u.OrganisasiKode = &orgKode.String
			}
			result[u.KodeUser] = u
		}
	}
	return result
}
