FROM alpine

WORKDIR /workspace/demo-gogo

COPY teach .

ADD config.yaml .
ADD log.json .


EXPOSE 9094

CMD ["./teach"]
