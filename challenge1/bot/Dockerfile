# Using Python on Debian Trixie (Testing)
FROM python:3.11-slim-trixie

WORKDIR /app

# 1. Install the EXACT package version you found in apt-cache
RUN apt-get update && apt-get install -y --allow-downgrades \
    chromium \
    fonts-liberation \
    fonts-unifont \
    && rm -rf /var/lib/apt/lists/*

# 2. Install Playwright Python library
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 3. Copy your test script
COPY test.py .

# Run the test
CMD ["python", "test.py"]
