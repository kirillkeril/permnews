from bs4 import BeautifulSoup
import requests
import time
from kafka import KafkaProducer
import os
from dotenv import load_dotenv

def main():
    load_dotenv()
    kfk = os.getenv('KAFKA_URL')
    producer = KafkaProducer(bootstrap_servers=[kfk],
                             value_serializer=lambda x:
                             dumps(x).encode('utf-8'))
    print('start')
    url = os.getenv('URL')
    while True:
        print('process')
        process(url)
        time.sleep(60 * 60 * 8)


def process(url):
    for i in range(100):
        url = 'https://www.afisha.ru/prm/{}/page{}'.format(url, i + 1)
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
            break
        soup = BeautifulSoup(res.text)
        list = soup.find_all("div", {'aria-label': 'Список'})[0]
        items = list.find_all("div", {'role': 'listitem', 'data-test': 'ITEM'})
        last_len = len(items)
        data = []
        for i in items:
            title = i.get('title')
            notice = i.find_all("div", {'data-test': 'ITEM-META ITEM-NOTICE'})[0].text
            link = "https://www.afisha.ru" + i.find_all('a', {'data-test': 'LINK ITEM-NAME ITEM-URL'})[0].get('href')
            data.append({"title": title, "notice": notice, "link": link})

        producer.send('numtest', value=data)
        time.sleep(20)


main()
