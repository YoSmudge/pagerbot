FROM bitnami/minideb:stretch
RUN install_packages ca-certificates
COPY dist/linux_amd64/pagerbot config.yml /
CMD ["/pagerbot"]
