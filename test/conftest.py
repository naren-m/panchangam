"""
Pytest configuration and fixtures for Panchangam API tests
"""
import pytest
import requests
import subprocess
import time
import os
from pathlib import Path
from typing import Generator
from urllib.parse import urlparse


@pytest.fixture(scope="session")
def api_base_url() -> str:
    """Base URL for the API gateway"""
    return os.getenv("PANCHANGAM_API_URL", "http://localhost:8080")


@pytest.fixture(scope="session")
def start_servers(api_base_url: str) -> Generator[None, None, None]:
    """Start the gRPC and Gateway servers for testing"""
    if os.getenv("SKIP_SERVER_START", "false").lower() == "true":
        yield
        return

    grpc_port = os.getenv("PANCHANGAM_GRPC_PORT", "50051")
    api_url = urlparse(api_base_url)
    http_port = str(api_url.port or (443 if api_url.scheme == "https" else 80))

    project_root = Path(__file__).resolve().parent.parent
    test_bin_dir = project_root / "tmp" / "test-bin"
    test_bin_dir.mkdir(parents=True, exist_ok=True)

    # Build servers from maintained command packages.
    subprocess.run(["go", "build", "-o", str(test_bin_dir / "grpc-server"), "./cmd/server"], check=True, cwd=project_root)
    subprocess.run(["go", "build", "-o", str(test_bin_dir / "gateway-server"), "./cmd/gateway"], check=True, cwd=project_root)

    # Start gRPC server
    grpc_proc = subprocess.Popen([
        str(test_bin_dir / "grpc-server"),
        f"--grpc-port={grpc_port}",
    ], cwd=project_root)
    time.sleep(2)

    # Start Gateway server
    gateway_proc = subprocess.Popen([
        str(test_bin_dir / "gateway-server"),
        f"--grpc-endpoint=localhost:{grpc_port}",
        f"--http-port={http_port}",
    ], cwd=project_root)
    time.sleep(3)

    # Verify servers are running
    max_retries = 10
    for i in range(max_retries):
        try:
            response = requests.get(f"{api_base_url}/api/v1/health")
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
