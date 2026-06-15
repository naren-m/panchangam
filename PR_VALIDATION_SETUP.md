# PR Validation Pipeline Setup - Issue #81

## Manual Setup Required

GitHub workflow files (`.github/workflows/*.yml`) cannot be automatically pushed via the API due to security restrictions. You need to manually add these files to enable the PR validation pipeline.

---

## Files Created

The following workflow files have been created locally and need to be committed manually:

### 1. **`.github/workflows/pr-validation.yml`** (New File)
**Purpose**: Comprehensive PR validation workflow for Issue #81

**What it does**:
- Runs automatically on every PR to `main` or `develop`
- Executes complete integration test suite (Issue #81)
- Validates all performance targets
- Posts detailed results as PR comment
- Provides "All Tests Passed" status for branch protection

### 2. **`.github/workflows/ci-cd.yml`** (Modified)
**Changes made**:
- Added `integration-test` job with complete Issue #81 test suite
- Includes Redis service for cache testing
- Updated dependencies for e2e-test job
- Enhanced PR comment to show integration test results

### 3. **`.github/workflows/README.md`** (New File)
**Purpose**: Documentation for all workflows
- Workflow descriptions and usage
- Integration test suite overview
- Branch protection setup guide
- Troubleshooting guide

---

## Quick Setup Steps

### Option 1: Commit Directly (Recommended)

```bash
# Navigate to repository
cd /home/user/panchangam

# Check the workflow files
git status

# Stage workflow files
git add .github/workflows/pr-validation.yml
git add .github/workflows/ci-cd.yml
git add .github/workflows/README.md

# Commit workflow files
git commit -m "feat: add comprehensive PR validation pipeline for Issue #81

Added automated PR validation workflow that runs complete integration
test suite on every pull request.

- New: pr-validation.yml - PR validation workflow
- Updated: ci-cd.yml - Added integration-test job
- New: README.md - Workflow documentation

Issue: #81"

# Push to your branch
git push origin claude/fix-open-issues-01GXhuVxtPN3eLiDDiVVre5i
```

### Option 2: Commit via GitHub UI

If the push fails due to permissions:

1. Go to your repository on GitHub
2. Navigate to the branch `claude/fix-open-issues-01GXhuVxtPN3eLiDDiVVre5i`
3. Click "Add file" → "Create new file"
4. For each workflow file:
   - Create `.github/workflows/pr-validation.yml`
   - Copy content from local file
   - Commit directly to the branch
5. Repeat for `ci-cd.yml` and `README.md`

---

## Setting Up Branch Protection

Once the workflows are pushed, set up branch protection:

### Via GitHub UI:

1. Go to **Repository Settings** → **Branches**
2. Click **Add rule** for `main` branch
3. Configure:
   - PASS Require status checks to pass before merging
   - PASS Require branches to be up to date before merging
   - **Required status checks**:
     - PASS `All Tests Passed` (from pr-validation.yml)
     - PASS `Code Quality Analysis`
     - PASS `Security Analysis`
   - PASS Require pull request reviews before merging
     - Minimum: 1 approval
   - PASS Dismiss stale pull request approvals when new commits are pushed
   - PASS Require conversation resolution before merging
4. Click **Create** or **Save changes**

### Via GitHub CLI:

```bash
gh api repos/naren-m/panchangam/branches/main/protection \
  --method PUT \
  -H "Accept: application/vnd.github+json" \
  --field required_status_checks='{"strict":true,"contexts":["All Tests Passed","Code Quality Analysis","Security Analysis"]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}'
```

---

## What Happens After Setup

### On Every Pull Request:

1. **PR Validation Workflow Triggers**
   ```
   PR Created/Updated
   ↓
   Validate PR Requirements
   ↓
   Start Services (Redis + Backend API)
   ↓
   Run Backend Go Integration Tests
   ↓
   Run Data Validation Tests
   ↓
   Run Performance Benchmarks
   ↓
   Run Cache Integration Tests
   ↓
   Run End-to-End Tests
   ↓
   Generate Coverage Report (90% target)
   ↓
   Post Results Comment on PR
   ↓
   All Tests Passed (or Failed)
   ```

2. **PR Gets Automated Comment**:
   ```markdown
   PR Validation Results

   ### Validation Checks
   PASS PR Requirements: success

   ### Integration Tests (Issue #81)
   PASS Status: success

   All tests passed! PASS Data Validation (Known Astronomical Events)
   PASS Performance Benchmarks (<500ms avg, 50 concurrent in <5s)
   PASS Cache Integration & Consistency
   PASS End-to-End Data Flow
   PASS Error Handling & Recovery (<3s)

   Test reports available in workflow artifacts.

   This PR meets all Issue #81 requirements and is ready for review! ```

3. **Merge Requirements**:
   - PASS All Tests Passed status check must pass
   - PASS At least 1 approval required
   - PASS All conversations resolved
   - PASS Branch up to date with base

---

## Integration Tests Validated

Every PR automatically validates:

| Requirement | Target | Test File |
|-------------|--------|-----------|
| API Response Time | <500ms average | `test_performance.py` |
| Concurrent Requests | 50 in <5 seconds | `test_performance.py` |
| Data Consistency | 100% accuracy | `test_data_validation.py` |
| Error Recovery | <3 seconds | `test_performance.py` |
| Cache Behavior | Hit/Miss consistency | `test_cache_integration.py` |
| End-to-End Flow | Complete validation | `test_e2e_integration.py` |
| Code Coverage | 90% target (CLAUDE.md) | All tests |

---

## Local Testing

Before creating a PR, developers can run tests locally:

```bash
cd test

# Run all Issue #81 tests
make test-issue-81

# Run specific test categories
make test-data-validation
make test-performance
make test-cache
make test-e2e

# Run with coverage
make test-coverage

# Run load tests
make test-load
```

---

## Documentation

After setup, refer to:

- **`.github/workflows/README.md`** - Workflow documentation
- **`test/INTEGRATION_TESTING.md`** - Integration test suite docs
- **`CLAUDE.md`** - Project guidelines (90% coverage requirement)
- **Issue #81** - Original requirements

---

## Verifying Setup

After pushing workflow files:

1. **Check Workflows Tab**:
   - Go to repository → **Actions** tab
   - You should see "PR Validation - Integration Tests" workflow

2. **Test with a PR**:
   - Create a test PR
   - Watch the workflow run automatically
   - Check for PR comment with results

3. **Verify Branch Protection**:
   - Try merging a PR without passing tests
   - Should be blocked with "Required status check not passing"

---

## Troubleshooting

### Workflow Not Running

**Problem**: Workflow doesn't trigger on PR
**Solution**:
- Check `.github/workflows/pr-validation.yml` exists in repository
- Verify workflow syntax: `yamllint .github/workflows/pr-validation.yml`
- Check repository Actions settings are enabled

### Tests Failing

**Problem**: Integration tests fail in CI
**Solution**:
- Check workflow logs in Actions tab
- Download test artifacts for detailed reports
- Run tests locally: `cd test && make test-issue-81`
- Review `test/INTEGRATION_TESTING.md` for troubleshooting

### Permission Errors

**Problem**: Cannot push workflow files
**Solution**:
- Push via GitHub UI (commit directly in browser)
- Or request repository admin to add workflows
- Or use personal access token with `workflow` scope

---

## Support

For issues with the PR validation pipeline:

1. Check workflow logs in GitHub Actions tab
2. Review test documentation: `test/INTEGRATION_TESTING.md`
3. Check `.github/workflows/README.md`
4. Open an issue with workflow logs and error details

---

## Summary

Once this PR validation pipeline is set up:

PASS **Every PR is automatically validated**
PASS **All Issue #81 requirements are tested**
PASS **Performance targets are enforced**
PASS **90% coverage requirement is checked**
PASS **Clear pass/fail status for merging**
PASS **Detailed test reports in PR comments**
PASS **No manual test execution needed**

This ensures code quality and prevents regressions!