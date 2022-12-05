FROM postgres:14

RUN sed -i -e"s/# ru_RU.UTF-8 UTF-8/ru_RU.UTF-8 UTF-8/" /etc/locale.gen \
    && locale-gen \
    && echo 'LANG="ru_RU.UTF-8"' > /etc/default/locale

ENV LC_ALL ru_RU.UTF-8
ENV LANG ru_RU.UTF-8