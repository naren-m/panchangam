#!/bin/bash
# Docker-based test execution script for Panchangam API

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Panchangam API Docker Test Runner"
    echo ""
    echo "Usage: $0 [OPTIONS] [TEST_TYPE]"
    echo ""
    echo "Test Types:"
    echo "  all           Run all tests (default)"
    echo "  smoke         Run smoke tests only"
    echo "  integration   Run integration tests only"
    echo "  performance   Run performance tests only"
    echo "  security      Run security tests only"
    echo ""
    echo "Options:"
    echo "  -h, --help        Show this help message"
    echo "  -v, --verbose     Verbose output"
    echo "  -c, --coverage    Generate coverage report"
    echo "  --html-report     Generate HTML reports"
    echo "  --build           Force rebuild Docker image"
    echo "  --cleanup         Clean up containers after execution"
    echo "  --interactive     Run container interactively"
    echo "  --shell           Open shell in test container"
    echo "  --logs            Show container logs"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Run all tests"
    echo "  $0 smoke -v                          # Run smoke tests with verbose output"
    echo "  $0 integration --coverage --html-report  # Integration tests with reports"
    echo "  $0 --shell                           # Open interactive shell"
    echo "  $0 --build --cleanup                 # Force rebuild and cleanup"
}

# Default values
TEST_TYPE="all"
VERBOSE=""
COVERAGE=""
HTML_REPORT=""
BUILD_FLAG=""
CLEANUP=false
INTERACTIVE=false
SHELL_MODE=false
SHOW_LOGS=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE="--verbose"
            shift
            ;;
        -c|--coverage)
            COVERAGE="--coverage"
            shift
            ;;
        --html-report)
            HTML_REPORT="--html-report"
            shift
            ;;
        --build)
            BUILD_FLAG="--build"
            shift
            ;;
        --cleanup)
            CLEANUP=true
            shift
            ;;
        --interactive)
            INTERACTIVE=true
            shift
            ;;
        --shell)
            SHELL_MODE=true
            shift
            ;;
        --logs)
            SHOW_LOGS=true
            shift
            ;;
        all|smoke|integration|performance|security)
            TEST_TYPE="$1"
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Function to check if Docker is available
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed or not in PATH"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed or not in PATH"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        print_error "Docker daemon is not running"
        exit 1
    fi
}

# Function to build Docker image
build_image() {
    print_status "Building Docker image..."
    if docker-compose build $BUILD_FLAG panchangam-tests; then
        print_success "Docker image built successfully"
    else
        print_error "Failed to build Docker image"
        exit 1
    fi
}

# Function to cleanup containers
cleanup_containers() {
    print_status "Cleaning up containers..."
    docker-compose down --remove-orphans
    docker-compose rm -f
    print_success "Cleanup completed"
}

# Function to show logs
show_container_logs() {
    print_status "Container logs:"
    docker-compose logs panchangam-tests
}

# Function to run shell in container
run_shell() {
    print_status "Opening interactive shell in test container..."
    docker-compose up -d panchangam-tests
    docker-compose exec panchangam-tests /bin/bash
}

# Function to run tests
run_tests() {
    local test_command="python run_tests.py --type $TEST_TYPE $VERBOSE $COVERAGE $HTML_REPORT"
    
    print_status "Starting test execution..."
    print_status "Test type: $TEST_TYPE"
    print_status "Command: $test_command"
    
    if [ "$INTERACTIVE" = true ]; then
        print_status "Running in interactive mode..."
        docker-compose run --rm panchangam-tests $test_command
    else
        print_status "Running in detached mode..."
        docker-compose up -d panchangam-tests
        docker-compose exec panchangam-tests $test_command
        local exit_code=$?
        
        if [ $exit_code -eq 0 ]; then
            print_success "Tests completed successfully"
        else
            print_error "Tests failed with exit code $exit_code"
            if [ "$SHOW_LOGS" = true ]; then
                show_container_logs
            fi
        fi
        
        return $exit_code
    fi
}

# Function to copy test reports
copy_reports() {
    print_status "Copying test reports..."
    
    # Create local directories if they don't exist
    mkdir -p ./reports ./htmlcov ./logs
    
    # Copy files from container
    docker-compose exec panchangam-tests bash -c "
        if [ -f /app/reports/report.html ]; then
            cp /app/reports/report.html /workspace/test/reports/
        fi
        if [ -d /app/htmlcov ]; then
            cp -r /app/htmlcov/* /workspace/test/htmlcov/ 2>/dev/null || true
        fi
        if [ -d /app/logs ]; then
            cp -r /app/logs/* /workspace/test/logs/ 2>/dev/null || true
        fi
    " 2>/dev/null || true
    
    # Check if reports were generated
    if [ -f "./reports/report.html" ]; then
        print_success "Test report available at: ./reports/report.html"
    fi
    
    if [ -f "./htmlcov/index.html" ]; then
        print_success "Coverage report available at: ./htmlcov/index.html"
    fi
}

# Main execution
main() {
    print_status "Panchangam API Docker Test Runner"
    print_status "================================="
    
    # Check prerequisites
    check_docker
    
    # Handle special modes
    if [ "$SHOW_LOGS" = true ]; then
        show_container_logs
        exit 0
    fi
    
    if [ "$SHELL_MODE" = true ]; then
        run_shell
        exit 0
    fi
    
    # Build image
    build_image
    
    # Ensure cleanup on exit
    if [ "$CLEANUP" = true ]; then
        trap cleanup_containers EXIT
    fi
    
    # Run tests
    if run_tests; then
        # Copy reports if they exist
        copy_reports
        print_success "Test execution completed successfully"
        exit 0
    else
        print_error "Test execution failed"
        exit 1
    fi
}

# Execute main function
main