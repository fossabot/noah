FROM ubuntu:16.04 as builder
RUN apt-get update
RUN apt-get --assume-yes install \
    libreadline6 \
    libreadline6-dev \
    git-all \
    build-essential \
    zlib1g-dev \
    libossp-uuid-dev \
    flex \
    bison \
    libxml2-utils \
    xsltproc
RUN git clone https://github.com/postgres/postgres.git
RUN mkdir /postbuild
RUN ./postgres/configure --prefix=/postbuild --with-ossp-uuid
RUN cd ./postgres
RUN make
RUN make world
RUN make install

FROM ubuntu:16.04
RUN mkdir /postgres
COPY --from=builder /postbuild /postgres
RUN mkdir /db
#COPY ./pg_hba.conf /db/pg_hba.conf
#COPY ./postgresql.conf /db/postgresql.conf
#RUN echo 12 >> /db/PG_VERSION
EXPOSE 5432
ENV LD_LIBRARY_PATH=/postgres/lib

RUN groupadd -g 999 postdb && \
    useradd -r -u 999 -g postdb postdb
RUN chown -R postdb /db
RUN chmod -R 750 /db
USER postdb
RUN /postgres/bin/initdb -D /db
CMD ["/postgres/bin/postgres", "-D", "/db"]