package service

import (
	"coal/model"
	"coal/util"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (us *UserService) Register(username, email, password string) error {
	// 检查用户是否已存在
	existingUser, _ := us.userRepo.GetByUsername(username)
	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	existingUser, _ = us.userRepo.GetByEmail(email)
	if existingUser != nil {
		return errors.New("邮箱已被使用")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	return us.userRepo.Create(user)
}

func (us *UserService) Login(username, password string) (string, error) {
	// 获取用户信息
	user, err := us.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return "", errors.New("用户名或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 生成JWT Token
	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *UserService) GetUserByID(id uint) (*model.User, error) {
	return us.userRepo.GetByID(id)
}

// 注意：实际项目中需要将此方法移到token相关的服务中
func generateToken(userID uint, username string) (string, error) {
	// 调用util包中的JWT工具函数生成token
	return util.GenerateToken(userID, username)
}
