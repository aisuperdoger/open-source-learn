package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	userRows                = "`id`,`name`,`password`,`mobile`,`gender`,`nickname`,`create_time`,`update_time`"
	userRowsExpectAutoSet   = "`name`,`password`,`mobile`,`gender`,`nickname`"
	userRowsWithPlaceHolder = "`name`=?,`password`=?,`mobile`=?,`gender`=?,`nickname`=?"

	cacheUserIdPrefix     = "cache:gozero:user:id:"
	cacheUserNamePrefix   = "cache:gozero:user:name:"
	cacheUserMobilePrefix = "cache:gozero:user:mobile:"
)

type (
	UserModel interface {
		Insert(data User) (int64, error)
		FindOne(id int64) (*User, error)
		FindOneByName(name string) (*User, error)
		FindOneByMobile(mobile string) (*User, error)
		Update(data User) error
		Delete(id int64) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id         int64     `db:"id"`
		Name       string    `db:"name"`
		Password   string    `db:"password"`
		Mobile     string    `db:"mobile"`
		Gender     string    `db:"gender"`
		Nickname   string    `db:"nickname"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user`",
	}
}

func (m *defaultUserModel) Insert(data User) (int64, error) {
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)

	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
		return conn.Exec(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname)
	}, userNameKey, userMobileKey)

	if err != nil {
		return 0, err
	}

	id, err := ret.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *defaultUserModel) FindOne(id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByName(name string) (*User, error) {
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, name)
	var resp User
	err := m.QueryRowIndex(&resp, userNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", userRows, m.table)
		if err := conn.QueryRow(&resp, query, name); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByMobile(mobile string) (*User, error) {
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
	var resp User
	err := m.QueryRowIndex(&resp, userMobileKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `mobile` = ? limit 1", userRows, m.table)
		if err := conn.QueryRow(&resp, query, mobile); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) Update(data User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)

	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
	}, userIdKey, userNameKey, userMobileKey)

	return err
}

func (m *defaultUserModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)

	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userIdKey, userNameKey, userMobileKey)

	return err
}

func (m *defaultUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
	return conn.QueryRow(v, query, primary)
}
