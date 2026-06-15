#!/usr/bin/env python3
"""Documentation state tests."""

import unittest
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]


def read_text(relative_path: str) -> str:
    return (PROJECT_ROOT / relative_path).read_text(encoding="utf-8")


class DocumentationStateTest(unittest.TestCase):
    def test_llm_docs_do_not_teach_removed_api_or_component_names(self):
        stale_terms = [
            "CalculatePanchangam",
            "NewPanchangamService",
            "PanchangamService",
            "CalculateDaily",
            "PanchangamDisplay",
            "DateNavigator",
            "TithiCard",
            "usePanchangamData",
            "gRPC-Gateway",
        ]

        for relative_path in [
            "llm/coding-standards.md",
            "llm/project-architecture.md",
            "llm/testing-guidelines.md",
        ]:
            content = read_text(relative_path)
            for stale_term in stale_terms:
                self.assertNotIn(stale_term, content, f"{relative_path} contains stale term")

    def test_llm_coding_standards_keep_abstractions_simple(self):
        coding_standards = read_text("llm/coding-standards.md")

        self.assertIn("NewPanchangamServer", coding_standards)
        self.assertIn("MonthNavigation", coding_standards)
        self.assertIn("usePanchangam", coding_standards)
        self.assertIn("Add interfaces only when a caller needs them", coding_standards)
        self.assertNotIn("SOLID Principles", coding_standards)
        self.assertNotIn("Dependency Injection", coding_standards)
        self.assertNotIn("options pattern", coding_standards)

    def test_llm_testing_guidelines_use_current_panchangam_service_api(self):
        testing_guidelines = read_text("llm/testing-guidelines.md")

        self.assertIn("NewPanchangamServer()", testing_guidelines)
        self.assertIn("ppb.GetPanchangamRequest", testing_guidelines)
        self.assertIn("server.Get(ctx, req)", testing_guidelines)
        self.assertIn("Panchangam.Get", testing_guidelines)

        for stale_api in [
            "NewPanchangamService",
            "PanchangamService",
            "CalculatePanchangam",
            "CalculateDaily",
        ]:
            self.assertNotIn(stale_api, testing_guidelines)

    def test_llm_testing_guidelines_use_current_frontend_examples(self):
        testing_guidelines = read_text("llm/testing-guidelines.md")

        for current_example in [
            "MonthNavigation",
            "SettingsPanel",
            "usePanchangam",
            "panchangamApiClient",
        ]:
            self.assertIn(current_example, testing_guidelines)

        for stale_example in [
            "PanchangamDisplay",
            "DateNavigator",
            "TithiCard",
            "usePanchangamData",
            "panchangamApi.getPanchangam",
            "@/services/panchangamApi",
            "PanchangamDisplay.test.tsx",
        ]:
            self.assertNotIn(stale_example, testing_guidelines)

    def test_service_docs_do_not_describe_calculated_response_as_placeholder_data(self):
        stale_claims = [
            "Service returns placeholder data instead of calculated values",
            "Service layer not integrated with astronomy calculation modules",
            "Functional tests validate placeholder data structure",
            "Replace placeholder data with calculated values",
            "replace placeholder data with real calculations",
            "Replaced all placeholder data with",
            "Service integration missing (confirmed by placeholder data)",
            "primary remaining work is integrating the service layer",
            "WARN Partial (placeholder data)",
            "Falls back to placeholder data",
            "Integration Tests Needed",
            "Integration tests: WARN **Missing**",
            "Performance tests: WARN **Missing**",
            "random error simulation",
            "simulation delays in service",
            "Missing real calculation integration tests",
            "timeout due to delays",
            "currently 61%",
        ]

        for relative_path in [
            "FEATURES.md",
            "FEATURE_COVERAGE_REPORT.md",
            "PROJECT_COMPLETION_SUMMARY.md",
            "ui/README.md",
        ]:
            content = read_text(relative_path)
            for stale_claim in stale_claims:
                self.assertNotIn(stale_claim, content, f"{relative_path} contains stale service claim")

    def test_developer_guide_uses_current_panchangam_get_response_shape(self):
        developer_guide = read_text("docs/DEVELOPER_GUIDE.md")

        self.assertIn("TestPanchangamGetIntegration", developer_guide)
        self.assertIn("client.Get(context.Background(), req)", developer_guide)
        self.assertIn("resp.PanchangamData", developer_guide)
        self.assertNotIn("TestPanchangamServiceIntegration", developer_guide)
        self.assertNotIn("resp.Data", developer_guide)


if __name__ == "__main__":
    unittest.main()
