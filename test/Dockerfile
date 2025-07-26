# Panchangam API Testing Framework Docker Image
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    golang-go \
    git \
    && rm -rf /var/lib/apt/lists/*

# Set Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Copy requirements first for better Docker layer caching
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy test files
COPY . .

# Create directories for test outputs
RUN mkdir -p /app/reports /app/htmlcov /app/logs

# Set environment variables
ENV PYTHONPATH=/app
ENV PYTEST_DISABLE_PLUGIN_AUTOLOAD=1

# Expose port for potential test server
EXPOSE 8080

# Default command
CMD ["python", "run_tests.py", "--help"]