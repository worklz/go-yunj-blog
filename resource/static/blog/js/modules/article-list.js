/**
 * 文章列表
 */
layui.define(['jquery', 'yunj'], function (exports) {
    let win = window;
    let doc = document;
    let $ = layui.jquery;

    class ArticleList {

        constructor(args) {
            args = this.objSupp(args, {
                listEl: null,
                page: 0,
                pageCount: 0,
                pageSize: 0,
                keywords: "",
                categoryIds: [],
                sortRule: ""
            });
            this.args = args;

            this.listBoxEl = null;
            this.listLoadingBoxEl = null;
            this.page = args.page;
            this.pageCount = args.pageCount;
            this.pageSize = args.pageSize;
            this.keywords = args.keywords;
            this.categoryIds = args.categoryIds;
            this.sortRule = args.sortRule;

            this.allowRequest = true;
            this.init();
        }

        init() {
            let that = this;
            that.initAttr();
            that.initPageingEl();
            if (that.pageCount > 0 && that.page >= that.pageCount) {
                that.allowRequest = false;
                that.loadingAll();
            } else {
                that.loading();
            }
            that.setEventBind();
        }

        initAttr() {
            let that = this;
            let listBoxEl = $(that.args.listEl);
            if (listBoxEl.length <= 0) throw new Error("列表容器缺失");
            that.listBoxEl = listBoxEl;
        }

        initPageingEl() {
            let that = this;
            if ($(".article-list-loading-box").length <= 0) {
                let listBoxEl = that.listBoxEl;
                listBoxEl.after(`<div class="article-list-loading-box"></div>`);
            }
            that.listLoadingBoxEl = $(".article-list-loading-box");
        }

        loading() {
            let that = this;
            if (!that.allowRequest) return;
            that.loadingStart();
            let apiUrl = "/blog/api/article/list";
            let requestData = {
                page: parseInt(that.page) + 1,
                pageSize: that.pageSize,
                keywords: that.keywords,
                categoryIds: that.categoryIds,
                sortRule: that.sortRule,
            };
            $.ajax({
                url: apiUrl,
                type: "post",
                data: JSON.stringify(requestData),
                contentType: 'application/json; charset=UTF-8',
                cache: false,
                dataType: "json",
                success: function (res) {
                    res.errcode ? that.loadingError(res.msg) : that.loadingEnd(res.data);
                },
                error: function (err) {
                    that.loadingError(err);
                }
            });
        }

        loadingStart() {
            let that = this;
            that.allowRequest = false;
            that.listLoadingBoxEl.html(`<div class="loading-img"><img src="/static/blog/imgs/loading.gif"></div>`);
            that.listLoadingBoxEl.show();
        }

        loadingEnd(data) {
            let that = this;
            let html = "";
            data.items.forEach(article => {
                let descHtml = "";
                if (article.hasOwnProperty("desc") && article.desc) {
                    descHtml = `<p class="desc">
                                    <textarea style="display:none;" name="markdown_doc">${article.desc}</textarea>
                                    <div class="article-desc-markdown" id="article_desc_markdown_${article.id}" style="display: none;"></div>
                                </p>`;
                }
                html += `<article class="excerpt">
                            <a class="focus" href="/article/${article.id}" title="${article.title}">
                                <img class="thumb" data-original="${article.cover}" src="${article.cover}" alt="${article.title}"style="display: inline;"></a>
                            <header>
                                <a class="cat" title="yunj">yunj<i></i></a>
                                <h2><a href="/article/${article.id}" title="${article.title}">${article.title}</a></h2>
                            </header>
                            <p class="meta">
                                <time class="time"><i class="yunj-icon yunj-icon-time"></i> ${article.display_create_time}</time>
                                <span class="views"><i class="yunj-icon yunj-icon-eye"></i> ${article.view_count}</span>
                            </p>
                            <p class="note">${article.title}</p>
                            ${descHtml}
                        </article>`;
            });
            that.listBoxEl.append(html);
            that.renderDesc();
            that.page = data.page;
            that.pageSize = data.pageSize;
            if (that.page === 1) that.pageCount = data.pageCount;
            if (that.page >= that.pageCount) {
                that.loadingAll();
                return;
            }
            that.listLoadingBoxEl.hide();
            that.allowRequest = true;
        }

        loadingError(error) {
            let that = this;
            that.listLoadingBoxEl.html(`<div class="loading-error">${error}</div>`);
            that.listLoadingBoxEl.show();
        }

        loadingAll() {
            let that = this;
            that.listLoadingBoxEl.html(`<div class="loading-all">已加载全部</div>`);
            that.listLoadingBoxEl.show();
        }

        renderDesc() {
            let that = this;
            that.listBoxEl.find(".excerpt>.desc-markdown").each(function () {
                let descMarkdownEl = $(this);
                let descEl = descMarkdownEl.next(".desc");

                let doc = descMarkdownEl.find("textarea[name=doc]").val();
                if (!doc) return true;
                let rawStyleEl = descMarkdownEl.find(".show");
                yunj.markdownToHtml(rawStyleEl.attr("id"), { markdown: doc }).then(res => {
                    let abstract = rawStyleEl.html().replace(/<em>/g, "tag_em_start").replace(/<\/em>/g, "tag_em_end").replace(/<\/?.+?>/g, "")
                        .replace(/tag_em_start/g, "<em>").replace(/tag_em_end/g, "</em>");
                    descEl.html(abstract);
                    descMarkdownEl.remove();
                });
            });

        }

        setEventBind() {
            let that = this;

            let winH = $(win).height();
            $(win).scroll(function () {
                let docH = $(document.body).height();
                let scrollTop = $(win).scrollTop();     //滚动条top
                let bottomRate = (docH - winH - scrollTop) / winH;
                if (bottomRate < 0.3) that.loading();
            });
        }

        varType(data) {
            return Object.prototype.toString.call(data).replace(/^\[object\s(.+)\]$/, '$1').toLowerCase();
        }

        isObj(obj) {
            return this.varType(obj) === 'object' && !obj.length;
        }

        objSupp(obj, ruleObj) {
            for (let attr in ruleObj) {
                if (attr in obj) {
                    if (this.isObj(ruleObj[attr]) && this.isObj(obj[attr])) {
                        ruleObj[attr] = this.objSupp(obj[attr], ruleObj[attr]);
                    } else {
                        ruleObj[attr] = obj[attr];
                    }
                }
            }
            return ruleObj;
        }

    }

    exports('ArticleList', ArticleList);
});