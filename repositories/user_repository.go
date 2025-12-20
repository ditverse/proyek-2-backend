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
		INSERT INTO users (nama, email, password_hash, role, organisasi_kode, no_hp)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING kode_user, created_at
	`
	var orgKode interface{}
	if user.OrganisasiKode != nil {
		orgKode = *user.OrganisasiKode
	}
	var noHP interface{}
	if user.NoHP != nil {
		noHP = *user.NoHP
	}
	err := r.DB.QueryRow(
		query,
		user.Nama,
		user.Email,
		user.PasswordHash,
		user.Role,
		orgKode,
		noHP,
	).Scan(&user.KodeUser, &user.CreatedAt)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT kode_user, nama, email, password_hash, role, organisasi_kode, no_hp, created_at
		FROM users
		WHERE email = $1
	`
	var orgKode sql.NullString
	var noHP sql.NullString
	err := r.DB.QueryRow(query, email).Scan(
		&user.KodeUser,
		&user.Nama,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&orgKode,
		&noHP,
		&user.CreatedAt,
	)
	if orgKode.Valid {
		user.OrganisasiKode = &orgKode.String
	}
	if noHP.Valid {
		user.NoHP = &noHP.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByID(kodeUser string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT kode_user, nama, email, password_hash, role, organisasi_kode, no_hp, created_at
		FROM users
		WHERE kode_user = $1
	`
	var orgKode sql.NullString
	var noHP sql.NullString
	err := r.DB.QueryRow(query, kodeUser).Scan(
		&user.KodeUser,
		&user.Nama,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&orgKode,
		&noHP,
		&user.CreatedAt,
	)
	if orgKode.Valid {
		user.OrganisasiKode = &orgKode.String
	}
	if noHP.Valid {
		user.NoHP = &noHP.String
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByRole(role models.RoleEnum) ([]models.User, error) {
	query := `
		SELECT kode_user, nama, email, role, organisasi_kode, no_hp, created_at
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
		var noHP sql.NullString
		err := rows.Scan(
			&u.KodeUser,
			&u.Nama,
			&u.Email,
			&u.Role,
			&orgKode,
			&noHP,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if orgKode.Valid {
			u.OrganisasiKode = &orgKode.String
		}
		if noHP.Valid {
			u.NoHP = &noHP.String
		}
		users = append(users, u)
	}
	return users, nil
}
