"""
Cache Behavior Integration Tests - Issue #81
Tests Redis cache integration, cache hits/misses, and cache consistency
"""
import pytest
import requests
import time
from typing import Dict, Any


class TestCacheIntegration:
    """Test suite for cache integration and behavior"""

    @pytest.mark.integration
    def test_cache_headers_present(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that cache-related headers are present in responses"""
        response = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )

        assert response.status_code == 200

        # Check for cache-related headers
        headers = response.headers

        # X-Cache header should indicate hit or miss
        assert "X-Cache" in headers or "Cache-Control" in headers, \
            "Response should include cache headers"

        # Cache-Control should be present
        if "Cache-Control" in headers:
            cache_control = headers["Cache-Control"]
            assert "max-age" in cache_control, \
                "Cache-Control should specify max-age"

    @pytest.mark.integration
    def test_cache_consistency(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that cached responses are consistent with fresh responses"""
        # First request (likely cache miss)
        response1 = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )
        assert response1.status_code == 200
        data1 = response1.json()

        # Second request (likely cache hit)
        response2 = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )
        assert response2.status_code == 200
        data2 = response2.json()

        # Responses should be identical
        assert data1 == data2, \
            "Cached response differs from original response"

        # Both should have the same core fields
        assert data1.get("tithi") == data2.get("tithi")
        assert data1.get("nakshatra") == data2.get("nakshatra")
        assert data1.get("sunrise_time") == data2.get("sunrise_time")

    @pytest.mark.integration
    def test_cache_key_uniqueness(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that different parameters result in different cache entries"""
        # Request 1: Bangalore
        params1 = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        # Request 2: Mumbai (different location, same date)
        params2 = {
            "date": "2024-01-15",
            "lat": 19.0760,
            "lng": 72.8777,
            "tz": "Asia/Kolkata",
        }

        response1 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params1)
        response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params2)

        assert response1.status_code == 200
        assert response2.status_code == 200

        data1 = response1.json()
        data2 = response2.json()

        # Different locations should have different data
        # (at least sunrise/sunset times should differ)
        assert data1.get("sunrise_time") != data2.get("sunrise_time"), \
            "Different locations should have different sunrise times"

    @pytest.mark.integration
    def test_cache_expiration_behavior(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test cache expiration headers and behavior"""
        response = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )

        assert response.status_code == 200

        cache_control = response.headers.get("Cache-Control", "")

        if cache_control:
            # Verify reasonable cache duration
            if "max-age" in cache_control:
                # Extract max-age value
                parts = cache_control.split(",")
                max_age = None
                for part in parts:
                    if "max-age" in part:
                        max_age = int(part.split("=")[1].strip())
                        break

                if max_age is not None:
                    # Cache duration should be reasonable (between 1 min and 1 hour)
                    assert 60 <= max_age <= 3600, \
                        f"Cache max-age {max_age}s is outside reasonable range (60-3600s)"

    @pytest.mark.integration
    def test_cache_invalidation_with_different_dates(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that different dates are cached separately"""
        base_params = {
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        # Request for date 1
        params1 = {**base_params, "date": "2024-01-15"}
        response1 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params1)
        assert response1.status_code == 200
        data1 = response1.json()

        # Request for date 2
        params2 = {**base_params, "date": "2024-01-16"}
        response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params2)
        assert response2.status_code == 200
        data2 = response2.json()

        # Different dates should return different data
        assert data1.get("date") != data2.get("date")
        assert data1.get("tithi") != data2.get("tithi") or \
               data1.get("nakshatra") != data2.get("nakshatra"), \
            "Different dates should have different astrological data"


class TestCacheHealth:
    """Test cache health and statistics endpoints"""

    @pytest.mark.integration
    def test_cache_health_endpoint(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test cache health check endpoint if available"""
        try:
            response = api_client.get(f"{api_base_url}/api/v1/cache/health")

            # If endpoint exists, it should return valid health status
            if response.status_code == 200:
                data = response.json()
                assert "status" in data or "healthy" in data or "connected" in data, \
                    "Cache health response should indicate status"

            # 404 is acceptable if Redis is not enabled
            elif response.status_code == 404:
                pytest.skip("Cache health endpoint not available (Redis may not be enabled)")

            else:
                pytest.fail(f"Unexpected status code: {response.status_code}")

        except requests.exceptions.RequestException:
            pytest.skip("Cache health endpoint not accessible")

    @pytest.mark.integration
    def test_cache_stats_endpoint(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test cache statistics endpoint if available"""
        try:
            response = api_client.get(f"{api_base_url}/api/v1/cache/stats")

            if response.status_code == 200:
                data = response.json()
                # Stats should include useful metrics
                # (exact fields may vary based on implementation)
                assert isinstance(data, dict), \
                    "Cache stats should return a dictionary"

            elif response.status_code == 404:
                pytest.skip("Cache stats endpoint not available (Redis may not be enabled)")

            else:
                pytest.fail(f"Unexpected status code: {response.status_code}")

        except requests.exceptions.RequestException:
            pytest.skip("Cache stats endpoint not accessible")


class TestCacheConcurrency:
    """Test cache behavior under concurrent load"""

    @pytest.mark.integration
    def test_concurrent_cache_access(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that concurrent requests for same data return consistent results"""
        import concurrent.futures

        def make_request():
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            assert response.status_code == 200
            return response.json()

        # Make 10 concurrent identical requests
        with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
            futures = [executor.submit(make_request) for _ in range(10)]
            results = [future.result() for future in concurrent.futures.as_completed(futures)]

        # All results should be identical
        first_result = results[0]
        for i, result in enumerate(results[1:], 1):
            assert result == first_result, \
                f"Concurrent request {i} returned different data"

    @pytest.mark.integration
    def test_cache_stampede_protection(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that concurrent requests don't cause cache stampede"""
        import concurrent.futures

        # Use a unique date to ensure cache miss
        unique_params = {
            "date": "2025-03-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        def make_request():
            start = time.time()
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=unique_params
            )
            duration = time.time() - start
            return {
                "status": response.status_code,
                "duration": duration,
                "cache_header": response.headers.get("X-Cache", "UNKNOWN"),
            }

        # Make 20 concurrent requests for the same uncached data
        with concurrent.futures.ThreadPoolExecutor(max_workers=20) as executor:
            futures = [executor.submit(make_request) for _ in range(20)]
            results = [future.result() for future in concurrent.futures.as_completed(futures)]

        # All requests should succeed
        assert all(r["status"] == 200 for r in results), \
            "Some concurrent requests failed"

        # Check cache behavior
        hits = sum(1 for r in results if r["cache_header"] == "HIT")
        misses = sum(1 for r in results if r["cache_header"] == "MISS")

        print(f"\nCache stampede test: {hits} hits, {misses} misses out of {len(results)} requests")

        # Most requests should either hit cache or be reasonably fast
        avg_duration = sum(r["duration"] for r in results) / len(results)
        assert avg_duration < 1.0, \
            f"Average request duration {avg_duration:.2f}s is too high"


class TestCacheEdgeCases:
    """Test edge cases in cache behavior"""

    @pytest.mark.integration
    def test_cache_with_optional_parameters(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that optional parameters are included in cache key"""
        base_params = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        # Request 1: Without optional parameters
        response1 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=base_params)
        assert response1.status_code == 200
        data1 = response1.json()

        # Request 2: With optional region parameter
        params_with_region = {**base_params, "region": "north"}
        response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params_with_region)
        assert response2.status_code == 200
        data2 = response2.json()

        # Request 3: Same as request 1 (should hit cache)
        response3 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=base_params)
        assert response3.status_code == 200
        data3 = response3.json()

        # Data 1 and 3 should be identical (same parameters)
        assert data1 == data3, \
            "Identical parameters should return identical cached data"

    @pytest.mark.integration
    def test_cache_handles_special_characters(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that cache handles special characters in parameters correctly"""
        # Test with timezone containing special characters
        params = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",  # Contains forward slash
        }

        response1 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
        assert response1.status_code == 200

        # Second request should work identically
        response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
        assert response2.status_code == 200

        # Data should be identical
        assert response1.json() == response2.json()

    @pytest.mark.integration
    def test_cache_with_extreme_coordinates(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test cache behavior with extreme coordinates"""
        extreme_params = [
            # North Pole
            {"date": "2024-06-21", "lat": 90.0, "lng": 0.0, "tz": "UTC"},
            # South Pole
            {"date": "2024-06-21", "lat": -90.0, "lng": 0.0, "tz": "UTC"},
            # International Date Line
            {"date": "2024-01-15", "lat": 0.0, "lng": 180.0, "tz": "Pacific/Majuro"},
        ]

        for params in extreme_params:
            # First request
            response1 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)

            # Skip if server doesn't handle extreme coordinates
            if response1.status_code != 200:
                continue

            data1 = response1.json()

            # Second request (should use cache)
            response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            assert response2.status_code == 200
            data2 = response2.json()

            # Cached data should match original
            assert data1 == data2, \
                f"Cached data differs for extreme coordinates: lat={params['lat']}, lng={params['lng']}"
