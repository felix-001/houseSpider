#!/usr/bin/env python3
# -*- coding: utf-8 -*-

r'''
    houseSpider.py
'''

from urllib import request,parse
import re

class Parser(object):
    def ParsePage():
        pass
    def ParseDetail():
        pass
    pass

class Writer(object):
    def write():
        pass
    def SetWriteMethod( method ):
        pass
    pass

class Spider(object):
    def GetPage(url):
        pass

    def GetPageNum():
        pass

    pass

class Logger(object):
    def log( val ):
        pass
    pass


def main():
    spider = Spider()
    logger = Logger()
    writer = Writer()
    writer.SetWriteMethod('excel')
    pageNum = spider.GetPageNum()
    logger.log( pageNum )
    for index in range( pageNum ):
        page = spider.GetPage( pageUrl )
        urls = parser.ParsePage( page )
        for url in urls:
            result = parser.ParseDetail()
            writer.write( result )


if __name__ == '__main__':
    main()
