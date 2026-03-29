FROM eclipse-temurin:17-jdk-jammy
RUN groupadd -r ctfgroup && useradd -r -g ctfgroup ctfuser
WORKDIR /app
COPY target/*.jar app.jar
RUN mkdir -p /tmp/ctf_payload && chown ctfuser:ctfgroup /tmp/ctf_payload

# The Prize
RUN echo "IFTxSecurinets{22cf023b5edf2c6f0914d8095cf3ac38}" > /thebigflag.txt
RUN chown root:root /thebigflag.txt && chmod 444 /thebigflag.txt
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
USER ctfuser
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
