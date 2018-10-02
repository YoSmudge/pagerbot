FROM bitnami/minideb:stretch
COPY dist/linux_amd64/pagerbot /
CMD ["/pagerbot"]
