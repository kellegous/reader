FROM kellegous/build:f1799259 AS build

ARG SHA
ARG BUILD_TIME

COPY . /src

RUN cd /src && make SHA=${SHA} BUILD_TIME=${BUILD_TIME} clean ALL

FROM lsiobase/debian:bookworm

COPY etc/setup.sh /setup.sh

RUN /setup.sh && rm /setup.sh

COPY --from=build /src/bin/reader /app/bin/reader

CMD ["/usr/bin/with-contenv", "/app/bin/reader", "--config-file=/data/reader.yaml"]
