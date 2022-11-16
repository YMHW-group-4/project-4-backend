


WORKDIR /crypto
COPY . .
RUN make build_arm64v8

FROM arm64v8/alpine:3.14


ENV ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.19
CMD ["/crypto"]
