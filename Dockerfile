FROM alpine:3.14

RUN apk update --quiet && apk upgrade && \
    apk add --no-cache bash git openssh tar

RUN git clone https://github.com/florianbeer/mella && \
    mv mella/mella /bin/

ADD drone-mella /bin/
ENTRYPOINT /bin/drone-mella
