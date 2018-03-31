FROM golang:1.9

WORKDIR /go/src/github.com/LeReverandNox/GuessWhat/src

RUN apt-get update
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash
RUN apt-get install -y nodejs
RUN npm install -g bower

ADD files /
RUN chmod +x /run.sh

ADD src .

RUN go-wrapper download
RUN go-wrapper install
RUN go build -o guesswhat

RUN go get github.com/pilu/fresh

RUN bower install --allow-root

EXPOSE 3000

ENV ENVIRONMENT prod

ENTRYPOINT ["/run.sh"]
