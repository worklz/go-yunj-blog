package key

import (
	"errors"
	"fmt"
	"time"

	adLocation "github.com/worklz/yunj-blog-go/app/enum/ad/location"
	"github.com/worklz/yunj-blog-go/app/enum/ad/status"
	articleStatus "github.com/worklz/yunj-blog-go/app/enum/article/status"
	linkStatus "github.com/worklz/yunj-blog-go/app/enum/link/status"
	"github.com/worklz/yunj-blog-go/app/enum/state"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/util"
)

type RedisKeyConst string

// 定义redis key常量
const (

	// 测试
	TEST RedisKeyConst = "test:"

	// 文章正常已发布的数量
	ARTICLE_STATE_NORMAL_STATUS_PUBLISH_COUNT RedisKeyConst = "article.state:normal.status:publish.count"

	// 文章的浏览数量
	ARTICLE_VIEW_COUNT RedisKeyConst = "article.view.count"

	// 前台菜单html结构
	INDEX_CATEGORY_MENU_HTML_LAYOUT RedisKeyConst = "index.category.menu.html.layout"

	// 前台的链接
	INDEX_LINK_ITEMS RedisKeyConst = "index.link.items"

	// 前台对应位置广告
	INDEX_AD_ITEMS_BY_LOCATION RedisKeyConst = "index.ad.items.by.location:"

	// 日志ip2long对应的ip_id
	LOG_IP_ID_BY_LOG2LONG RedisKeyConst = "log.ip.id.by.log2long:"

	// 日志16位MD5的user agent对应的user_agent id
	LOG_UA_ID_BY_MD5_16 RedisKeyConst = "log.ua.id.by.md5.16:"

	// 日志16位MD5的pageUrl对应的page id
	LOG_PAGE_ID_BY_MD5_16 RedisKeyConst = "log.page.id.by.md5.16:"
)

// 所有常量属性
var AllConstAttrs = map[RedisKeyConst]interface{}{
	TEST: map[string]interface{}{
		"desc": "测试",
		"keySuffix": func(rk *RedisKey) (string, error) {
			if args, ok := rk.Args.(map[string]interface{}); ok {
				if id, exists := args["id"]; exists {
					if suffix, err := util.ToString(id); err == nil {
						return suffix, nil
					} else {
						return "", err
					}
				}
			}
			return "", errors.New("args参数异常")
		},
		"value": func(rk *RedisKey) (interface{}, error) {
			return time.Now().Format("2006-01-02 15:04:05"), nil
		},
	},
	ARTICLE_STATE_NORMAL_STATUS_PUBLISH_COUNT: map[string]interface{}{
		"desc": "文章正常已发布的数量",
		"value": func(rk *RedisKey) (interface{}, error) {
			var count int
			err := global.MySQL.Model(&model.Article{}).Where("state = ? AND status = ?", state.NORMAL, articleStatus.PUBLISH).Count(&count).Error
			return count, err
		},
	},
	ARTICLE_VIEW_COUNT: map[string]interface{}{
		"desc": "文章的浏览数量",
		"value": func(rk *RedisKey) (interface{}, error) {
			type SumRes struct {
				TotalViewCount int64
			}
			var sumRes SumRes
			err := global.MySQL.Model(&model.Article{}).Select("SUM(view_count) as total_view_count").Scan(&sumRes).Error
			return sumRes.TotalViewCount, err
		},
	},
	INDEX_CATEGORY_MENU_HTML_LAYOUT: map[string]interface{}{
		"desc": "前台菜单html结构",
	},
	INDEX_LINK_ITEMS: map[string]interface{}{
		"desc": "前台的链接",
		"value": func(rk *RedisKey) (interface{}, error) {
			var links []model.Link
			err := global.MySQL.
				Model(&model.Link{}).
				Where("status = ? and state = ?", linkStatus.PUBLISH, state.NORMAL).
				Order("sort asc").
				Select("id,link,`desc`").
				Find(&links).
				Error
			if err != nil {
				global.Logger.Error("获取前台的链接数据异常！", err)
			}
			return links, nil
		},
	},
	INDEX_AD_ITEMS_BY_LOCATION: map[string]interface{}{
		"desc": "前台对应位置广告",
		"keySuffix": func(rk *RedisKey) (string, error) {
			location, err := rk.GetArgs()
			if err != nil {
				return "", err
			}
			if suffix, err := util.ToString(location); err == nil {
				return suffix, nil
			} else {
				return "", fmt.Errorf("args参数转换为string类型错误！%v", err)
			}
		},
		"value": func(rk *RedisKey) (interface{}, error) {
			var ads []model.Ad
			location, err := rk.GetArgs()
			if err != nil {
				return ads, err
			}
			err = global.MySQL.
				Model(&model.Ad{}).
				Where("location = ? and status = ? and state = ?", location, status.PUBLISH, state.NORMAL).
				Order("sort asc").
				Select("id,cover,link,`desc`").
				Find(&ads).
				Error
			if err != nil {
				global.Logger.Error("获取前台对应位置广告数据异常！", err)
			}
			return ads, nil
		},
	},
	LOG_IP_ID_BY_LOG2LONG: map[string]interface{}{
		"desc": "日志ip2long对应的ip_id",
		"keySuffix": func(rk *RedisKey) (string, error) {
			ip2long, err := rk.GetArgs()
			if err != nil {
				return "", err
			}
			if suffix, err := util.ToString(ip2long); err == nil {
				return suffix, nil
			} else {
				return "", fmt.Errorf("args参数转换为string类型错误！%v", err)
			}
		},
		"value": func(rk *RedisKey) (interface{}, error) {
			var logIp model.LogIp
			ip2long, err := rk.GetArgs()
			if err != nil {
				return 0, err
			}
			err = global.MySQL.Where("ip2long = ?", ip2long).First(&logIp).Error
			if err != nil {
				global.Logger.Error("获取日志ip2long对应的ip_id数据异常！", err)
			}
			return logIp.Id, nil
		},
	},
	LOG_UA_ID_BY_MD5_16: map[string]interface{}{
		"desc": "日志16位MD5的user agent对应的user_agent id",
		"keySuffix": func(rk *RedisKey) (string, error) {
			md5_16, err := rk.GetArgs()
			if err != nil {
				return "", err
			}
			if suffix, err := util.ToString(md5_16); err == nil {
				return suffix, nil
			} else {
				return "", fmt.Errorf("args参数转换为string类型错误！%v", err)
			}
		},
		"value": func(rk *RedisKey) (interface{}, error) {
			var logUserAgent model.LogUserAgent
			md5_16, err := rk.GetArgs()
			if err != nil {
				return 0, err
			}
			err = global.MySQL.Where("md5_16 = ?", md5_16).First(&logUserAgent).Error
			if err != nil {
				global.Logger.Error("获取日志16位MD5的user agent对应的user_agent id数据异常！", err)
			}
			return logUserAgent.Id, nil
		},
	},
	LOG_PAGE_ID_BY_MD5_16: map[string]interface{}{
		"desc": "日志16位MD5的pageUrl对应的page id",
		"keySuffix": func(rk *RedisKey) (string, error) {
			md5_16, err := rk.GetArgs()
			if err != nil {
				return "", err
			}
			if suffix, err := util.ToString(md5_16); err == nil {
				return suffix, nil
			} else {
				return "", fmt.Errorf("args参数转换为string类型错误！%v", err)
			}
		},
		"value": func(rk *RedisKey) (interface{}, error) {
			var logPage model.LogPage
			md5_16, err := rk.GetArgs()
			if err != nil {
				return 0, err
			}
			err = global.MySQL.Where("md5_16 = ?", md5_16).First(&logPage).Error
			if err != nil {
				global.Logger.Error("日志16位MD5的pageUrl对应的page id数据异常！", err)
			}
			return logPage.Id, nil
		},
	},
}

// 定义所有常量获取参数的规则
var AllConstGetArgs = map[RedisKeyConst]func(*RedisKey, ...interface{}) (interface{}, error){
	// 前台对应位置广告参数获取校验
	INDEX_AD_ITEMS_BY_LOCATION: func(rk *RedisKey, param ...interface{}) (interface{}, error) {
		location, ok := rk.Args.(adLocation.AdLocationConst)
		if !ok {
			return nil, fmt.Errorf("%s配置的args参数异常！", rk.StringValue())
		}
		return location, nil
	},
	// 日志ip2long对应的ip_id
	LOG_IP_ID_BY_LOG2LONG: func(rk *RedisKey, param ...interface{}) (interface{}, error) {
		ip2long, ok := rk.Args.(uint32)
		if !ok {
			return nil, fmt.Errorf("%s配置的args参数异常！", rk.StringValue())
		}
		return ip2long, nil
	},
	// 日志16位MD5的user agent对应的user_agent id
	LOG_UA_ID_BY_MD5_16: func(rk *RedisKey, param ...interface{}) (interface{}, error) {
		md5_16, ok := rk.Args.(string)
		if !ok || md5_16 == "" {
			return nil, fmt.Errorf("%s配置的args参数异常！", rk.StringValue())
		}
		return md5_16, nil
	},
	// 日志16位MD5的pageUrl对应的page id
	LOG_PAGE_ID_BY_MD5_16: func(rk *RedisKey, param ...interface{}) (interface{}, error) {
		md5_16, ok := rk.Args.(string)
		if !ok || md5_16 == "" {
			return nil, fmt.Errorf("%s配置的args参数异常！", rk.StringValue())
		}
		return md5_16, nil
	},
}
