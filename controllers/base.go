package controllers

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

//Response struct
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data    interface{} `json:"data"`
}

//Response struct of Error
type ErrResponse struct {
	Errcode int         `json:"errcode"`
	Errmsg  interface{} `json:"errmsg"`
}
