import sys
import asyncio
import logging
import codecs
import config
import proxymanager
import httpreq

from functools import partial
from bs4 import BeautifulSoup
from datetime import datetime, timedelta

dbg_get_url_err_count = 0

def log_init():
    logging.basicConfig(level=logging.INFO, format='%(asctime)s|%(filename)s:%(lineno)d|%(levelname)s: %(message)s')
    sys.stdout = codecs.getwriter("utf-8")(sys.stdout.detach())

def gen_search_params(group_id, q):
    params = {
            'cat': 1013,
            'sort': 'time',
            }
    params.update(dict(
        group=group_id,
        q=q
        ))
    return params

def get_douban_html(url, group, q):
    params = None

    if group is not None:
        group_id, group_alias = group
        params = gen_search_params(group_id, q)
    return httpreq.get_html_with_proxy(url, "www.douban.com", params )

def parse_search_result_tr(tr):
    time_tag = tr.find('td', {'class': 'td-time'})
    time_str = time_tag['title']
    last_update_time = datetime.strptime(time_str, '%Y-%m-%d %H:%M:%S')
    # 只要3天内的
    if datetime.now() - last_update_time > timedelta(days=3):
        return None

    a_tag = tr.td.a
    url = a_tag.get('href', '')
    title = a_tag.text
    return (url, title)

def parse_search_result(html):
    posts = []

    soup = BeautifulSoup(html, 'html.parser')
    results = soup.find_all('tr', {'class': 'pl'})
    for index, result in enumerate(results):
        rv = parse_search_result_tr(result)
        if rv is not None:
            url, title = rv
            posts.append( (url, title) )
    return posts

async def search_douban_group_by_query(group, q, posts):
    event_loop = asyncio.get_event_loop()
    html = await event_loop.run_in_executor(
        None,
        partial( get_douban_html, url='http://www.douban.com/group/search', group=group, q=q) 
    )
    if html is None:
        logging.info("fetch group:%s q:%s error", str(group), str(q))
        return None
    res = parse_search_result( html )
    posts.extend(res)

async def search_douban_groups():
    all_posts = []
    tasks = []

    event_loop = asyncio.get_event_loop()
    for group in config.GROUPS:
        for q in config.KEYWORDS:
            tasks.append(event_loop.create_task(search_douban_group_by_query(group, q, all_posts)))
    done, pending = await asyncio.wait(tasks)
    assert bool(pending) is False
    return all_posts

def filter_posts(posts):
    for post in posts[:]:
        url,title = post
        for keyword in config.TITLE_FILTER_KEYWRODS:
            if title.find(keyword) > 0:
                posts.remove(post)
                break
    logging.info("len:%d", len(posts))

async def get_douban_posts():
    posts = await search_douban_groups()
    logging.info("len:%d", len(posts))
    filter_posts( posts )
    logging.info("len:%d", len(posts))
    for url,title in posts:
        print('- [ ', title, ' ](', url, ')')
#    logging.info("posts:%s", str(posts))

def main():
    log_init()
    proxymanager.get_high_anonymous_proxy_list()
    event_loop = asyncio.get_event_loop()
    event_loop.run_until_complete(get_douban_posts())

if __name__ == '__main__':
    main()

