#!/usr/bin/env python3
"""Deployment health-check wiring tests."""

import unittest
import os
import subprocess
import tempfile
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]


def read_text(relative_path: str) -> str:
    return (PROJECT_ROOT / relative_path).read_text(encoding="utf-8")


def text_between(content: str, start: str, end: str) -> str:
    return content.split(start, 1)[1].split(end, 1)[0]


class DeploymentHealthConfigTest(unittest.TestCase):
    def test_kubernetes_base_uses_real_probes(self):
        manifest = read_text("deployments/k8s/base/backend-deployment.yaml")

        self.assertNotIn("--health-check", manifest)
        self.assertIn('command: ["/panchangam-grpc"]', manifest)
        self.assertIn('args: ["--grpc-port=50052"]', manifest)
        self.assertIn('args: ["--grpc-endpoint=grpc-server:50052", "--http-port=8080"]', manifest)
        self.assertGreaterEqual(manifest.count("grpc:\n              port: 50052"), 2)
        self.assertGreaterEqual(manifest.count("path: /api/v1/health"), 2)

    def test_generated_deploy_manifest_uses_real_probes(self):
        script = read_text("scripts/deploy.sh")

        self.assertNotIn("--health-check", script)
        self.assertIn('args: ["--grpc-endpoint=localhost:50052", "--http-port=8080"]', script)
        self.assertIn('command: ["/panchangam-grpc"]', script)
        self.assertIn('args: ["--grpc-port=50052"]', script)
        self.assertGreaterEqual(script.count("grpc:\n            port: 50052"), 2)

    def test_backend_dockerfile_builds_canonical_server(self):
        dockerfile = read_text("docker/Dockerfile.backend")

        self.assertNotIn("./server", dockerfile)
        self.assertIn("./cmd/server", dockerfile)
        self.assertIn('CMD ["/panchangam-gateway"]', dockerfile)
        self.assertIn('HEALTHCHECK', dockerfile)
        self.assertIn('"/panchangam-gateway", "--health-check"', dockerfile)

    def test_nginx_content_security_policy_does_not_allow_eval(self):
        nginx = read_text("deployments/nginx/nginx.conf")
        csp_line = next(line for line in nginx.splitlines() if "Content-Security-Policy" in line)

        self.assertNotIn("'unsafe-eval'", csp_line)
        self.assertIn("script-src 'self' 'unsafe-inline'", csp_line)

    def test_deploy_scripts_use_real_production_domains(self):
        for relative_path in ["scripts/deploy.sh", "scripts/validate-deployment.sh"]:
            script = read_text(relative_path)

            self.assertNotIn("panchangam.example.com", script)
            self.assertIn("panchangam.app", script)
            self.assertIn("api.panchangam.app", script)

    def test_staging_smoke_tests_fall_back_when_service_ip_is_missing(self):
        script = read_text("scripts/deploy.sh")

        self.assertIn("service_ip=$(kubectl get svc panchangam-backend", script)
        self.assertIn('if [ -n "$service_ip" ]; then', script)
        self.assertIn('service_url="http://$service_ip"', script)
        self.assertIn('service_url="http://localhost:8080"', script)
        self.assertNotIn('service_url="http://$(kubectl get svc panchangam-backend', script)

    def test_root_deploy_script_requires_only_tools_it_uses(self):
        script = read_text("scripts/deploy.sh")

        self.assertIn('if [ "$DRY_RUN" = "false" ] && ! command -v docker', script)
        self.assertNotIn("if ! command -v docker", script)
        self.assertIn("command -v kubectl", script)
        self.assertNotIn("command -v helm", script)

    def test_dry_run_does_not_claim_images_were_validated(self):
        script = read_text("scripts/deploy.sh")
        validate_images = text_between(script, "validate_images() {", "\n# Generate Kubernetes manifests")

        self.assertIn('if [ "$DRY_RUN" = "true" ]; then', validate_images)
        dry_run_block = text_between(validate_images, 'if [ "$DRY_RUN" = "true" ]; then', "\n    fi")
        self.assertIn('log_info "DRY RUN: Would validate container images"', dry_run_block)
        self.assertIn("return", dry_run_block)
        self.assertLess(
            validate_images.index('if [ "$DRY_RUN" = "true" ]; then'),
            validate_images.index('log_success "Container images validated"'),
        )

    def test_generated_manifest_dir_uses_unique_temp_directory(self):
        script = read_text("scripts/deploy.sh")

        self.assertIn('manifest_dir=$(mktemp -d "${TMPDIR:-/tmp}/panchangam-deploy-${VERSION}.XXXXXX")', script)
        self.assertNotIn('local manifest_dir="/tmp/panchangam-deploy-${VERSION}"', script)
        self.assertNotIn('mkdir -p "$manifest_dir"', script)

    def test_generate_manifests_keeps_stdout_to_manifest_path_only(self):
        script = read_text("scripts/deploy.sh")
        generate_manifests = text_between(script, "generate_manifests() {", "\n# Deploy to Kubernetes")

        self.assertIn('log_info "Generating Kubernetes manifests..." >&2', generate_manifests)
        self.assertNotIn('log_info "Generating Kubernetes manifests..."\n', generate_manifests)
        self.assertIn('echo "$manifest_dir"', generate_manifests)

    def test_generate_manifests_names_environment_specific_values_once(self):
        script = read_text("scripts/deploy.sh")
        generate_manifests = text_between(script, "generate_manifests() {", "\n# Deploy to Kubernetes")

        self.assertIn('local replicas="2"', generate_manifests)
        self.assertIn('local log_level="info"', generate_manifests)
        self.assertIn('replicas="3"', generate_manifests)
        self.assertIn('log_level="warn"', generate_manifests)
        self.assertIn("replicas: ${replicas}", generate_manifests)
        self.assertIn('value: "${log_level}"', generate_manifests)
        self.assertNotIn('replicas: $( [ "$ENVIRONMENT" = "production" ]', generate_manifests)
        self.assertNotIn('value: "$( [ "$ENVIRONMENT" = "production" ]', generate_manifests)

    def test_generated_manifest_dir_is_cleaned_on_exit(self):
        script = read_text("scripts/deploy.sh")
        main = text_between(script, "main() {", "\n# Run main function")

        self.assertIn('MANIFEST_DIR=""', script)
        self.assertIn('manifest_dir=$(generate_manifests)', main)
        self.assertIn('MANIFEST_DIR="$manifest_dir"', main)
        self.assertIn('trap \'rm -rf -- "$MANIFEST_DIR"\' EXIT', main)
        self.assertNotIn('trap \'rm -rf "$manifest_dir"\' EXIT', main)
        self.assertNotIn('# Cleanup\n        rm -rf "$manifest_dir"', main)

    def test_dry_run_does_not_apply_namespace(self):
        script = read_text("scripts/deploy.sh")
        deploy_to_kubernetes = text_between(script, "deploy_to_kubernetes() {", "\n# Run smoke tests")
        dry_run_block = text_between(deploy_to_kubernetes, 'if [ "$DRY_RUN" = "true" ]; then', "\n    fi")

        self.assertIn('kubectl apply -f "$manifest_dir" --dry-run=client', dry_run_block)
        self.assertIn("return", dry_run_block)
        self.assertNotIn("kubectl apply -f -", dry_run_block)
        self.assertLess(
            deploy_to_kubernetes.index('if [ "$DRY_RUN" = "true" ]; then'),
            deploy_to_kubernetes.index('kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -'),
        )

    def test_production_confirmation_reports_non_interactive_input(self):
        script = read_text("scripts/deploy.sh")
        main = text_between(script, "main() {", "\n# Run main function")

        self.assertIn("if ! read -r confirmation; then", main)
        self.assertIn('log_error "Production confirmation requires input"', main)
        self.assertNotIn("        read -r confirmation\n", main)

    def test_compose_health_checks_use_real_binary_checks(self):
        compose = read_text("docker-compose.prod.yml")

        self.assertIn('command: ["/panchangam-grpc", "--grpc-port=50052"]', compose)
        self.assertIn('test: ["CMD", "/panchangam-grpc", "--health-check", "--grpc-port=50052"]', compose)
        self.assertIn('test: ["CMD", "/panchangam-gateway", "--health-check", "--http-port=8080"]', compose)

    def test_makefile_docker_targets_use_project_compose_file(self):
        makefile = read_text("Makefile")

        self.assertIn("docker-compose -f docker-compose.prod.yml up --build", makefile)
        self.assertIn("docker-compose -f docker-compose.prod.yml down", makefile)
        self.assertIn(
            "docker-compose -f docker-compose.prod.yml up --force-recreate --remove-orphans --detach",
            makefile,
        )
        self.assertIn("Panchangam services are running.", makefile)
        self.assertIn("http://localhost:8080/api/v1/health", makefile)
        self.assertNotIn("OpenTelemetry Demo is running.", makefile)
        self.assertNotIn("http://192.168.68.73:16686/", makefile)
        self.assertNotIn("docker-compose up --build", makefile)
        self.assertNotIn("docker compose up --force-recreate", makefile)

    def test_project_does_not_ship_stale_default_compose_file(self):
        self.assertFalse((PROJECT_ROOT / "docker-compose.yml").exists())

        for relative_path in ["DEPLOYMENT.md", "PROJECT_COMPLETION_SUMMARY.md"]:
            content = read_text(relative_path)
            self.assertIn("docker-compose -f docker-compose.prod.yml up -d", content)
            self.assertNotIn("docker-compose up -d", content)
            self.assertNotIn("**docker-compose.yml:**", content)

    def test_project_uses_focused_dockerfiles_not_root_supervisor_image(self):
        self.assertFalse((PROJECT_ROOT / "Dockerfile").exists())

        makefile = read_text("Makefile")
        compose = read_text("docker-compose.prod.yml")
        test_compose = read_text("test/docker-compose.yml")
        test_docs = read_text("test/DOCKER_SETUP.md")

        self.assertIn("docker/Dockerfile.backend", makefile)
        self.assertIn("ui/Dockerfile", makefile)
        self.assertIn("dockerfile: docker/Dockerfile.backend", compose)
        self.assertIn("dockerfile: Dockerfile", test_compose)
        self.assertIn("test/\n├── Dockerfile", test_docs)

        for content in (makefile, compose):
            self.assertNotIn("-f Dockerfile", content)
            self.assertNotIn("supervisord", content)

        self.assertNotIn("Panchangam Multi-Stage Dockerfile", test_docs)

    def test_deployment_guide_points_to_tracked_deployment_artifacts(self):
        deployment_guide = read_text("DEPLOYMENT.md")

        self.assertIn("docker/Dockerfile.backend", deployment_guide)
        self.assertIn("ui/Dockerfile", deployment_guide)
        self.assertIn("docker-compose.prod.yml", deployment_guide)
        self.assertIn("deployments/k8s/overlays/production", deployment_guide)

        self.assertNotIn("**Dockerfile (gRPC Server):**", deployment_guide)
        self.assertNotIn("See `k8s/` directory", deployment_guide)
        self.assertNotIn("kubectl apply -f k8s/", deployment_guide)

    def test_operator_docs_use_gateway_health_endpoint(self):
        for relative_path in ["deployments/README.md", "deployments/RUNBOOK.md"]:
            content = read_text(relative_path)

            self.assertNotIn("localhost:8080/health", content)
            self.assertIn("localhost:8080/api/v1/health", content)

    def test_deployments_do_not_set_unused_port_env_vars(self):
        exact_stale_lines = {
            "deployments/k8s/base/backend-deployment.yaml": [
                '- name: PORT\n              value: "50052"',
                '- name: GRPC_PORT',
                '- name: HTTP_PORT',
                '- name: GRPC_ENDPOINT',
            ],
            "scripts/deploy.sh": [
                '- name: HTTP_PORT',
                '- name: GRPC_ENDPOINT',
                '- name: PORT\n          value: "50052"',
            ],
            "docker-compose.prod.yml": [
                "- PORT=50052",
                "- GRPC_PORT=50052",
                "- HTTP_PORT=8080",
                "- GRPC_ENDPOINT=grpc-server:50052",
            ],
        }

        for relative_path, stale_lines in exact_stale_lines.items():
            content = read_text(relative_path)
            for stale_line in stale_lines:
                self.assertNotIn(stale_line, content)

    def test_backend_deployments_use_gateway_redis_addr_only(self):
        exact_stale_lines = {
            "deployments/k8s/base/backend-deployment.yaml": [
                "- name: REDIS_HOST",
                "- name: REDIS_PORT",
                "- name: DB_HOST",
                "- name: DB_PORT",
                "- name: DB_NAME",
                "- name: DB_USER",
                "- name: DB_PASSWORD",
            ],
            "docker-compose.prod.yml": [
                "- REDIS_HOST=redis",
                "- REDIS_PORT=6379",
                "- DB_HOST=postgres",
                "- DB_PORT=5432",
                "- DB_NAME=${POSTGRES_DB:-panchangam}",
                "- DB_USER=${POSTGRES_USER:-panchangam}",
                "- DB_PASSWORD=${POSTGRES_PASSWORD:-panchangam123}",
            ],
        }

        for relative_path, stale_lines in exact_stale_lines.items():
            content = read_text(relative_path)
            for stale_line in stale_lines:
                self.assertNotIn(stale_line, content)

        manifest = read_text("deployments/k8s/base/backend-deployment.yaml")
        self.assertIn("- name: REDIS_ADDR", manifest)
        self.assertIn("key: REDIS_ADDR", manifest)

        compose = read_text("docker-compose.prod.yml")
        self.assertIn("- REDIS_ADDR=redis:6379", compose)

        configmap = read_text("deployments/k8s/base/configmap.yaml")
        self.assertIn('REDIS_ADDR: "redis:6379"', configmap)
        self.assertNotIn("REDIS_HOST", configmap)
        self.assertNotIn("REDIS_PORT", configmap)
        self.assertNotIn("DB_HOST", configmap)
        self.assertNotIn("DB_PORT", configmap)
        self.assertNotIn("ALLOW_ALL_ORIGINS", configmap)

        compose = read_text("docker-compose.prod.yml")
        self.assertNotIn("- ALLOW_ALL_ORIGINS=", compose)

    def test_compose_omits_unused_database_and_backup_services(self):
        compose = read_text("docker-compose.prod.yml")
        prometheus = read_text("deployments/prometheus/prometheus.yml")
        alerts = read_text("deployments/prometheus/alerts/application.yml")

        stale_compose_lines = [
            "  postgres:",
            "  postgres-replica:",
            "  backup:",
            "POSTGRES_",
            "PGHOST=postgres",
            "postgres_data:",
            "postgres_replica_data:",
            "backup_data:",
            "./deployments/postgres/init",
            "./deployments/backup/backup.sh",
        ]
        for stale_line in stale_compose_lines:
            self.assertNotIn(stale_line, compose)

        self.assertNotIn("job_name: 'postgres'", prometheus)
        self.assertNotIn("targets: ['postgres:5432']", prometheus)
        self.assertNotIn("panchangam_database", alerts)
        self.assertNotIn("PostgreSQLDown", alerts)
        self.assertNotIn("pg_stat_", alerts)

    def test_deploy_commands_do_not_reference_unused_database_operations(self):
        makefile = read_text("Makefile")
        deploy_script = read_text("deployments/scripts/deploy.sh")

        stale_makefile_lines = [
            "# Database Operations",
            "# Backup & Recovery",
            "migrate-up:",
            "migrate-down:",
            "migrate-create:",
            "migrate-version:",
            "backup:",
            "restore:",
            "./deployments/migrations/migrate.sh",
            "./deployments/backup/backup.sh",
            "./deployments/backup/restore.sh",
            "deployments/{postgres/init",
            ",backup,nginx}",
        ]
        for stale_line in stale_makefile_lines:
            self.assertNotIn(stale_line, makefile)

        stale_deploy_lines = [
            "run_migrations",
            "Running database migrations",
            "postgres-primary",
            "DB_PASSWORD",
            "backend-gateway",
            "./migrate.sh up",
        ]
        for stale_line in stale_deploy_lines:
            self.assertNotIn(stale_line, deploy_script)

    def test_secondary_deploy_script_uses_strict_mode_and_safe_env_loading(self):
        deploy_script = read_text("deployments/scripts/deploy.sh")

        self.assertIn("set -euo pipefail", deploy_script)
        self.assertNotIn("set -e\n", deploy_script)
        self.assertNotIn("export $(cat", deploy_script)
        self.assertNotIn("grep -v '^#' | xargs", deploy_script)
        self.assertNotIn("set -a", deploy_script)
        self.assertNotIn('. ".env.${ENVIRONMENT}"', deploy_script)
        self.assertIn("while IFS= read -r line || [ -n \"$line\" ]; do", deploy_script)
        self.assertIn('key="${line%%=*}"', deploy_script)
        self.assertIn('export "$line"', deploy_script)

    def test_secondary_deploy_script_supports_rollback_command(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            fake_compose = Path(temp_dir) / "docker-compose"
            fake_compose.write_text(
                "#!/bin/sh\n"
                "printf 'fake docker-compose %s\\n' \"$*\"\n",
                encoding="utf-8",
            )
            fake_compose.chmod(0o755)

            env = {
                **os.environ,
                "PATH": f"{temp_dir}{os.pathsep}/usr/bin:/bin:/usr/sbin:/sbin",
            }
            result = subprocess.run(
                [
                    "bash",
                    str(PROJECT_ROOT / "deployments/scripts/deploy.sh"),
                    "rollback",
                    "docker-compose",
                ],
                cwd=PROJECT_ROOT,
                env=env,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                check=False,
            )

        self.assertEqual(0, result.returncode, result.stdout)
        self.assertIn("Rolling back deployment", result.stdout)
        self.assertIn("fake docker-compose -f docker-compose.prod.yml down", result.stdout)
        self.assertIn("Rollback completed", result.stdout)
        self.assertNotIn("Invalid environment: rollback", result.stdout)

    def test_secondary_deploy_script_reports_missing_rollback_tools(self):
        env = {
            **os.environ,
            "PATH": "/usr/bin:/bin:/usr/sbin:/sbin",
        }
        cases = [
            ("kubernetes", "kubectl is not installed"),
            ("docker-compose", "docker-compose is not installed"),
        ]

        for platform, expected_error in cases:
            with self.subTest(platform=platform):
                result = subprocess.run(
                    [
                        "bash",
                        str(PROJECT_ROOT / "deployments/scripts/deploy.sh"),
                        "rollback",
                        platform,
                    ],
                    cwd=PROJECT_ROOT,
                    env=env,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.STDOUT,
                    text=True,
                    check=False,
                )

                self.assertEqual(1, result.returncode, result.stdout)
                self.assertIn(expected_error, result.stdout)
                self.assertNotIn("command not found", result.stdout)
                self.assertNotIn("Rolling back deployment", result.stdout)

    def test_secondary_deploy_script_names_kubernetes_deployments_once(self):
        deploy_script = read_text("deployments/scripts/deploy.sh")

        self.assertIn("KUBERNETES_DEPLOYMENTS=(", deploy_script)
        for deployment in [
            "panchangam-grpc",
            "panchangam-gateway",
            "panchangam-frontend",
        ]:
            self.assertEqual(1, deploy_script.count(deployment), deployment)

        self.assertIn('for deployment in "${KUBERNETES_DEPLOYMENTS[@]}"; do', deploy_script)
        self.assertIn('kubectl rollout status "deployment/$deployment"', deploy_script)
        self.assertIn('kubectl rollout undo "deployment/$deployment"', deploy_script)
        self.assertNotIn("kubectl rollout status deployment/panchangam-", deploy_script)
        self.assertNotIn("kubectl rollout undo deployment/panchangam-", deploy_script)

    def test_project_omits_unused_database_tooling_and_operator_docs(self):
        for relative_path in [
            "deployments/postgres",
            "deployments/migrations",
            "deployments/backup",
        ]:
            self.assertFalse((PROJECT_ROOT / relative_path).exists())

        stale_doc_lines = [
            "PostgreSQL",
            "Database Operations",
            "Database Connection Issues",
            "Database Upgrade Procedure",
            "Backup & Recovery",
            "Backup Procedures Tested",
            "deployments/postgres",
            "deployments/migrations",
            "deployments/backup",
            "postgres-primary",
            "postgres-replica",
            "POSTGRES_",
            "pg_dump",
            "pg_stat_",
        ]
        for relative_path in [
            "deployments/README.md",
            "deployments/RUNBOOK.md",
            "deployments/IMPLEMENTATION_SUMMARY.md",
            "deployments/grafana/dashboards/overview.json",
        ]:
            content = read_text(relative_path)
            for stale_line in stale_doc_lines:
                self.assertNotIn(stale_line, content)

    def test_design_docs_match_current_runtime_storage(self):
        stale_lines = [
            "PostgreSQL",
            "POSTGRES",
            "postgres:",
            "Configuration DB",
            "Config DB",
            "Database Architecture",
            "Database Layer",
            "pgbouncer",
            "read_replicas",
        ]
        for relative_path in [
            "docs/design/README.md",
            "docs/design/high-level-architecture.md",
            "docs/design/everyday-integration-architecture.md",
        ]:
            content = read_text(relative_path)
            for stale_line in stale_lines:
                self.assertNotIn(stale_line, content)

    def test_kubernetes_base_omits_unused_database_resources(self):
        self.assertFalse((PROJECT_ROOT / "deployments/k8s/base/postgres-deployment.yaml").exists())

        kustomization = read_text("deployments/k8s/base/kustomization.yaml")
        configmap = read_text("deployments/k8s/base/configmap.yaml")
        secret = read_text("deployments/k8s/base/secret.yaml")

        stale_lines = [
            "postgres-deployment.yaml",
            "DB_NAME",
            "DB_USER",
            "DB_PASSWORD",
            "POSTGRES_PASSWORD",
            "Database configuration",
            "Database credentials",
            "postgres-primary",
            "postgres-pvc",
            "POSTGRES_",
        ]
        for stale_line in stale_lines:
            self.assertNotIn(stale_line, kustomization)
            self.assertNotIn(stale_line, configmap)
            self.assertNotIn(stale_line, secret)

    def test_compose_backend_dependencies_match_runtime_needs(self):
        compose = read_text("docker-compose.prod.yml")
        grpc_service = text_between(compose, "  # Backend gRPC service", "  # Backend HTTP Gateway")
        gateway_service = text_between(compose, "  # Backend HTTP Gateway", "  # Frontend React application")

        self.assertNotIn("postgres:", grpc_service)
        self.assertNotIn("redis:", grpc_service)
        self.assertIn("jaeger:", grpc_service)

        self.assertIn("grpc-server:", gateway_service)
        self.assertIn("redis:", gateway_service)
        self.assertNotIn("postgres:", gateway_service)

    def test_redis_health_checks_authenticate(self):
        manifest = read_text("deployments/k8s/base/redis-deployment.yaml")
        self.assertNotIn("- redis-cli\n                - ping", manifest)
        self.assertGreaterEqual(manifest.count('REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli ping'), 2)

        compose = read_text("docker-compose.prod.yml")
        redis_service = text_between(compose, "  # Redis for caching", "  # ====================================\n  # Application Services")
        self.assertIn("- REDIS_PASSWORD=${REDIS_PASSWORD:-panchangam}", redis_service)
        self.assertIn(
            'test: ["CMD-SHELL", "REDISCLI_AUTH=\\"$${REDIS_PASSWORD}\\" redis-cli --raw incr ping"]',
            redis_service,
        )
        self.assertNotIn('test: ["CMD", "redis-cli", "--raw", "incr", "ping"]', redis_service)

    def test_operator_redis_commands_authenticate(self):
        runbook = read_text("deployments/RUNBOOK.md")

        self.assertNotIn("-- redis-cli\n", runbook)
        self.assertNotIn("\n    redis-cli INFO stats", runbook)
        self.assertIn('REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli', runbook)
        self.assertIn('REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli INFO stats', runbook)

    def test_operator_cache_clear_matches_cache_key_prefix(self):
        runbook = read_text("deployments/RUNBOOK.md")

        self.assertNotIn("panchangam:cache:*", runbook)
        self.assertNotIn("KEYS panchangam:", runbook)
        self.assertNotIn("DEL panchangam:", runbook)
        self.assertIn('REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli FLUSHALL', runbook)
        self.assertIn('redis-cli --scan --pattern "panchangam:*"', runbook)
        self.assertIn("xargs -r redis-cli DEL", runbook)

    def test_api_health_routes_use_gateway_health_endpoint(self):
        runbook = read_text("deployments/RUNBOOK.md")
        self.assertNotIn("api.panchangam.app/health", runbook)
        self.assertIn("api.panchangam.app/api/v1/health", runbook)

        deploy_script = read_text("deployments/scripts/deploy.sh")
        self.assertNotIn('${API_URL}/health', deploy_script)
        self.assertIn('${API_URL}/api/v1/health', deploy_script)

        load_balancer = read_text("deployments/nginx/lb.conf")
        self.assertNotIn("proxy_pass http://backend/health;", load_balancer)
        self.assertIn("location /api/v1/health", load_balancer)
        self.assertIn("proxy_pass http://backend/api/v1/health;", load_balancer)

        ingress = read_text("deployments/k8s/base/ingress.yaml")
        self.assertNotIn("path: /health", ingress)
        self.assertIn("path: /api/v1/health", ingress)


if __name__ == "__main__":
    unittest.main()
