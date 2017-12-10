FROM debian:stretch-slim as extract
ARG WITH_JMX
COPY datadog-agent*_amd64.deb /
WORKDIR /output

# Get tini and check gpg signature
ENV TINI_VERSION v0.16.1
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-amd64 tini
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-amd64.asc tini.asc
RUN apt-get update \
 && apt-get install --no-install-recommends -y dirmngr gpg \
 && gpg --keyserver ipv4.pool.sks-keyservers.net --recv-keys 595E85A6B1B4779EA4DAAEC70B588DFF0527A9B7 \
 && gpg --verify tini.asc \
 && rm tini.asc \
 && chmod +x tini

# Get Consul tools
ENV CONSUL_VERSION=0.9.3
ENV CONSULTEMPLATE_VERSION=0.19.3
ENV HASHI_KEY 91A6E7F85D05C65630BEF18951852D87348FFC4C

ADD https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip consul_${CONSUL_VERSION}_linux_amd64.zip
ADD https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_SHA256SUMS consul_${CONSUL_VERSION}_SHA256SUMS
ADD https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_SHA256SUMS.sig consul_${CONSUL_VERSION}_SHA256SUMS.sig

ADD https://releases.hashicorp.com/consul-template/${CONSULTEMPLATE_VERSION}/consul-template_${CONSULTEMPLATE_VERSION}_linux_amd64.zip consul-template_${CONSULTEMPLATE_VERSION}_linux_amd64.zip
ADD https://releases.hashicorp.com/consul-template/${CONSULTEMPLATE_VERSION}/consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS
ADD https://releases.hashicorp.com/consul-template/${CONSULTEMPLATE_VERSION}/consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS.sig consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS.sig

RUN mkdir -p usr/local/bin \
 && apt-get update \
 && apt-get install --no-install-recommends unzip \
 && gpg --keyserver keys.gnupg.net --receive-keys $HASHI_KEY \
 && gpg --batch --verify consul_${CONSUL_VERSION}_SHA256SUMS.sig consul_${CONSUL_VERSION}_SHA256SUMS \
 && grep consul_${CONSUL_VERSION}_linux_amd64.zip consul_${CONSUL_VERSION}_SHA256SUMS | sha256sum -c \
 && unzip -d usr/local/bin consul_${CONSUL_VERSION}_linux_amd64.zip \
 && rm consul_${CONSUL_VERSION}_linux_amd64.zip \
 && rm consul_${CONSUL_VERSION}_SHA256SUMS \
 && rm consul_${CONSUL_VERSION}_SHA256SUMS.sig \
 && chmod +x usr/local/bin/consul \
 && gpg --batch --verify consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS.sig consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS \
 && grep consul-template_${CONSULTEMPLATE_VERSION}_linux_amd64.zip consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS | sha256sum -c \
 && unzip -d usr/local/bin consul-template_${CONSULTEMPLATE_VERSION}_linux_amd64.zip \
 && rm consul-template_${CONSULTEMPLATE_VERSION}_linux_amd64.zip \
 && rm consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS \
 && rm consul-template_${CONSULTEMPLATE_VERSION}_SHA256SUMS.sig \
 && chmod +x usr/local/bin/consul-template

# Extract and cleanup:
#   - unused systemd unit
#   - GPL sources for embedded software  # FIXME: move upstream
#   - docs and manpages                  # FIXME: move upstream
#   - static libraries                   # FIXME: move upstream
#   - jmxfetch on nojmx build
RUN dpkg -x /datadog-agent*_amd64.deb . \
 && echo "Cleanup and trimming" \
 && rm -rf usr/bin etc/init lib \
    opt/datadog-agent/sources \
    opt/datadog-agent/embedded/share/doc \
    opt/datadog-agent/embedded/share/man \
 && find opt/datadog-agent/ -iname "*.a" -delete \
 && if [ -z "$WITH_JMX" ]; then rm -rf opt/datadog-agent/bin/agent/dist/jmx; fi \
 && echo "Linking embedded ssl certificates" \
 && ln -s /opt/datadog-agent/embedded/ssl etc/ssl

# Configuration:
#   - empty config file to remove logline
#   - enable docker check by default
#   - enable process agent
COPY datadog*.yaml etc/datadog-agent/
RUN mv etc/datadog-agent/conf.d/docker.d/conf.yaml.example etc/datadog-agent/conf.d/docker.d/conf.yaml.default \
 && sed -i "s/enabled: false/enabled: true/" etc/datadog-agent/conf.d/process_agent.yaml.default
COPY entrypoint.sh .


FROM debian:stretch-slim
LABEL maintainer "Datadog <package@datadoghq.com>"
ARG WITH_JMX
ENV DOCKER_DD_AGENT=yes \
    PATH=/opt/datadog-agent/bin/agent/:/opt/datadog-agent/embedded/bin/:$PATH \
    CURL_CA_BUNDLE=/opt/datadog-agent/embedded/ssl/certs/cacert.pem

# Install openjdk-8-jre-headless on jmx flavor
RUN if [ -n "$WITH_JMX" ]; then apt-get update \
 && mkdir /usr/share/man/man1 \
 && apt-get install --no-install-recommends -y openjdk-8-jre-headless \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*; fi

# Copy agent from extract stage
COPY --from=extract /output/ /

# Expose DogStatsD and trace-agent ports
EXPOSE 8125/udp 8126/tcp

# Extra conf.d and checks.d
VOLUME ["/conf.d", "/checks.d"]

# Placeholder healthcheck, to be improved
HEALTHCHECK --interval=5m --timeout=20s --retries=1 \
  CMD ["/opt/datadog-agent/bin/agent/agent", "status"]

ENTRYPOINT ["/tini", "--", "/entrypoint.sh"]
CMD ["/opt/datadog-agent/bin/agent/agent", "start"]
