FROM alpine

WORKDIR /workspace/demo-gogo

COPY teach .

ADD conf .


EXPOSE 9094

CMD ["./teach"]
