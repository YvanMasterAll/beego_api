package models

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
)

type Paylist struct {
	Id		int		`json:"id" orm:"column(payid);pk;unique"`
	User 		*int 		`json:"user" orm:"column(user);size(11)"`
	Product		*int		`json:"product" orm:"column(product);size(11)"`
	Createtime    time.Time    	`json:"create_at" orm:"column(createtime);auto_now;type(datetime)"`
}

func (t *Paylist) TableName() string {
	return "paylist"
}

func init() {
	orm.RegisterModel(new(Paylist))
}

// GetAllPaylist retrieves all paylist matches certain condition. Returns empty list if no records exist
func GetAllPaylist(query map[string]string, fields []string, sortby []string, order []string,
offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Paylist))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Paylist
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

//GetPaylistByUserID
func GetPaylistByUserID(userid int, productid int) (ml []interface{}, err error){
	o := orm.NewOrm()
	qs := o.QueryTable(new(Paylist))

	if productid != -1 {
		qs.Filter("product", productid)
	}
	qs = qs.Filter("user", userid)

	var l []Paylist

	if _,err := qs.All(&l); err == nil {
		for _, v := range l {
			ml = append(ml, v)
		}
		return ml, nil
	}

	err = errors.New("Bad Query Result")
	return nil, err
}

