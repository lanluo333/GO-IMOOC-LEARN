package repositories

import (
	"database/sql"
	"fmt"
	"imooc-shop/common"
	"imooc-shop/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

func NewOrderManagerRepository(table string, sql *sql.DB) *OrderManagerRepository {
	return &OrderManagerRepository{
		table:     table,
		mysqlConn: sql,
	}
}

type OrderManagerRepository struct {
	table 	string
	mysqlConn	*sql.DB
}

func (o *OrderManagerRepository) Conn() error  {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}

	if o.table == "" {
		o.table = "order"
	}

	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (id int64, err error)  {
	if err = o.Conn(); err != nil {
		return
	}

	sql := "INSERT `" + o.table + "` SET userId=?, productId=?, orderStatus=?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	fmt.Println(sql)
	fmt.Println(errStmt)
	if errStmt != nil {
		return id, err
	}

	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return id, errResult
	}

	return result.LastInsertId()
}

func (o *OrderManagerRepository) Delete(id int64) bool  {
	if err := o.Conn(); err != nil {
		return false
	}

	sql := "DELETE from " + o.table + "where id = ?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return false
	}

	_, errResult := stmt.Exec(id)
	if errResult != nil {
		return false
	}

	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) error  {
	if errConn := o.Conn(); errConn != nil {
		return errConn
	}

	sql := "UPDATE " + o.table + " set userId = ?, productId=?, orderStatus=? where id = " +
		strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return errStmt
	}

	_, errResult := stmt.Exec(sql, order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return errResult
	}

	return nil
}

func (o *OrderManagerRepository) SelectByKey(id int64) (order *datamodels.Order, err error)  {
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}

	sql := "SELECT * from " + o.table + "where id = " + strconv.FormatInt(id, 10)
	row, errRow := o.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}

	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order, err error)  {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	sql := "SELECT * FROM " + o.table
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}

	return orderArray, nil
}

func (o *OrderManagerRepository) SelectAllWithInfo() (orderMap map[int]map[string]string, err error)  {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	sql := "select o.ID,p.productName,o.orderStatus from imooc.order as o left join product as p on o.productID = p.id"
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}

	return common.GetResultRows(rows), err

}

