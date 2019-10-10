#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
豆瓣租房爬虫

Install:
    pip install beautifulsoup4

Usage:
    python douban_zufang.py
"""

from __future__ import print_function
import asyncio
import copy
from datetime import datetime, timedelta
from functools import partial

from bs4 import BeautifulSoup
import requests


# 有些帖子已经看过了，手动把 url 加入黑名单
post_url_black_list = []
try:
    with open('url_blacklist.txt', 'r') as f:
        for line in f.readline():
            url = f.readline().split('\n')[0]
            post_url_black_list.append(url)
except IOError:
    pass


expected_groups = [
    (26926, u'北京租房豆瓣'),
    (279962, u'北京租房（非中介）'),
    (262626, u'北京无中介租房（寻天使投资）'),
    (35417, u'北京租房'),
    (56297, u'北京个人租房 （真房源|无中介）'),
    (257523, u'北京租房房东联盟(中介勿扰) '),
]
expected_query_strs = (u'立水桥', u'天通苑', u'霍营', u'回龙观东大街', )


group_search_url = 'http://www.douban.com/group/search'
default_params = {
    'cat': 1013,  # 不要问我有啥用，不知道...
    'sort': 'time',
}
default_headers = {
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) '
                  'AppleWebKit/537.36 (KHTML, like Gecko) '
                  'Chrome/61.0.3163.100 Safari/537.36',
}


def gen_search_params(group_id, q):
    params = copy.deepcopy(default_params)
    params.update(dict(
        group=group_id,
        q=q
    ))
    return params


def parse_html_tr(tr):
    time_tag = tr.find('td', {'class': 'td-time'})
    time_str = time_tag['title']
    last_update_time = datetime.strptime(time_str, '%Y-%m-%d %H:%M:%S')
    if datetime.now() - last_update_time > timedelta(days=7):
        return None

    a_tag = tr.td.a
    url = a_tag.get('href', '')
    title = a_tag.text
    return (url, title)


async def search_group(group, q, all_posts):
    event_loop = asyncio.get_event_loop()
    posts = []

    group_id, group_alias = group
    params = gen_search_params(group_id, q)
    response = await event_loop.run_in_executor(
        None,
        partial(requests.get, url=group_search_url, params=params, headers=default_headers)
    )
    if response.status_code == 200:
        html_doc = response.content
        soup = BeautifulSoup(html_doc, 'html.parser')
        results = soup.find_all('tr', {'class': 'pl'})
        for index, result in enumerate(results):
            rv = parse_html_tr(result)
            if rv is not None:
                url, title = rv
                if url in post_url_black_list:
                    continue
                if q in title:
                    posts.append((url, title))

    else:
        print(u'爬虫出现了未知问题...')

    all_posts.extend(posts)
    #print(u'在 『{group_alias}』 查找到了 {count} 个包含 『{q}』 的帖子'
    #      .format(group_alias=group[1], count=len(posts), q=q))


async def search_groups():
    """查找符合要求的帖子

    1. 包含我们查询的字符
    2. 时间是 7 天以内的
    """
    event_loop = asyncio.get_event_loop()

    all_posts = []
    tasks = []
    for group in expected_groups:
        for q in expected_query_strs:
            tasks.append(event_loop.create_task(search_group(group, q, all_posts)))
    done, pending = await asyncio.wait(tasks)
    assert bool(pending) is False
    return all_posts


def uniq_and_sort(posts):
    url_title_map = dict(posts)
    post_weight_map = {}  # {url: weight}
    title_set = set()

    for post in posts:
        url, title = post

        if any([u'马泉营' in title,
                u'孙河' in title,
                u'顺义' in title,
                u'石门' in title,
                u'特价房源' in title,
                u'限女' in title,
                u'求租' in title,
                u'咨询电话' in title,
                u'无中介费' in title,
                u'月付' in title,
                u'随时看房' in title,
                u'押一付一' in title,
                u'押一付' in title,
                u'主卧' in title,
                u'征室友' in title,
                u'同床' in title,
                u'次卧' in title,
                u'女生' in title,
                u'随时入住' in title,
                u'阳隔' in title,
                u'法信' in title]):
            continue

        if title in title_set:
            continue
        else:
            title_set.add(title)

        if url in post_weight_map:
            post_weight_map[url] += 1
        else:
            post_weight_map[url] = 1

    sorted_result = sorted(
        [(_url, weight) for (_url, weight) in post_weight_map.items()],
        key=lambda each: each[1],
        reverse=True
    )
    for url, weight in sorted_result:
        title = url_title_map[url]
        yield (url, title)


async def main():
    posts = await search_groups()
    for url, title in uniq_and_sort(posts):
        print('- [ ', title, ' ](', url, ')')


if __name__ == '__main__':
    event_loop = asyncio.get_event_loop()
    event_loop.run_until_complete(main())
