FROM bravissimolabs/alpine-git

RUN apk update --quiet && apk add tar && \
    git clone https://github.com/florianbeer/mella && \
    mv mella/mella /bin/

ADD drone-mella /bin/
ENTRYPOINT /bin/drone-mella
