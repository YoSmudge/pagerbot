FROM bitnami/minideb:stretch
COPY dist/linux_amd64/pagerbot config.yml /
CMD ["/pagerbot"]
