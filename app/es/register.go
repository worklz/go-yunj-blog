package es

import (
	"fmt"

	"github.com/worklz/yunj-blog-go/app/es/index"
)

// 所有索引
var indexs = []IndexInterface{
	&index.Article{},
}

// 注册全部索引
// boot时调用
func Register() {
	var err error
	var query *Query
	for _, index := range indexs {
		query = NewQuery(index)
		if err = query.IndexRegister(); err != nil {
			panic(fmt.Sprintf("ES索引%s注册失败！%v", query.IndexName(), err))
		}
	}
}
