#!/usr/bin/env python3
"""Docker test runner wiring tests."""

import unittest
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]


def read_text(relative_path: str) -> str:
    return (PROJECT_ROOT / relative_path).read_text(encoding="utf-8")


class DockerTestRunnerConfigTest(unittest.TestCase):
    def test_runner_builds_image_with_argv_array(self):
        script = read_text("test/run-docker-tests.sh")

        self.assertIn("local build_args=(", script)
        self.assertIn('docker-compose build "${build_args[@]}"', script)
        self.assertNotIn("docker-compose build $BUILD_FLAG", script)

    def test_runner_executes_run_tests_as_argv_array(self):
        script = read_text("test/run-docker-tests.sh")

        self.assertIn("local test_command=(", script)
        self.assertIn('docker-compose run --rm panchangam-tests "${test_command[@]}"', script)
        self.assertIn('docker-compose exec panchangam-tests "${test_command[@]}"', script)
        self.assertNotIn('local test_command="', script)


if __name__ == "__main__":
    unittest.main()
