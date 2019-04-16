# coding: utf8
USER_AGENT = (
    "Mozilla/5.0 (Windows NT 6.3; WOW64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/40.0.2214.93 Safari/537.36"
)

RULES = {
    # 每个帖子项
    "topic_item": "//table[@class='olt']/tr",
    "url_list": "//table[@class='olt']/tr/td[@class='title']/a/@href",
    # 列表元素
    "title": "td[@class='title']/a/@title",
    "author": "td[@nowrap='nowrap'][1]/a/text()",
    "reply": "td[@nowrap='nowrap'][2]/text()",
    "last_reply_time": "td[@class='time']/text()",
    "url": "td[@class='title']/a/@href",
    # 帖子详情
    "detail_title_sm": "//td[@class='tablecc']/text()",
    # 完整标题
    "detail_title_lg": "//div[@id='content']/h1/text()",
    "create_time": "//span[@class='color-green']/text()",
    "detail_author": "//span[@class='from']/a/text()",
    "content": "//div[@class='topic-richtext']/p/text()",
    "content2": "//div[@class='topic-richtext']/p/text()",
}

MAX_COUNT = 5
WATCH_INTERVAL = 10 * 60

GROUP_URLS = [
    # 北京租房豆瓣
    "http://www.douban.com/group/26926/",
    # 北京租房（非中介）
    "http://www.douban.com/group/279962/",
    # 北京租房房东联盟(中介勿扰)
    "http://www.douban.com/group/257523/",
    # 北京租房
    "http://www.douban.com/group/beijingzufang/",
    # 北京租房小组
    "http://www.douban.com/group/xiaotanzi/",
    # 北京无中介租房
    "http://www.douban.com/group/zhufang/",
    # 北漂爱合租
    "http://www.douban.com/group/aihezu/",
    # 北京同志们来租房
    "http://www.douban.com/group/325060/",
    # 北京个人租房
    "http://www.douban.com/group/opking/",
    # 北京租房小组!
    "http://www.douban.com/group/374051/",
]

GROUP_SUFFIX = "discussion?start=%d"
COMMUNITY_LIST = [
	"像素",
	"民族家园",
]

REGION_LIST = [
	"黄渠",
	"褡裢坡",
	"常营",
	"草房",
	"物资学院",
	"四惠",
	"四惠东",
	"高碑店",
	"传媒大学",
	"双桥",
        "管庄",
        "八里桥",
        "潘家园",
        "双井",
        "劲松",
        "梨园",
        "通州北苑",
        "果园",
        "九棵树",
        "通州北关",
]

FILTER_LIST = [
	"无中介费",
	"个人房源",
	"房源",
]

POOL_SIZE = 100

