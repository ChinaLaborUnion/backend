FROM alpine:3.6

RUN sed -i 's/dl-cdn\.alpinelinux\.org/mirrors\.aliyun\.com/g' /etc/apk/repositories
RUN apk update --no-cache
ENV TZ=Asia/Shanghai
RUN apk update \
    && apk add tzdata \
    && echo "${TZ}" > /etc/timezone \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && rm /var/cache/apk/*
RUN mkdir /main


COPY main /main/
COPY conf.ini /main/
COPY apiclient_cert.p12 /main/apiclient_cert.p12
COPY apiclient_cert.pem /main/apiclient_cert.pem
COPY apiclient_key.pem /main/apiclient_key.pem

COPY alipayCertPublicKey_RSA2.crt /main/alipayCertPublicKey_RSA2.crt
COPY alipayRootCert.crt /main/alipayRootCert.crt
COPY appCertPublicKey_2021002134605526.crt /main/appCertPublicKey_2021002134605526.crt

WORKDIR /main
ENTRYPOINT ["./main"]