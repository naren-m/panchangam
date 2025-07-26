"""
Health check endpoint tests
"""
import pytest
import requests
from typing import Dict, Any


@pytest.mark.smoke
def test_health_check_endpoint(api_client: requests.Session, api_base_url: str):
    """Test that health check endpoint returns proper response"""
    response = api_client.get(f"{api_base_url}/api/v1/health")
    
    assert response.status_code == 200
    assert response.headers["content-type"] == "application/json"
    
    data = response.json()
    assert data["status"] == "healthy"
    assert data["service"] == "panchangam-gateway"
    assert "timestamp" in data
    assert "version" in data


@pytest.mark.smoke
def test_health_check_with_request_id(api_client: requests.Session, api_base_url: str):
    """Test health check preserves custom request ID"""
    request_id = "test-health-123"
    headers = {"X-Request-Id": request_id}
    
    response = api_client.get(f"{api_base_url}/api/v1/health", headers=headers)
    
    assert response.status_code == 200
    assert response.headers.get("X-Request-Id") == request_id


@pytest.mark.performance
def test_health_check_response_time(api_client: requests.Session, api_base_url: str):
    """Test health check response time is acceptable"""
    import time
    
    start_time = time.time()
    response = api_client.get(f"{api_base_url}/api/v1/health")
    end_time = time.time()
    
    response_time = (end_time - start_time) * 1000  # Convert to milliseconds
    
    assert response.status_code == 200
    assert response_time < 100  # Should respond within 100ms


@pytest.mark.integration
def test_health_check_cors_headers(api_client: requests.Session, api_base_url: str):
    """Test CORS headers are present for allowed origins"""
    headers = {"Origin": "http://localhost:5173"}
    
    response = api_client.get(f"{api_base_url}/api/v1/health", headers=headers)
    
    assert response.status_code == 200
    assert "Access-Control-Allow-Origin" in response.headers
    assert response.headers["Access-Control-Allow-Origin"] == "http://localhost:5173"