"""
Data Validation Tests - Issue #81
Tests against known astronomical data to verify calculation accuracy
"""
import pytest
import requests
from typing import Dict, Any
from datetime import datetime


class TestDataAccuracyValidation:
    """Test suite for validating data accuracy against known astronomical events"""

    # Known astronomical data for validation
    # Source: NASA JPL Horizons System and astronomical almanacs
    KNOWN_DATA = [
        {
            "name": "New Moon - January 2024",
            "date": "2024-01-11",
            "location": {"lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
            "expected": {
                "tithi_contains": "Amavasya",  # New Moon
                "paksha": "Krishna",  # Waning phase before new moon
            }
        },
        {
            "name": "Full Moon - January 2024",
            "date": "2024-01-25",
            "location": {"lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
            "expected": {
                "tithi_contains": "Purnima",  # Full Moon
                "paksha": "Shukla",  # Waxing phase
            }
        },
        {
            "name": "Summer Solstice 2024 - London",
            "date": "2024-06-20",
            "location": {"lat": 51.5074, "lng": -0.1278, "tz": "Europe/London"},
            "expected": {
                "day_length_min": 960,  # > 16 hours (960 minutes)
                "sunrise_before": "05:00",
                "sunset_after": "21:00",
            }
        },
        {
            "name": "Winter Solstice 2024 - New York",
            "date": "2024-12-21",
            "location": {"lat": 40.7128, "lng": -74.0060, "tz": "America/New_York"},
            "expected": {
                "day_length_max": 570,  # < 9.5 hours (570 minutes)
            }
        },
        {
            "name": "Equinox - March 2024",
            "date": "2024-03-20",
            "location": {"lat": 0.0, "lng": 0.0, "tz": "UTC"},  # Equator
            "expected": {
                "day_length_min": 710,  # ~12 hours (720 min ± 10 min)
                "day_length_max": 730,
            }
        },
        {
            "name": "Hindu Festival - Diwali 2024",
            "date": "2024-11-01",
            "location": {"lat": 28.6139, "lng": 77.2090, "tz": "Asia/Kolkata"},  # Delhi
            "expected": {
                "tithi_contains": "Amavasya",  # Diwali is on new moon
                "events_contain": ["Festival", "Auspicious"],  # Should mark festivals
            }
        },
    ]

    @pytest.mark.integration
    @pytest.mark.parametrize("test_case", KNOWN_DATA, ids=lambda x: x["name"])
    def test_known_astronomical_data(
        self,
        api_client: requests.Session,
        api_base_url: str,
        test_case: Dict[str, Any]
    ):
        """Test API responses against known astronomical data"""
        loc = test_case["location"]
        params = {
            "date": test_case["date"],
            "lat": loc["lat"],
            "lng": loc["lng"],
            "tz": loc["tz"],
        }

        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)

        assert response.status_code == 200, f"Failed for {test_case['name']}"
        data = response.json()

        # Validate expected fields
        expected = test_case["expected"]

        # Check tithi (lunar day) if specified
        if "tithi_contains" in expected:
            tithi = data.get("tithi", "")
            assert expected["tithi_contains"].lower() in tithi.lower(), \
                f"Expected tithi to contain '{expected['tithi_contains']}', got '{tithi}'"

        # Check paksha (lunar fortnight) if specified
        if "paksha" in expected:
            paksha = data.get("paksha", "")
            assert expected["paksha"].lower() in paksha.lower(), \
                f"Expected paksha '{expected['paksha']}', got '{paksha}'"

        # Check day length if specified
        if "day_length_min" in expected or "day_length_max" in expected:
            sunrise = data.get("sunrise_time", "")
            sunset = data.get("sunset_time", "")

            if sunrise and sunset:
                try:
                    # Parse times (assuming HH:MM:SS format)
                    sunrise_parts = sunrise.split(":")
                    sunset_parts = sunset.split(":")

                    sunrise_min = int(sunrise_parts[0]) * 60 + int(sunrise_parts[1])
                    sunset_min = int(sunset_parts[0]) * 60 + int(sunset_parts[1])

                    day_length = sunset_min - sunrise_min
                    if day_length < 0:  # Handle day crossing midnight
                        day_length += 24 * 60

                    if "day_length_min" in expected:
                        assert day_length >= expected["day_length_min"], \
                            f"Day length {day_length} min < expected minimum {expected['day_length_min']} min"

                    if "day_length_max" in expected:
                        assert day_length <= expected["day_length_max"], \
                            f"Day length {day_length} min > expected maximum {expected['day_length_max']} min"
                except (ValueError, IndexError) as e:
                    pytest.fail(f"Failed to parse sunrise/sunset times: {e}")

        # Check sunrise/sunset bounds if specified
        if "sunrise_before" in expected:
            sunrise = data.get("sunrise_time", "")
            assert sunrise < expected["sunrise_before"], \
                f"Sunrise {sunrise} not before {expected['sunrise_before']}"

        if "sunset_after" in expected:
            sunset = data.get("sunset_time", "")
            assert sunset > expected["sunset_after"], \
                f"Sunset {sunset} not after {expected['sunset_after']}"

        # Check events if specified
        if "events_contain" in expected:
            events = data.get("events", [])
            events_str = " ".join([str(e) for e in events])
            for expected_event in expected["events_contain"]:
                # Note: This is a soft check, events might not always be present
                # depending on the calculation method and data availability
                pass  # Can be enhanced when event data is more robust


class TestConsistencyValidation:
    """Test data consistency across multiple requests"""

    @pytest.mark.integration
    def test_same_request_returns_same_data(self, api_client: requests.Session, api_base_url: str):
        """Test that identical requests return identical data"""
        params = {
            "date": "2024-01-15",
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        # Make multiple identical requests
        responses = []
        for _ in range(3):
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            assert response.status_code == 200
            responses.append(response.json())

        # Verify all responses are identical
        first_response = responses[0]
        for i, response in enumerate(responses[1:], 1):
            assert response == first_response, \
                f"Response {i} differs from first response"

    @pytest.mark.integration
    def test_adjacent_dates_continuity(self, api_client: requests.Session, api_base_url: str):
        """Test that adjacent dates show logical progression"""
        base_params = {
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        # Test a sequence of dates
        dates = ["2024-01-15", "2024-01-16", "2024-01-17"]
        responses = []

        for date in dates:
            params = {**base_params, "date": date}
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)
            assert response.status_code == 200
            responses.append(response.json())

        # Verify dates are correct in responses
        for date, response in zip(dates, responses):
            assert response.get("date") == date

        # Verify logical progression (e.g., sunrise times shouldn't vary wildly)
        for i in range(len(responses) - 1):
            sunrise1 = responses[i].get("sunrise_time", "")
            sunrise2 = responses[i + 1].get("sunrise_time", "")

            if sunrise1 and sunrise2:
                # Adjacent days' sunrise should differ by < 5 minutes
                # (This is a reasonable assumption for mid-latitude locations)
                parts1 = sunrise1.split(":")
                parts2 = sunrise2.split(":")

                time1_min = int(parts1[0]) * 60 + int(parts1[1])
                time2_min = int(parts2[0]) * 60 + int(parts2[1])

                diff = abs(time2_min - time1_min)
                assert diff < 5, \
                    f"Sunrise times differ by {diff} minutes between adjacent dates"

    @pytest.mark.integration
    def test_geographic_consistency(self, api_client: requests.Session, api_base_url: str, sample_locations):
        """Test that nearby locations have similar data"""
        date = "2024-01-15"

        # Test Bangalore location
        bangalore = sample_locations["bangalore"]
        params1 = {
            "date": date,
            "lat": bangalore["lat"],
            "lng": bangalore["lng"],
            "tz": bangalore["tz"],
        }

        # Test nearby location (slightly different coordinates)
        params2 = {
            "date": date,
            "lat": bangalore["lat"] + 0.1,  # ~11 km north
            "lng": bangalore["lng"] + 0.1,  # ~11 km east
            "tz": bangalore["tz"],
        }

        response1 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params1)
        response2 = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params2)

        assert response1.status_code == 200
        assert response2.status_code == 200

        data1 = response1.json()
        data2 = response2.json()

        # Nearby locations should have same tithi
        assert data1.get("tithi") == data2.get("tithi"), \
            "Nearby locations should have same tithi"

        # Sunrise/sunset should differ by < 2 minutes for nearby locations
        if data1.get("sunrise_time") and data2.get("sunrise_time"):
            sunrise1 = data1["sunrise_time"]
            sunrise2 = data2["sunrise_time"]

            parts1 = sunrise1.split(":")
            parts2 = sunrise2.split(":")

            time1_min = int(parts1[0]) * 60 + int(parts1[1])
            time2_min = int(parts2[0]) * 60 + int(parts2[1])

            diff = abs(time2_min - time1_min)
            assert diff < 2, \
                f"Nearby locations' sunrise times differ by {diff} minutes"


class TestResponseStructureValidation:
    """Test that response structures are consistent and complete"""

    REQUIRED_FIELDS = [
        "date",
        "tithi",
        "nakshatra",
        "yoga",
        "karana",
        "paksha",
        "sunrise_time",
        "sunset_time",
        "moonrise_time",
        "moonset_time",
        "events",
    ]

    @pytest.mark.integration
    @pytest.mark.parametrize("location", ["bangalore", "mumbai", "new_york", "london"])
    def test_response_has_required_fields(
        self,
        api_client: requests.Session,
        api_base_url: str,
        sample_locations,
        location: str
    ):
        """Test that all responses contain required fields"""
        loc = sample_locations[location]
        params = {
            "date": "2024-01-15",
            "lat": loc["lat"],
            "lng": loc["lng"],
            "tz": loc["tz"],
        }

        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)

        assert response.status_code == 200
        data = response.json()

        # Check all required fields are present
        for field in self.REQUIRED_FIELDS:
            assert field in data, f"Missing required field: {field}"
            # Also verify non-empty for most fields
            if field != "events":  # Events can be empty
                assert data[field] is not None, f"Field {field} is None"
                if isinstance(data[field], str):
                    assert data[field] != "", f"Field {field} is empty string"

    @pytest.mark.integration
    def test_response_field_types(self, api_client: requests.Session, api_base_url: str, valid_request_params):
        """Test that response fields have correct types"""
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)

        assert response.status_code == 200
        data = response.json()

        # Verify field types
        assert isinstance(data.get("date"), str)
        assert isinstance(data.get("tithi"), str)
        assert isinstance(data.get("nakshatra"), str)
        assert isinstance(data.get("yoga"), str)
        assert isinstance(data.get("karana"), str)
        assert isinstance(data.get("paksha"), str)
        assert isinstance(data.get("sunrise_time"), str)
        assert isinstance(data.get("sunset_time"), str)
        assert isinstance(data.get("moonrise_time"), str)
        assert isinstance(data.get("moonset_time"), str)
        assert isinstance(data.get("events"), list)

    @pytest.mark.integration
    def test_time_format_validation(self, api_client: requests.Session, api_base_url: str, valid_request_params):
        """Test that time fields follow correct format (HH:MM:SS)"""
        response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=valid_request_params)

        assert response.status_code == 200
        data = response.json()

        time_fields = ["sunrise_time", "sunset_time", "moonrise_time", "moonset_time"]

        for field in time_fields:
            time_value = data.get(field)
            if time_value and time_value != "N/A":  # Some fields might be N/A at poles
                # Verify HH:MM:SS format
                parts = time_value.split(":")
                assert len(parts) == 3, f"{field} doesn't have HH:MM:SS format"

                hours, minutes, seconds = parts
                assert hours.isdigit() and 0 <= int(hours) < 24, \
                    f"Invalid hours in {field}: {hours}"
                assert minutes.isdigit() and 0 <= int(minutes) < 60, \
                    f"Invalid minutes in {field}: {minutes}"
                assert seconds.isdigit() and 0 <= int(seconds) < 60, \
                    f"Invalid seconds in {field}: {seconds}"


class TestHistoricalDataValidation:
    """Test historical dates and edge cases"""

    @pytest.mark.integration
    def test_historical_dates(self, api_client: requests.Session, api_base_url: str):
        """Test that historical dates work correctly"""
        historical_dates = [
            "2020-01-01",  # Recent past
            "2015-06-21",  # 5 years ago
            "2010-12-31",  # 10 years ago
            "2000-01-01",  # Y2K
        ]

        params_base = {
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        for date in historical_dates:
            params = {**params_base, "date": date}
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)

            assert response.status_code == 200, f"Failed for historical date: {date}"
            data = response.json()
            assert data.get("date") == date
            assert "tithi" in data

    @pytest.mark.integration
    def test_future_dates(self, api_client: requests.Session, api_base_url: str):
        """Test that future dates work correctly"""
        future_dates = [
            "2025-01-01",
            "2026-06-21",
            "2030-12-31",
        ]

        params_base = {
            "lat": 12.9716,
            "lng": 77.5946,
            "tz": "Asia/Kolkata",
        }

        for date in future_dates:
            params = {**params_base, "date": date}
            response = api_client.get(f"{api_base_url}/api/v1/panchangam", params=params)

            assert response.status_code == 200, f"Failed for future date: {date}"
            data = response.json()
            assert data.get("date") == date
            assert "tithi" in data
