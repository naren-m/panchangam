#!/usr/bin/env python3
"""Shell script wiring tests."""

import re
import subprocess
import unittest
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]


def read_text(relative_path: str) -> str:
    return (PROJECT_ROOT / relative_path).read_text(encoding="utf-8")


class IntegrationScriptConfigTest(unittest.TestCase):
    def test_integration_runner_uses_configurable_http_base_url(self):
        script = read_text("test/integration_test.sh")

        self.assertIn('HTTP_PORT=${HTTP_PORT:-8080}', script)
        self.assertIn('HTTP_BASE_URL="http://localhost:${HTTP_PORT}"', script)
        self.assertIn("${HTTP_BASE_URL}/api/v1/health", script)
        self.assertNotIn("http://localhost:8080", script)

    def test_integration_runner_does_not_use_eval(self):
        script = read_text("test/integration_test.sh")

        self.assertNotIn("eval ", script)
        self.assertNotIn("local command=", script)
        self.assertIn('response=$("$@" 2>&1) || true', script)
        self.assertIn('"$check_function" "$response"', script)

    def test_integration_runner_uses_set_e_safe_counters(self):
        script = read_text("test/integration_test.sh")

        self.assertNotIn("TESTS_PASSED++", script)
        self.assertNotIn("TESTS_FAILED++", script)
        self.assertIn("TESTS_PASSED=$((TESTS_PASSED + 1))", script)
        self.assertIn("TESTS_FAILED=$((TESTS_FAILED + 1))", script)

    def test_integration_runner_uses_stop_script_for_cleanup(self):
        script = read_text("test/integration_test.sh")

        self.assertIn("./scripts/stop-servers.sh", script)
        self.assertNotIn("pkill -f grpc-server", script)
        self.assertNotIn("pkill -f gateway-server", script)


class ValidateDeploymentScriptConfigTest(unittest.TestCase):
    def test_deployment_scripts_use_set_e_safe_counters(self):
        scripts = [
            read_text("scripts/validate-deployment.sh"),
            read_text("scripts/deploy.sh"),
        ]

        for script in scripts:
            self.assertNotIn("++))", script)

        self.assertIn("passed=$((passed + 1))", scripts[0])
        self.assertIn("successful_requests=$((successful_requests + 1))", scripts[0])
        self.assertIn("attempt=$((attempt + 1))", scripts[1])

    def test_deployment_scripts_quote_numeric_test_variables(self):
        scripts = {
            "scripts/validate-deployment.sh": read_text("scripts/validate-deployment.sh"),
            "scripts/deploy.sh": read_text("scripts/deploy.sh"),
            "test/run-docker-tests.sh": read_text("test/run-docker-tests.sh"),
        }

        for path, script in scripts.items():
            unquoted_tests = re.findall(r"\[\s+\$[A-Za-z_][A-Za-z0-9_]*", script)
            self.assertEqual([], unquoted_tests, path)

        self.assertIn(
            'generate_report "$health_status" "$api_status" "$performance_status" "$error_status" "$frontend_status"',
            scripts["scripts/validate-deployment.sh"],
        )
        self.assertIn('if [ "${#missing_tools[@]}" -ne 0 ]; then', scripts["scripts/deploy.sh"])
        self.assertNotIn('if [ ${#missing_tools[@]} -ne 0 ]; then', scripts["scripts/deploy.sh"])

    def test_deploy_script_reports_missing_option_values(self):
        script = read_text("scripts/deploy.sh")

        expected_errors = [
            "Error: --environment requires a value",
            "Error: --version requires a value",
            "Error: --registry requires a value",
            "Error: --namespace requires a value",
        ]
        for expected_error in expected_errors:
            self.assertIn(expected_error, script)

        self.assertIn('if [ "$#" -lt 2 ] || [[ "$2" == -* ]]; then', script)

    def test_deploy_script_rejects_unsafe_version_tag_before_prerequisites(self):
        result = subprocess.run(
            [
                "bash",
                str(PROJECT_ROOT / "scripts/deploy.sh"),
                "--dry-run",
                "--version",
                "../bad",
            ],
            cwd=PROJECT_ROOT,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            check=False,
        )

        self.assertEqual(1, result.returncode)
        self.assertIn(
            "Error: Version may only contain letters, numbers, dots, dashes, and underscores",
            result.stdout,
        )
        self.assertNotIn("Checking prerequisites", result.stdout)

    def test_deploy_script_rejects_invalid_namespace_before_prerequisites(self):
        result = subprocess.run(
            [
                "bash",
                str(PROJECT_ROOT / "scripts/deploy.sh"),
                "--dry-run",
                "--version",
                "v1.2.3",
                "--namespace",
                "../bad",
            ],
            cwd=PROJECT_ROOT,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            check=False,
        )

        self.assertEqual(1, result.returncode)
        self.assertIn(
            "Error: Namespace must be a lowercase Kubernetes DNS label",
            result.stdout,
        )
        self.assertNotIn("Checking prerequisites", result.stdout)

    def test_deploy_script_rejects_empty_registry_before_prerequisites(self):
        result = subprocess.run(
            [
                "bash",
                str(PROJECT_ROOT / "scripts/deploy.sh"),
                "--dry-run",
                "--version",
                "v1.2.3",
                "--registry",
                "",
            ],
            cwd=PROJECT_ROOT,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            check=False,
        )

        self.assertEqual(1, result.returncode)
        self.assertIn("Error: Registry must not be empty", result.stdout)
        self.assertNotIn("Checking prerequisites", result.stdout)

    def test_validate_deployment_script_reports_missing_option_values(self):
        script = read_text("scripts/validate-deployment.sh")

        expected_errors = [
            "Error: --environment requires a value",
            "Error: --url requires a value",
            "Error: --timeout requires a value",
        ]
        for expected_error in expected_errors:
            self.assertIn(expected_error, script)

        self.assertIn('if [ "$#" -lt 2 ] || [[ "$2" == -* ]]; then', script)

    def test_validate_deployment_script_uses_portable_response_body_parsing(self):
        script = read_text("scripts/validate-deployment.sh")

        self.assertNotIn("head -n -1", script)
        self.assertNotIn('echo "$response"', script)
        self.assertNotIn('echo -e "\\n000"', script)
        self.assertIn("sed '$d'", script)
        self.assertIn("printf '%s\\n' \"$response\" | tail -n 1", script)
        self.assertIn("printf '\\n000\\n'", script)

    def test_validate_deployment_script_uses_curl_timing_for_performance_checks(self):
        script = read_text("scripts/validate-deployment.sh")

        self.assertNotIn("date +%s.%3N", script)
        self.assertNotIn("end_time - $start_time", script)
        self.assertIn("%{time_total}", script)

    def test_validate_deployment_script_rejects_unknown_environment(self):
        result = subprocess.run(
            [
                "bash",
                str(PROJECT_ROOT / "scripts/validate-deployment.sh"),
                "-e",
                "qa",
                "-u",
                "http://127.0.0.1:1",
                "-t",
                "0",
            ],
            cwd=PROJECT_ROOT,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            check=False,
        )

        self.assertEqual(1, result.returncode)
        self.assertIn("Error: Environment must be 'staging' or 'production'", result.stdout)
        self.assertNotIn("Validating health check endpoint", result.stdout)

    def test_validate_deployment_script_rejects_non_numeric_timeout(self):
        result = subprocess.run(
            [
                "bash",
                str(PROJECT_ROOT / "scripts/validate-deployment.sh"),
                "-t",
                "soon",
                "-u",
                "http://127.0.0.1:1",
            ],
            cwd=PROJECT_ROOT,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            check=False,
        )

        self.assertEqual(1, result.returncode)
        self.assertIn("Error: --timeout must be a whole number of seconds", result.stdout)
        self.assertNotIn("Validating health check endpoint", result.stdout)


if __name__ == "__main__":
    unittest.main()
