package controllers

import (
	"github.com/prime_beego/beego_api/models"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	ErrNicknameOrPasswd = ErrResponse{403001, "Error Nickname or Password"}
	ErrGenerateToken = ErrResponse{403002, "Error Generate Token"}
	ErrParseToken = ErrResponse{403003, "Error Parse Token"}
	ErrVerifyToken = ErrResponse{403004, "Error Verify Token"}
	ErrAccessDeny = ErrResponse{403004, "Error Access"}
)

type LoginToken struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

//Secret Key
const SecretKey = "jwtprime"

// UserController operations for User
type UserController struct {
	beego.Controller
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Logout", c.Logout)
	c.Mapping("Login", c.LoginHandle)
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Auth", c.AuthVerify)
	c.Mapping("Deny", c.AccessDeny)
}

// @Title 用户登陆
// @Description 用户登陆
// @Param	nickname		formdata 	string		true		"user login nickname"
// @Param	password		formdata 	string		true		"user login password"
// @Failure 404 no enough input
// @Failure 401 No Admin
// @router /login [post]
func (c *UserController) LoginHandle() {
	nickname := c.GetString("nickname")
	password := c.GetString("password")

	user, ok := models.CheckUser(nickname, password)
	if !ok {
		c.Data["json"] = ErrNicknameOrPasswd
		c.ServeJSON()
		return
	}

	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil || tokenString == "" {
		c.Data["json"] = ErrGenerateToken
		c.ServeJSON()
		return
	}

	//login success & response with token
	c.Data["json"] = Response{
		Code: 0,
		Msg: "success",
		Data: LoginToken{
			User: user,
			Token: tokenString,
		},
	}
	c.ServeJSON()
}

// @Title 登陆验证
// @Description 登陆验证
// @Success 200 {object}
// @Failure 401 unauthorized
// @router /auth [get]
func (c *UserController) AuthVerify() {
	bytoken := strings.TrimSpace(c.Ctx.Request.Header.Get("Authorization"))
	token, err := jwt.Parse(bytoken, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Data["json"] = ErrParseToken
		c.ServeJSON()
	}else {
		if token.Valid 	{
			c.Data["json"] = Response{
				Code: 0,
				Msg: "success",
			}
			c.ServeJSON()
		}else {
			c.Data["json"] = ErrVerifyToken
			c.ServeJSON()
		}
	}
}

// @Title 登出
// @Description 登出
// @Success 200 {object}
// @Failure 401 unauthorized
// @router /logout [get]
func (c *UserController) Logout() {
	c.Data["json"] = Response{
		Code: 0,
		Msg: "success",
	}
	c.ServeJSON()
}

// @Title 权限验证失败
// @Description 权限验证失败
// @Success 200 {object}
// @Failure 401 unauthorized
// @router /deny [get]
func (c *UserController) AccessDeny() {
	c.Data["json"] = ErrAccessDeny
	c.ServeJSON()
}

// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	var v models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddUser(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
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
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.User{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateUserById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteUser(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
