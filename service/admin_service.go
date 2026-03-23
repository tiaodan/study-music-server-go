package service

import (
	"crypto/md5"
	"encoding/hex"
	"study-music-server-go/common"
	"study-music-server-go/mapper"
)

type AdminService struct {
	adminMapper *mapper.AdminMapper
}

func NewAdminService() *AdminService {
	return &AdminService{
		adminMapper: mapper.NewAdminMapper(),
	}
}

func (s *AdminService) Login(username, password string) *common.Response {
	admin, err := s.adminMapper.FindByUsername(username)
	if err != nil {
		return common.Error("管理员不存在")
	}

	// Verify password
	h := md5.New()
	h.Write([]byte(password + common.SALT))
	pwd := hex.EncodeToString(h.Sum(nil))

	if admin.Password != pwd {
		return common.Error("密码错误")
	}

	return common.SuccessWithData("登录成功", admin)
}
