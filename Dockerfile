FROM alpine

WORKDIR /workspace/demo-gogo

COPY demo-gogo .

ADD config.yaml .
ADD log.json .

EXPOSE 9094

CMD ["./demo-gogo"]
