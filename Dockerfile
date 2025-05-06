FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y iperf3
RUN rm -rf /var/cache/apt/archives /var/lib/apt/lists/*