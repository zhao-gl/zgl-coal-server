package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"coal/config"
	"coal/controller"
	"coal/model"
	"coal/service"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (ur *userRepository) Create(user *model.User) error {
	return ur.db.Create(user).Error
}

func (ur *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := ur.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) Update(user *model.User) error {
	return ur.db.Save(user).Error
}

func (ur *userRepository) Delete(id uint) error {
	return ur.db.Delete(&model.User{}, id).Error
}

func main() {
	// 加载 .env 文件（优先加载，若文件不存在则忽略，不影响生产环境）
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到 .env 文件，使用系统环境变量或默认值")
	}
	// 初始化数据库连接
	dbConfig := config.GetDatabaseConfig()
	dsn := dbConfig.ConnectionString()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}

	// 自动迁移数据库模式
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 初始化数据仓库和服务
	userRepo := &userRepository{db: db}
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// 设置路由
	r := mux.NewRouter()
	r.HandleFunc("/register", userController.Register).Methods("POST")
	r.HandleFunc("/login", userController.Login).Methods("POST")

	// 需要身份验证的路由
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/profile", userController.GetProfile).Methods("GET")

	// 启动服务器
	log.Println("服务器启动在端口 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
