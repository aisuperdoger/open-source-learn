package main

import (
	"fmt"
	"log"

	"gorm-best-practices/config"
	"gorm-best-practices/models"
	"gorm-best-practices/repository"
	"gorm-best-practices/service"

	"gorm.io/gorm"
)

func main() {

	// 初始化数据库连接
	dbConfig := config.DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "zwl1819",
		Name:     "testdb",
	}

	// 连接数据库
	db := config.InitDatabase(dbConfig)
	defer config.CloseDatabase()

	// 自动迁移数据库表
	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.Post{}, &models.UserProfile{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migration completed successfully")

	// 创建仓库和服务实例
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// 演示GORM最佳实践
	demonstrateGORMBestPractices(userService, userRepo, db)
}

func demonstrateGORMBestPractices(userService service.UserService, userRepo repository.UserRepository, db *gorm.DB) {
	fmt.Println("\n=== GORM最佳实践演示 ===")

	// 1. 创建用户
	fmt.Println("\n1. 创建用户:")
	user := &models.User{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "password123",
		Age:      25,
		Status:   "active",
	}

	if err := userService.Register(user); err != nil {
		// 如果用户已存在，获取现有用户
		existingUser, getErr := userRepo.GetByUsername("john_doe")
		if getErr != nil {
			log.Fatal("Failed to get existing user:", getErr)
		}
		user = existingUser
		fmt.Printf("用户已存在，ID: %d\n", user.ID)
	} else {
		fmt.Printf("用户创建成功，ID: %d\n", user.ID)
	}

	// 2. 查询用户
	fmt.Println("\n2. 查询用户:")
	foundUser, err := userService.GetProfile(user.ID)
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}
	fmt.Printf("用户信息: ID=%d, Username=%s, Email=%s, Age=%d\n",
		foundUser.ID, foundUser.Username, foundUser.Email, foundUser.Age)

	// 3. 更新用户
	fmt.Println("\n3. 更新用户:")
	foundUser.Age = 26
	if err := userService.UpdateProfile(foundUser); err != nil {
		log.Fatal("Failed to update user:", err)
	}
	fmt.Println("用户更新成功")

	// 4. 查询用户列表
	fmt.Println("\n4. 查询用户列表:")
	users, err := userService.ListUsers(1, 10)
	if err != nil {
		log.Fatal("Failed to list users:", err)
	}
	fmt.Printf("查询到 %d 个用户:\n", len(users))
	for _, u := range users {
		fmt.Printf("  - ID: %d, Username: %s, Email: %s\n", u.ID, u.Username, u.Email)
	}

	// 5. 获取用户总数
	fmt.Println("\n5. 获取用户总数:")
	count, err := userService.GetUserCount()
	if err != nil {
		log.Fatal("Failed to get user count:", err)
	}
	fmt.Printf("用户总数: %d\n", count)

	// 6. 演示预加载关联数据
	fmt.Println("\n6. 预加载关联数据:")
	var userWithPosts models.User
	if err := db.Preload("Posts").First(&userWithPosts, user.ID).Error; err != nil {
		log.Fatal("Failed to preload user posts:", err)
	}
	fmt.Printf("用户 %s 的文章数量: %d\n", userWithPosts.Username, len(userWithPosts.Posts))

	// 7. 演示事务处理
	fmt.Println("\n7. 事务处理:")
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}

	// 在事务中创建订单
	order := &models.Order{
		UserID:  user.ID,
		Amount:  99.99,
		Status:  "paid",
		Product: "GORM最佳实践课程",
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		log.Fatal("Failed to create order:", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}
	fmt.Printf("订单创建成功，ID: %d\n", order.ID)

	// 8. 演示原生SQL查询
	fmt.Println("\n8. 原生SQL查询:")
	var userCount int64
	db.Raw("SELECT COUNT(*) FROM users").Scan(&userCount)
	fmt.Printf("通过原生SQL查询到的用户数量: %d\n", userCount)

	fmt.Println("\n=== 演示完成 ===")

	err = userService.DeleteAccount(user.ID)
	if err != nil {
		log.Fatal("Failed to delete user:", err)
	}
}
