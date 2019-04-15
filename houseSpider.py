#!/usr/bin/env python3
# -*- coding: utf-8 -*-

r'''
    houseSpider.py
'''

from urllib import request,parse
import requests
import re
import sys
import codecs
from lxml import etree
import time
import csv
import logging
import datetime
from utils import ProxyManager
import ast

from config import (
        USER_AGENT, RULES, MAX_COUNT, GROUP_SUFFIX, REGION_LIST,FILTER_LIST,COMMUNITY_LIST, GROUP_URLS
        )

class Parser(object):
    def __init__( self, proxy_manager=None ):
        self.rules = RULES
        self.proxy_manager = proxy_manager

    def set_proxy_manager( self, proxy_manager=None):
        self.proxy_manager = proxy_manager

    def get_proxy_list( self ):
        r = self.GetPage('https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list')
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
                    #kwargs["proxies"] = {
                        #proxy }
                    kwargs["proxies"] = proxy.copy()
                resp = requests.get(url, **kwargs)
                if resp.status_code != 200:
                    raise HTTPError(resp.status_code, url)
                break
            except Exception as exc:
                logging.warn("%s %d failed! %s proxy is %s \n", url, i, str(exc), proxy )
                time.sleep(2)
                continue
        if resp is None:
            raise URLFetchError(url)
        return resp.content.decode("utf8")

    def _get_detail_info(self, url):
        """获取帖子详情
        @html, str, 页面
        """
        html = self.GetPage( url )
        if u"机器人" in html:
            logging.warn("%s 403.html", url)
            return None
        topic = {}
        title = self.extract(self.rules["detail_title_sm"], html) \
            or self.extract(self.rules["detail_title_lg"], html)
        if title is None: 
            return None
        topic["标题"] = title.strip().encode("utf8")
        logging.info(" [ %s ] ", topic["标题"].decode("utf8") )
        topic["创建时间"] = self.extract(
            self.rules["create_time"], html).strip().encode('utf8')
        topic["发布者"] = self.extract(
            self.rules["detail_author"], html).strip().encode('utf8')
        topic["描述"] = '\n'.join(
            self.extract(self.rules["content"], html, multi=True) \
            or self.extract(self.rules["content2"], html, multi=True) ).encode('utf8')
        topic["链接"] = url.encode('utf8')
        return topic

    def _get_page_info(self, topic_list):
        """获取每一页的帖子基本信息
        @topic_list, list, 当前页的帖子项
        """
        topics = []
        # 第一行是标题头,舍掉
        for topic_item in topic_list[1:]:
            topic = {}
            topic["title"] = self.extract(self.rules["title"], topic_item)
            logging.info(" [ %s ] ", topic["title"] )
            topic["author"] = self.extract(self.rules["author"], topic_item)
            topic["reply"] = self.extract(self.rules["reply"], topic_item) or 0
            topic["last_reply_time"] = self.extract(
                self.rules["last_reply_time"], topic_item)
            topic["url"] = self.extract(self.rules["url"], topic_item)
            now = time.time()
            topic["got_time"] = now
            topic["last_update_time"] = now
            topics.append(topic)
        return topics        

    def extract(self, regx, body, multi=False):
        """解析元素,xpath语法
        @regx, str, 解析表达式
        @body, unicode or element, 网页源码或元素
        @multi, bool, 是否取多个
        """
        if isinstance(body, str):
            body = etree.HTML(body)
        res = body.xpath(regx)
        #logging.info("type(res) = %s", type(res))
        #logging.info("res = %s\n", res )
        if multi:
            return res
        return res[0] if res else None

    def CrawlPage( self, html ):
        topic_urls = self.extract(
            self.rules["url_list"], html, multi=True)
        topic_list = self.extract( self.rules["topic_item"], html, multi=True)

        return topic_urls

    def filter( self, topic ):
        valid = False
        if topic is None:
            logging.info("topic is none type")
            return True
        for region in REGION_LIST:
            #if topic["描述"].deocde('utf8') == NoneType:
            #    return True
            if region in topic["描述"].decode('utf8'):
                valid = True
                break
        if valid == False:
            for community in COMMUNITY_LIST:
                if community in topic["描述"].decode('utf8'):
                    valid = True
                    break
        if valid == True:
            for f in FILTER_LIST:
                if f in topic["描述"].decode('utf8'):
                    return True 
        if valid == True:
            now = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()) 
            diff = self.getTimeDiff( now, topic["创建时间"].decode('utf8') )
            if diff > 4.0 :
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

class Spider(object):

    def GetPage(self, url):
        kwargs = {
            "headers": {
                "User-Agent": USER_AGENT,
                "Referer": "http://www.douban.com/"
            },
            "timeout":10,
        }
        try :
            r = requests.get(url, **kwargs )
        except requests.exceptions.RequestException as e:
            print(e)
            return None
        if r.status_code == requests.codes.ok:
            return r.content.decode("utf8")
        else:
            return None

    def GetPageNum( self ):
        r = self.GetPage( self.root_urls[0] )
        logger.log( r.content )


    def SetRootUrls( self, urls ):
        self.root_urls = urls

    pass

    


def main():
    val = 1
    count = 0
    logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(filename)s[line:%(lineno)d] - %(levelname)s: %(message)s')
    first = 1

    sys.stdout = codecs.getwriter("utf-8")(sys.stdout.detach())
    parser = Parser( )

    parser.get_proxy_list()
    proxy_manager = ProxyManager("./proxy_list.txt", 30)
    parser.set_proxy_manager(proxy_manager)
    csv_file = open( 'house.csv', 'w', encoding='utf-8' )

    for url in GROUP_URLS:
        for i in range( MAX_COUNT ):
            base_url = url +  GROUP_SUFFIX + str(i*25)
            page = parser.GetPage( base_url )
            urls = parser.CrawlPage( page )
            for topic_url in urls:
                val = val + 1
                logging.info("progress [ %03d/%03d ] %s", count, val, topic_url )
                topic = parser._get_detail_info( topic_url )
                if parser.filter( topic ) == False:
                    logging.info("valid [ %03d ] %s", count, topic_url )
                    w = csv.DictWriter(csv_file, topic.keys() )
                    if first == 1:
                        w.writeheader()
                        first = 0
                    w.writerow({k:v.decode('utf8') for k,v in topic.items()})
                    count = count + 1
                time.sleep( 1 )

    csv_file.close()



if __name__ == '__main__':
    main()
