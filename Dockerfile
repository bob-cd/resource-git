FROM python:3

WORKDIR /opt

ADD server.py .
ADD requirements.txt .

RUN pip install -r requirements.txt

ENTRYPOINT ["python3", "server.py"]
