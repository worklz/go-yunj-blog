layui.use(["jquery", "yunj"], function () {
    let win = window;
    let doc = document;
    let $ = layui.jquery;
    let serverPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwijlkdji1/KBWWxxoYil
NOgfTCd+wHazHFlLR/YVwMeuyiahOiDHvOlHOdu3y1hoRbRxyFdIU1C4ci324s/P
g505pmDmVZL/ajP19HSq301l3sn2dWGLuUJ6VisEy3k4GjDx1bvPkmm+vt8Nh3GY
k+qv2dEjJEicrNLmwbwKmKx9zvfBu6xlr/U68BjRMXmeVWMLCs6mHSmeFkXwURu0
dciT7Q72rEdLiY93Nw5qAHR0sV5ulwwMWiY0DmMXTDkgEp4ZoAK3fyMqw08yUvHV
1lhsHcFfQNLwHb/rGe9WCl7ylxjCzO/H7VTWmOrHdHCgzkhL4h5DdyG8vcQn07RG
6QIDAQAB
-----END PUBLIC KEY-----`;

    // 服务端数据加密
    win.serverDataEncrypt = async function (data) {
        let key = yunj.randStr(16);
        let iv = yunj.randStr(16);
        let encrypted = await yunj.aesEncrypt(data, key, iv);
        key = await yunj.rsaEncrypt(key, serverPublicKey);
        iv = await yunj.rsaEncrypt(iv, serverPublicKey);
        return {key: key, iv: iv, encrypted: encrypted};
    };

    // 缓存
    win.cache = {
        set(key, value) {
            localStorage.setItem(key, value);
        },
        get(key) {
            return localStorage.getItem(key);
        },
        clear(key) {
            localStorage.removeItem(key);
        }
    };

    /**
     * 顶部固定
     */
    $.fn.navSmartFloat = function () {
        let position = function (element) {
            let top = element.position().top;
            let pos = element.css("position");
            $(win).scroll(function () {
                let scrolls = $(this).scrollTop();
                if (scrolls > top) {
                    $('.header-topbar').fadeOut(0);
                    if (win.XMLHttpRequest) {
                        element.css({
                            position: "fixed",
                            top: 0,
                            zIndex: 1
                        }).addClass("shadow");
                        let h = element.height() + 'px';
                        $(".replace-header").css({height: h, display: 'block'});
                    } else {
                        $(".replace-header").hide();
                        element.css({top: scrolls});
                    }
                } else {
                    $(".replace-header").hide();
                    $('.header-topbar').fadeIn(500);
                    element.css({
                        position: pos,
                        top: top
                    }).removeClass("shadow")
                }
            })
        };
        return $(this).each(function () {
            position($(this))
        })
    };
    $("#navbar").navSmartFloat();

    /**
     * 回到顶部
     */
    $("#gotop").hide();
    $(win).scroll(function () {
        if ($(win).scrollTop() > 100) {
            $("#gotop").fadeIn()
        } else {
            $("#gotop").fadeOut()
        }
    });
    $("#gotop").click(function () {
        $('html,body').animate({'scrollTop': 0}, 500)
    });

    /**
     * 懒加载图片
     */
    $("img.lazy").lazyload({
        //图片显示时淡入效果
        effect: "fadeIn",
        //没有加载图片时的临时占位符
        placeholder: "/static/blog/imgs/lazyload.png",
        //图片在距离屏幕 200 像素时提前加载.
        threshold: 200,
    });

    /**
     * 侧边栏下拉固定显示
     */
    $(win).scroll(function () {
        let sidebar = $('.sidebar');
        if (sidebar.length <= 0) return;
        let sidebarHeight = sidebar.height();
        let windowScrollTop = $(win).scrollTop();
        let fixedEl = $('.fixed');
        windowScrollTop > sidebarHeight - 60 ? fixedEl.css({
            'position': 'fixed',
            'top': '75px',
            'width': '360px'
        }) : fixedEl.removeAttr("style");

    });

    function keywordsSearch() {
        let keywords = $('input[name=keywords]:visible').val();
        if (keywords.length <= 1) return;
        location.href = "/search?keywords=" + encodeURIComponent(keywords);
    }

    $(doc).on("click", ".btn-search", function (e) {
        keywordsSearch();
        e.stopPropagation(e);
    });

    $(doc).on("keyup", function (e) {
        if (e.keyCode === 13) keywordsSearch();
        e.stopPropagation();
    });

    class GUID {
        constructor() {
            this.cacheKey = "guid";
            this.checkApiUrl = "/blog/api/guid/check";
            this.validApiUrl = "/blog/api/guid/valid";
            this.init();
        }

        init() {
            let that = this;
            let guid = that.get() || "";
            that.check(guid).then(res => {
                let isValid = false;
                switch (res.data.guid) {
                    case "1":
                        break;
                    case "2":
                        isValid = true;
                        break;
                    default:
                        guid = res.data.guid;
                        isValid = true;
                }
                if (isValid && guid) that.valid(guid);
            });
        }

        check(guid) {
            let that = this;
            return new Promise(resolve => {
                yunj.request({
                    url:that.checkApiUrl, 
                    data:JSON.stringify({guid: guid}), 
                    type:"post",
                    contentType: 'application/json; charset=UTF-8',
                }).then(res => {
                    return resolve(res);
                });
            });
        }

        valid(guid) {
            let that = this;
            yunj.request({
                url:that.validApiUrl, 
                data:JSON.stringify({guid: guid}), 
                type:"post",
                contentType: 'application/json; charset=UTF-8',
            }).then(res => {
                that.set(guid);
            });
        }

        get() {
            let that = this;
            return cache.get(that.cacheKey);
        }

        set(guid) {
            let that = this;
            cache.set(that.cacheKey, guid);
        }
    }

    class Log {

        constructor() {
            this.sendApiUrl = "/blog/api/log/record";
            this.init();
        }

        init() {
            let that = this;
            that.setEventBind();
            that.viewLoad();
        }

        // 日志体生成
        generate(args) {
            let that = this;
            return yunj.objSupp(args, {
                page_url: yunj.url(),
                page_view_id:win.pageViewId,
                type: 11,
                referer: "",
                title: ""
            });
        }

        // 推送日志到后端
        send(log){
            let that = this;
            let guid = blogGUID.get() || "";
            if (!guid) return;
            log = that.generate(log);
            log.guid = guid;
            yunj.request({
                url:that.sendApiUrl, 
                data:JSON.stringify(log), 
                type:"post",
                contentType: 'application/json; charset=UTF-8',
            });
        }

        setEventBind(){
            let that = this;

            // 监听页面关闭
            win.onbeforeunload = function (e) {
                that.viewUnload();
            };
        }

        // 页面加载
        viewLoad() {
            let that = this;
            let log = {
                type: 11,
                referer: doc.referrer || "",
                title: doc.title || ""
            };
            that.send(log);
        }

        // 页面卸载
        viewUnload(){
            let that = this;
            let log = {
                type: 22,
                referer: doc.referrer || "",
                title: doc.title || ""
            };
            that.send(log);
        }

    }

    class IndexCommon {

        constructor() {
            this.init();
        }

        init() {
            let that = this;
            win.blogGUID = new GUID();
            win.blogLog = new Log();
            that.setEventBind();
        }

        setEventBind() {
            let that = this;

        }

    }

    win.blogCommon = new IndexCommon();

});