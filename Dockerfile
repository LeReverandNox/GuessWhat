FROM golang:1.9

WORKDIR /go/src/github.com/LeReverandNox/GuessWhat/src

ADD files /
RUN chmod +x /run.sh

ADD src .

RUN go-wrapper download
RUN go-wrapper install

RUN go get github.com/pilu/fresh

VOLUME /go/src/github.com/LeReverandNox/GuessWhat/src

EXPOSE 3000

ENV ENVIRONMENT prod

ENTRYPOINT ["/run.sh"]