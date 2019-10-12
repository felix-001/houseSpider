# coding: utf8
DEFAULT_USER_AGENT = (
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"
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


GROUPS = [
    (26926,     u'北京租房豆瓣'),
    (279962,    u'北京租房（非中介）'),
    (262626,    u'北京无中介租房（寻天使投资）'),
    (35417,     u'北京租房'),
    (56297,     u'北京个人租房 （真房源|无中介）'),
    (257523,    u'北京租房房东联盟(中介勿扰) '),
]

KEYWORDS = ( u'天通苑', 
             u'霍营', 
             u'回龙观东大街',
)

TITLE_FILTER_KEYWRODS = [
        u'限女',
        u'求租',
        u'咨询电话',
        u'无中介费',
        u'月付',
        u'随时看房',
        u'押一付一',
        u'押一付',
        u'主卧',
        u'征室友',
        u'同床',
        u'次卧',
        u'女生',
        u'随时入住',
        u'阳隔',
        u'精选房源'
]
