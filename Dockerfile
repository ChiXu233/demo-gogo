FROM alpine

WORKDIR /workspace/demo-gogo

COPY demo-gogo .

ADD conf .


EXPOSE 9094

CMD ["./demo-gogo"]
