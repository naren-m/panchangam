"""
Pytest configuration and fixtures for Panchangam API tests
"""
import pytest
import requests
import subprocess
import time
import os
from typing import Generator


@pytest.fixture(scope="session")
def api_base_url() -> str:
    """Base URL for the API gateway"""
    return os.getenv("PANCHANGAM_API_URL", "http://localhost:8080")


@pytest.fixture(scope="session")
def start_servers() -> Generator[None, None, None]:
    """Start the gRPC and Gateway servers for testing"""
    if os.getenv("SKIP_SERVER_START", "false").lower() == "true":
        yield
        return
    
    # Build servers (from parent directory)
    subprocess.run(["go", "build", "-o", "grpc-server", "../server/server.go"], check=True, cwd="../")
    subprocess.run(["go", "build", "-o", "gateway-server", "../cmd/gateway/main.go"], check=True, cwd="../")
    
    # Start gRPC server
    grpc_proc = subprocess.Popen(["../grpc-server"], cwd="../")
    time.sleep(2)
    
    # Start Gateway server
    gateway_proc = subprocess.Popen([
        "../gateway-server",
        "--grpc-endpoint=localhost:50052",
        "--http-port=8080"
    ], cwd="../")
    time.sleep(3)
    
    # Verify servers are running
    max_retries = 10
    for i in range(max_retries):
        try:
            response = requests.get("http://localhost:8080/api/v1/health")
            if response.status_code == 200:
                break
        except requests.exceptions.ConnectionError:
            if i == max_retries - 1:
                raise
            time.sleep(1)
    
    yield
    
    # Cleanup
    gateway_proc.terminate()
    grpc_proc.terminate()
    gateway_proc.wait(timeout=5)
    grpc_proc.wait(timeout=5)


@pytest.fixture
def api_client(api_base_url: str, start_servers) -> requests.Session:
    """HTTP client for API requests"""
    session = requests.Session()
    session.headers.update({
        "Accept": "application/json",
        "Content-Type": "application/json"
    })
    session.base_url = api_base_url
    return session


@pytest.fixture
def sample_locations():
    """Sample locations for testing"""
    return {
        "bangalore": {
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
            "name": "Bangalore"
        },
        "mumbai": {
            "lat": 19.0760,
            "lng": 72.8777,
            "tz": "Asia/Kolkata",
            "name": "Mumbai"
        },
        "new_york": {
            "lat": 40.7128,
            "lng": -74.0060,
            "tz": "America/New_York",
            "name": "New York"
        },
        "london": {
            "lat": 51.5074,
            "lng": -0.1278,
            "tz": "Europe/London",
            "name": "London"
        }
    }


@pytest.fixture
def valid_request_params(sample_locations):
    """Valid request parameters for testing"""
    location = sample_locations["bangalore"]
    return {
        "date": "2024-01-15",
        "lat": location["lat"],
        "lng": location["lng"],
        "tz": location["tz"]
    }


def pytest_configure(config):
    """Configure pytest with custom markers"""
    config.addinivalue_line(
        "markers", "smoke: mark test as smoke test"
    )
    config.addinivalue_line(
        "markers", "integration: mark test as integration test"
    )
    config.addinivalue_line(
        "markers", "performance: mark test as performance test"
    )
    config.addinivalue_line(
        "markers", "security: mark test as security test"
    )