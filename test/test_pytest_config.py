#!/usr/bin/env python3
"""Pytest fixture wiring tests."""

import unittest
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]


def read_text(relative_path: str) -> str:
    return (PROJECT_ROOT / relative_path).read_text(encoding="utf-8")


class PytestConfigTest(unittest.TestCase):
    def test_server_health_check_uses_configured_api_base_url(self):
        conftest = read_text("test/conftest.py")

        self.assertIn("def start_servers(api_base_url: str)", conftest)
        self.assertIn("from urllib.parse import urlparse", conftest)
        self.assertIn('grpc_port = os.getenv("PANCHANGAM_GRPC_PORT", "50051")', conftest)
        self.assertIn("api_url = urlparse(api_base_url)", conftest)
        self.assertIn('f"--grpc-port={grpc_port}"', conftest)
        self.assertIn('f"--grpc-endpoint=localhost:{grpc_port}"', conftest)
        self.assertIn('f"--http-port={http_port}"', conftest)
        self.assertIn('requests.get(f"{api_base_url}/api/v1/health"', conftest)
        self.assertNotIn('"--grpc-endpoint=localhost:50052"', conftest)
        self.assertNotIn('requests.get("http://localhost:8080/api/v1/health"', conftest)


if __name__ == "__main__":
    unittest.main()
