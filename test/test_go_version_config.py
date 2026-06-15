import re
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def go_major_minor() -> str:
    match = re.search(r"^go\s+(\d+\.\d+)", (ROOT / "go.mod").read_text(), re.MULTILINE)
    if not match:
        raise AssertionError("go.mod should declare a Go version")
    return match.group(1)


class GoVersionConfigTest(unittest.TestCase):
    def test_build_config_uses_go_mod_major_minor(self):
        version = go_major_minor()

        checks = {
            "Makefile": f"GO_VERSION := {version}",
            ".github/workflows/ci-cd.yml": f"GO_VERSION: '{version}'",
            "cmd/panchangam-cli/local_commands.go": f'"go_version":  "{version}+"',
        }
        for name, expected in checks.items():
            content = (ROOT / name).read_text()
            self.assertIn(expected, content, f"{name} should use Go {version}")

    def test_developer_docs_use_go_mod_major_minor(self):
        version = go_major_minor()

        for name in ["docs/README.md", "docs/DEVELOPER_GUIDE.md"]:
            content = (ROOT / name).read_text()
            self.assertIn(f"Go {version}+", content, f"{name} should document Go {version}+")


if __name__ == "__main__":
    unittest.main()
