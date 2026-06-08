import json
import re
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def package_scripts() -> set[str]:
    package_json = json.loads((ROOT / "ui/package.json").read_text())
    return set(package_json.get("scripts", {}))


class FrontendScriptConfigTest(unittest.TestCase):
    def test_frontend_package_metadata_names_the_project(self):
        package_json = json.loads((ROOT / "ui/package.json").read_text())
        package_lock = json.loads((ROOT / "ui/package-lock.json").read_text())

        self.assertEqual("panchangam-ui", package_json["name"])
        self.assertEqual("panchangam-ui", package_lock["name"])
        self.assertEqual("panchangam-ui", package_lock["packages"][""]["name"])
        self.assertNotEqual("vite-react-typescript-starter", package_json["name"])

    def test_makefile_and_ci_only_call_defined_frontend_scripts(self):
        scripts = package_scripts()

        for name in ["Makefile", ".github/workflows/ci-cd.yml"]:
            content = (ROOT / name).read_text()
            for script in re.findall(r"npm run ([A-Za-z0-9:_-]+)", content):
                self.assertIn(script, scripts, f"{name} calls missing ui npm script {script}")

    def test_ui_readme_only_documents_defined_frontend_scripts(self):
        scripts = package_scripts()
        content = (ROOT / "ui/README.md").read_text()

        for script in re.findall(r"npm run ([A-Za-z0-9:_-]+)", content):
            self.assertIn(script, scripts, f"ui/README.md documents missing npm script {script}")

    def test_playwright_only_collects_e2e_specs(self):
        content = (ROOT / "ui/playwright.config.js").read_text()

        self.assertIn("testDir: './e2e'", content)
        self.assertIn("testMatch: '**/*.spec.ts'", content)


if __name__ == "__main__":
    unittest.main()
