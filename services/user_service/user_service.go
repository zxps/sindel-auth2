package user_service

import (
	"auth2/models"
	"auth2/storages"
	"auth2/utils"
)

type Service struct {
	storage *storages.UserStorage
}

func New(storage *storages.UserStorage) *Service {
	return &Service{storage: storage}
}

func (u *Service) UpdatePassword(username string, hashedPassword []byte) error {
	return u.storage.UpdatePassword(username, hashedPassword)
}

func (u *Service) Scan(opts *storages.ScanOptions) ([]models.User, error) {
	return u.storage.Scan(opts)
}

func (u *Service) GetById(id int) (*models.User, error) {
	return u.storage.GetById(id)
}

func (u *Service) GetByEmail(email string) (*models.User, error) {
	return u.storage.GetByEmail(email)
}

func (u *Service) GetByUsername(username string) (*models.User, error) {
	return u.storage.GetByUsername(username)
}

func (u *Service) GetByConfirmationToken(token string) (*models.User, error) {
	return u.storage.FindByColumn("confirmation_token", token)
}

func (u *Service) CreateUser(user *models.User) error {
	user.Created = utils.StorageTimestamp()
	user.Updated = utils.StorageTimestamp()

	return u.storage.Save(user)
}

func (u *Service) UpdateUser(user *models.User) error {
	user.Updated = utils.StorageTimestamp()
	return u.storage.Save(user)
}

func (u *Service) DeleteUser(user *models.User) error {
	return u.storage.Delete(user.Id)
}

func (u *Service) GetTotalCount() int {
	totalCount, _ := u.storage.GetTotalCount()
	return totalCount
}

func (u *Service) IsExists(id int) bool {
	_, err := u.storage.GetById(id)

	return err == nil
}

func (u *Service) IsExistsByEmail(email string) bool {
	_, err := u.storage.GetByEmail(email)

	return err == nil
}

func (u *Service) IsExistsByUsername(username string) bool {
	_, err := u.storage.GetByUsername(username)

	return err == nil
}

func (u *Service) GetRole(id string) *models.Role {
	return u.storage.GetRole(id)
}

func (u *Service) GetRoles(ids []string) []*models.Role {
	return u.storage.GetRoles(ids)
}
