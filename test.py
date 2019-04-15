import requests
import logging
import gevent
from gevent import monkey; monkey.patch_socket()
from gevent.queue import Queue
from gevent.pool import Pool
import time

kwargs = {
	"headers": {
		"User-Agent": "Mozilla/5.0 (Windows NT 6.3; WOW64) "
			"AppleWebKit/537.36 (KHTML, like Gecko) "
			"Chrome/40.0.2214.93 Safari/537.36",
		"Refer":"www.douban.com",
	},
        "timeout":10,
        "proxies":{ 
            "http":"181.113.17.230:53281"
            }
} 

url = "http://www.baidu.com"

def proxy_test():
    try:
        r = requests.get(url, **kwargs )
        logging.warn( r )
    except Exception as exc:
        logging.warn("%s %d failed! %s  \n", url, 0, str(exc) )

page_queue = Queue()
pool = Pool()
all_greenlet = []

def init_urls():
    for i in range( 1000 ):
        url = "www.baidu.com" + str(i);
        print("init_urls url = ", url )
        page_queue.put( url )
        time.sleep(1)


def consume():
    while 1:
        url = page_queue.get()
        print("url = ", url )
        gevent.sleep(1)

def gevent_queue_test():
    all_greenlet.append( gevent.spawn(init_urls) )
    all_greenlet.append( gevent.spawn(consume) )
    gevent.joinall( all_greenlet )

def GetPage(self, url, timeout=10, retury_num=10):
    """发起HTTP请求
    @url, str, URL
    @timeout, int, 超时时间
    @retury_num, int, 重试次数
    """
    kwargs = {
            "headers": {
                "User-Agent": USER_AGENT,
                "Referer": url
                },
            }
    kwargs["timeout"] = timeout
    resp = None
    for i in range(retury_num):
        try:
            resp = requests.get(url, **kwargs)
            if resp.status_code != 200:
                logging.warn("return code not 200")
            break
        except Exception as exc:
            logging.warn("%s %d failed! %s \n", url, i, str(exc) )
            time.sleep(2)
            continue
    logging.warn( resp.status_code )


def gevent_test():
    url_list = [ "http://www.baidu.com",
            "http://www.google.com",
            "http://www.douban.com",
            "http://www.qq.com",
            "http://www.dog.com",
            "http://www.sougou.com",
            "http://www.xicidaili.com",
            "http://www.renren.com"
            ]

    for url in url_list:




if __name__ == '__main__':
    #proxy_test()
    #gevent_queue_test()
    gevent_test()
