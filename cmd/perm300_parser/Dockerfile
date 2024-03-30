FROM python:3.9-slim as builder

WORKDIR /app

COPY . .

RUN pip install beautifulsoup4
RUN pip install kafka-python
RUN pip install requests
RUN pip install python-dotenv
RUN pip install requests-toolbelt

CMD ["python", "-u", "app.py"]