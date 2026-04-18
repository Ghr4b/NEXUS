FROM eclipse-temurin:17-jdk-jammy
RUN groupadd -r ctfgroup && useradd -r -g ctfgroup ctfuser
WORKDIR /app
COPY target/*.jar app.jar
RUN mkdir -p /tmp/ctf_payload && chown ctfuser:ctfgroup /tmp/ctf_payload

# The Prize
RUN RANDOM_FILE="flag_$(tr -dc 'a-z0-9' < /dev/urandom | head -c 8).txt" && \
    echo "cybersphere{f1rst_T1me_wasS00Nic3!H4d2}" > "/$RANDOM_FILE" && \
    echo "Created flag file: $RANDOM_FILE"
RUN chown root:root /$RANDOM_FILE && chmod 444 /$RANDOM_FILE
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER ctfuser
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
