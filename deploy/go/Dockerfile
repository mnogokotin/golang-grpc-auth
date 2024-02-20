FROM golang:1.22.0-alpine3.19 as builder

# install alpine packages
RUN apk update \
    && apk add make

# copy project files
COPY ./. /app
WORKDIR /app

# build app
RUN go build ./cmd/sso

# run migrations
#RUN make migrate-up


FROM scratch
#COPY --from=builder /app/.env /app/
COPY --from=builder /app /app/
CMD ["/app/sso"]
