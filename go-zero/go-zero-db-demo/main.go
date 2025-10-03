package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"go-zero-db-demo/model"
)

func main() {
	// 数据库连接配置
	dataSource := "root:password@tcp(127.0.0.1:3306)/test"

	// 创建数据库连接
	conn := sqlx.NewSqlConn("mysql", dataSource)

	// 缓存配置
	cacheConf := cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: "127.0.0.1:6379",
				Type: "node",
			},
			Weight: 100,
		},
	}

	// 创建用户模型
	userModel := model.NewUserModel(conn, cacheConf)

	// 创建用户
	user := model.User{
		Name:     fmt.Sprintf("user_%d", time.Now().Unix()),
		Password: "password123",
		Mobile:   fmt.Sprintf("138%08d", time.Now().Unix()%100000000),
		Gender:   "male",
		Nickname: "test_user",
	}

	fmt.Println("创建用户...")
	id, err := userModel.Insert(user)
	if err != nil {
		log.Fatal("创建用户失败:", err)
	}
	fmt.Printf("用户创建成功，ID: %d\n", id)

	// 查询用户
	fmt.Println("\n查询用户...")
	foundUser, err := userModel.FindOne(id)
	if err != nil {
		log.Fatal("查询用户失败:", err)
	}
	fmt.Printf("查询到用户: %+v\n", foundUser)

	// 根据用户名查询
	fmt.Println("\n根据用户名查询...")
	userByName, err := userModel.FindOneByName(user.Name)
	if err != nil {
		log.Fatal("根据用户名查询失败:", err)
	}
	fmt.Printf("根据用户名查询到用户: %+v\n", userByName)

	// 更新用户
	fmt.Println("\n更新用户...")
	foundUser.Nickname = "updated_user"
	err = userModel.Update(*foundUser)
	if err != nil {
		log.Fatal("更新用户失败:", err)
	}
	fmt.Println("用户更新成功")

	// 再次查询确认更新
	fmt.Println("\n确认更新...")
	updatedUser, err := userModel.FindOne(id)
	if err != nil {
		log.Fatal("查询更新后的用户失败:", err)
	}
	fmt.Printf("更新后的用户: %+v\n", updatedUser)

	// 删除用户
	fmt.Println("\n删除用户...")
	err = userModel.Delete(id)
	if err != nil {
		log.Fatal("删除用户失败:", err)
	}
	fmt.Println("用户删除成功")

	// 确认删除
	fmt.Println("\n确认删除...")
	_, err = userModel.FindOne(id)
	if err != nil {
		if err == model.ErrNotFound {
			fmt.Println("用户已成功删除")
		} else {
			log.Fatal("查询已删除用户时发生错误:", err)
		}
	} else {
		fmt.Println("用户删除失败")
	}
}
