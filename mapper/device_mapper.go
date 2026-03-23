package mapper

import (
	"study-music-server-go/models"
)

type DeviceMapper struct{}

func NewDeviceMapper() *DeviceMapper {
	return &DeviceMapper{}
}

func (*DeviceMapper) Add(device *models.Device) error {
	return DB.Create(device).Error
}

func (*DeviceMapper) FindAll() ([]models.Device, error) {
	var devices []models.Device
	err := DB.Find(&devices).Error
	return devices, err
}

func (*DeviceMapper) FindById(id uint) (*models.Device, error) {
	var device models.Device
	err := DB.First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (*DeviceMapper) FindByType(deviceType string) ([]models.Device, error) {
	var devices []models.Device
	err := DB.Where("type = ?", deviceType).Find(&devices).Error
	return devices, err
}

func (*DeviceMapper) FindDefaultByType(deviceType string) (*models.Device, error) {
	var device models.Device
	err := DB.Where("type = ? AND is_default = ?", deviceType, true).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (*DeviceMapper) FindByName(name string) (*models.Device, error) {
	var device models.Device
	err := DB.Where("name = ?", name).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (*DeviceMapper) Update(device *models.Device) error {
	return DB.Save(device).Error
}

func (*DeviceMapper) Delete(id uint) error {
	return DB.Delete(&models.Device{}, id).Error
}
