package boot

import (
	"fmt"

	"github.com/worklz/yunj-blog-go/app/es"
	"github.com/worklz/yunj-blog-go/pkg/global"

	"github.com/olivere/elastic/v7"
)

// 初始化ES连接
func InitEs() {
	if !global.Config.Elasticsearch.Enable {
		return
	}
	client, err := elastic.NewClient(
		elastic.SetURL(global.Config.Elasticsearch.Hosts...),
		elastic.SetSniff(false), // 禁用 sniffing
	)
	if err != nil {
		panic(fmt.Sprintf("初始化ES连接失败！%v", err))
	}
	global.EsClient = client
	// 注册全部索引
	es.Register()
	// todo
	// err = es.Reset()
	// fmt.Println("ES初始化", err)
}
