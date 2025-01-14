package baidu

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/util"

	"github.com/sirupsen/logrus"
)

// 百度收录推送
func SiteMapSend(urls []string) {
	// 将url地址转换为字符串，并用换行符分隔
	urlsData := strings.Join(urls, "\r\n")
	// 请求地址
	apiUrl := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", global.Config.Baidu.SiteMapSite, global.Config.Baidu.SiteMapToken)

	// 创建http客户端
	client := &http.Client{}
	// 创建post请求
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(urlsData))
	if err != nil {
		global.Logger.WithError(err).Error("百度收录POST请求创建失败！")
		return
	}
	// 设置请求头
	request.Header.Set("Content-Type", "text/plain")
	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		global.Logger.WithError(err).Error("百度收录发送请求失败！")
		return
	}
	// 关闭响应体
	defer response.Body.Close()

	// 读取响应内容
	body := new(bytes.Buffer)
	_, err = body.ReadFrom(response.Body)
	if err != nil {
		global.Logger.WithError(err).Error("百度收录读取响应内容失败！")
		return
	}
	res, _ := util.JsonTo[map[string]interface{}](body.String())
	// 日志记录
	global.Logger.WithFields(logrus.Fields{
		"reqUrls": urls,
		"res":     res,
		"resBody": body.String(),
	}).Info("百度收录！")

	// 打印响应内容
	fmt.Printf("百度站点收录\r\n推送：%s\r\n响应：%s\r\n", urlsData, body.String())
}
