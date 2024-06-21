FROM alpine

WORKDIR /workspace/demo-gogo

COPY demo-gogogo .

ADD conf .


EXPOSE 9094

CMD ["./demo-gogo"]
