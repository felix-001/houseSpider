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
	"龙跃苑",
	"龙腾苑",
	"龙锦苑",
	"龙博苑",
	"龙泽苑",
	"龙华园",
	"龙禧苑",
	"龙泽苑西区",
	"龙泽苑东区",
	"龙回苑",
	"龙城嘉园",
	"龙冠润景",
	"新龙城",
	"天龙苑",
	"矩阵小区",
	"天慧园",
	"云趣园",
	"风雅园",
	"佰嘉城",
	"流星花园",
	"和谐家园",
	"天露园",
	"通达园",
	"田园风光雅园",
	"三合庄园",
	"天鑫家园",
	"宽HOUSE",
	"慧华苑",
	"北店嘉园",
	"东亚上北",
	"上坡佳园",
	"东村家园",
	"骊龙园",
	"北回归线",
	"龙兴园北区",
	"北京人家",
	"北郊农场",
	"冠庭园",
]

REGION_LIST = [
	"西二旗",
	"龙泽",
	"回龙观",
	"上地",
	"回龙观东大街",
	"回龙观西大街",
	"五道口",
	"平西府",
	"育新",
]

FILTER_LIST = [
	"无中介费",
	"个人房源",
	"房源",
]

POOL_SIZE = 100

