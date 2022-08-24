FROM golang:1.19

RUN mkdir -p /app
RUN mkdir -p /app/data
ADD repository /app
RUN chmod +x /app/repository
CMD /app/repository /app/data