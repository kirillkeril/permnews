import requests
import time
from kafka import KafkaProducer
import os
from dotenv import load_dotenv
import json
import asyncio

async def main():
    load_dotenv()
    types = os.getenv('EVENT_TYPES').split(';') # spectacle;all-events-cinema
    kfk = os.getenv('KAFKA_URL')
    t = int(os.getenv('TIMEOUT'))
    topic = os.getenv('KAFKA_TOPIC')
    pagesCount = int(os.getenv('PAGES'))
    if t is None:
        t = 60*60*1
    else:
        t = int(t)
    producer = KafkaProducer(bootstrap_servers=[kfk], value_serializer=lambda m: json.dumps(m).encode('ascii'))
    print('start')
    while True:
        await asyncio.gather(
            *map(lambda x: process_afisha(x, pagesCount, producer, topic), types)
        )
        print('fall asleep for {} seconds'.format(t))
        time.sleep(t)


async def process_afisha(type, pagesCount, producer, topic):
    print('processing https://afisha.yandex.ru/api/events/selection/{}'.format(type))
    for i in range(1, pagesCount):
        offset = 12 * (i - 1)
        url = 'https://afisha.yandex.ru/api/events/selection/{}?limit=12&offset={}&city=perm'.format(type, offset)
        res = requests.get(url, headers={
            'Accept':	'application/json, text/javascript, */*; q=0.01',
            'Accept-Encoding':	'gzip, deflate, br',
            'Accept-Language':	'ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3',
            'Cache-Control':	'no-cache',
            'Connection':	'keep-alive',
            'Host':	'afisha.yandex.ru',
            'Pragma':	'no-cache',
            'Referer':	'https://afisha.yandex.ru/perm/selections/{}?page={}'.format(type, i),
            'Sec-Fetch-Dest':	'empty',
            'Sec-Fetch-Mode':	'no-cors',
            'Sec-Fetch-Site':	'same-origin',
            'User-Agent':	'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0',
            'x-force-cors-preflight':	'1',
            'X-Parent-Request-Id':	'1711781556060751-13537792777347991917',
            'X-Requested-With':	'XMLHttpRequest',
            'X-Retpath-Y':	'https://afisha.yandex.ru/perm/selections/{}?page={}'.format(type, i),
        })
        if not res.ok:
            print(url)
            break
        data = res.json()
        events = []
        for i in data['data']:
            link = "afisha.yandex.ru{}".format(i['event']['url'])
            notice = "с {} по {}, {}".format(i['scheduleInfo']['dateStarted'], i['scheduleInfo']['dateEnd'], i['scheduleInfo']['placePreview'])
            image = None
            try:
                image = i['event']['image']['sizes']['eventCover']['url']
            except TypeError:
                print(i['event'])
            title = i['event']['title']
            events.append({"title": title, "notice": notice, "link": link, "image": image, "source": "afisha.yandex.ru"})
        
        producer.send(topic, events).add_callback(on_send_success).add_errback(on_send_error)
        await asyncio.sleep(5)


def on_send_success(record_metadata):
    print(record_metadata.topic)
    print(record_metadata.partition)
    print(record_metadata.offset)

def on_send_error(excp):
    print('error: {}'.format(excp))

if __name__ == "__main__":
    asyncio.run(main())

