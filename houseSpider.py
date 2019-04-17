#!/usr/bin/env python3
# -*- coding: utf-8 -*-

r'''
    houseSpider.py
'''

from urllib import request,parse
import gevent
from gevent.pool import Pool
from gevent.queue import Queue
from gevent import monkey; monkey.patch_all()
import requests
import re
import sys
import codecs
import random
from lxml import etree
import time
import csv
import logging
import datetime
from utils import ProxyManager
import ast
from utils import Timer
from progress import ShowProcess
import os

from config import (
        USER_AGENT, RULES, MAX_COUNT, GROUP_SUFFIX, REGION_LIST,FILTER_LIST,COMMUNITY_LIST, GROUP_URLS, POOL_SIZE,WATCH_INTERVAL
        )

class HTTPError(Exception):

    """ HTTP状态码不是200异常 """

    def __init__(self, status_code, url):
        self.status_code = status_code
        self.url = url

    def __str__(self):
        return "%s HTTP %s" % (self.url, self.status_code)


class URLFetchError(Exception):

    """ HTTP请求结果为空异常 """

    def __init__(self, url):
        self.url = url

    def __str__(self):
        return "%s fetch failed!" % self.self.url

class Parser(object):
    def __init__( self,  filestr="house.csv", proxy_manager=None ):
        self.rules = RULES
        self.proxy_manager = proxy_manager
        self.pool = Pool(size=POOL_SIZE)
        self.page_queue = Queue()
        self.topic_queue = Queue()
        self.group_list = GROUP_URLS
        self.valid_topic = 0
        self.total_crawed_topic = 0
        self.csv_file = open( filestr, 'w', encoding='utf-8' )
        self.interval = WATCH_INTERVAL
        self.count = 0
        self.first = 1
        self.total = len(self.group_list) * MAX_COUNT * 25 
        self.process_bar = ShowProcess( self.total, 'OK')
#        self.f = open('./output.txt', 'w+')
        logging.info("total : %d", self.total )

    def set_proxy_manager( self, proxy_manager=None):
        self.proxy_manager = proxy_manager

    def _init_page_tasks(self, group_url):
        """初始化页面任务
        @group_url, str, 小组URL
        """
        for page in range(MAX_COUNT):
            base_url = "%s%s" % (group_url, GROUP_SUFFIX)
            url = base_url % (page * 25)
            #logging.info("url is %s", url )
            self.page_queue.put(url)

    def _page_loop(self):
        """page loop
        """
        while 1:
            #logging.info( "page_loop")
            page_url = self.page_queue.get(block=True)
            #logging.info("page_url = %s", page_url )
#            gevent.sleep(1)
            self.pool.spawn(self.CrawlPage, page_url)

    def run( self ):
        all_greenlet = []
        # 定时爬取
        for group_url in self.group_list:
            # timer = Timer(random.randint(0, self.interval), self.interval)
            """
            当一个greenlet遇到IO操作时，比如访问网络，就自动切换到其他的greenlet，
            等到IO操作完成，再在适当的时候切换回来继续执行。由于IO操作非常耗时，经常使程序处于等待状态，
            有了gevent为我们自动切换协程，就保证总有greenlet在运行，而不是等待IO。
            """
            timer = Timer(random.randint(0, 2), self.interval)
            greenlet = gevent.spawn(
                timer.run, self._init_page_tasks, group_url)
            all_greenlet.append(greenlet)
        all_greenlet.append(gevent.spawn(self._page_loop))
        all_greenlet.append(gevent.spawn(self._topic_loop))
        gevent.joinall(all_greenlet)

    def get_proxy_list( self ):
        r = self.GetPage('https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list', 1, 1 )
        if r is None:
            return None

        f = open('proxy.list', 'w')
        f.write(r)
        f.close()

        f2 = open('proxy.list','r')
        txt = open('proxy_list.txt', 'w')
        for line in open('proxy.list'):
            try :
                line = f2.readline()
                item = ast.literal_eval( line )
                url = "{ \"" + item["type"] + "\"" + ":" + "\"" + item["type"] + "://" + item["host"] + ":" + str(item["port"]) + "\"}" +"\n"
                txt.write(url)
            except Exception as e:
                logging.info(e)
                logging.info(line)
                logging.info(item)
                logging.info(url)
        txt.close()
        f2.close()


    def GetPage(self, url, timeout=10, retury_num=10):
        """发起HTTP请求
        @url, str, URL
        @timeout, int, 超时时间
        @retury_num, int, 重试次数
        """
        kwargs = {
            "headers": {
                "User-Agent": USER_AGENT,
                "Referer": "http://www.douban.com/"
            },
        }
        kwargs["timeout"] = timeout
        resp = None
        for i in range(retury_num):
            try:
                # 是否启动代理
                if self.proxy_manager is not None:
                    proxy = self.proxy_manager.get_proxy()
                    kwargs["proxies"] = proxy.copy()
                resp = requests.get(url, **kwargs)
                if resp.status_code != 200:
                    raise HTTPError(resp.status_code, url)
                break
            except Exception as exc:
                """
                if self.proxy_manager is not None:
                    logging.warn("%s %d failed! %s proxy is %s \n", url, i, str(exc), proxy )
                else:
                    logging.warn("%s %d failed! %s \n", url, i, str(exc) )
                """
                time.sleep(2)
                continue

        if resp is None:
            logging.info("fetch url : %s error", url )
            return None
        return resp.content.decode("utf8")

    def _get_detail_info(self, url):
        """获取帖子详情
        @html, str, 页面
        """
        html = self.GetPage( url )
        if html is None :
            logging.info("fetch url : %s error", url )
            return None

        if u"机器人" in html:
            logging.warn("%s 403.html", url)
            return None

        #logging.info("progress [ valid: %03d crawed: %03d ] %s", self.valid_topic, self.total_crawed_topic, url )
        self.process_bar.show_process( self.total_crawed_topic )
        topic = {}
        title = self.extract(self.rules["detail_title_sm"], html) \
            or self.extract(self.rules["detail_title_lg"], html)
        if title is None: 
            return None
        #logging.info("tile = %s", title.strip() )
        topic["标题"] = title.strip().encode("utf8")
        topic["创建时间"] = self.extract(
            self.rules["create_time"], html).strip().encode('utf8')
        topic["发布者"] = self.extract(
            self.rules["detail_author"], html).strip().encode('utf8')
        topic["描述"] = '\n'.join(
            self.extract(self.rules["content"], html, multi=True) \
            or self.extract(self.rules["content2"], html, multi=True) ).encode('utf8')
        topic["链接"] = url.encode('utf8')
        if self.filter( topic ) == False:
            logging.info("valid url : %s", topic["链接"].decode('utf8') )
#            self.f.write(topic)
            w = csv.DictWriter(self.csv_file, topic.keys() )
            if self.first == 1:
                w.writeheader()
                self.first = 0
            #logging.info("writing to csv")
            w.writerow({k:v.decode('utf8') for k,v in topic.items()})
            self.valid_topic += 1
        self.total_crawed_topic += 1
        return topic

    def extract(self, regx, body, multi=False):
        if isinstance(body, str):
            body = etree.HTML(body)
        res = body.xpath(regx)
        #logging.info("type(res) = %s", type(res))
        #logging.info("res = %s\n", res )
        if multi:
            return res
        return res[0] if res else None

    def _topic_loop(self):
        """topic loop
        """
        while 1:
            topic_url = self.topic_queue.get(block=True)
            #logging.info("topic loop topic_url = %s [%03d] ", topic_url, self.count )
            self.count += 1
            self.pool.spawn(self._get_detail_info, topic_url)

    def _init_topic_tasks(self, topic_urls):
        for url in topic_urls:
            self.topic_queue.put(url)

    def CrawlPage( self, url ):
        html = self.GetPage( url )
        if html is None:
            return None
        topic_list = self.extract( self.rules["topic_item"], html, multi=True)
        url_list = []
        for topic in topic_list[1:]:
            if topic[0][0].tag is None:
                logging.info("check tag error, url : %s", url )
                continue
            if topic[0][0].text is None:
                logging.info("check text error, url : %s", url )
                continue
            if topic[0][0].get('href') is None:
                logging.info("check href error, url : %s", url )
                continue
            for region in REGION_LIST:
                if region in topic[0][0].text:
                    url_list.append( topic[0][0].get('href'))

        self._init_topic_tasks( url_list )
        return url_list

    def filter( self, topic ):
        valid = False
        if topic is None:
            logging.info("topic is none type")
            return True

        for f in FILTER_LIST:
            if f in topic["描述"].decode('utf8'):
                return True 
        if valid == True:
            now = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()) 
            diff = self.getTimeDiff( now, topic["创建时间"].decode('utf8') )
            if diff > 10.0 :
                #logging.info("diff : %f, %s", diff, topic["链接"] )
                return True;

        if valid == False:
            return True

        return False

    def getTimeDiff(self, timeStra, timeStrb ):
        if timeStra<=timeStrb:
            return 0
        ta = time.strptime(timeStra, "%Y-%m-%d %H:%M:%S")
        tb = time.strptime(timeStrb, "%Y-%m-%d %H:%M:%S")
        y,m,d,H,M,S = ta[0:6]
        dataTimea=datetime.datetime(y,m,d,H,M,S)
        y,m,d,H,M,S = tb[0:6]
        dataTimeb=datetime.datetime(y,m,d,H,M,S)
        secondsDiff=(dataTimea-dataTimeb).total_seconds()
        return (secondsDiff/60/60/24)

class Writer(object):
    def write():
        pass
    def SetWriteMethod( self,  method ):
        pass
    pass

    
def main():
    logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(filename)s[line:%(lineno)d] - %(levelname)s: %(message)s')
    sys.stdout = codecs.getwriter("utf-8")(sys.stdout.detach())

    parser = Parser( )
    logging.info("get prox list...")
    parser.get_proxy_list()
    logging.info("processing proxy list")
    proxy_manager = ProxyManager("./proxy_list.txt", 30)
    parser.set_proxy_manager(proxy_manager)
    logging.info("spider run...")
    parser.run()



if __name__ == '__main__':
    main()

