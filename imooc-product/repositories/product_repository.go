package repositories

import (
	"database/sql"
	"imooc-shop/common"
	"imooc-shop/datamodels"
	"strconv"
)

// 第一步，先开发对应的接口
// 第二步，实现定义的接口

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product)(int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
	SubProductNum(productId int64) error
}

type ProductManager struct {
	table		string
	mysqlConn	*sql.DB
}

func NewProductManager(table string, db *sql.DB) *ProductManager {
	return &ProductManager{
		table:     table,
		mysqlConn: db,
	}
}

func (p *ProductManager) Conn()(err error)  {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}

	return nil
}

// 插入
func (p *ProductManager) Insert(product *datamodels.Product)(productId int64, err error) {
	// 判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}

	sql := "INSERT product SET productName=?, productNum=?, productImg=?, productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}

	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}

	return result.LastInsertId()
}

func (p *ProductManager) Delete(productId int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}

	sql := "DELETE from product WHERE ID=?"
	stm, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}

	_, errStm := stm.Exec(productId)
	if errStm != nil {
		return false
	}

	return true
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	// 判断连接是否存在
	if err := p.Conn(); err != nil {
		return nil
	}

	sql := "UPDATE product set productName=?, productNum=?, productImage=?, productUrl=? where ID="+
		strconv.FormatInt(product.ID, 10)

	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return errSql
	}

	_, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}

	return nil
}

// 根据商品id查询对应的商品
func (p *ProductManager) SelectByKey(productId int64) (productResult *datamodels.Product, err error)  {
	// 1.判断连接是否存在
	if err := p.Conn();err != nil {
		return &datamodels.Product{}, err
	}

	sql2 := "SELECT * FROM " + p.table + " WHERE id = " + strconv.FormatInt(productId, 10)
	row, errRow := p.mysqlConn.Query(sql2)
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}

	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return
}

func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error)  {
	// 1.判断连接是否存在
	if err := p.Conn();err != nil {
		return nil, err
	}

	sql := "SELECT * FROM " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0  {
		return nil,nil
	}

	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}

	return productArray, nil
}

func (p *ProductManager) SubProductNum(productId int64) error  {
	if err := p.Conn();err != nil {
		return err
	}

	sql := "UPDATE " + p.table + " SET productNum=productNum-1 where id="+strconv.FormatInt(productId, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	return err
}

