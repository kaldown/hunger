FROM golang:latest

COPY ./build/generated/quiz.pb.go /usr/src/hunger/build/generated/
COPY ./build/bin/server /usr/src/hunger/build/bin/

WORKDIR /usr/src/hunger/build/bin/

CMD ./server
