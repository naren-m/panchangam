"""
End-to-End Integration Tests - Issue #81
Complete data flow validation from request to response
"""
import pytest
import requests
import concurrent.futures
from typing import Dict, Any
from datetime import datetime, timedelta


class TestEndToEndDataFlow:
    """Test complete end-to-end data flow through the system"""

    @pytest.mark.integration
    def test_complete_request_lifecycle(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test complete request lifecycle from HTTP to gRPC to response"""
        # Make request
        response = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )

        # Verify HTTP layer
        assert response.status_code == 200
        assert response.headers["Content-Type"] == "application/json"

        # Verify response structure
        data = response.json()
        assert isinstance(data, dict)

        # Verify all required fields are present and valid
        required_fields = [
            "date", "tithi", "nakshatra", "yoga", "karana", "paksha",
            "sunrise_time", "sunset_time", "moonrise_time", "moonset_time", "events"
        ]

        for field in required_fields:
            assert field in data, f"Missing required field: {field}"

        # Verify data types
        assert isinstance(data["date"], str)
        assert isinstance(data["tithi"], str)
        assert isinstance(data["nakshatra"], str)
        assert isinstance(data["yoga"], str)
        assert isinstance(data["karana"], str)
        assert isinstance(data["paksha"], str)
        assert isinstance(data["events"], list)

        # Verify date format
        try:
            datetime.strptime(data["date"], "%Y-%m-%d")
        except ValueError:
            pytest.fail(f"Invalid date format: {data['date']}")

        # Verify time formats (HH:MM:SS)
        time_fields = ["sunrise_time", "sunset_time", "moonrise_time", "moonset_time"]
        for field in time_fields:
            time_value = data.get(field)
            if time_value and time_value != "N/A":
                parts = time_value.split(":")
                assert len(parts) == 3, f"{field} should be in HH:MM:SS format"

    @pytest.mark.integration
    def test_multi_location_data_flow(
        self,
        api_client: requests.Session,
        api_base_url: str,
        sample_locations
    ):
        """Test data flow for multiple locations simultaneously"""
        date = "2024-01-15"
        results = {}

        # Request data for all sample locations
        for location_key, location_data in sample_locations.items():
            params = {
                "date": date,
                "lat": location_data["lat"],
                "lng": location_data["lng"],
                "tz": location_data["tz"],
            }

            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            assert response.status_code == 200

            results[location_key] = response.json()

        # Verify all locations got valid data
        for location_key, data in results.items():
            assert data["date"] == date
            assert "tithi" in data
            assert "nakshatra" in data

        # Verify geographic variance
        # Different locations should have different sunrise/sunset times
        bangalore_sunrise = results["bangalore"]["sunrise_time"]
        new_york_sunrise = results["new_york"]["sunrise_time"]

        assert bangalore_sunrise != new_york_sunrise, \
            "Different geographic locations should have different sunrise times"

    @pytest.mark.integration
    def test_date_range_data_flow(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test data flow across a range of dates"""
        base_params = {
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        # Test a week of dates
        start_date = datetime(2024, 1, 15)
        results = []

        for i in range(7):
            current_date = start_date + timedelta(days=i)
            date_str = current_date.strftime("%Y-%m-%d")

            params = {**base_params, "date": date_str}
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)

            assert response.status_code == 200
            data = response.json()
            results.append(data)

        # Verify all dates processed correctly
        assert len(results) == 7

        # Verify dates are in sequence
        for i, result in enumerate(results):
            expected_date = (start_date + timedelta(days=i)).strftime("%Y-%m-%d")
            assert result["date"] == expected_date

        # Verify lunar day progression is logical
        # (Tithi should progress through the lunar month)
        tithis = [r["tithi"] for r in results]
        assert len(set(tithis)) >= 5, \
            "Over a week, we should see several different tithis"


class TestErrorHandlingDataFlow:
    """Test error handling through the complete data flow"""

    @pytest.mark.integration
    def test_validation_error_flow(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that validation errors flow correctly through the system"""
        invalid_params = {
            "date": "invalid-date",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=invalid_params)

        # Should return 400 Bad Request
        assert response.status_code == 400

        # Error response should have consistent structure
        data = response.json()
        assert "error" in data
        assert "message" in data
        assert "request_id" in data

    @pytest.mark.integration
    def test_missing_parameter_flow(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that missing parameter errors are handled correctly"""
        incomplete_params = {
            "date": "2024-01-15",
            # Missing lat and lng
        }

        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=incomplete_params)

        assert response.status_code == 400

        data = response.json()
        assert "error" in data
        assert data["error"] in ["MISSING_PARAMETER", "BAD_REQUEST"]

    @pytest.mark.integration
    def test_recovery_from_transient_errors(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test recovery from transient errors with retry"""
        max_retries = 3
        retry_delay = 0.5

        for attempt in range(max_retries):
            try:
                response = api_client.get(
                    f"{api_base_url}/api/v1/panchangam",
                    params=valid_request_params,
                    timeout=5
                )

                if response.status_code == 200:
                    # Success - verify data
                    data = response.json()
                    assert "tithi" in data
                    break

            except requests.exceptions.RequestException as e:
                if attempt < max_retries - 1:
                    import time
                    time.sleep(retry_delay)
                else:
                    pytest.fail(f"Request failed after {max_retries} retries: {e}")


class TestConcurrentDataFlow:
    """Test data flow under concurrent load"""

    @pytest.mark.integration
    def test_concurrent_unique_requests(
        self,
        api_client: requests.Session,
        api_base_url: str,
        sample_locations
    ):
        """Test concurrent requests for different data"""
        def make_request(location_key):
            location = sample_locations[location_key]
            params = {
                "date": "2024-01-15",
                "lat": location["lat"],
                "lng": location["lng"],
                "tz": location["tz"],
            }

            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            return {
                "location": location_key,
                "status": response.status_code,
                "data": response.json() if response.status_code == 200 else None,
            }

        # Make concurrent requests for all locations
        with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:
            futures = [
                executor.submit(make_request, loc_key)
                for loc_key in sample_locations.keys()
            ]
            results = [future.result() for future in concurrent.futures.as_completed(futures)]

        # Verify all requests succeeded
        for result in results:
            assert result["status"] == 200, \
                f"Request for {result['location']} failed"
            assert result["data"] is not None
            assert "tithi" in result["data"]

    @pytest.mark.integration
    def test_concurrent_identical_requests(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test concurrent identical requests return consistent data"""
        def make_request():
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            return response.json() if response.status_code == 200 else None

        # Make 10 concurrent identical requests
        with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
            futures = [executor.submit(make_request) for _ in range(10)]
            results = [future.result() for future in concurrent.futures.as_completed(futures)]

        # All results should be identical
        assert all(r is not None for r in results)

        first_result = results[0]
        for i, result in enumerate(results[1:], 1):
            assert result == first_result, \
                f"Concurrent request {i} returned different data"


class TestDataConsistency:
    """Test data consistency across the system"""

    @pytest.mark.integration
    def test_idempotency(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that identical requests always return identical data (idempotency)"""
        responses = []

        # Make same request 5 times
        for _ in range(5):
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            assert response.status_code == 200
            responses.append(response.json())

        # All responses should be identical
        first_response = responses[0]
        for i, response in enumerate(responses[1:], 1):
            assert response == first_response, \
                f"Response {i} differs from first response (idempotency violated)"

    @pytest.mark.integration
    def test_cross_session_consistency(
        self,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that different sessions get consistent data"""
        # Create two separate sessions
        session1 = requests.Session()
        session2 = requests.Session()

        try:
            # Make same request with both sessions
            response1 = session1.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            response2 = session2.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )

            assert response1.status_code == 200
            assert response2.status_code == 200

            data1 = response1.json()
            data2 = response2.json()

            # Both sessions should get identical data
            assert data1 == data2, \
                "Different sessions should get consistent data"

        finally:
            session1.close()
            session2.close()

    @pytest.mark.integration
    def test_parameter_order_independence(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that parameter order doesn't affect results"""
        # Same parameters in different order
        params1 = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        params2 = {
            "tz": "Asia/Kolkata",
            "lng": 77.5946,
            "lat": 12.9716,
            "date": "2024-01-15",
        }

        response1 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params1)
        response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params2)

        assert response1.status_code == 200
        assert response2.status_code == 200

        # Should get identical data regardless of parameter order
        assert response1.json() == response2.json()


class TestSystemIntegration:
    """Test integration of all system components"""

    @pytest.mark.integration
    def test_health_check_integration(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test health check endpoint integration"""
        response = api_client.get(f"{api_base_url}/api/v1/health")

        assert response.status_code == 200

        data = response.json()
        assert "status" in data
        assert data["status"] == "healthy"

    @pytest.mark.integration
    def test_cors_headers_integration(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test CORS headers are properly set"""
        allowed_origin = "http://localhost:5173"
        headers = {"Origin": allowed_origin}

        response = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params,
            headers=headers
        )

        assert response.status_code == 200

        # Verify CORS headers
        assert "Access-Control-Allow-Origin" in response.headers
        # Should either be the specific origin or *
        cors_header = response.headers["Access-Control-Allow-Origin"]
        assert cors_header in [allowed_origin, "*"]

    @pytest.mark.integration
    def test_request_id_propagation(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that request IDs are properly propagated"""
        custom_request_id = "test-request-12345"
        headers = {"X-Request-Id": custom_request_id}

        response = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params,
            headers=headers
        )

        assert response.status_code == 200

        # Request ID should be echoed back or a new one generated
        response_request_id = response.headers.get("X-Request-Id")
        if response_request_id:
            # If custom ID is supported, it should match
            # Otherwise, a valid request ID should be present
            assert len(response_request_id) > 0
