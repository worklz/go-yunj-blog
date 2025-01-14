package save

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/worklz/yunj-blog-go/app/enum/redis/key"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/app/param"
	"github.com/worklz/yunj-blog-go/app/service/blog/api/log"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/util"

	"github.com/jinzhu/gorm"
)

// 日志保存
type logSave struct {
	log.Service
}

var LogSave *logSave

// 执行
func (ls *logSave) Handler(param param.LogRecord) error {
	ipId, err := ls.getIpId(param.Ip)
	if err != nil {
		return err
	}
	userAgentId, err := ls.getUserAgentId(param.UserAgent)
	if err != nil {
		return err
	}
	pageId, err := ls.getPageId(param.PageUrl)
	if err != nil {
		return err
	}
	logId, err := util.SnowflakeId()
	if err != nil {
		return fmt.Errorf("log_id生成失败！%v", err)
	}
	guidInt64, _ := strconv.ParseInt(param.Guid, 10, 64)
	pageViewId, _ := strconv.ParseInt(param.PageViewId, 10, 64)
	logAttrDesc, _ := json.Marshal(map[string]interface{}{
		"referer": param.Referer,
		"title":   param.Title,
	})
	// 组装数据
	var log model.Log
	log.Id = logId
	log.Guid = guidInt64
	log.PageId = pageId
	log.PageViewId = pageViewId
	log.Type = param.Type
	log.UaId = userAgentId
	log.IpId = ipId
	log.CreateTime = param.CreateTime
	var logAttr model.LogAttr
	logAttr.LogId = logId
	logAttr.Desc = string(logAttrDesc)

	// 事务
	err = global.MySQL.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&log).Error; err != nil {
			return err
		}
		if err := tx.Save(&logAttr).Error; err != nil {
			return err
		}
		// 如果所有操作都成功，则返回 nil 以提交事务
		return nil
	})
	// 检查事务结果
	if err != nil {
		return fmt.Errorf("日志数据保存失败！%v", err)
	}
	return nil
}

// 获取ip在数据库中的id
func (ls *logSave) getIpId(ip string) (int, error) {
	log2long, err := util.IpToUint32(ip)
	if err != nil {
		return 0, fmt.Errorf("IP地址转换异常！%s", err)
	}
	redisKey := key.New(key.LOG_IP_ID_BY_LOG2LONG).SetArgs(log2long)
	ipId, err := redisKey.GetCacheInt(0)
	if err != nil {
		return 0, fmt.Errorf("日志IP地址缓存异常！%s", err)
	}
	if ipId <= 0 {
		var logIp model.LogIp
		logIp.Ip = ip
		logIp.Ip2long = log2long
		logIp.CreateTime = time.Now().Unix()
		err = global.MySQL.Create(&logIp).Error
		if err != nil {
			return 0, fmt.Errorf("日志IP地址写入数据库异常！%s", err)
		}
		ipId = logIp.Id
		redisKey.SetCache(ipId, 0)
	}
	return ipId, nil
}

// 获取user_agent在数据库中的id
func (ls *logSave) getUserAgentId(userAgent string) (int, error) {
	md5_16 := util.Md5(userAgent)[8:24]
	redisKey := key.New(key.LOG_UA_ID_BY_MD5_16).SetArgs(md5_16)
	userAgentId, err := redisKey.GetCacheInt(0)
	if err != nil {
		return 0, fmt.Errorf("日志user_agent缓存异常！%s", err)
	}
	if userAgentId <= 0 {
		var logUserAgent model.LogUserAgent
		logUserAgent.Md5_16 = md5_16
		logUserAgent.Content = userAgent
		logUserAgent.CreateTime = time.Now().Unix()
		err = global.MySQL.Create(&logUserAgent).Error
		if err != nil {
			return 0, fmt.Errorf("日志user_agent写入数据库异常！%s", err)
		}
		userAgentId = logUserAgent.Id
		redisKey.SetCache(userAgentId, 0)
	}
	return userAgentId, nil
}

// 获取page在数据库中的id
func (ls *logSave) getPageId(url string) (int, error) {
	md5_16 := util.Md5(url)[8:24]
	redisKey := key.New(key.LOG_PAGE_ID_BY_MD5_16).SetArgs(md5_16)
	pageId, err := redisKey.GetCacheInt(0)
	if err != nil {
		return 0, fmt.Errorf("日志page缓存异常！%s", err)
	}
	if pageId <= 0 {
		var logPage model.LogPage
		logPage.Md5_16 = md5_16
		logPage.Url = url
		err = global.MySQL.Create(&logPage).Error
		if err != nil {
			return 0, fmt.Errorf("日志page写入数据库异常！%s", err)
		}
		pageId = logPage.Id
		redisKey.SetCache(pageId, 0)
	}
	return pageId, nil
}
