
import requests
import logging
import config
import proxymanager
import codecs
import sys

"""
referer的作用，防盗链，比如www.google.com页面上放了一个baidu的链接
那么，点击这个链接，访问百度时，这个referer就是www.google.com，表示
是从哪点过来的，如果一个网站只允许在自己的网站访问本网站的图片，就
可以判断referer是不是自己的域名，如果不是，就不允许访问，就是防盗链
"""
def get_html_by_url(url, params=None, referer='www.baidu.com', proxies=None, timeout=10 ):
    headers= {
            'User-Agent':config.DEFAULT_USER_AGENT,
            "Referer":referer
    }
    resp = None

    try:
        resp = requests.get(url=url, params=params, proxies=proxies, headers=headers, timeout=timeout )
    except Exception as exc:
        #logging.info("fetch url:%s fail, params:%s error:%s", url, str(params), str(exc) )
        return None

    if resp.status_code == 200:
        return resp.content.decode('utf8')
    else:
        #logging.error("check status_code error:%d", resp.status_code )
        return None

def get_html_with_proxy(url, referer, params=None, timeout=10):
    for i in range(5):
        proxy = proxymanager.get_one_proxy()
        if proxy is None:
            logging.error("get proxy error")
            return None
        proxymanager.set_proxy_used(proxy)
        proxy_url = proxy['type'] + '://'+proxy['host']+':'+str(proxy['port'])
        proxies = { proxy['type']: proxy_url}
        resp = get_html_by_url(url, params, referer, proxies, timeout)
        if resp is not None:
            proxymanager.set_proxy_unused(proxy)
            break
        else:
            #logging.info("remove not available proxy:%s", str(proxy))
            pass
    return resp

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO, format='%(asctime)s|%(filename)s:%(lineno)d|%(levelname)s: %(message)s')
    sys.stdout = codecs.getwriter("utf-8")(sys.stdout.detach())
    params = {'cat': 1013, 'sort': 'time', 'group': 257523, 'q': '天通苑'}
    proxies = {'https':'https://167.71.97.196:3128'}
    s = get_html_by_url('http://www.douban.com/group/search', params, 'www.douban.com', proxies)
    print(len(s))

