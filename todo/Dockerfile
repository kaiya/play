FROM alpine

EXPOSE 8080
ARG SQLPASS
WORKDIR /app/work
COPY ./todo /app/work/todo

ENTRYPOINT ["sh", "-c", "./todo"]
