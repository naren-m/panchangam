"""
Load Testing with Locust - Issue #81
Performance and load testing for Panchangam API
"""
import random
from datetime import datetime, timedelta
from locust import HttpUser, task, between, events
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class PanchangamUser(HttpUser):
    """Simulates a user accessing the Panchangam API"""

    # Wait between 1 and 3 seconds between tasks
    wait_time = between(1, 3)

    # Sample locations for testing
    locations = [
        {"name": "Bangalore", "lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
        {"name": "Mumbai", "lat": 19.0760, "lng": 72.8777, "tz": "Asia/Kolkata"},
        {"name": "Delhi", "lat": 28.6139, "lng": 77.2090, "tz": "Asia/Kolkata"},
        {"name": "Chennai", "lat": 13.0827, "lng": 80.2707, "tz": "Asia/Kolkata"},
        {"name": "Kolkata", "lat": 22.5726, "lng": 88.3639, "tz": "Asia/Kolkata"},
        {"name": "Hyderabad", "lat": 17.3850, "lng": 78.4867, "tz": "Asia/Kolkata"},
        {"name": "New York", "lat": 40.7128, "lng": -74.0060, "tz": "America/New_York"},
        {"name": "London", "lat": 51.5074, "lng": -0.1278, "tz": "Europe/London"},
        {"name": "Tokyo", "lat": 35.6762, "lng": 139.6503, "tz": "Asia/Tokyo"},
        {"name": "Sydney", "lat": -33.8688, "lng": 151.2093, "tz": "Australia/Sydney"},
    ]

    def on_start(self):
        """Called when a simulated user starts"""
        logger.info(f"User {self.environment.runner.user_count} started")

    @task(10)
    def get_panchangam_current_date(self):
        """Most common: Get panchangam for current date (weighted 10x)"""
        location = random.choice(self.locations)
        today = datetime.now().strftime("%Y-%m-%d")

        params = {
            "date": today,
            "lat": location["lat"],
            "lng": location["lng"],
            "tz": location["tz"],
        }

        with self.client.get(
            "/api/v1/panchangam",
            params=params,
            catch_response=True,
            name="/api/v1/panchangam [current_date]"
        ) as response:
            if response.status_code == 200:
                data = response.json()
                # Validate response has required fields
                if "tithi" in data and "nakshatra" in data and "sunrise_time" in data:
                    response.success()
                else:
                    response.failure("Missing required fields in response")
            else:
                response.failure(f"Failed with status code: {response.status_code}")

    @task(5)
    def get_panchangam_random_future_date(self):
        """Get panchangam for random future date (weighted 5x)"""
        location = random.choice(self.locations)

        # Random date within next 30 days
        days_ahead = random.randint(1, 30)
        future_date = (datetime.now() + timedelta(days=days_ahead)).strftime("%Y-%m-%d")

        params = {
            "date": future_date,
            "lat": location["lat"],
            "lng": location["lng"],
            "tz": location["tz"],
        }

        with self.client.get(
            "/api/v1/panchangam",
            params=params,
            catch_response=True,
            name="/api/v1/panchangam [future_date]"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Failed with status code: {response.status_code}")

    @task(3)
    def get_panchangam_random_past_date(self):
        """Get panchangam for random past date (weighted 3x)"""
        location = random.choice(self.locations)

        # Random date within past 365 days
        days_ago = random.randint(1, 365)
        past_date = (datetime.now() - timedelta(days=days_ago)).strftime("%Y-%m-%d")

        params = {
            "date": past_date,
            "lat": location["lat"],
            "lng": location["lng"],
            "tz": location["tz"],
        }

        with self.client.get(
            "/api/v1/panchangam",
            params=params,
            catch_response=True,
            name="/api/v1/panchangam [past_date]"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Failed with status code: {response.status_code}")

    @task(2)
    def get_panchangam_with_optional_params(self):
        """Get panchangam with optional parameters (weighted 2x)"""
        location = random.choice(self.locations)
        today = datetime.now().strftime("%Y-%m-%d")

        regions = ["global", "north", "south"]
        methods = ["traditional", "modern"]
        locales = ["en", "hi", "ta", "te"]

        params = {
            "date": today,
            "lat": location["lat"],
            "lng": location["lng"],
            "tz": location["tz"],
            "region": random.choice(regions),
            "method": random.choice(methods),
            "locale": random.choice(locales),
        }

        with self.client.get(
            "/api/v1/panchangam",
            params=params,
            catch_response=True,
            name="/api/v1/panchangam [with_optional_params]"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Failed with status code: {response.status_code}")

    @task(1)
    def health_check(self):
        """Health check endpoint (weighted 1x)"""
        with self.client.get(
            "/api/v1/health",
            catch_response=True,
            name="/api/v1/health"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Health check failed with status: {response.status_code}")


class StressTestUser(HttpUser):
    """Aggressive user for stress testing"""

    wait_time = between(0.1, 0.5)  # Much faster requests

    locations = PanchangamUser.locations

    @task
    def rapid_fire_requests(self):
        """Make rapid fire requests to stress test the system"""
        location = random.choice(self.locations)
        date = datetime.now().strftime("%Y-%m-%d")

        params = {
            "date": date,
            "lat": location["lat"],
            "lng": location["lng"],
            "tz": location["tz"],
        }

        self.client.get("/api/v1/panchangam", params=params)


class CacheTestUser(HttpUser):
    """User specifically for testing cache behavior"""

    wait_time = between(0.5, 1.5)

    # Use fixed location and date to test cache hits
    fixed_location = {"lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"}
    fixed_date = "2024-01-15"

    @task(20)
    def get_cached_panchangam(self):
        """Request same data repeatedly to test cache hits"""
        params = {
            "date": self.fixed_date,
            "lat": self.fixed_location["lat"],
            "lng": self.fixed_location["lng"],
            "tz": self.fixed_location["tz"],
        }

        with self.client.get(
            "/api/v1/panchangam",
            params=params,
            catch_response=True,
            name="/api/v1/panchangam [cache_test]"
        ) as response:
            if response.status_code == 200:
                # Check if response was cached
                cache_header = response.headers.get("X-Cache", "UNKNOWN")
                if cache_header == "HIT":
                    logger.debug("Cache HIT")
                elif cache_header == "MISS":
                    logger.debug("Cache MISS")
                response.success()
            else:
                response.failure(f"Failed with status code: {response.status_code}")

    @task(1)
    def check_cache_health(self):
        """Check cache health endpoint"""
        try:
            response = self.client.get("/api/v1/cache/health", catch_response=True)
            if response.status_code == 200:
                response.success()
            elif response.status_code == 404:
                # Cache might not be enabled
                response.success()
            else:
                response.failure(f"Unexpected status: {response.status_code}")
        except Exception as e:
            logger.debug(f"Cache health check failed: {e}")


# Event handlers for tracking performance metrics
@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    """Called when the test starts"""
    logger.info("=" * 80)
    logger.info("Load test starting - Issue #81 Performance Requirements:")
    logger.info("  - API Response Time: <500ms average")
    logger.info("  - Concurrent Requests: 50 requests in <5 seconds")
    logger.info("  - Data Consistency: 100% accuracy verification")
    logger.info("  - Error Recovery: <3 seconds for retry scenarios")
    logger.info("=" * 80)


@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    """Called when the test stops"""
    logger.info("=" * 80)
    logger.info("Load test completed")

    # Get stats
    stats = environment.stats

    # Calculate summary
    total_requests = stats.total.num_requests
    total_failures = stats.total.num_failures
    success_rate = ((total_requests - total_failures) / total_requests * 100) if total_requests > 0 else 0

    avg_response_time = stats.total.avg_response_time
    median_response_time = stats.total.median_response_time
    max_response_time = stats.total.max_response_time

    logger.info(f"Total Requests: {total_requests}")
    logger.info(f"Total Failures: {total_failures}")
    logger.info(f"Success Rate: {success_rate:.2f}%")
    logger.info(f"Average Response Time: {avg_response_time:.2f}ms")
    logger.info(f"Median Response Time: {median_response_time:.2f}ms")
    logger.info(f"Max Response Time: {max_response_time:.2f}ms")

    # Check against Issue #81 requirements
    logger.info("-" * 80)
    logger.info("Validation against Issue #81 Requirements:")

    # Requirement 1: Average response time <500ms
    if avg_response_time < 500:
        logger.info(f"✓ Average response time ({avg_response_time:.2f}ms) meets <500ms target")
    else:
        logger.warning(f"✗ Average response time ({avg_response_time:.2f}ms) exceeds 500ms target")

    # Requirement 2: Success rate should be high
    if success_rate >= 99:
        logger.info(f"✓ Success rate ({success_rate:.2f}%) is excellent")
    elif success_rate >= 95:
        logger.info(f"⚠ Success rate ({success_rate:.2f}%) is acceptable but could be better")
    else:
        logger.warning(f"✗ Success rate ({success_rate:.2f}%) is below acceptable threshold")

    logger.info("=" * 80)


# Custom shape for testing concurrent load (Issue #81: 50 concurrent requests)
from locust import LoadTestShape


class ConcurrentLoadShape(LoadTestShape):
    """
    Custom load shape to test specific concurrent load scenarios
    Tests Issue #81 requirement: 50 concurrent requests in <5 seconds
    """

    stages = [
        # Stage 1: Ramp up to 10 users in 10 seconds
        {"duration": 10, "users": 10, "spawn_rate": 1},
        # Stage 2: Spike to 50 users for 5 seconds (Issue #81 requirement)
        {"duration": 15, "users": 50, "spawn_rate": 10},
        # Stage 3: Sustain 50 users for 10 seconds
        {"duration": 25, "users": 50, "spawn_rate": 1},
        # Stage 4: Ramp down
        {"duration": 35, "users": 10, "spawn_rate": 5},
        # Stage 5: Cool down
        {"duration": 45, "users": 0, "spawn_rate": 5},
    ]

    def tick(self):
        """
        Returns the desired user count and spawn rate for the current time
        """
        run_time = self.get_run_time()

        for stage in self.stages:
            if run_time < stage["duration"]:
                return (stage["users"], stage["spawn_rate"])

        return None  # End of test
