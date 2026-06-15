#!/usr/bin/env python3
"""
Panchangam API Test Runner
Comprehensive test execution script for the Panchangam API gateway
"""
import os
import sys
import subprocess
import argparse
import time
from pathlib import Path


def run_command(cmd, capture_output=True, env=None):
    """Run a command and return the result"""
    try:
        result = subprocess.run(
            cmd,
            capture_output=capture_output,
            env=env,
            text=True,
            check=False
        )
        return result.returncode, result.stdout, result.stderr
    except Exception as e:
        return 1, "", str(e)


def check_dependencies():
    """Check if required dependencies are installed"""
    print("Checking dependencies...")

    # Check Python version
    if sys.version_info < (3, 8):
        print("FAIL Python 3.8+ required")
        return False

    # Check if pytest is installed
    returncode, _, _ = run_command([sys.executable, "-m", "pytest", "--version"])
    if returncode != 0:
        print("FAIL pytest not installed. Run: pip install -r requirements.txt")
        return False

    # Check if Go is available (for building servers)
    returncode, _, _ = run_command(["go", "version"])
    if returncode != 0:
        print("FAIL Go not installed or not in PATH")
        return False

    print("PASS All dependencies available")
    return True


def test_binary_dir(project_root):
    """Directory for binaries built by the test harness"""
    return project_root / "tmp" / "test-bin"


def server_build_commands(project_root):
    """Build commands for the local test servers"""
    bin_dir = test_binary_dir(project_root)
    return [
        ("gRPC server", ["go", "build", "-o", str(bin_dir / "grpc-server"), "./cmd/server"]),
        ("Gateway server", ["go", "build", "-o", str(bin_dir / "gateway-server"), "./cmd/gateway"]),
    ]


def build_servers():
    """Build the gRPC and Gateway servers"""
    print("Building servers...")

    # Change to project root
    project_root = Path(__file__).parent.parent
    os.chdir(project_root)
    test_binary_dir(project_root).mkdir(parents=True, exist_ok=True)

    for label, command in server_build_commands(project_root):
        returncode, stdout, stderr = run_command(command)
        if returncode != 0:
            print(f"FAIL Failed to build {label}: {stderr}")
            return False

    print("PASS Servers built successfully")
    return True


def install_dependencies():
    """Install Python test dependencies"""
    print("Installing Python dependencies...")

    test_dir = Path(__file__).parent
    requirements_file = test_dir / "requirements.txt"

    if not requirements_file.exists():
        print("FAIL requirements.txt not found")
        return False

    returncode, stdout, stderr = run_command([sys.executable, "-m", "pip", "install", "-r", str(requirements_file)])
    if returncode != 0:
        print(f"FAIL Failed to install dependencies: {stderr}")
        return False

    print("PASS Dependencies installed successfully")
    return True


def run_tests(test_type="all", verbose=False, coverage=False, html_report=False, markers=None):
    """Run pytest tests with specified options"""
    test_dir = Path(__file__).parent
    os.chdir(test_dir)

    # Build pytest command
    cmd_parts = [sys.executable, "-m", "pytest"]

    # Add test selection based on type
    if test_type == "smoke":
        cmd_parts.extend(["-m", "smoke"])
    elif test_type == "integration":
        cmd_parts.extend(["-m", "integration"])
    elif test_type == "performance":
        cmd_parts.extend(["-m", "performance"])
    elif test_type == "security":
        cmd_parts.extend(["-m", "security"])
    elif markers:
        cmd_parts.extend(["-m", markers])

    # Add verbosity
    if verbose:
        cmd_parts.append("-v")
    else:
        cmd_parts.append("-q")

    # Add coverage
    if coverage:
        cmd_parts.extend(["--cov=../gateway", "--cov-report=term-missing"])
        if html_report:
            cmd_parts.append("--cov-report=html:htmlcov")

    # Add HTML report
    if html_report:
        cmd_parts.extend(["--html=report.html", "--self-contained-html"])

    # Add current directory
    cmd_parts.append(".")

    print(f"Running tests: {' '.join(cmd_parts)}")

    # Set environment variables
    env = os.environ.copy()
    env["PYTHONPATH"] = str(test_dir.parent)

    # Run tests
    start_time = time.time()
    returncode, stdout, stderr = run_command(cmd_parts, capture_output=False, env=env)
    end_time = time.time()

    duration = end_time - start_time

    if returncode == 0:
        print(f"PASS Tests completed successfully in {duration:.2f}s")
        if html_report:
            print(f"HTML report generated: {test_dir}/report.html")
        if coverage and html_report:
            print(f"Coverage report: {test_dir}/htmlcov/index.html")
    else:
        print(f"FAIL Tests failed after {duration:.2f}s")
        print(f"Exit code: {returncode}")

    return returncode == 0


def main():
    """Main test runner function"""
    parser = argparse.ArgumentParser(description="Panchangam API Test Runner")

    parser.add_argument(
        "--type",
        choices=["all", "smoke", "integration", "performance", "security"],
        default="all",
        help="Type of tests to run"
    )

    parser.add_argument(
        "--markers",
        help="Custom pytest markers to run (e.g., 'smoke and not performance')"
    )

    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Verbose output"
    )

    parser.add_argument(
        "--coverage", "-c",
        action="store_true",
        help="Generate coverage report"
    )

    parser.add_argument(
        "--html-report",
        action="store_true",
        help="Generate HTML test and coverage reports"
    )

    parser.add_argument(
        "--install-deps",
        action="store_true",
        help="Install Python dependencies before running tests"
    )

    parser.add_argument(
        "--build-servers",
        action="store_true",
        help="Build Go servers before running tests"
    )

    parser.add_argument(
        "--skip-server-start",
        action="store_true",
        help="Skip automatic server startup (assumes servers are already running)"
    )

    args = parser.parse_args()

    print("Panchangam API Test Runner")
    print("=" * 50)

    # Check dependencies
    if not check_dependencies():
        sys.exit(1)

    # Install dependencies if requested
    if args.install_deps:
        if not install_dependencies():
            sys.exit(1)

    # Build servers if requested
    if args.build_servers:
        if not build_servers():
            sys.exit(1)

    # Set environment variable if skipping server start
    if args.skip_server_start:
        os.environ["SKIP_SERVER_START"] = "true"

    # Run tests
    success = run_tests(
        test_type=args.type,
        verbose=args.verbose,
        coverage=args.coverage,
        html_report=args.html_report,
        markers=args.markers
    )

    if not success:
        sys.exit(1)

    print("\nTest execution completed!")


if __name__ == "__main__":
    main()
