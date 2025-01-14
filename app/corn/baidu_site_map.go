package corn

import (
	"strconv"
	"time"

	"github.com/worklz/yunj-blog-go/app/enum/ad/status"
	"github.com/worklz/yunj-blog-go/app/enum/state"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/pkg/baidu"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/util"
)

// 实现TaskInterface接口
type BaiduSiteMap struct{}

// 名称
func (t *BaiduSiteMap) Name() string {
	return "百度站点收录"
}

// 超时时间
func (t *BaiduSiteMap) Timeout() time.Duration {
	return 10 * time.Second
}

// 任务处理函数
func (t *BaiduSiteMap) Handler() {
	blogUrl := global.Config.App.Url

	// 要收录的地址
	urls := []string{
		blogUrl,
		blogUrl + "/about-us",
		blogUrl + "/category",
	}
	// 所有发布的文章
	var articleIds []int
	global.MySQL.Model(&model.Article{}).Where("status = ? and state = ?", status.PUBLISH, state.NORMAL).Pluck("id", &articleIds)
	for _, articleId := range articleIds {
		urls = append(urls, blogUrl+"/article/"+strconv.Itoa(articleId))
	}
	// 因为百度收录推送有限制，所以次数随机获取5条进行推送
	urls = util.RandomSlice(urls, 5)
	// 推送
	baidu.SiteMapSend(urls)
}
