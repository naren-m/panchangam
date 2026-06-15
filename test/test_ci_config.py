#!/usr/bin/env python3
"""CI and startup-instruction wiring tests."""

import unittest
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]


def read_text(relative_path: str) -> str:
    return (PROJECT_ROOT / relative_path).read_text(encoding="utf-8")


def text_between(content: str, start: str, end: str) -> str:
    return content.split(start, 1)[1].split(end, 1)[0]


class CiConfigTest(unittest.TestCase):
    def test_e2e_workflow_starts_backend_in_steps_not_service_container(self):
        workflow = read_text(".github/workflows/ci-cd.yml")
        e2e_job = text_between(workflow, "  e2e-test:", "  # ===========================\n  # Build & Package Jobs")

        self.assertNotIn("    services:\n      backend:", e2e_job)
        self.assertNotIn("image: golang:1.21", e2e_job)
        self.assertIn("      - name: Start backend services", e2e_job)
        self.assertIn("./bin/panchangam-grpc &", e2e_job)

    def test_e2e_workflow_builds_canonical_backend_commands(self):
        workflow = read_text(".github/workflows/ci-cd.yml")
        build_backend = text_between(workflow, "      - name: Build backend", "      - name: Build frontend")

        self.assertIn("mkdir -p bin", build_backend)
        self.assertIn("go build -o bin/panchangam-gateway ./cmd/gateway", build_backend)
        self.assertIn("go build -o bin/panchangam-grpc ./cmd/server", build_backend)
        self.assertNotIn("./cmd/grpc-server", build_backend)

    def test_user_facing_startup_help_uses_canonical_server_command(self):
        network_error = read_text("ui/src/components/common/Error/NetworkError.tsx")

        self.assertIn("go run ./cmd/server", network_error)
        self.assertIn("go run ./cmd/gateway", network_error)
        self.assertNotIn("cmd/grpc-server", network_error)

    def test_makefile_builds_one_canonical_grpc_binary(self):
        makefile = read_text("Makefile")
        build_backend = text_between(makefile, "# Build backend binaries", "# Build frontend")

        self.assertIn("go build $(LDFLAGS) -o bin/panchangam-grpc ./cmd/server", build_backend)
        self.assertEqual(1, build_backend.count("./cmd/server"))
        self.assertNotIn("panchangam-server", build_backend)

    def test_start_script_uses_canonical_grpc_binary(self):
        script = read_text("scripts/start-servers.sh")

        self.assertIn("go build -o ./bin/panchangam-grpc ./cmd/server", script)
        self.assertIn("./bin/panchangam-grpc --grpc-port=${GRPC_PORT}", script)
        self.assertNotIn("panchangam-server", script)

    def test_deployment_doc_uses_canonical_grpc_binary(self):
        deployment_doc = read_text("DEPLOYMENT.md")

        self.assertIn("go build -o bin/panchangam-grpc ./cmd/server", deployment_doc)
        self.assertIn("./bin/panchangam-grpc", deployment_doc)
        self.assertIn("ExecStart=/opt/panchangam/bin/panchangam-grpc --grpc-port=50051", deployment_doc)
        self.assertIn("docker/Dockerfile.backend", deployment_doc)
        self.assertNotIn("panchangam-server", deployment_doc)

    def test_runnable_go_commands_use_package_paths(self):
        checks = {
            ".github/workflows/ci-cd.yml": ["go run ./cmd/test-service"],
            "ui/src/components/Settings/ApiHealthCheck.tsx": ["go run ./cmd/gateway"],
        }

        for relative_path, expected_commands in checks.items():
            content = read_text(relative_path)
            for expected_command in expected_commands:
                self.assertIn(expected_command, content)

            self.assertNotRegex(content, r"go run cmd/[A-Za-z0-9_-]+/main\.go")
            self.assertNotRegex(content, r"go build -o [^\n]+ cmd/[A-Za-z0-9_-]+/main\.go")

    def test_user_docs_use_runnable_go_package_paths(self):
        checks = {
            "cmd/panchangam-cli/README.md": [
                "go build -o panchangam-cli ./cmd/panchangam-cli",
                "go run ./cmd/panchangam-cli [command]",
            ],
            "cmd/sunrise-demo/README.md": [
                "go run ./cmd/sunrise-demo -location london",
                "go build -o sunrise-demo ./cmd/sunrise-demo",
            ],
            "ui/README.md": [
                "Go 1.23+",
                "go run ./cmd/gateway",
            ],
        }

        for relative_path, expected_commands in checks.items():
            content = read_text(relative_path)
            for expected_command in expected_commands:
                self.assertIn(expected_command, content)

            self.assertNotRegex(content, r"go run cmd/[A-Za-z0-9_-]+/main\.go")
            self.assertNotRegex(content, r"go build -o [^\n]+ cmd/[A-Za-z0-9_-]+/main\.go")
            self.assertNotIn("Go 1.21+", content)

        architecture_doc = read_text("llm/project-architecture.md")
        self.assertIn("gRPC Server (`cmd/server`)", architecture_doc)
        self.assertIn("│   ├── server/", architecture_doc)
        self.assertIn("# Main gRPC server", architecture_doc)
        self.assertNotIn("├── main.go", architecture_doc)
        self.assertNotIn("\n├── client/", architecture_doc)
        self.assertNotIn("\n├── server/", architecture_doc)
        self.assertNotIn("grpc-server/", architecture_doc)
        self.assertNotIn("cmd/grpc-server", architecture_doc)
        for current_gateway_file in [
            "gateway/server.go",
            "gateway/panchangam_handler.go",
            "gateway/current_tithi_handler.go",
            "gateway/sky_view_handler.go",
            "gateway/middleware.go",
        ]:
            self.assertIn(current_gateway_file, architecture_doc)
        self.assertIn("cmd/server/main.go", architecture_doc)
        self.assertIn("services/panchangam/server.go", architecture_doc)
        self.assertIn("services/panchangam/service.go", architecture_doc)
        self.assertNotIn("`handler.go`", architecture_doc)
        self.assertNotIn("`interceptors.go`", architecture_doc)
        self.assertNotIn("└── handlers/", architecture_doc)
        self.assertNotIn("gRPC-Gateway", architecture_doc)
        self.assertNotIn("Future routes", architecture_doc)
        self.assertIn("GET  /api/v1/health", architecture_doc)
        self.assertIn("GET  /api/v1/tithi/current", architecture_doc)
        self.assertIn("GET  /api/v1/sky-view", architecture_doc)
        self.assertNotIn("GET  /health", architecture_doc)
        self.assertNotIn("GET  /api/v1/muhurta", architecture_doc)
        self.assertNotIn("GET  /api/v1/festivals", architecture_doc)
        self.assertNotIn("GET  /api/v1/skyview", architecture_doc)
        for current_frontend_area in [
            "Calendar/",
            "DayDetail/",
            "EclipticBeltVisualization/",
            "SkyVisualization/",
            "LocationPicker/",
            "Settings/",
            "useProgressivePanchangam.ts",
            "skyViewApi.ts",
            "locationService.ts",
        ]:
            self.assertIn(current_frontend_area, architecture_doc)
        for stale_frontend_area in [
            "PanchangamDisplay.tsx",
            "TithiCard.tsx",
            "SkyView/",
            "Common/",
            "Layout/",
            "useSkyViewData.ts",
            "useLocation.ts",
            "panchangamService.ts",
            "skyviewService.ts",
        ]:
            self.assertNotIn(stale_frontend_area, architecture_doc)
        self.assertIn("service Panchangam {", architecture_doc)
        self.assertIn("rpc Get(GetPanchangamRequest) returns (GetPanchangamResponse);", architecture_doc)
        self.assertIn("message GetPanchangamRequest", architecture_doc)
        self.assertIn("message GetPanchangamResponse", architecture_doc)
        self.assertIn("gRPC Panchangam.Get()", architecture_doc)
        self.assertNotIn("service PanchangamService", architecture_doc)
        self.assertNotIn("PanchangamService", architecture_doc)
        self.assertNotIn("CalculatePanchangam(", architecture_doc)
        self.assertNotIn("CalculateMuhurta", architecture_doc)
        self.assertNotIn("GetFestivals", architecture_doc)
        for stale_pattern in [
            "Repository Pattern",
            "Factory Pattern",
            "Strategy Pattern",
            "Dependency Injection",
            "Authentication (Future)",
            "Authorization",
        ]:
            self.assertNotIn(stale_pattern, architecture_doc)
        self.assertIn("Keep code paths direct", architecture_doc)

    def test_stale_aaa_interceptor_package_is_absent(self):
        self.assertFalse((PROJECT_ROOT / "aaa").exists())

        for relative_path in [
            "docs/design/high-level-architecture.md",
            "docs/design/README.md",
        ]:
            self.assertNotIn("Custom AAA interceptors", read_text(relative_path))

    def test_ci_workflow_does_not_claim_placeholder_deployments(self):
        workflow = read_text(".github/workflows/ci-cd.yml")

        self.assertIn("name: CI Pipeline", workflow)
        self.assertNotIn("CI/CD Pipeline", workflow)
        self.assertNotIn("deploy_environment:", workflow)
        self.assertNotIn("deploy-staging:", workflow)
        self.assertNotIn("deploy-production:", workflow)
        self.assertNotIn("notify-completion:", workflow)
        self.assertNotIn("pull-requests: write", workflow)
        self.assertNotIn("checks: write", workflow)
        self.assertNotIn("Add actual deployment logic", workflow)
        self.assertNotIn("Add smoke test logic", workflow)
        self.assertNotIn("Add production smoke test logic", workflow)
        self.assertNotIn("github.rest.repos.createDeploymentStatus", workflow)
        self.assertNotIn("context.payload.deployment.id", workflow)
        self.assertNotIn("panchangam-staging.example.com", workflow)
        self.assertNotIn("panchangam.example.com", workflow)

    def test_makefile_ci_helpers_do_not_push_or_deploy(self):
        makefile = read_text("Makefile")

        self.assertIn("# Local Verification Commands", makefile)
        ci_helpers = text_between(makefile, "# Local Verification Commands", "# Utilities")

        self.assertIn("ci-lint:", ci_helpers)
        self.assertIn("ci-test:", ci_helpers)
        self.assertIn("ci-build:", ci_helpers)
        self.assertIn("Running local lint gate", ci_helpers)
        self.assertIn("Running local test gate", ci_helpers)
        self.assertIn("Running local build gate", ci_helpers)
        self.assertNotIn("CI/CD Pipeline Commands", makefile)
        self.assertNotIn("Running CI", ci_helpers)
        self.assertNotIn("ci-deploy:", ci_helpers)
        self.assertNotIn("docker-push", ci_helpers)
        self.assertNotIn("deploy-staging", ci_helpers)
        self.assertNotIn("Running CI deployment", ci_helpers)

    def test_image_jobs_only_push_outside_pull_requests(self):
        workflow = read_text(".github/workflows/ci-cd.yml")
        package_jobs = workflow.split("  # Build & Package Jobs", 1)[1]

        self.assertIn("if: github.event_name != 'pull_request'", package_jobs)
        self.assertIn("push: ${{ github.event_name != 'pull_request' }}", package_jobs)
        self.assertIn("- name: Build Docker image", package_jobs)
        self.assertNotIn("- name: Build and push Docker image", package_jobs)
        self.assertNotIn("push: true", package_jobs)

    def test_image_jobs_wait_for_e2e_and_do_not_expose_unused_outputs(self):
        workflow = read_text(".github/workflows/ci-cd.yml")
        backend_job = text_between(workflow, "  build-backend:", "  build-frontend:")
        frontend_job = workflow.split("  build-frontend:", 1)[1]

        self.assertIn("needs: [code-quality, security-scan, backend-test, e2e-test]", backend_job)
        self.assertIn("needs: [code-quality, frontend-test, e2e-test]", frontend_job)

        for job in (backend_job, frontend_job):
            self.assertNotIn("id: build", job)
            self.assertNotIn("outputs:", job)
            self.assertNotIn("image-digest:", job)
            self.assertNotIn("steps.build.outputs.digest", job)


if __name__ == "__main__":
    unittest.main()
