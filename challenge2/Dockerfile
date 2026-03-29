FROM eclipse-temurin:17-jdk-jammy
RUN groupadd -r ctfgroup && useradd -r -g ctfgroup ctfuser
WORKDIR /app
COPY target/*.jar app.jar
RUN mkdir -p /tmp/ctf_payload && chown ctfuser:ctfgroup /tmp/ctf_payload

# The Prize
RUN echo "CTF{redos_to_jndi_rce_master}" > /flag.txt
RUN chown root:root /flag.txt && chmod 444 /flag.txt
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER ctfuser
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]