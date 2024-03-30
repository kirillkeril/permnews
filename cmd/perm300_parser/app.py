from bs4 import BeautifulSoup
import requests
import time
from kafka import KafkaProducer
import os
from dotenv import load_dotenv
import json
import asyncio
from requests_toolbelt.multipart.encoder import MultipartEncoder

def get_category(n):
    if n == 0:
        return "theatre"
    elif n == 1:
        return "music"
    elif n == 2:
        return "cinema"
    elif n == 4:
        return "exhibition"
    elif n == 5:
        return "festival"
    elif n == 7:
        return "other"
    elif n == 8:
        return "excursions"
    elif n == 10:
        return "sports"
    elif n == 11:
        return "education"

async def main():
    load_dotenv()
    c = list(map(int, os.getenv('CATEGORIES').split(";"))) # 0 - theatre; 1 - music; 2 - cinema; 4 - exhibition; 5 - festival; 7 - other; 8 - excursions; 10 - sports; 11 - obrazovanie;  
    kfk = os.getenv('KAFKA_URL')
    t = int(os.getenv('TIMEOUT'))
    topic = os.getenv('KAFKA_TOPIC')
    if t is None:
        t = 60*60*1
    else:
        t = int(t)
    producer = KafkaProducer(bootstrap_servers=[kfk], value_serializer=lambda m: json.dumps(m).encode('ascii'))
    print('start')
    while True:
        await asyncio.gather(
            *map(lambda category: process_afisha(category, producer, topic), c)
        )
        await process_afisha(c, producer, topic)
        print('fall asleep for {} seconds'.format(t))
        time.sleep(t)


async def process_afisha(category, producer, topic):
    shortUrl = "https://perm-300.ru"
    url = shortUrl + "/events8"

    formData = MultipartEncoder(fields={
        'tv67[]': '{}'.format(category),
        'mode': 'eventis'
    })

    response = requests.post(url, data=formData, headers={
        'Content-Type': formData.content_type,
        'Accept':'*/*',
        'Accept-Encoding':'gzip, deflate, br',
        'Accept-Language':'ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3',
        'Connection':'keep-alive',
        'Cookie':'evoosd7pn=2813a76c298aa522654a12ecfb069a2e; view1=1; _ym_uid=1711732612377140203; _ym_d=1711732612; _ym_isad=2; view7=1; cookie=1; view1710=1; view1715=1; _ym_visorc=w',
        'Host':'perm-300.ru',
        'Origin':'https://perm-300.ru',
        'Referer':'https://perm-300.ru/afisha/',
        'Sec-Fetch-Dest':'empty',
        'Sec-Fetch-Mode':'cors',
        'Sec-Fetch-Site':'same-origin',
        'User-Agent':'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0',
        'X-Requested-With':'XMLHttpRequest',
    })
    print(formData, category)
    if not response.ok:
        print(response.status_code, url, category)
        pass

    htmlContent = response.content

    soup = BeautifulSoup(htmlContent, "html.parser")

    events = soup.find_all("a", class_="event")

    data = []
    for event in events:
        eventDate = event.find("div", class_="eventdate").text.strip()
        eventTitle = event.find("div", class_="eventtitle").text.strip()
        link = shortUrl + event['href']

        styleString = event.find("div", class_="eventimg")['style']
        image = shortUrl + '/' + styleString.split("url('", 1)[-1].split("')", 1)[0].strip()
        desc = ""
        try:
            desc = event.find("div", class_="shadow").text.strip()
        except Exception:
            print('no desc')
        data.append({
            "image": image,
            "link:": link,
            "source": shortUrl,
            "notice": eventDate,
            "title": eventTitle,
            "desc": desc,
            "category": get_category(category)
        })
        print({
            "image": image,
            "link:": link,
            "source": shortUrl,
            "notice": eventDate,
            "title": eventTitle,
            "desc": desc,
            "category": get_category(category)
        })


    producer.send(topic, data).add_callback(on_send_success).add_errback(on_send_error)
    print("sent {} items in data".format(len(data)))
    await asyncio.sleep(1)

def on_send_success(record_metadata):
    print(record_metadata.topic)
    print(record_metadata.partition)
    print(record_metadata.offset)

def on_send_error(excp):
    print('error: {}'.format(excp))

if __name__ == "__main__":
    asyncio.run(main())

