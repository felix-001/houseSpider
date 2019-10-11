import ast
import httpreq
import logging
import time

g_proxy_list = []

def get_high_anonymous_proxy_list():
    proxies = httpreq.get_html_by_url('https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list')
    if proxies is None:
        return None
    # 使用'\n' 切分字符串,得到一个list
    lines = proxies.splitlines(False)
    for line in lines:
        # 只保留高密代理,https的一定是高密的
        if ( line.find('https') > 0 ):
            try:
                proxy = ast.literal_eval(line)
                proxy['used'] = False
                g_proxy_list.append( proxy )
            except Exception as e:
                pass

def get_unused_proxy():
    for proxy in g_proxy_list:
        if ( proxy['used'] == False):
            return proxy
    return None

def get_one_proxy():
    for i in range(10):
        proxy = get_unused_proxy()
        if (proxy):
            return proxy
        logging.warn("no available proxy,sleep 1s to retry")
        time.sleep(1)
    logging.error("no available proxy, retry 10 times")
    return None

def set_proxy_unused(proxy):
    proxy['used'] = False

def set_proxy_used(proxy):
    proxy['used'] = True

def remove_proxy(proxy):
    g_proxy_list.remove(proxy)

