FROM timescale/timescaledb:latest-pg15

# Install build dependencies
RUN apk add --no-cache \
    git \
    build-base \
    postgresql15-dev \
    rust \
    cargo \
    clang15 \
    llvm15

# Install pgvectorscale
RUN git clone https://github.com/timescale/pgvectorscale.git \
    && cd pgvectorscale/pgvectorscale \
    && cargo build --release \
    && cp target/release/libvectorscale.so "$(pg_config --pkglibdir)" \
    && cp ../extension/vectorscale* "$(pg_config --sharedir)/extension/" \
    && cd ../.. \
    && rm -rf pgvectorscale