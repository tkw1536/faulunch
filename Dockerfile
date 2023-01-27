FROM alpine as permission

# create www-data
RUN set -x ; \
  addgroup -g 82 -S www-data ; \
  adduser -u 82 -D -S -G www-data www-data && exit 0 ; exit 1
RUN mkdir /data/ && chown -R www-data:www-data /data/
RUN apk --no-cache add ca-certificates

# build the server
FROM golang as build

# build the app
ADD . /app/
WORKDIR /app/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/faulunch /app/cmd/serve/

# add it into a scratch image
FROM scratch
WORKDIR /

# add the user
COPY --from=permission /etc/passwd /etc/passwd
COPY --from=permission /etc/group /etc/group
COPY --from=permission /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# add the app
COPY --from=build /app/faulunch /faulunch

# and set the entry command
EXPOSE 8080
USER www-data:www-data
COPY --from=permission --chown=www-data:www-data /data/ /data/
VOLUME /data/
CMD ["/faulunch", "-sync", "12h", "-addr", "0.0.0.0:8080", "/data/data.sql"]
