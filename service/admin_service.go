package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/utils"

	"golang.org/x/crypto/bcrypt"
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

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return common.Error("密码错误")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		return common.Error("生成token失败")
	}

	return common.SuccessWithData("登录成功", map[string]string{
		"token": token,
	})
}

// CheckLoginStatus 检查登录状态
func (s *AdminService) CheckLoginStatus(token string) *common.Response {
	if token == "" {
		return common.Error("未登录")
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		return common.Error("token无效或已过期")
	}

	return common.SuccessWithData("已登录", map[string]interface{}{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"role":     claims.Role,
	})
}
