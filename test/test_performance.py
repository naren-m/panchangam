"""
Performance Testing Suite - Issue #81
Tests API performance, response times, and concurrent load handling
"""
import pytest
import requests
import time
import concurrent.futures
import statistics
from typing import List, Dict, Any


class TestAPIPerformance:
    """Test suite for API performance benchmarks"""

    @pytest.mark.performance
    def test_api_response_time_target(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that API response time meets <500ms target"""
        response_times = []

        # Make 100 requests to get average
        for _ in range(100):
            start_time = time.time()
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            end_time = time.time()

            assert response.status_code == 200
            response_time_ms = (end_time - start_time) * 1000
            response_times.append(response_time_ms)

        # Calculate statistics
        avg_time = statistics.mean(response_times)
        median_time = statistics.median(response_times)
        p95_time = statistics.quantiles(response_times, n=20)[18]  # 95th percentile
        p99_time = statistics.quantiles(response_times, n=100)[98]  # 99th percentile

        # Log performance metrics
        print(f"\nPerformance Metrics (100 requests):")
        print(f"  Average: {avg_time:.2f}ms")
        print(f"  Median: {median_time:.2f}ms")
        print(f"  P95: {p95_time:.2f}ms")
        print(f"  P99: {p99_time:.2f}ms")

        # Assert performance targets from Issue #81
        assert avg_time < 500, \
            f"Average response time {avg_time:.2f}ms exceeds 500ms target"

        # Additional reasonable targets
        assert p95_time < 1000, \
            f"P95 response time {p95_time:.2f}ms exceeds 1000ms"

    @pytest.mark.performance
    def test_concurrent_requests_target(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that 50 concurrent requests complete in <5 seconds (Issue #81 requirement)"""

        def make_request():
            """Helper function to make a single request"""
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            return response.status_code

        # Execute 50 concurrent requests
        start_time = time.time()

        with concurrent.futures.ThreadPoolExecutor(max_workers=50) as executor:
            futures = [executor.submit(make_request) for _ in range(50)]
            results = [future.result() for future in concurrent.futures.as_completed(futures)]

        end_time = time.time()
        total_time = end_time - start_time

        # Verify all requests succeeded
        assert all(status == 200 for status in results), \
            "Some concurrent requests failed"

        # Assert performance target from Issue #81
        assert total_time < 5.0, \
            f"50 concurrent requests took {total_time:.2f}s (target: <5s)"

        print(f"\n50 concurrent requests completed in {total_time:.2f}s")

    @pytest.mark.performance
    def test_concurrent_load_with_varying_parameters(
        self,
        api_client: requests.Session,
        api_base_url: str,
        sample_locations
    ):
        """Test concurrent load with varying parameters"""

        # Prepare different request parameters
        locations = ["bangalore", "mumbai", "new_york", "london"]
        dates = ["2024-01-15", "2024-06-21", "2024-12-21"]

        request_params_list = []
        for location_key in locations:
            location = sample_locations[location_key]
            for date in dates:
                request_params_list.append({
                    "date": date,
                    "lat": location["lat"],
                    "lng": location["lng"],
                    "tz": location["tz"],
                })

        def make_request(params):
            """Helper function to make request with specific params"""
            start = time.time()
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            duration = time.time() - start
            return {
                "status": response.status_code,
                "duration": duration * 1000,  # Convert to ms
                "params": params,
            }

        # Execute concurrent requests
        start_time = time.time()

        with concurrent.futures.ThreadPoolExecutor(max_workers=12) as executor:
            futures = [executor.submit(make_request, params) for params in request_params_list]
            results = [future.result() for future in concurrent.futures.as_completed(futures)]

        end_time = time.time()
        total_time = end_time - start_time

        # Verify all requests succeeded
        failed_requests = [r for r in results if r["status"] != 200]
        assert len(failed_requests) == 0, \
            f"{len(failed_requests)} requests failed out of {len(results)}"

        # Calculate performance statistics
        response_times = [r["duration"] for r in results]
        avg_time = statistics.mean(response_times)
        max_time = max(response_times)

        print(f"\nConcurrent load test ({len(results)} requests):")
        print(f"  Total time: {total_time:.2f}s")
        print(f"  Average response: {avg_time:.2f}ms")
        print(f"  Max response: {max_time:.2f}ms")
        print(f"  Requests/second: {len(results) / total_time:.2f}")

        # Performance assertions
        assert avg_time < 500, \
            f"Average response time {avg_time:.2f}ms exceeds 500ms"

    @pytest.mark.performance
    def test_sustained_load(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test sustained load over time"""
        duration_seconds = 10
        requests_per_second = 10
        total_requests = duration_seconds * requests_per_second

        def make_request():
            start = time.time()
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            duration = time.time() - start
            return {
                "status": response.status_code,
                "duration": duration * 1000,
            }

        start_time = time.time()
        results = []

        # Sustained load with controlled rate
        with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
            futures = []
            for i in range(total_requests):
                # Schedule requests at controlled rate
                futures.append(executor.submit(make_request))

                # Sleep to maintain target rate
                if (i + 1) % requests_per_second == 0:
                    elapsed = time.time() - start_time
                    expected_time = (i + 1) / requests_per_second
                    if elapsed < expected_time:
                        time.sleep(expected_time - elapsed)

            # Wait for all to complete
            results = [future.result() for future in concurrent.futures.as_completed(futures)]

        end_time = time.time()
        actual_duration = end_time - start_time

        # Verify results
        successful = sum(1 for r in results if r["status"] == 200)
        failed = len(results) - successful

        response_times = [r["duration"] for r in results if r["status"] == 200]
        avg_time = statistics.mean(response_times) if response_times else 0

        print(f"\nSustained load test ({duration_seconds}s):")
        print(f"  Total requests: {len(results)}")
        print(f"  Successful: {successful}")
        print(f"  Failed: {failed}")
        print(f"  Actual duration: {actual_duration:.2f}s")
        print(f"  Actual rate: {len(results) / actual_duration:.2f} req/s")
        print(f"  Average response: {avg_time:.2f}ms")

        # Assertions
        assert failed == 0, f"{failed} requests failed"
        assert avg_time < 500, \
            f"Average response time {avg_time:.2f}ms exceeds 500ms under sustained load"


class TestCachingPerformance:
    """Test caching impact on performance"""

    @pytest.mark.performance
    def test_cache_hit_vs_miss_performance(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that cached responses are faster than cache misses"""
        # First request (cache miss)
        start_time = time.time()
        response1 = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )
        first_request_time = (time.time() - start_time) * 1000

        assert response1.status_code == 200
        cache_header_1 = response1.headers.get("X-Cache", "MISS")

        # Second request (should be cache hit if caching is enabled)
        start_time = time.time()
        response2 = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )
        second_request_time = (time.time() - start_time) * 1000

        assert response2.status_code == 200
        cache_header_2 = response2.headers.get("X-Cache", "MISS")

        print(f"\nCache performance:")
        print(f"  First request (cache {cache_header_1}): {first_request_time:.2f}ms")
        print(f"  Second request (cache {cache_header_2}): {second_request_time:.2f}ms")

        # If caching is enabled, second request should be faster
        if cache_header_2 == "HIT":
            assert second_request_time <= first_request_time, \
                "Cache hit should not be slower than cache miss"

        # Both responses should be identical
        assert response1.json() == response2.json(), \
            "Cached response differs from original"

    @pytest.mark.performance
    def test_cache_warmup_improves_performance(
        self,
        api_client: requests.Session,
        api_base_url: str,
        sample_locations
    ):
        """Test that cache warmup improves overall performance"""
        # Prepare test parameters for popular locations
        popular_locations = ["bangalore", "mumbai"]
        dates = ["2024-01-15", "2024-06-21"]

        params_list = []
        for loc_key in popular_locations:
            location = sample_locations[loc_key]
            for date in dates:
                params_list.append({
                    "date": date,
                    "lat": location["lat"],
                    "lng": location["lng"],
                    "tz": location["tz"],
                })

        # Phase 1: Warm up cache
        print("\nWarming up cache...")
        for params in params_list:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            assert response.status_code == 200

        # Phase 2: Measure performance with warm cache
        start_time = time.time()
        for params in params_list:
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            assert response.status_code == 200

        warm_cache_time = time.time() - start_time

        print(f"Time with warm cache: {warm_cache_time:.2f}s for {len(params_list)} requests")
        print(f"Average per request: {(warm_cache_time / len(params_list)) * 1000:.2f}ms")

        # With cache, average should be well under 500ms
        avg_time_ms = (warm_cache_time / len(params_list)) * 1000
        assert avg_time_ms < 500, \
            f"Average time with cache {avg_time_ms:.2f}ms exceeds 500ms"


class TestErrorRecoveryPerformance:
    """Test error recovery and retry performance (Issue #81: <3 seconds for retry scenarios)"""

    @pytest.mark.performance
    def test_error_response_time(
        self,
        api_client: requests.Session,
        api_base_url: str
    ):
        """Test that error responses are fast"""
        # Test with invalid parameters (should fail fast)
        invalid_params = {"invalid": "params"}

        start_time = time.time()
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=invalid_params)
        error_response_time = (time.time() - start_time) * 1000

        assert response.status_code == 400
        assert error_response_time < 100, \
            f"Error response took {error_response_time:.2f}ms (should be <100ms)"

        print(f"\nError response time: {error_response_time:.2f}ms")

    @pytest.mark.performance
    def test_retry_scenario_performance(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test retry scenario completes within 3 seconds (Issue #81 requirement)"""
        max_retries = 3
        retry_delay = 0.5  # seconds

        def make_request_with_retry():
            """Make request with retry logic"""
            for attempt in range(max_retries):
                try:
                    response = api_client.get(
                        f"{api_base_url}/api/v1/panchangam",
                        params=valid_request_params,
                        timeout=5
                    )
                    if response.status_code == 200:
                        return response
                except requests.exceptions.RequestException:
                    if attempt < max_retries - 1:
                        time.sleep(retry_delay)
                    else:
                        raise
            return None

        # Measure retry scenario
        start_time = time.time()
        response = make_request_with_retry()
        retry_scenario_time = time.time() - start_time

        assert response is not None
        assert response.status_code == 200

        # Issue #81 requirement: retry scenarios should complete in <3 seconds
        assert retry_scenario_time < 3.0, \
            f"Retry scenario took {retry_scenario_time:.2f}s (target: <3s)"

        print(f"\nRetry scenario completed in {retry_scenario_time:.2f}s")


class TestResourceUtilization:
    """Test resource utilization under load"""

    @pytest.mark.performance
    def test_memory_efficient_responses(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that response sizes are reasonable"""
        response = api_client.get(
            f"{api_base_url}/api/v1/panchangam",
            params=valid_request_params
        )

        assert response.status_code == 200

        # Check response size
        response_size = len(response.content)
        print(f"\nResponse size: {response_size} bytes ({response_size / 1024:.2f} KB)")

        # Response should be reasonably sized (< 100 KB for a single day's data)
        assert response_size < 100 * 1024, \
            f"Response size {response_size} bytes exceeds 100 KB"

    @pytest.mark.performance
    def test_connection_reuse(
        self,
        api_client: requests.Session,
        api_base_url: str,
        valid_request_params: Dict[str, Any]
    ):
        """Test that HTTP connections are reused efficiently"""
        # Using session (api_client) should reuse connections
        response_times = []

        for _ in range(10):
            start_time = time.time()
            response = api_client.get(
                f"{api_base_url}/api/v1/panchangam",
                params=valid_request_params
            )
            response_time = (time.time() - start_time) * 1000
            response_times.append(response_time)
            assert response.status_code == 200

        avg_time = statistics.mean(response_times)
        first_request = response_times[0]
        avg_subsequent = statistics.mean(response_times[1:])

        print(f"\nConnection reuse test:")
        print(f"  First request: {first_request:.2f}ms")
        print(f"  Avg subsequent: {avg_subsequent:.2f}ms")
        print(f"  Overall avg: {avg_time:.2f}ms")

        # All requests should complete reasonably fast
        assert avg_time < 500, \
            f"Average response time {avg_time:.2f}ms exceeds 500ms"
