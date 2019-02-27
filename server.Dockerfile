FROM alpine:3.7

COPY ./build/generated/quiz.pb.go /app/build/generated
COPY ./build/bin/server /app/server

CMD ./app/server
