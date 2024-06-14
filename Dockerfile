FROM kellegous/build:0d98364e as build

COPY . /src

RUN cd /src && make clean ALL

FROM lsiobase/debian:bookworm

COPY etc/setup.sh /setup.sh

RUN /setup.sh && rm /setup.sh

COPY --from=build /src/bin/reader /app/bin/reader

CMD ["/usr/bin/with-contenv", "/app/bin/reader", "--config-file=/data/reader.yaml"]
