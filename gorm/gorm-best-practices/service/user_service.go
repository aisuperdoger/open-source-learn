package service

import (
	"errors"

	"gorm-best-practices/models"
	"gorm-best-practices/repository"

	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	Register(user *models.User) error
	Login(username, password string) (*models.User, error)
	GetProfile(id uint) (*models.User, error)
	UpdateProfile(user *models.User) error
	DeleteAccount(id uint) error
	ListUsers(page, pageSize int) ([]models.User, error)
	GetUserCount() (int64, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Register 用户注册
func (s *userService) Register(user *models.User) error {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(user.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// 创建用户
	return s.userRepo.Create(user)
}

// Login 用户登录
func (s *userService) Login(username, password string) (*models.User, error) {
	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// 这里应该验证密码，实际项目中需要加密处理
	// 为简化示例，我们直接比较明文密码
	if user.Password != password {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(user *models.User) error {
	// 检查用户是否存在
	existingUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}

	// 更新用户信息
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.Age = user.Age

	return s.userRepo.Update(existingUser)
}

// DeleteAccount 删除账户
func (s *userService) DeleteAccount(id uint) error {
	return s.userRepo.Delete(id)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(page, pageSize int) ([]models.User, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

// GetUserCount 获取用户总数
func (s *userService) GetUserCount() (int64, error) {
	return s.userRepo.Count()
}
