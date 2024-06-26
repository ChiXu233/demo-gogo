FROM alpine

WORKDIR /workspace/teach

COPY teach .

ADD conf .


EXPOSE 9094

CMD ["./teach"]
