#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import logging
import time

USER_AGENT = (
    "Mozilla/5.0 (Windows NT 6.3; WOW64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/40.0.2214.93 Safari/537.36"
)

class HttpUtils(object):
    def __init__( self, referer='' ):
        self.referer = referer
        self.proxies = None

    def get_html( self, url, timeout=10, retry=10 ):
        kwargs = {}
        headers = { 'User-Agent' : USER_AGENT, 'Referer' : self.referer }
        kwargs['headers'] = headers
        kwargs['timeout'] = timeout
        if self.proxies is not None:
            kwargs['proxies'] = self.proxies
        for i in range( retry ):
            try:
                resp = requests.get( url, **kwargs )
                if resp.status_code != 200:
                    logging.error('GET url %s status code : %d', url, resp.status_code )
                    continue

            except Exception as exc:
                logging.error('GET url %s failed, i = %d, error : %s', url, i, str(exc) )
                time.sleep( 2 )
                continue
        if resp is None:
            logging.error('GET url %s failed', url )
            return None
        return resp.content.decode("utf8")

    def set_referer( self, referer ):
        self.referer = referer
    
    def set_proxies( self, proxies ):
        self.proxies = proxies

class ProxyManager( object ):
    def __init__( self, httpUtil ):
        self.httpUtil = httpUtil
        self.initial_proxy_url = 'https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list'

    def get_initial_proxies( self ):
        self.httpUtil.get_html( self.initial_proxy_url )


def main():
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(filename)s[line:%(lineno)d] - %(levelname)s: %(message)s')
    httpUtil = HttpUtils()
    resp = httpUtil.get_html( 'http://www.baidu.com' )
    logging.info("resp is %s", resp )

if __name__ == '__main__':
    main()

