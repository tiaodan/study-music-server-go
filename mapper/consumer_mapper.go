package mapper

import (
	"study-music-server-go/models"
)

type ConsumerMapper struct{}

func NewConsumerMapper() *ConsumerMapper {
	return &ConsumerMapper{}
}

func (*ConsumerMapper) Add(consumer *models.Consumer) error {
	return DB.Create(consumer).Error
}

func (*ConsumerMapper) FindByUsername(username string) (*models.Consumer, error) {
	var consumer models.Consumer
	err := DB.Where("username = ?", username).First(&consumer).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (*ConsumerMapper) FindByEmail(email string) (*models.Consumer, error) {
	var consumer models.Consumer
	err := DB.Where("email = ?", email).First(&consumer).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (*ConsumerMapper) FindById(id uint) (*models.Consumer, error) {
	var consumer models.Consumer
	err := DB.First(&consumer, id).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (*ConsumerMapper) FindAll() ([]models.Consumer, error) {
	var consumers []models.Consumer
	err := DB.Find(&consumers).Error
	return consumers, err
}

func (*ConsumerMapper) Update(consumer *models.Consumer) error {
	return DB.Save(consumer).Error
}

func (*ConsumerMapper) Delete(id uint) error {
	return DB.Delete(&models.Consumer{}, id).Error
}

func (*ConsumerMapper) UpdatePassword(id uint, password string) error {
	return DB.Model(&models.Consumer{}).Where("id = ?", id).Update("password", password).Error
}
