FROM amd64/golang:latest
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go build -o linux_sindel_auth2 .
#RUN adduser -S -D -H -h /app appuser
#USER appuser
