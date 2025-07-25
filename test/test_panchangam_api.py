"""
Panchangam API endpoint tests
"""
import pytest
import requests
from typing import Dict, Any
from datetime import datetime, timedelta


class TestPanchangamAPI:
    """Test suite for Panchangam API endpoint"""

    @pytest.mark.smoke
    def test_valid_panchangam_request(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test valid panchangam request returns proper data"""
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)
        
        assert response.status_code == 200
        assert response.headers["content-type"] == "application/json"
        
        data = response.json()
        
        # Verify required fields
        assert data["date"] == valid_request_params["date"]
        assert "tithi" in data
        assert "nakshatra" in data
        assert "yoga" in data
        assert "karana" in data
        assert "sunrise_time" in data
        assert "sunset_time" in data
        assert "events" in data
        assert isinstance(data["events"], list)

    @pytest.mark.integration
    def test_multiple_locations(self, api_client: requests.Session, api_base_url: str, sample_locations: Dict[str, Any]):
        """Test panchangam calculation for different locations"""
        date = "2024-01-15"
        
        for location_name, location in sample_locations.items():
            params = {
                "date": date,
                "lat": location["lat"],
                "lng": location["lng"],
                "tz": location["tz"]
            }
            
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            
            assert response.status_code == 200, f"Failed for {location_name}"
            
            data = response.json()
            assert data["date"] == date
            
            # Sunrise/sunset times should be different for different locations
            assert "sunrise_time" in data
            assert "sunset_time" in data

    @pytest.mark.integration
    def test_date_range(self, api_client: requests.Session, api_base_url: str, sample_locations: Dict[str, Any]):
        """Test panchangam calculation for different dates"""
        location = sample_locations["bangalore"]
        base_date = datetime(2024, 1, 15)
        
        for i in range(7):  # Test 7 consecutive days
            test_date = base_date + timedelta(days=i)
            params = {
                "date": test_date.strftime("%Y-%m-%d"),
                "lat": location["lat"],
                "lng": location["lng"],
                "tz": location["tz"]
            }
            
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            
            assert response.status_code == 200
            data = response.json()
            assert data["date"] == test_date.strftime("%Y-%m-%d")

    @pytest.mark.integration
    def test_optional_parameters(self, api_client: requests.Session, api_base_url: str):
        """Test request with optional parameters"""
        params = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
            "region": "Karnataka",
            "method": "Drik",
            "locale": "en"
        }
        
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
        
        assert response.status_code == 200
        data = response.json()
        assert data["date"] == "2024-01-15"

    @pytest.mark.performance
    def test_response_time(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test API response time is within acceptable limits"""
        import time
        
        # Warm up request
        api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)
        
        # Measure actual request
        start_time = time.time()
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)
        end_time = time.time()
        
        response_time = (end_time - start_time) * 1000  # Convert to milliseconds
        
        assert response.status_code == 200
        assert response_time < 500  # Should respond within 500ms

    @pytest.mark.performance
    def test_concurrent_requests(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test API can handle concurrent requests"""
        import concurrent.futures
        import time
        
        def make_request():
            return api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)
        
        start_time = time.time()
        
        # Make 10 concurrent requests
        with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
            futures = [executor.submit(make_request) for _ in range(10)]
            responses = [future.result() for future in concurrent.futures.as_completed(futures)]
        
        end_time = time.time()
        total_time = (end_time - start_time) * 1000
        
        # All requests should succeed
        for response in responses:
            assert response.status_code == 200
        
        # Total time should be reasonable (not much slower than single request)
        assert total_time < 2000  # 10 concurrent requests should complete within 2 seconds

    @pytest.mark.integration
    def test_cache_headers(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test that proper cache headers are set"""
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)
        
        assert response.status_code == 200
        assert "Cache-Control" in response.headers
        assert "max-age=300" in response.headers["Cache-Control"]

    @pytest.mark.integration
    def test_request_id_tracking(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test request ID tracking works correctly"""
        custom_id = "test-panchangam-456"
        headers = {"X-Request-Id": custom_id}
        
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params, headers=headers)
        
        assert response.status_code == 200
        assert response.headers.get("X-Request-Id") == custom_id
        
        # Test auto-generated request ID
        response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)
        assert response2.status_code == 200
        assert "X-Request-Id" in response2.headers
        assert response2.headers["X-Request-Id"].startswith("req_")

    @pytest.mark.security
    def test_cors_configuration(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test CORS configuration for different origins"""
        # Test allowed origin
        allowed_origins = [
            "http://localhost:5173",
            "http://localhost:3000"
        ]
        
        for origin in allowed_origins:
            headers = {"Origin": origin}
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params, headers=headers)
            
            assert response.status_code == 200
            assert response.headers.get("Access-Control-Allow-Origin") == origin
            assert "X-Request-Id" in response.headers.get("Access-Control-Expose-Headers", "")

    @pytest.mark.integration
    def test_method_not_allowed(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test that only GET method is allowed"""
        # Test POST method
        response = api_client.post(f"{api_base_url}/api/v1/panchangam", json=valid_request_params)
        assert response.status_code == 405
        
        # Test PUT method
        response = api_client.put(f"{api_base_url}/api/v1/panchangam", json=valid_request_params)
        assert response.status_code == 405
        
        # Test DELETE method
        response = api_client.delete(f"{api_base_url}/api/v1/panchangam")
        assert response.status_code == 405