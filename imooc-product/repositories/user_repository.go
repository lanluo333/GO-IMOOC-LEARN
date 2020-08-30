package repositories

import (
	"database/sql"
	"errors"
	"imooc-shop/common"
	"imooc-shop/datamodels"
	"strconv"
)

type IUserRepository interface {
	Conn()	error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

func NewUserRepository(table string, db *sql.DB) *UserManagerRepository {
	return &UserManagerRepository{
		table:     table,
		mysqlConn: db,
	}
}

type UserManagerRepository struct {
	table 	string
	mysqlConn	*sql.DB
}

func (u *UserManagerRepository) Conn() (err error)  {
	if u.mysqlConn == nil {
		mysql, errMysql := common.NewMysqlConn()
		if errMysql != nil {
			return errMysql
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}

	return 

}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error)  {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空")
	}
	if err=u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	sql := "select * from " + u.table + " where userName = ?"
	rows, errRows := u.mysqlConn.Query(sql, userName)
	defer rows.Close()

	if errRows != nil {
		return &datamodels.User{}, errRows
	}

	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)

	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error)  {
	if err = u.Conn(); err != nil {
		return
	}

	sql := "insert into " + u.table + " set nickName=?,userName=?,passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return userId, errStmt
	}

	result, errResult := stmt.Exec(user.NickName, user.UserName, user.HasPassword)
	if errResult != nil {
		return userId, errResult
	}

	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectById(userId int64) (user *datamodels.User, err error)  {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	sql := "select * from " + u.table + " where ID = " + strconv.FormatInt(userId, 10)
	row, errRow := u.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.User{}, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}

