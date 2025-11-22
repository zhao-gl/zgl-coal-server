package controller

import (
	"encoding/json"
	"net/http"

	"coal/service"
	"coal/util"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "用户名、邮箱和密码不能为空", http.StatusBadRequest)
		return
	}

	if err := uc.userService.Register(req.Username, req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "用户注册成功"})
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求数据", http.StatusBadRequest)
		return
	}

	token, err := uc.userService.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "登录成功",
		"token":   token,
	})
}

func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取用户ID (通常从JWT token中解析)
	claims, err := util.ParseTokenFromHeader(r)
	if err != nil {
		http.Error(w, "无效的访问令牌", http.StatusUnauthorized)
		return
	}

	if claims == nil {
		http.Error(w, "访问令牌缺失", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID
	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
