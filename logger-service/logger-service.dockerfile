FROM alpine:latest

RUN mkdir /app

COPY logAPP /app 

CMD ["/app/logAPP"]