package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/redis"
	"github.com/worklz/yunj-blog-go/pkg/util"

	"github.com/jinzhu/gorm"
)

const (
	// guid暂存过期时间
	GUID_TEMP_EXPIRE_TIME int = 10
	// guid暂存标识
	GUID_TEMP_FLAG string = "1"
)

type guid struct {
	Service
}

var Guid *guid

// 校验参数
type GuidCheckParams struct {
	Guid string `json:"guid"`
}

// 验证参数
type GuidValidParams struct {
	Guid string `json:"guid"`
}

// 校验id数据并返回int64类型id
func (g *guid) CheckId(id interface{}) (int64, error) {
	if idInt64, ok := id.(int64); ok {
		return idInt64, nil
	}
	idStr, ok := id.(string)
	if !ok {
		return 0, errors.New("错误格式guid")
	}
	idInt64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("错误格式guid！%v", err)
	}
	return idInt64, nil
}

// 校验有效性
func (g *guid) Check(id interface{}) (int64, error) {
	idInt64, err := g.CheckId(id)
	if err != nil {
		return 0, err
	}
	if idInt64 > 0 {
		// 判断是否存在
		exists, err := g.IsExist(idInt64)
		if err != nil {
			return 0, err
		}
		if exists {
			return 1, nil
		}
		// 判断是否临时存在
		tempExists, err := g.IsTempExist(idInt64)
		if err != nil {
			return 0, err
		}
		if tempExists {
			return 2, nil
		}
	}
	return 0, nil
}

// 是否存在
func (g *guid) IsExist(id interface{}) (bool, error) {
	idInt64, err := g.CheckId(id)
	if err != nil {
		return false, err
	}

	if idInt64 <= 0 {
		return false, nil
	}
	// 数据库层面判断
	var guid model.Guid
	err = global.MySQL.Where("id = ?", idInt64).First(&guid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, fmt.Errorf("guid数据查询失败！%v", err)
		}
	}
	return guid.Id > 0, nil
}

// 是否临时存在
func (g *guid) IsTempExist(id interface{}) (bool, error) {
	idInt64, err := g.CheckId(id)
	if err != nil {
		return false, err
	}

	if idInt64 <= 0 {
		return false, nil
	}
	idStr, err := util.ToString(idInt64)
	if err != nil {
		return false, fmt.Errorf("guid是否临时存在判断异常！%v", err)
	}
	exists, err := redis.Exists(idStr)
	if err != nil {
		return false, fmt.Errorf("guid是否临时存在redis判断异常！%v", err)
	}
	if !exists {
		return false, nil
	}
	flag, err := redis.GetString(idStr)
	if err != nil {
		return false, fmt.Errorf("guid是否临时存在redis获取异常！%v", err)
	}
	return flag == GUID_TEMP_FLAG, nil
}

// 是否临时存在
func (g *guid) Generate() (int64, error) {
	id, err := util.SnowflakeId()
	if err != nil {
		return 0, fmt.Errorf("临时guid生成失败！%v", err)
	}
	idStr, err := util.ToString(id)
	if err != nil {
		return 0, fmt.Errorf("临时guid生成异常！%v", err)
	}
	err = redis.SetEx(idStr, GUID_TEMP_EXPIRE_TIME, GUID_TEMP_FLAG)
	if err != nil {
		return 0, fmt.Errorf("临时guid生成缓存异常！%v", err)
	}
	return id, nil
}

// 保存
func (g *guid) Save(id interface{}) error {
	idInt64, err := g.CheckId(id)
	if err != nil {
		return err
	}
	// 保存数据
	guid := model.Guid{Id: idInt64, CreateTime: time.Now().Unix()}
	// 删除临时缓存
	idStr, err := util.ToString(id)
	if err != nil {
		return fmt.Errorf("guid保存参数无效！%v", err)
	}
	err = redis.Del(idStr)
	if err != nil {
		return fmt.Errorf("guid保存缓存删除无效！%v", err)
	}
	// todo 保存在ES
	// 保存数据库
	if err = global.MySQL.Save(&guid).Error; err != nil {
		return fmt.Errorf("guid保存数据库失败！%v", err)
	}
	return nil
}
