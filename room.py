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
from lxml import etree

words_too_less_cnt = 0
no_need_keyword_cnt = 0
filter_keyword_cnt = 0
total_posts = 0
total_posts_search = 0
cur = 0
cur_search = 0

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

def get_douban_html(url, group=None, q=None):
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
    global cur_search

    show_progress( cur_search, total_posts_search )
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
    show_progress( cur_search, total_posts_search )
    cur_search += 1

async def search_douban_groups():
    global total_posts_search
    all_posts = []
    tasks = []

    total_posts_search = len(config.GROUPS) * len(config.KEYWORDS)
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

def check_post_valid(post):
    global words_too_less_cnt
    global no_need_keyword_cnt
    global filter_keyword_cnt

    found = 0

    if len(post['content']) < 200:
        words_too_less_cnt += 1
        return False
    if len(config.KEYWORDS) > 0:
        for keyword in config.KEYWORDS:
            if keyword.encode('utf8') in post['content']:
                found = 1
                break
    if found == 0:
        no_need_keyword_cnt += 1
        return False

    for keyword in config.TITLE_FILTER_KEYWRODS:
        if keyword.encode('utf8') in post['content']:
            filter_keyword_cnt += 1
            return False
        if keyword.encode('utf8') in post['title']:
            filter_keyword_cnt += 1
            return False
    return True

def extract( regx, body, multi=False):
    if isinstance(body, str):
        body = etree.HTML(body)
        res = body.xpath(regx)
        if multi:
            return res
        return res[0] if res else None

def parse_detail_page(html):
    post = {}

    title = extract(config.RULES["detail_title_sm"], html) or extract(config.RULES["detail_title_lg"], html)
    if title is None:
        logging.error("parse title error")
        return None
    post["title"] = title.strip().encode("utf8")
    post["create_time"] = extract( config.RULES["create_time"], html).strip().encode('utf8')
    post["author"] = extract( config.RULES["detail_author"], html).strip().encode('utf8')
    post["content"] = '\n'.join( extract(config.RULES["content"], html, multi=True) \
            or extract(config.RULES["content2"], html, multi=True) ).encode('utf8')
    return post

def show_progress( cur, total):
    percent = cur*100.0/total
    progress = '[ ' + '%.2f' % percent + '% ] ' + str(cur)+'/'+str(total)+'\r'
    sys.stdout.write(progress)
    sys.stdout.flush()
    if cur >= total:
        print('Done.')

async def get_douban_detail_page(url, posts):
    global cur

    event_loop = asyncio.get_event_loop()
    show_progress(cur, total_posts)
    html = await event_loop.run_in_executor(
        None,
        partial( get_douban_html, url ) 
    )
    if html is None:
        logging.info("fetch detail page:%s error", url )
        return None
    post = parse_detail_page( html )
    if check_post_valid(post):
        res = ( post['title'], url )
        posts.append(res)
    show_progress(cur, total_posts)
    cur += 1

async def get_douban_posts():
    global total_posts
    tasks = []
    all_posts = []

    logging.info("search all the groups...")
    posts = await search_douban_groups()
    logging.info("total crawled %d posts", len(posts))
    filter_posts( posts )
    logging.info("after filter, levae %d posts", len(posts))
    logging.info("search all the post detail...")
    total_posts = len(posts)
    event_loop = asyncio.get_event_loop()
    for url,title in posts:
        tasks.append(event_loop.create_task(get_douban_detail_page(url, all_posts)))
    done, pending = await asyncio.wait(tasks)
    assert bool(pending) is False
    logging.info("final leave %d posts, words_too_less_cnt : %d no_need_keyword_cnt : %d filter_keyword_cnt:%d",
            len(all_posts), words_too_less_cnt, no_need_keyword_cnt, filter_keyword_cnt )
    f = open('douban.md', 'w')
    for title,url in all_posts:
        md = '- [ '+ title.decode('utf8')+ ' ](' + url + ')' + '\n'
        f.write(md)
    f.close()

def main():
    log_init()
    logging.info("start...")
    proxymanager.get_high_anonymous_proxy_list()
    event_loop = asyncio.get_event_loop()
    event_loop.run_until_complete(get_douban_posts())

if __name__ == '__main__':
    main()

