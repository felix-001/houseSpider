
import requests
import gevent
from gevent import monkey; monkey.patch_all()
from gevent import socket
import logging
import time
from gevent.pool import Pool

USER_AGENT = (
    "Mozilla/5.0 (Windows NT 6.3; WOW64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/40.0.2214.93 Safari/537.36"
)

def get_html( url ):
    kwargs = {
        "headers": {
            "User-Agent": USER_AGENT,
        },
    }
    kwargs["timeout"] = 5
    try:
        logging.warn("start url = %s", url )
        resp = requests.get( url, **kwargs )
        logging.warn( "url = %s, resp = %d", url, resp.status_code )
    except Exception as exc:
        logging.warn("%s %s failed!\n", url,  str(exc) )
    logging.warn("exit get_html, %s", url )


def gevent_test():
    urls = ['https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list', 'http://www.baidu.com', 'http://www.example.com', 'http://www.python.org', 'http://www.google.com', 'http://www.sougou.com',
            'http://www.renren.com', ]
    jobs = [gevent.spawn( get_html, url) for url in urls]
    print(jobs)
    gevent.joinall(jobs, timeout=2)

def gevent_test2():
    urls = ['http://www.baidu.com', 'http://www.example.com', 'http://www.python.org', 'http://www.google.com', 'http://www.sougou.com',
            'http://www.renren.com']


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(filename)s[line:%(lineno)d] - %(levelname)s: %(message)s')
    gevent_test()
    while 1:
        time.sleep(2)
