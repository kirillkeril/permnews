from bs4 import BeautifulSoup
import requests
import time
from kafka import KafkaProducer
import os
from dotenv import load_dotenv
import json
import asyncio

async def main():
    load_dotenv()
    u = os.getenv('URLS').split(';')
    kfk = os.getenv('KAFKA_URL')
    t = int(os.getenv('TIMEOUT'))
    pagesCount = int(os.getenv('PAGES'))
    topic = os.getenv('KAFKA_TOPIC')
    if t is None:
        t = 60*60*1
    else:
        t = int(t)
    producer = KafkaProducer(bootstrap_servers=[kfk], value_serializer=lambda m: json.dumps(m).encode('ascii'))
    print('start')
    while True:
        await asyncio.gather(
            *map(lambda x: process_afisha(x, pagesCount, producer, topic), u)
        )
        print('fall asleep for {} seconds'.format(t))
        time.sleep(t)


async def process_afisha(u, pagesCount, producer, topic):
    print('processing https://www.afisha.ru/prm/{}'.format(u))
    for i in range(pagesCount):
        url = 'https://www.afisha.ru/prm/{}/page{}'.format(u, i + 1)
        res = requests.get(url,
                           headers={
                               'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
                               'Accept-Encoding': 'gzip, deflate, br',
                               'Accept-Language': 'ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3',
                               'Cache-Control': 'no-cache',
                               'Connection': 'keep-alive',
                               'Host': 'www.afisha.ru',
                               'Pragma': 'no-cache',
                               'Sec-Fetch-Dest': 'document',
                               'Sec-Fetch-Mode': 'navigate',
                               'Sec-Fetch-Site': 'same-origin',
                               'Upgrade-Insecure-Requests': '1',
                               'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0'
                           }
                           )
        if not res.ok:
            print(res.status_code, url)
            break
        soup = BeautifulSoup(res.text, features="html.parser")
        list = soup.find_all("div", {'aria-label': 'Список'})[0]
        items = list.find_all("div", {'role': 'listitem', 'data-test': 'ITEM'})
        data = []
        for i in items:
            title = i.get('title')
            notice = i.find_all("div", {'data-test': 'ITEM-META ITEM-NOTICE'})[0].text
            link = "https://www.afisha.ru" + i.find_all('a', {'data-test': 'LINK ITEM-NAME ITEM-URL'})[0].get('href')
            source = "https://www.afisha.ru"
            image = i.find_all("img", { 'data-test': "IMAGE ITEM-IMAGE" })[0].get('src')
            data.append({"title": title, "notice": notice, "link": link, "source": source, "image": image})

        producer.send(topic, data).add_callback(on_send_success).add_errback(on_send_error)
        print("sent {} items in data".format(len(data)))
        await asyncio.sleep(20)

def on_send_success(record_metadata):
    print(record_metadata.topic)
    print(record_metadata.partition)
    print(record_metadata.offset)

def on_send_error(excp):
    print('error: {}'.format(excp))

if __name__ == "__main__":
    asyncio.run(main())
