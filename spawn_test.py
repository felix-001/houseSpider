from gevent import monkey; monkey.patch_all();import gevent
import requests

def subtask( i, maintask ):
    print('enter subtask ', i, ' ', maintask )
    try:
        requests.get('http://www.google.com')
    except Exception as exc:
        print('...')
    print('leave subtask ', i, ' ', maintask )

def maintask( i ):
    tasks = []
    print('enter maintask ', i )
    try:
        requests.get('http://www.google.com')
    except Exception as exc:
        print('...')

    for v in range(500):
        task = gevent.spawn( subtask, v, i )
        tasks.append( task )
    gevent.joinall( tasks )
    print('leave maintask ', i )

def main():
    mains = []
    for i in range(500):
        maint = gevent.spawn( maintask, i )
        mains.append( maint )
    gevent.joinall( mains )

if __name__ == '__main__':
    main()
