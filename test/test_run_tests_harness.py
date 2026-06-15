#!/usr/bin/env python3
"""Unit tests for the Python test runner wiring."""

import sys
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))

import run_tests


class ServerBuildPathTest(unittest.TestCase):
    def test_server_build_commands_use_test_bin_directory(self):
        project_root = Path("/repo")

        commands = run_tests.server_build_commands(project_root)

        self.assertEqual(
            commands,
            [
                ("gRPC server", ["go", "build", "-o", "/repo/tmp/test-bin/grpc-server", "./cmd/server"]),
                ("Gateway server", ["go", "build", "-o", "/repo/tmp/test-bin/gateway-server", "./cmd/gateway"]),
            ],
        )


class DependencyCommandTest(unittest.TestCase):
    def test_check_dependencies_uses_argv_commands(self):
        original_run_command = run_tests.run_command
        calls = []

        def fake_run_command(cmd, capture_output=True, env=None):
            calls.append(cmd)
            return 0, "", ""

        run_tests.run_command = fake_run_command
        try:
            self.assertTrue(run_tests.check_dependencies())
        finally:
            run_tests.run_command = original_run_command

        self.assertEqual(
            calls,
            [
                [sys.executable, "-m", "pytest", "--version"],
                ["go", "version"],
            ],
        )

    def test_install_dependencies_uses_argv_command(self):
        original_run_command = run_tests.run_command
        calls = []

        def fake_run_command(cmd, capture_output=True, env=None):
            calls.append(cmd)
            return 0, "", ""

        run_tests.run_command = fake_run_command
        try:
            self.assertTrue(run_tests.install_dependencies())
        finally:
            run_tests.run_command = original_run_command

        self.assertEqual(
            calls,
            [[sys.executable, "-m", "pip", "install", "-r", str(Path(run_tests.__file__).parent / "requirements.txt")]],
        )


class RunTestsEnvironmentTest(unittest.TestCase):
    def test_run_tests_passes_project_root_pythonpath(self):
        original_run_command = run_tests.run_command
        calls = []

        def fake_run_command(cmd, capture_output=True, env=None):
            calls.append((cmd, capture_output, env))
            return 0, "", ""

        run_tests.run_command = fake_run_command
        try:
            self.assertTrue(run_tests.run_tests())
        finally:
            run_tests.run_command = original_run_command

        self.assertEqual(len(calls), 1)
        _, _, env = calls[0]
        self.assertIsNotNone(env)
        self.assertEqual(env["PYTHONPATH"], str(Path(run_tests.__file__).parent.parent))

    def test_run_tests_passes_pytest_args_as_list(self):
        original_run_command = run_tests.run_command
        calls = []

        def fake_run_command(cmd, capture_output=True, env=None):
            calls.append(cmd)
            return 0, "", ""

        run_tests.run_command = fake_run_command
        try:
            self.assertTrue(run_tests.run_tests(markers="smoke and not performance"))
        finally:
            run_tests.run_command = original_run_command

        self.assertEqual(
            calls[0],
            [sys.executable, "-m", "pytest", "-m", "smoke and not performance", "-q", "."],
        )


if __name__ == "__main__":
    unittest.main()
