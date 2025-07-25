"""
Error handling and edge case tests for Panchangam API
"""
import pytest
import requests
from typing import Dict, Any


class TestErrorHandling:
    """Test suite for API error handling scenarios"""

    @pytest.mark.integration
    def test_missing_required_parameters(self, api_client: requests.Session, api_base_url: str):
        """Test missing required parameters return proper 400 errors"""
        test_cases = [
            # Missing date
            {"lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
            # Missing lat
            {"date": "2024-01-15", "lng": 77.5946, "tz": "Asia/Kolkata"},
            # Missing lng
            {"date": "2024-01-15", "lat": 12.9716, "tz": "Asia/Kolkata"},
            # Missing tz
            {"date": "2024-01-15", "lat": 12.9716, "lng": 77.5946},
            # Empty parameters
            {},
        ]
        
        for params in test_cases:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            
            assert response.status_code == 400
            assert response.headers["content-type"] == "application/json"
            
            data = response.json()
            assert "error" in data
            assert "message" in data
            assert "request_id" in data
            assert isinstance(data["error"], str)

    @pytest.mark.integration
    def test_invalid_parameter_types(self, api_client: requests.Session, api_base_url: str):
        """Test invalid parameter types return proper validation errors"""
        test_cases = [
            # Invalid date format
            {"date": "invalid-date", "lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
            {"date": "2024-13-45", "lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
            # Invalid latitude (out of range)
            {"date": "2024-01-15", "lat": 95.0, "lng": 77.5946, "tz": "Asia/Kolkata"},
            {"date": "2024-01-15", "lat": -95.0, "lng": 77.5946, "tz": "Asia/Kolkata"},
            # Invalid longitude (out of range)
            {"date": "2024-01-15", "lat": 12.9716, "lng": 185.0, "tz": "Asia/Kolkata"},
            {"date": "2024-01-15", "lat": 12.9716, "lng": -185.0, "tz": "Asia/Kolkata"},
            # Invalid timezone
            {"date": "2024-01-15", "lat": 12.9716, "lng": 77.5946, "tz": "Invalid/Timezone"},
            # Non-numeric lat/lng
            {"date": "2024-01-15", "lat": "not-a-number", "lng": 77.5946, "tz": "Asia/Kolkata"},
            {"date": "2024-01-15", "lat": 12.9716, "lng": "not-a-number", "tz": "Asia/Kolkata"},
        ]
        
        for params in test_cases:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            
            assert response.status_code == 400
            data = response.json()
            assert "error" in data
            assert "message" in data
            assert data["error"] in ["INVALID_ARGUMENT", "BAD_REQUEST"]

    @pytest.mark.integration
    def test_empty_parameter_values(self, api_client: requests.Session, api_base_url: str):
        """Test empty parameter values are handled correctly"""
        test_cases = [
            {"date": "", "lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
            {"date": "2024-01-15", "lat": "", "lng": 77.5946, "tz": "Asia/Kolkata"},
            {"date": "2024-01-15", "lat": 12.9716, "lng": "", "tz": "Asia/Kolkata"},
            {"date": "2024-01-15", "lat": 12.9716, "lng": 77.5946, "tz": ""},
        ]
        
        for params in test_cases:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            
            assert response.status_code == 400
            data = response.json()
            assert "error" in data
            assert "message" in data

    @pytest.mark.integration
    def test_extreme_coordinates(self, api_client: requests.Session, api_base_url: str):
        """Test extreme but valid coordinates"""
        extreme_cases = [
            # North/South poles
            {"date": "2024-01-15", "lat": 90.0, "lng": 0.0, "tz": "UTC"},
            {"date": "2024-01-15", "lat": -90.0, "lng": 0.0, "tz": "UTC"},
            # International Date Line
            {"date": "2024-01-15", "lat": 0.0, "lng": 180.0, "tz": "Pacific/Majuro"},
            {"date": "2024-01-15", "lat": 0.0, "lng": -180.0, "tz": "Pacific/Majuro"},
            # Equator
            {"date": "2024-01-15", "lat": 0.0, "lng": 0.0, "tz": "UTC"},
        ]
        
        for params in extreme_cases:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            
            # These should either succeed (200) or fail gracefully (400/500)
            assert response.status_code in [200, 400, 500]
            
            data = response.json()
            if response.status_code == 200:
                # Verify required fields are present
                assert "date" in data
                assert "tithi" in data
            else:
                # Verify error structure
                assert "error" in data
                assert "message" in data

    @pytest.mark.integration
    def test_invalid_optional_parameters(self, api_client: requests.Session, api_base_url: str):
        """Test invalid optional parameters are handled gracefully"""
        base_params = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata"
        }
        
        invalid_optional_cases = [
            # Invalid region
            {**base_params, "region": "InvalidRegion123"},
            # Invalid method
            {**base_params, "method": "InvalidMethod"},
            # Invalid locale
            {**base_params, "locale": "xx"},
            # Multiple invalid optionals
            {**base_params, "region": "Invalid", "method": "Invalid", "locale": "invalid"},
        ]
        
        for params in invalid_optional_cases:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            
            # Should either succeed (ignoring invalid optionals) or return 400
            assert response.status_code in [200, 400]
            
            data = response.json()
            if response.status_code == 400:
                assert "error" in data
                assert "message" in data

    @pytest.mark.integration
    def test_malformed_request_headers(self, api_client: requests.Session, api_base_url: str, valid_request_params: Dict[str, Any]):
        """Test malformed request headers are handled gracefully"""
        # Test with malformed Accept header
        malformed_headers = {
            "Accept": "invalid/content-type",
            "Content-Type": "text/plain",
            "Authorization": "Bearer invalid-token",
            "X-Custom-Header": "x" * 1000,  # Very long header
        }
        
        response = api_client.get(
            f"{api_base_url}/api/v1/panchangam", 
            params=valid_request_params,
            headers=malformed_headers
        )
        
        # Should handle gracefully - either succeed or return proper error
        assert response.status_code in [200, 400, 406]
        
        data = response.json()
        if response.status_code != 200:
            assert "error" in data

    @pytest.mark.integration
    def test_request_size_limits(self, api_client: requests.Session, api_base_url: str):
        """Test request size limits are enforced"""
        # Test with extremely long parameter values
        long_value = "x" * 10000  # 10KB string
        
        params = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
            "region": long_value,  # Very long parameter
        }
        
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
        
        # Should handle large requests gracefully
        assert response.status_code in [200, 400, 413, 414]  # 413=Payload Too Large, 414=URI Too Long

    @pytest.mark.integration
    def test_concurrent_error_scenarios(self, api_client: requests.Session, api_base_url: str):
        """Test error handling under concurrent load"""
        import concurrent.futures
        
        def make_invalid_request():
            # Make request with missing parameters
            return api_client.get(f"{api_base_url}/api/v1/panchangam", params={"invalid": "params"})
        
        # Make 10 concurrent invalid requests
        with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
            futures = [executor.submit(make_invalid_request) for _ in range(10)]
            responses = [future.result() for future in concurrent.futures.as_completed(futures)]
        
        # All should return consistent error responses
        for response in responses:
            assert response.status_code == 400
            data = response.json()
            assert "error" in data
            assert "message" in data
            assert "request_id" in data

    @pytest.mark.integration
    def test_error_response_consistency(self, api_client: requests.Session, api_base_url: str):
        """Test error responses have consistent structure"""
        error_scenarios = [
            # Missing parameters
            {"params": {"invalid": "params"}},
            # Invalid date
            {"params": {"date": "invalid", "lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"}},
            # Out of range coordinates
            {"params": {"date": "2024-01-15", "lat": 100, "lng": 77.5946, "tz": "Asia/Kolkata"}},
        ]
        
        for scenario in error_scenarios:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=scenario["params"])
            
            assert response.status_code == 400
            assert response.headers["content-type"] == "application/json"
            
            data = response.json()
            
            # Verify consistent error structure
            required_fields = ["error", "message", "request_id"]
            for field in required_fields:
                assert field in data, f"Missing {field} in error response"
            
            # Verify field types
            assert isinstance(data["error"], str)
            assert isinstance(data["message"], str)
            assert isinstance(data["request_id"], str)
            
            # Verify error code is valid
            valid_error_codes = ["INVALID_ARGUMENT", "BAD_REQUEST", "INTERNAL_SERVER_ERROR"]
            assert data["error"] in valid_error_codes

    @pytest.mark.performance
    def test_error_handling_performance(self, api_client: requests.Session, api_base_url: str):
        """Test error handling doesn't significantly impact performance"""
        import time
        
        # Test invalid request performance
        start_time = time.time()
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params={"invalid": "params"})
        end_time = time.time()
        
        error_response_time = (end_time - start_time) * 1000
        
        assert response.status_code == 400
        assert error_response_time < 100  # Error responses should be fast (<100ms)

    @pytest.mark.security
    def test_error_information_disclosure(self, api_client: requests.Session, api_base_url: str):
        """Test error messages don't disclose sensitive information"""
        # Test various error scenarios
        error_scenarios = [
            {"params": {"invalid": "params"}},
            {"params": {"date": "2024-01-15", "lat": "sql_injection_attempt", "lng": 77.5946, "tz": "Asia/Kolkata"}},
        ]
        
        sensitive_info_patterns = [
            "internal error",
            "stack trace",
            "database",
            "sql",
            "password",
            "token",
            "key",
            "secret",
            "connection",
            "localhost",
            "127.0.0.1",
        ]
        
        for scenario in error_scenarios:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=scenario["params"])
            
            assert response.status_code == 400
            data = response.json()
            
            # Check error message doesn't contain sensitive information
            error_message = data.get("message", "").lower()
            for pattern in sensitive_info_patterns:
                assert pattern not in error_message, f"Error message contains sensitive info: {pattern}"
            
            # Error messages should be user-friendly, not technical
            assert len(error_message) > 10  # Should have meaningful message
            assert not error_message.startswith("panic:")  # No Go panic messages
            assert "grpc" not in error_message.lower()  # No internal protocol details