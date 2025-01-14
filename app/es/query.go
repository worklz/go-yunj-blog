package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/util"

	"github.com/olivere/elastic/v7"
)

type IndexInterface interface {
	Name() string                 // 索引名称
	Mapping() string              // 映射jso结构
	GetSyncDbDatas() (any, error) // 获取要同步的数据库数据
}

type Query struct {
	Index           IndexInterface
	indexName       string   // 索引名称
	mappingPropKeys []string // 映射字段的keys
}

// 创建一个新的查询query
func NewQuery(index IndexInterface) *Query {
	query := &Query{Index: index}
	return query
}

// 获取索引名称
func (q *Query) IndexName() string {
	if q.indexName != "" {
		return q.indexName
	}
	IndexPrefix := global.Config.Elasticsearch.IndexPrefix
	indexName := q.Index.Name()
	if !strings.HasPrefix(indexName, IndexPrefix) {
		indexName = IndexPrefix + indexName
	}
	q.indexName = indexName
	return q.indexName
}

// 索引注册
// 判断索引是否存在，不存在则创建
func (q *Query) IndexRegister() error {
	// 判断索引是否存在
	exists, err := q.IndexExists()
	if err != nil {
		return err
	}
	if !exists {
		// 不存在则重置索引，并推送所有文档
		if err := q.IndexReset(time.Now().Format("20060102150405")); err != nil {
			return err
		}
	}
	return nil
}

// 判断索引是否存在
func (q *Query) IndexExists() (bool, error) {
	ctx := context.Background()
	exists, err := global.EsClient.IndexExists(q.IndexName()).Do(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// 索引创建，并设置映射
func (q *Query) IndexCreate() error {
	ctx := context.Background()
	_, err := global.EsClient.CreateIndex(q.IndexName()).BodyString(q.Index.Mapping()).Do(ctx)
	return err
}

// 索引删除
func (q *Query) IndexDelete() error {
	ctx := context.Background()
	_, err := global.EsClient.DeleteIndex(q.IndexName()).Do(ctx)
	return err
}

// 索引重置
// 重新创建索引和映射结构，并推送所有文档
// @param bkMark 备份标识
func (q *Query) IndexReset(bkMark string) error {
	// 判断索引是否存在
	exists, err := q.IndexExists()
	if err != nil {
		return err
	}
	if exists {
		// 存在则备份索引
		// 构建 _reindex 请求体
		indexName := q.IndexName()
		targetIndexName := indexName + "_bk_" + bkMark
		// 执行 _reindex 请求
		// reindex 请求在es中会发起一个异步任务，此处需要指定阻塞完成后才返回结果
		_, err := global.EsClient.Reindex().SourceIndex(indexName).DestinationIndex(targetIndexName).WaitForCompletion(true).Do(context.Background())
		if err != nil {
			return err
		}
		// 删除旧索引
		err = q.IndexDelete()
		if err != nil {
			return err
		}
	}
	// 创建新索引
	err = q.IndexCreate()
	if err != nil {
		return err
	}
	// 推送数据（后面改为队列处理，或者协程处理）
	dbDatas, err := q.Index.GetSyncDbDatas()
	if err != nil {
		return fmt.Errorf("获取ES要同步的数据库数据异常！%v", err)
	}
	err = q.DocBatchSave(dbDatas)
	return err
}

// 获取文档映射属性字段的keys
func (q *Query) GetMappingPropKeys() ([]string, error) {
	if q.mappingPropKeys != nil {
		return q.mappingPropKeys, nil
	}
	// 定义一个结构体来解析定义的mapping json 数据
	type Mapping struct {
		Mappings struct {
			Properties map[string]map[string]any `json:"properties"`
		} `json:"mappings"`
	}

	var mapping Mapping
	if err := json.Unmarshal([]byte(q.Index.Mapping()), &mapping); err != nil {
		return nil, fmt.Errorf("mapping json 数据解析异常: %v", err)
	}

	keys := []string{}
	for k := range mapping.Mappings.Properties {
		keys = append(keys, k)
	}
	q.mappingPropKeys = keys
	return q.mappingPropKeys, nil
}

// 获取映射文档
// 将数据库的数据转换成es文档
func (q *Query) GetMappingDoc(dbData any) (map[string]any, error) {
	var data map[string]any
	data, ok := dbData.(map[string]any)
	if !ok {
		dbDataJson, err := util.ToJson(dbData)
		if err != nil {
			return nil, fmt.Errorf("获取映射文档的数据库数据转换异常1！%v", err)
		}
		data, err = util.JsonTo[map[string]any](dbDataJson)
		if err != nil {
			return nil, fmt.Errorf("获取映射文档的数据库数据转换异常2！%v", err)
		}
	}

	keys, err := q.GetMappingPropKeys()
	if err != nil {
		return nil, fmt.Errorf("获取文档映射属性字段的keys！%v", err)
	}
	doc := map[string]any{}
	for _, k := range keys {
		v, exists := data[k]
		if exists {
			doc[k] = v
		}
	}
	return doc, nil
}

// 获取映射文档切片
// 将数据库的数据转换成es文档切片
func (q *Query) GetMappingDocs(dbDatas any) ([]map[string]any, error) {
	var datas []map[string]any
	datas, ok := dbDatas.([]map[string]any)
	if !ok {
		dbDatasJson, err := util.ToJson(dbDatas)
		if err != nil {
			return nil, fmt.Errorf("获取映射文档切片的数据库数据转换异常1！%v", err)
		}
		datas, err = util.JsonTo[[]map[string]any](dbDatasJson)
		if err != nil {
			return nil, fmt.Errorf("获取映射文档切片的数据库数据转换异常2！%v", err)
		}
	}
	docs := []map[string]any{}
	for _, data := range datas {
		doc, err := q.GetMappingDoc(data)
		if err != nil {
			return nil, err
		}
		if len(doc) > 0 {
			docs = append(docs, doc)
		}
	}
	return docs, nil
}

// 文档批量保存
// 有则更新，无则新增
// @param dbDatas 数据库数据切片，如：[]model.Article...
func (q *Query) DocBatchSave(dbDatas any) error {
	docs, err := q.GetMappingDocs(dbDatas)
	if err != nil {
		return fmt.Errorf("文档批量保存错误1！%v", err)
	}
	if len(docs) <= 0 {
		return errors.New("文档批量保存错误2！请传入合法数据")
	}
	docIds := make([]string, 0, len(docs))
	for _, doc := range docs {
		docId, exists := doc["id"]
		if !exists {
			return fmt.Errorf("文档批量保存错误3！文档缺少id字段！")
		}
		docIdStr, err := util.ToString(docId)
		if err != nil {
			return fmt.Errorf("文档批量保存错误，id转换为字符串异常！%v", err)
		}
		docIds = append(docIds, docIdStr)
	}
	if len(docIds) <= 0 {
		return errors.New("文档批量保存错误4！文档缺少id字段！")
	}
	// 获取新增和更新的文档
	hasByteDocs, err := q.GetByteDocsByIds(docIds)
	if err != nil {
		return fmt.Errorf("文档批量保存错误5！%v", err)
	}
	bulkRequest := global.EsClient.Bulk()
	for _, doc := range docs {
		docId := doc["id"]
		docIdStr, err := util.ToString(docId)
		if err != nil {
			return fmt.Errorf("文档批量保存错误6！%v", err)
		}
		var esDoc map[string]any
		for _, hasByteDoc := range hasByteDocs {
			hasDoc, err := util.JsonTo[map[string]any](string(hasByteDoc))
			if err != nil {
				continue
			}
			hasDocId, exists := hasDoc["id"]
			if !exists {
				continue
			}
			if docId == hasDocId {
				esDoc = hasDoc
				break
			}
		}
		var request elastic.BulkableRequest
		if esDoc == nil {
			// 新增
			request = elastic.NewBulkIndexRequest().Index(q.IndexName()).Doc(doc).Id(docIdStr)
		} else {
			// 更新，正常情况下doc只包含想要更新的字段，其他字段将保持不变
			request = elastic.NewBulkUpdateRequest().Index(q.IndexName()).Doc(doc).Id(docIdStr)
		}
		bulkRequest.Add(request)
	}
	// 执行批量操作
	bulkRes, err := bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}
	// 判断失败
	failItems := bulkRes.Failed()
	if len(failItems) > 0 {
		return errors.New(failItems[0].Error.Reason)
	}
	return nil
}

// 文档保存
// 有则更新，无则新增
// @param dbData 数据库数据，如：model.Article
// 使用示例
// err = es.NewQuery(&index.Article{}).DocSave(article)
// fmt.Println("3333 $#v", err)
func (q *Query) DocSave(dbData any) error {
	doc, err := q.GetMappingDoc(dbData)
	if err != nil {
		return fmt.Errorf("ES文档数据获取错误！%v", err)
	}
	docId, exists := doc["id"]
	if !exists {
		return fmt.Errorf("ES文档数据缺少id字段！")
	}
	docIdStr, err := util.ToString(docId)
	if err != nil {
		return fmt.Errorf("ES文档id转换为字符串异常！%v", err)
	}

	exists, err = q.DocExists(docIdStr)
	if err != nil {
		return fmt.Errorf("ES文档数据获取错误！%v", err)
	}
	if exists {
		// 存在，更新
		_, err := global.EsClient.Update().Index(q.IndexName()).Doc(doc).Id(docIdStr).Do(context.Background())
		if err != nil {
			return fmt.Errorf("ES文档数据更新错误！%v", err)
		}
	} else {
		// 不存在，新增
		_, err := global.EsClient.Index().Index(q.IndexName()).BodyJson(doc).Id(docIdStr).Do(context.Background())
		if err != nil {
			return fmt.Errorf("ES文档数据新增错误！%v", err)
		}
	}
	return nil
}

// 文档批量保存
// 有则更新，无则新增
// @param ids 要删除的数据ids
// err := es.NewQuery(&index.Article{}).DocBatchDelete([]string{"1", "4567894654"})
// fmt.Println("1111 $#v", err)
func (q *Query) DocBatchDelete(ids []string) error {
	if len(ids) <= 0 {
		return nil
	}
	bulkRequest := global.EsClient.Bulk()
	for _, id := range ids {
		request := elastic.NewBulkDeleteRequest().Index(q.IndexName()).Id(id)
		bulkRequest.Add(request)
	}
	// 执行批量操作
	bulkRes, err := bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}
	// 判断失败
	failItems := bulkRes.Failed()
	if len(failItems) > 0 {
		for _, failItem := range failItems {
			if failItem.Error != nil && failItem.Error.Reason != "" {
				return fmt.Errorf("ES文档数据删除错误！%v", failItem.Error.Reason)
			}
		}
	}
	return nil
}

// 文档保存
// 有则更新，无则新增
// @param id 文档id
// 使用示例
// err = es.NewQuery(&index.Article{}).DocDelete(article)
// fmt.Println("3333 $#v", err)
func (q *Query) DocDelete(id any) error {
	idStr, err := util.ToString(id)
	if err != nil {
		return fmt.Errorf("ES文档id转换为字符串异常！%v", err)
	}
	_, err = global.EsClient.Delete().Index(q.IndexName()).Id(idStr).Do(context.Background())
	if err != nil && !elastic.IsNotFound(err) {
		return fmt.Errorf("ES文档删除异常！%v", err)
	}
	return nil
}

// 根据id判断文档数据是否存在
func (q *Query) DocExists(id string) (bool, error) {
	res, err := global.EsClient.Get().Index(q.IndexName()).Id(id).Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("根据id获取对应文档查询失败！%v", err)
	}
	return res.Found, nil
}

// 根据id获取对应文档数据
// 使用示例：
// byteDoc, err := es.NewQuery(&index.Article{}).GetByteDocById("1")
// fmt.Println("1111 $#v", err)
// doc, err := util.JsonTo[map[string]any](string(byteDoc))
// fmt.Println("2222 $#v $#v", err, doc)
func (q *Query) GetByteDocById(id string) (json.RawMessage, error) {
	res, err := global.EsClient.Get().Index(q.IndexName()).Id(id).Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("根据id获取对应文档查询失败！%v", err)
	}
	return res.Source, nil
}

// 根据ids批量获取对应字段文档切片数据
func (q *Query) GetByteDocsByIds(ids []string) ([]json.RawMessage, error) {
	// 准备 MultiGetRequest 请求
	items := make([]*elastic.MultiGetItem, 0, len(ids))
	for _, id := range ids {
		item := elastic.NewMultiGetItem().Index(q.IndexName()).Id(id)
		items = append(items, item)
	}
	res, err := global.EsClient.MultiGet().Add(items...).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("根据ids批量获取对应文档查询失败！%v", err)
	}
	docs := make([]json.RawMessage, 0, len(ids))
	for _, v := range res.Docs {
		if !v.Found { // 跳过未找到的文档
			continue
		}
		docs = append(docs, v.Source)
	}
	return docs, nil
}

// 文档数据搜索
// 使用示例如下：
//
//	res, err := es.NewQuery(&index.Article{}).DocSearch(`{
//		"query": {
//			"match_all": {}
//		},
//		"highlight": {
//			"fields": {
//			"*":{}
//			}
//		}
//	}`)
//
// fmt.Printf("4444 %v \r\n", err)
// fmt.Printf("4444 %v \r\n", res.Hits.TotalHits.Value)
func (q *Query) DocSearch(dslBody any) (*elastic.SearchResult, error) {
	res, err := global.EsClient.Search().Index(q.IndexName()).Source(dslBody).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ES文档搜索异常！%v", err)
	}
	return res, nil
}
