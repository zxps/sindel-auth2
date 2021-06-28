package storages

import (
	"auth2/models"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	availableUserRoles = map[string]*models.Role{
		models.RoleUser: {
			Id:   models.RoleUser,
			Name: "Пользователь",
		},
		models.RoleSuperAdmin: {
			Id:   models.RoleSuperAdmin,
			Name: "Суперадмин",
		},
	}
)

type UserStorage struct {
	table   string
	dialect goqu.DialectWrapper
	db      *sqlx.DB
}

func NewUserStorage(
	table string,
	db *sqlx.DB,
	dialect goqu.DialectWrapper,
) *UserStorage {
	return &UserStorage{
		table:   table,
		dialect: dialect,
		db:      db,
	}
}

func (s *UserStorage) FindIdByColumn(column string, value interface{}) (int, error) {
	var id int

	err := s.db.Get(&id, fmt.Sprintf("SELECT id FROM "+s.table+"WHERE %s=?", column), value)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserStorage) FindByColumn(column string, value interface{}) (*models.User, error) {
	user := models.User{}

	err := s.db.Get(&user, fmt.Sprintf("SELECT * FROM "+s.table+"WHERE %s=?", column), value)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) GetByUsername(username string) (*models.User, error) {
	user := models.User{}
	err := s.db.Get(&user, "SELECT * FROM "+s.table+" WHERE username = ? LIMIT 1", username)
	return &user, err
}

func (s *UserStorage) GetByEmail(email string) (*models.User, error) {
	user := models.User{}
	err := s.db.Get(&user, "SELECT * FROM "+s.table+" WHERE email = ? LIMIT 1", email)
	return &user, err
}

type ScanOptions struct {
	Limit  uint
	Offset uint
	IsDesc *bool
}

func (s *UserStorage) Scan(opts *ScanOptions) ([]models.User, error) {
	var users []models.User

	query := s.dialect.From(s.table)
	query = query.Limit(opts.Limit)
	query = query.Offset(opts.Offset)

	if opts.IsDesc != nil && *opts.IsDesc {
		query = query.Order(goqu.I("id").Desc())
	}

	sql, _, _ := query.ToSQL()

	err := s.db.Select(&users, sql)
	if err != nil {
		logrus.Warning(err)
	}

	return users, err
}

func (s *UserStorage) GetTotalCount() (int, error) {
	var count int

	err := s.db.Get(&count, "SELECT COUNT(*) FROM "+s.table)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UserStorage) GetById(id int) (*models.User, error) {
	user := &models.User{}

	err := s.db.Get(user, "SELECT * FROM "+s.table+" WHERE id = ? LIMIT 1", id)

	return user, err
}

func (s UserStorage) Save(user *models.User) error {
	var sql string
	var err error

	userId := user.Id

	if userId > 0 {
		sql, _, err = s.dialect.Update(s.table).Set(user).Where(goqu.Ex{"id": userId}).Limit(1).ToSQL()
	} else {
		sql, _, err = s.dialect.Insert(s.table).Rows(user).ToSQL()
	}

	if err != nil {
		return err
	}

	result, err := s.db.Exec(sql)

	if err == nil {
		lastInsertID, err := result.LastInsertId()
		if err == nil {
			user.Id = int(lastInsertID)
		}
	}

	if userId > 0 {
		user.Id = userId
	}

	return err
}

func (s *UserStorage) Delete(id int) error {
	sql, _, err := s.dialect.Delete(s.table).Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(sql)

	return err
}

func (s *UserStorage) UpdatePassword(username string, password []byte) error {
	user, err := s.GetByUsername(username)
	if err != nil {
		return err
	}

	_, err = s.db.NamedExec(
		"UPDATE "+s.table+" SET password = :password WHERE username = :username",
		map[string]interface{}{
			"password": password,
			"username": user.Username,
		},
	)

	return err
}

func (s *UserStorage) GetRole(id string) *models.Role {
	role, isExists := availableUserRoles[id]

	if !isExists {
		return nil
	}

	return role
}

func (s *UserStorage) GetRoles(ids []string) []*models.Role {
	var roles []*models.Role

	for _, id := range ids {
		role := s.GetRole(id)
		if role == nil {
			continue
		}

		roles = append(roles, role)
	}

	return roles
}
