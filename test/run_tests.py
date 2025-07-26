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


def run_command(cmd, capture_output=True):
    """Run a shell command and return the result"""
    try:
        result = subprocess.run(
            cmd, 
            shell=True, 
            capture_output=capture_output, 
            text=True,
            check=False
        )
        return result.returncode, result.stdout, result.stderr
    except Exception as e:
        return 1, "", str(e)


def check_dependencies():
    """Check if required dependencies are installed"""
    print("🔍 Checking dependencies...")
    
    # Check Python version
    if sys.version_info < (3, 8):
        print("❌ Python 3.8+ required")
        return False
    
    # Check if pytest is installed
    returncode, _, _ = run_command("python -m pytest --version")
    if returncode != 0:
        print("❌ pytest not installed. Run: pip install -r requirements.txt")
        return False
    
    # Check if Go is available (for building servers)
    returncode, _, _ = run_command("go version")
    if returncode != 0:
        print("❌ Go not installed or not in PATH")
        return False
    
    print("✅ All dependencies available")
    return True


def build_servers():
    """Build the gRPC and Gateway servers"""
    print("🏗️ Building servers...")
    
    # Change to project root
    project_root = Path(__file__).parent.parent
    os.chdir(project_root)
    
    # Build gRPC server
    returncode, stdout, stderr = run_command("go build -o grpc-server ./server/server.go")
    if returncode != 0:
        print(f"❌ Failed to build gRPC server: {stderr}")
        return False
    
    # Build Gateway server
    returncode, stdout, stderr = run_command("go build -o gateway-server ./cmd/gateway/main.go")
    if returncode != 0:
        print(f"❌ Failed to build Gateway server: {stderr}")
        return False
    
    print("✅ Servers built successfully")
    return True


def install_dependencies():
    """Install Python test dependencies"""
    print("📦 Installing Python dependencies...")
    
    test_dir = Path(__file__).parent
    requirements_file = test_dir / "requirements.txt"
    
    if not requirements_file.exists():
        print("❌ requirements.txt not found")
        return False
    
    returncode, stdout, stderr = run_command(f"pip install -r {requirements_file}")
    if returncode != 0:
        print(f"❌ Failed to install dependencies: {stderr}")
        return False
    
    print("✅ Dependencies installed successfully")
    return True


def run_tests(test_type="all", verbose=False, coverage=False, html_report=False, markers=None):
    """Run pytest tests with specified options"""
    test_dir = Path(__file__).parent
    os.chdir(test_dir)
    
    # Build pytest command
    cmd_parts = ["python", "-m", "pytest"]
    
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
    
    cmd = " ".join(cmd_parts)
    print(f"🧪 Running tests: {cmd}")
    
    # Set environment variables
    env = os.environ.copy()
    env["PYTHONPATH"] = str(test_dir.parent)
    
    # Run tests
    start_time = time.time()
    returncode, stdout, stderr = run_command(cmd, capture_output=False)
    end_time = time.time()
    
    duration = end_time - start_time
    
    if returncode == 0:
        print(f"✅ Tests completed successfully in {duration:.2f}s")
        if html_report:
            print(f"📊 HTML report generated: {test_dir}/report.html")
        if coverage and html_report:
            print(f"📈 Coverage report: {test_dir}/htmlcov/index.html")
    else:
        print(f"❌ Tests failed after {duration:.2f}s")
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
    
    print("🚀 Panchangam API Test Runner")
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
    
    print("\n🎉 Test execution completed!")


if __name__ == "__main__":
    main()