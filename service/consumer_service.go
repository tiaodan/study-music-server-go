package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
	"study-music-server-go/utils"
)

type ConsumerService struct {
	consumerMapper *mapper.ConsumerMapper
}

func NewConsumerService() *ConsumerService {
	return &ConsumerService{
		consumerMapper: mapper.NewConsumerMapper(),
	}
}

func (s *ConsumerService) AddUser(req *models.ConsumerRequest) *common.Response {
	// Check if username exists
	existing, _ := s.consumerMapper.FindByUsername(req.Username)
	if existing != nil {
		return common.Warning("用户名已存在")
	}

	// Hash password
	h := md5.New()
	h.Write([]byte(req.Password + common.SALT))
	password := hex.EncodeToString(h.Sum(nil))

	consumer := &models.Consumer{
		Username:    req.Username,
		Password:    password,
		Sex:         req.Sex,
		PhoneNum:    req.PhoneNum,
		Email:       req.Email,
		Birth:       req.Birth,
		Introduction: req.Introduction,
		Location:    req.Location,
		Avator:      req.Avator,
	}

	err := s.consumerMapper.Add(consumer)
	if err != nil {
		return common.Error("注册失败")
	}
	return common.Success("注册成功")
}

func (s *ConsumerService) LoginStatus(req *models.ConsumerRequest) *common.Response {
	consumer, err := s.consumerMapper.FindByUsername(req.Username)
	if err != nil {
		return common.Error("用户名或密码错误")
	}

	// Verify password
	h := md5.New()
	h.Write([]byte(req.Password + common.SALT))
	password := hex.EncodeToString(h.Sum(nil))

	if consumer.Password != password {
		return common.Error("用户名或密码错误")
	}

	return common.SuccessWithData("登录成功", consumer)
}

func (s *ConsumerService) LoginEmailStatus(req *models.ConsumerRequest) *common.Response {
	consumer, err := s.consumerMapper.FindByEmail(req.Email)
	if err != nil {
		return common.Error("邮箱或密码错误")
	}

	// Verify password
	h := md5.New()
	h.Write([]byte(req.Password + common.SALT))
	password := hex.EncodeToString(h.Sum(nil))

	if consumer.Password != password {
		return common.Error("邮箱或密码错误")
	}

	return common.SuccessWithData("登录成功", consumer)
}

func (s *ConsumerService) AllUser() *common.Response {
	consumers, err := s.consumerMapper.FindAll()
	if err != nil {
		return common.Error("获取用户列表失败")
	}
	return common.SuccessWithData("获取成功", consumers)
}

func (s *ConsumerService) UserOfId(id uint) *common.Response {
	consumer, err := s.consumerMapper.FindById(id)
	if err != nil {
		return common.Error("用户不存在")
	}
	return common.SuccessWithData("获取成功", consumer)
}

func (s *ConsumerService) DeleteUser(id uint) *common.Response {
	err := s.consumerMapper.Delete(id)
	if err != nil {
		return common.Error("删除失败")
	}
	return common.Success("删除成功")
}

func (s *ConsumerService) UpdateUserMsg(req *models.ConsumerRequest) *common.Response {
	consumer, err := s.consumerMapper.FindById(req.ID)
	if err != nil {
		return common.Error("用户不存在")
	}

	consumer.Username = req.Username
	consumer.Sex = req.Sex
	consumer.PhoneNum = req.PhoneNum
	consumer.Email = req.Email
	consumer.Birth = req.Birth
	consumer.Introduction = req.Introduction
	consumer.Location = req.Location

	err = s.consumerMapper.Update(consumer)
	if err != nil {
		return common.Error("更新失败")
	}
	return common.Success("更新成功")
}

func (s *ConsumerService) UpdatePassword(req *models.ConsumerRequest) *common.Response {
	consumer, err := s.consumerMapper.FindById(req.ID)
	if err != nil {
		return common.Error("用户不存在")
	}

	// Verify old password
	h := md5.New()
	h.Write([]byte(req.OldPassword + common.SALT))
	oldPassword := hex.EncodeToString(h.Sum(nil))

	if consumer.Password != oldPassword {
		return common.Error("原密码错误")
	}

	// Update password
	h.Reset()
	h.Write([]byte(req.Password + common.SALT))
	newPassword := hex.EncodeToString(h.Sum(nil))

	err = s.consumerMapper.UpdatePassword(req.ID, newPassword)
	if err != nil {
		return common.Error("密码更新失败")
	}
	return common.Success("密码更新成功")
}

func (s *ConsumerService) UpdatePasswordByEmail(email, password string) error {
	h := md5.New()
	h.Write([]byte(password + common.SALT))
	pwd := hex.EncodeToString(h.Sum(nil))

	consumer, err := s.consumerMapper.FindByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.consumerMapper.UpdatePassword(consumer.ID, pwd)
}

func (s *ConsumerService) UpdateUserAvator(filePath string, id uint) *common.Response {
	consumer, err := s.consumerMapper.FindById(id)
	if err != nil {
		return common.Error("用户不存在")
	}

	// Save file and get path
	filename := utils.SaveFile(filePath, "avatorImages")
	if filename == "" {
		return common.Error("头像上传失败")
	}

	consumer.Avator = common.AVATOR_IMAGES_PATH + filename
	err = s.consumerMapper.Update(consumer)
	if err != nil {
		return common.Error("更新头像失败")
	}
	return common.SuccessWithData("更新成功", consumer.Avator)
}

func (s *ConsumerService) FindByEmail(email string) (*models.Consumer, error) {
	return s.consumerMapper.FindByEmail(email)
}
