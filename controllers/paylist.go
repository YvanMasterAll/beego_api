package controllers

import (
	"github.com/prime_beego/beego_api/models"
	"strings"

	"github.com/astaxie/beego"
)

var (
	ErrDataNotFound = ErrResponse{403001, "No Data Selected"}
	ErrQueryIsBad = ErrResponse{403002, "Bad Query"}
	ErrParamIsBad = ErrResponse{403003, "Bad Param"}
	ErrQueryResultIsWrong = ErrResponse{403004, "Bad Query Result"}
)

// PaylistController operations for Paylist
type PaylistController struct {
	beego.Controller
}

// URLMapping ...
func (c *PaylistController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GetPaylist", c.GetPaylistByUserID)
}

// @Title Get All
// @Description get Paylist
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Paylist
// @Failure 403
// @router / [get]
func (c *PaylistController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				//c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.Ctx.ResponseWriter.WriteHeader(403)
				c.Data["json"] = ErrQueryIsBad
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllPaylist(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// @Title 查询用户订单
// @Description 查询指定用户所有的订单
// @Param	userid		formData 	int	true 		"用户 ID"
// @Param	productid	formData 	int	false		"产品 ID"
// @Success 200 {object}
// @Failure 403 参数错误：缺失或格式错误
// @router /getpaylist [post]
func (c *PaylistController) GetPaylistByUserID() {
	userid := -1
	productid := -1

	//userid
	if v, err := c.GetInt("userid"); err == nil {
		userid = v
	}
	//productid
	if v, err := c.GetInt("productid"); err == nil {
		productid = v
	}

	if userid == -1 {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = ErrParamIsBad
		c.ServeJSON()
		return
	}

	l, err := models.GetPaylistByUserID(userid, productid)
	if err != nil {
		c.Data["json"] = ErrQueryResultIsWrong
	}else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}
