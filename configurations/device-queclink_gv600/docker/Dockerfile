FROM golang:1.10.2

ADD . /go/src/bitbucket.org/positrace/queclink55

RUN git config --global http.sslverify false

RUN go get -u github.com/kardianos/govendor

WORKDIR /go/src/bitbucket.org/positrace/queclink55/src/serviceq600

RUN go build -o /app/queclink600

WORKDIR /go/src/bitbucket.org/positrace/queclink55/
RUN ln -fs /config/credentials.json /app/credentials.json
RUN ln -fs /config/ReportConfiguration.xml /app/ReportConfiguration.xml 

WORKDIR /app

CMD ["./queclink600"]

