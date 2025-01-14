package blog

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/worklz/yunj-blog-go/app/enum/article/sortrule"
	"github.com/worklz/yunj-blog-go/app/enum/article/status"
	"github.com/worklz/yunj-blog-go/app/enum/redis/key"
	"github.com/worklz/yunj-blog-go/app/enum/state"
	"github.com/worklz/yunj-blog-go/app/es"
	"github.com/worklz/yunj-blog-go/app/es/index"
	"github.com/worklz/yunj-blog-go/app/model"
	"github.com/worklz/yunj-blog-go/app/service"
	"github.com/worklz/yunj-blog-go/pkg/global"
	"github.com/worklz/yunj-blog-go/pkg/util"
)

type article struct {
	Service
}

var Article *article

// 正常已发布文章数量
func (a *article) TotalByNormalPublish() int64 {
	res, err := key.New(key.ARTICLE_STATE_NORMAL_STATUS_PUBLISH_COUNT).GetCacheInt64(0)
	if err == nil {
		return res
	}
	global.Logger.Error("获取正常已发布文章数量缓存异常！", err)
	return 0
}

// 文章总浏览数量
func (a *article) ViewTotal() int64 {
	res, err := key.New(key.ARTICLE_VIEW_COUNT).GetCacheInt64(0)
	if err == nil {
		return res
	}
	global.Logger.Error("获取文章总浏览数量缓存异常！", err)
	return 0
}

// 热门文章
func (a *article) HotItems() []model.Article {
	var articles []model.Article
	err := global.MySQL.
		Model(&model.Article{}).
		Where("status = ? and state = ?", status.PUBLISH, state.NORMAL).
		Limit(8).
		Order("view_count desc").
		Select("id,category_id,title,cover,view_count,create_time").
		Find(&articles).
		Error
	if err != nil {
		global.Logger.Error("热门文章数据获取异常！", err)
	}
	// 补充封面图
	service.Article.ItemsAppendCover(articles)
	// 格式化输出时间
	a.FormatCreateTime(articles, "2006-01-02 15:04")
	return articles
}

// 格式化输出时间
func (a *article) FormatCreateTime(articles []model.Article, format string) {
	for i, _ := range articles {
		a.FormatSingleCreateTime(&articles[i], format)
	}
}

// 格式化单个输出时间
func (a *article) FormatSingleCreateTime(article *model.Article, format string) {
	article.DisplayCreateTime = time.Unix(article.CreateTime, 0).Format(format)
}

// 分页列表参数
type ArticlePageListParams struct {
	Page        int                           `json:"page" validate:"required,gt=0" message:"[page]参数必须，且为正整数"`
	PageSize    int                           `json:"pageSize" validate:"required,gt=0,lte=20" message:"[pageSize]参数必须，且须在给定范围内"`
	Keywords    string                        `json:"keywords" validate:"omitempty,min=2" message:"请输入至少2个字符关键词"`
	CategoryIds []int                         `json:"categoryIds" validate:"omitempty,nonnegativeIntSlice" message:"[categoryIds]参数需为非负整数数组"`
	SortRule    sortrule.ArticleSortRuleConst `json:"sortRule" validate:"omitempty,oneof=recent hot keywords_score" message:"[sortRule]参数需在给定范围内"`
}

// 分页列表
func (a *article) PageList(paramArgs ...ArticlePageListParams) map[string]interface{} {
	// 获取参数
	params := a.getPageListParams(paramArgs...)

	// 列表数据查询
	var articleCount int64
	var articles []model.Article
	if global.Config.Elasticsearch.Enable {
		articleCount, articles = a.PageListByEs(params)
	} else {
		articleCount, articles = a.PageListByDb(params)
	}

	// 补充封面图
	service.Article.ItemsAppendCover(articles)
	// 格式化输出时间
	a.FormatCreateTime(articles, "2006-01-02 15:04")

	res := map[string]interface{}{
		"page":       params.Page,
		"pageSize":   params.PageSize,
		"pageCount":  0,
		"itemsCount": articleCount,
		"items":      articles,
	}
	if params.Page == 1 {
		res["pageCount"] = math.Ceil(float64(articleCount) / float64(params.PageSize))
	}
	return res
}

// 分页列表请求参数获取
func (a *article) getPageListParams(paramArgs ...ArticlePageListParams) ArticlePageListParams {
	var params ArticlePageListParams
	if len(paramArgs) <= 0 {
		params = ArticlePageListParams{
			Page:        1,
			PageSize:    8,
			Keywords:    "",
			CategoryIds: []int{},
			SortRule:    "",
		}
	} else {
		params = paramArgs[0]
	}
	if params.SortRule == "" {
		params.SortRule = sortrule.RECENT
	}
	return params
}

// 分页列表字段
func (a *article) PageListFields() []string {
	return []string{"id", "category_id", "title", "keywords", "cover", "view_count", "create_time"}
}

// ES分页查询
func (a *article) PageListByEs(params ArticlePageListParams) (int64, []model.Article) {
	var err error

	dslBody := map[string]any{
		"size": params.PageSize,
		"from": (params.Page - 1) * params.PageSize,
	}
	a.setPageListEsDslQuery(dslBody, params)
	a.setPageListEsDslSort(dslBody, params)
	a.setPageListEsDslSource(dslBody, params)

	res, err := es.NewQuery(&index.Article{}).DocSearch(dslBody)
	if err != nil {
		global.Logger.WithError(err).Error("ES分页查询异常！%v", err)
		return 0, nil
	}
	// 查询数量
	articleCount := res.Hits.TotalHits.Value
	// 查询数据
	articles := make([]model.Article, 0, len(res.Hits.Hits))
	var article model.Article
	for _, hit := range res.Hits.Hits {
		article, err = util.JsonTo[model.Article](string(hit.Source))
		if err != nil {
			global.Logger.WithError(err).Error("ES分页查询source数据转换异常！%v", err)
			continue
		}
		// 高亮
		if hit.Highlight != nil {
			for k, v := range hit.Highlight {
				switch k {
				case "title":
					article.Title = strings.Join(v, "")
				case "keywords":
					article.Keywords = strings.Join(v, "")
				case "content":
					article.Desc = strings.Join(v, "")
				}
			}
		}
		articles = append(articles, article)
	}

	return articleCount, articles
}

// 设置ES分页查询条件
func (a *article) setPageListEsDslQuery(dslBody map[string]any, params ArticlePageListParams) {
	mustItems := []map[string]any{
		{
			"term": map[string]any{
				"status": status.PUBLISH,
			},
		},
		{
			"term": map[string]any{
				"state": state.NORMAL,
			},
		},
	}

	if params.Keywords != "" {
		mustItems = append(mustItems, map[string]any{
			"multi_match": map[string]any{
				"query": params.Keywords,
				"fields": []string{
					"title",
					"keywords",
					"content",
				},
			},
		})
	}
	if len(params.CategoryIds) > 0 {
		categoryIds := params.CategoryIds
		var mustItem map[string]any
		if len(categoryIds) == 1 {
			mustItem = map[string]any{
				"term": map[string]any{
					"category_id": categoryIds[0],
				},
			}
		} else {
			mustItem = map[string]any{
				"terms": map[string]any{
					"category_id": categoryIds,
				},
			}
		}
		mustItems = append(mustItems, mustItem)
	}
	dslBody["query"] = map[string]any{
		"bool": map[string]any{
			"must": mustItems,
		},
	}
}

// 设置ES分页查询排序
func (a *article) setPageListEsDslSort(dslBody map[string]any, params ArticlePageListParams) {
	switch params.SortRule {
	case sortrule.RECENT:
		dslBody["sort"] = []map[string]any{
			{
				"create_time": map[string]any{
					"order": "desc",
				},
			},
		}
		break
	case sortrule.HOT:
		dslBody["sort"] = []map[string]any{
			{
				"view_count": map[string]any{
					"order": "desc",
				},
			},
		}
		break
	case sortrule.KEYWORDS_SCORE:
		dslBody["highlight"] = map[string]any{
			"fields": map[string]any{
				"title":    map[string]any{},
				"keywords": map[string]any{},
				"content":  map[string]any{},
			},
		}
		break
	}
}

// 设置ES分页查询字段
func (a *article) setPageListEsDslSource(dslBody map[string]any, params ArticlePageListParams) {
	dslBody["_source"] = a.PageListFields()
}

// 数据库分页查询
func (a *article) PageListByDb(params ArticlePageListParams) (int64, []model.Article) {
	// 组装条件
	query := global.MySQL.
		Model(&model.Article{}).
		Where("status = ? and state = ?", status.PUBLISH, state.NORMAL).
		Offset(params.PageSize * (params.Page - 1)).
		Limit(params.PageSize).Select(strings.Join(a.PageListFields(), ","))
	// 其他条件
	if params.Keywords != "" {
		keywords := "%" + params.Keywords + "%"
		query = query.Where("(title like ? or keywords like ?)", keywords, keywords)
	}
	if len(params.CategoryIds) > 0 {
		categoryIds := params.CategoryIds
		if len(categoryIds) == 1 {
			query = query.Where("category_id = ?", categoryIds[0])
		} else {
			query = query.Where("category_id in (?)", categoryIds)
		}
	}

	// 查询数量
	var articleCount int64
	if params.Page == 1 {
		query.Count(&articleCount)
	}

	// 查询数据
	var articles []model.Article
	// 排序
	if params.SortRule == sortrule.HOT {
		query = query.Order("view_count desc")
	} else {
		query = query.Order("create_time desc")
	}

	err := query.Find(&articles).Error
	if err != nil {
		global.Logger.Error("文章数据库分页查询异常！", err)
	}
	return articleCount, articles
}

// 详情
func (a *article) Detail(id int) (model.Article, error) {
	// todo 此处后期需嵌入ES查询
	article := a.DetailByDb(id)
	if article.Id == 0 {
		return article, errors.New("文章不存在！")
	}
	if article.State != state.NORMAL || article.Status != status.PUBLISH {
		return article, errors.New("文章数据异常！")
	}
	// 补充扩展属性
	service.Article.ItemAttrAppend(&article)
	// 补充显示时间
	a.FormatSingleCreateTime(&article, "2006-01-02 15:04")
	// 关联分类数据
	if article.CategoryId > 0 {
		article.RelatedCategorys = service.Category.GetRelatedCategorysById(article.CategoryId)
	} else {
		article.RelatedCategorys = []model.Category{service.Category.DefaultCategory()}
	}
	// 标签
	article.Tags = strings.Split(article.Keywords, ",")
	// 推荐文章
	err := global.MySQL.Model(&model.Article{}).
		Where("category_id = ? and status = ? and state = ?", article.CategoryId, status.PUBLISH, state.NORMAL).
		Order("view_count desc").Limit(8).
		Find(&article.RecommendArticles).Error
	if err != nil {
		global.Logger.Error("文章详情获取推荐文章异常！", err)
	}
	// 数量自增
	service.Article.ViewCountInc(&article)

	return article, nil
}

// 数据库详情查询
func (a *article) DetailByDb(id int) model.Article {
	article, err := service.Article.ItemById(id)
	if err != nil {
		global.Logger.Error("文章数据库详情查询异常！", err)
		return article
	}
	return article
}
