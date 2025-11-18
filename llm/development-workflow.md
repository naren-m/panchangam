# Development Workflow

This document outlines the development workflow, Git practices, and collaboration guidelines for the Panchangam project.

## Branching Strategy

### Branch Naming Convention

All branches must follow this pattern:
```
<type>/<issue-number>-<brief-description>

Examples:
feature/123-add-muhurta-calculator
fix/456-tithi-calculation-bug
refactor/789-optimize-ephemeris-cache
docs/101-update-api-documentation
```

**Branch Types**:
- `feature/` - New features or enhancements
- `fix/` - Bug fixes
- `refactor/` - Code refactoring without feature changes
- `docs/` - Documentation updates
- `test/` - Test additions or modifications
- `chore/` - Maintenance tasks (dependencies, tooling)

### Branch Lifecycle

1. **Create Branch from Main**
   ```bash
   git checkout main
   git pull origin main
   git checkout -b feature/123-add-nakshatra-visualization
   ```

2. **Regular Syncing**
   ```bash
   # Sync with main regularly
   git fetch origin
   git rebase origin/main

   # Or merge if preferred
   git merge origin/main
   ```

3. **Keep Branches Short-lived**
   - Target: Complete work within 2-3 days
   - Maximum: 1 week before merging or closing
   - Break large features into smaller branches

4. **Clean Up After Merge**
   ```bash
   # After PR is merged
   git checkout main
   git pull origin main
   git branch -d feature/123-add-nakshatra-visualization
   ```

## Commit Guidelines

### Commit Message Format

Follow the Conventional Commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Example**:
```
feat(astronomy): add Nakshatra calculation with Swiss Ephemeris

- Implement CalculateNakshatra function
- Add support for all 27 nakshatras
- Include pada (quarter) calculations
- Add comprehensive test coverage

Closes #123
```

### Commit Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, missing semicolons)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements

### Commit Best Practices

**DO**:
- Write clear, descriptive commit messages
- Keep commits focused and atomic
- Reference issue numbers in commit messages
- Test before committing
- Commit working code

**DON'T**:
- Commit broken code
- Use generic messages ("fix bug", "update code")
- Commit large, unrelated changes together
- Commit secrets, credentials, or sensitive data
- Commit generated files (unless necessary)

### Atomic Commits

```bash
# Good: Separate logical changes
git add astronomy/nakshatra.go astronomy/nakshatra_test.go
git commit -m "feat(astronomy): add Nakshatra calculation"

git add ui/src/components/NakshatraDisplay.tsx
git commit -m "feat(ui): add Nakshatra display component"

# Bad: Mixing unrelated changes
git add .
git commit -m "various updates"
```

## Pull Request Process

### Before Creating PR

**Checklist**:
- [ ] All tests pass locally (`make test` and `npm test`)
- [ ] Code coverage ≥ 90%
- [ ] Code follows style guidelines (linters pass)
- [ ] No console.log or debug statements
- [ ] Documentation updated if needed
- [ ] Branch is up-to-date with main
- [ ] Commit history is clean

### Creating a Pull Request

1. **Push Your Branch**
   ```bash
   git push -u origin feature/123-add-nakshatra-visualization
   ```

2. **Create PR on GitHub**
   - Use descriptive title matching commit convention
   - Include issue number in title and description
   - Fill out PR template completely

3. **PR Title Format**
   ```
   feat: Add Nakshatra visualization (#123)
   fix: Correct Tithi calculation for edge cases (#456)
   ```

### PR Description Template

```markdown
## Description
Brief description of changes and motivation.

## Issue Reference
Closes #123

## Type of Change
- [ ] Bug fix (non-breaking change fixing an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Breaking change (fix or feature causing existing functionality to change)
- [ ] Documentation update

## Changes Made
- Detailed list of changes
- What was added
- What was modified
- What was removed

## Testing Performed
- Unit tests added/updated
- Integration tests run
- Manual testing performed
- Edge cases tested

## Screenshots (if applicable)
[Add screenshots for UI changes]

## Code Coverage
- Overall coverage: 92.5%
- New code coverage: 95.0%

## Checklist
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my code
- [ ] I have commented complex areas of code
- [ ] I have updated documentation as needed
- [ ] My changes generate no new warnings
- [ ] I have added tests with ≥90% coverage
- [ ] All tests pass locally
- [ ] Dependent changes have been merged

## Additional Notes
Any additional context or considerations.
```

### PR Review Process

**As Author**:
1. Respond to feedback promptly
2. Make requested changes
3. Mark conversations as resolved when addressed
4. Keep PR updated with main branch
5. Be open to suggestions

**As Reviewer**:
1. Review within 24-48 hours
2. Check for:
   - Code correctness
   - Test coverage
   - Style compliance
   - Documentation
   - Performance implications
   - Security concerns
3. Provide constructive feedback
4. Approve when satisfied

### PR Merge Requirements

Before merging:
- [ ] At least one approval (two for breaking changes)
- [ ] All CI checks pass
- [ ] No merge conflicts
- [ ] Code coverage meets 90% threshold
- [ ] Documentation updated
- [ ] No requested changes pending

### Merge Strategy

**Preferred**: Squash and Merge
- Keeps main branch history clean
- Combines all commits into one
- Use for feature branches

**Alternative**: Rebase and Merge
- Preserves commit history
- Use for important features with logical commit progression

**Avoid**: Merge Commit
- Creates merge bubbles
- Makes history harder to follow

## Code Review Guidelines

### What to Look For

**Functionality**:
- Does the code do what it's supposed to?
- Are edge cases handled?
- Is error handling appropriate?

**Code Quality**:
- Is the code readable and maintainable?
- Are names descriptive?
- Is complexity minimized?
- Are functions/methods appropriately sized?

**Testing**:
- Are there sufficient tests?
- Do tests cover edge cases?
- Is coverage ≥90%?
- Are tests meaningful (not just for coverage)?

**Performance**:
- Are there obvious performance issues?
- Is caching used appropriately?
- Are database queries optimized?

**Security**:
- Input validation present?
- No hardcoded secrets?
- Proper error handling (no information leakage)?
- SQL injection prevention?

**Documentation**:
- Are complex parts commented?
- Is public API documented?
- Are README/docs updated?

### Providing Feedback

**Good Feedback**:
```
❌ "This is wrong"
✅ "Consider using CalculateNakshatra() instead of inline calculation here to improve maintainability and reusability."

❌ "Needs tests"
✅ "Could you add a test case for when moonLongitude > 360°? This edge case isn't currently covered."

❌ "Bad naming"
✅ "The variable name 'x' isn't very descriptive. Consider renaming to 'tithiProgress' to make the code more readable."
```

**Categories**:
- **Required**: Must be addressed before merge
- **Suggestion**: Nice to have, not blocking
- **Question**: Seeking clarification
- **Nit**: Minor style/formatting issue

## Development Best Practices

### Local Development Setup

1. **Clone Repository**
   ```bash
   git clone https://github.com/naren-m/panchangam.git
   cd panchangam
   ```

2. **Install Dependencies**
   ```bash
   # Backend
   go mod download

   # Frontend
   cd ui
   npm install
   cd ..
   ```

3. **Run Tests**
   ```bash
   # Backend tests
   make test

   # Frontend tests
   cd ui && npm test
   ```

4. **Start Development Servers**
   ```bash
   # Terminal 1: Backend
   make run_server

   # Terminal 2: Frontend
   cd ui && npm run dev
   ```

### Before Starting Work

1. **Check for existing issues**
   - Search for duplicates
   - Comment if you plan to work on an issue

2. **Create or assign issue**
   - Describe what you'll implement
   - Get feedback on approach
   - Assign to yourself

3. **Create branch**
   - Branch from latest main
   - Use proper naming convention
   - Include issue number

### During Development

**Regular Tasks**:
- Commit frequently (atomic commits)
- Write tests alongside code
- Run tests locally before pushing
- Sync with main regularly
- Update documentation as you go

**Communication**:
- Comment on issue with progress updates
- Ask questions early if blocked
- Share WIP for early feedback

### Pre-Push Checklist

```bash
# 1. Run all tests
make test
cd ui && npm test

# 2. Check coverage
make test-coverage
cd ui && npm run test:coverage

# 3. Run linters
cd ui && npm run lint

# 4. Format code
gofmt -w .
cd ui && npm run lint -- --fix

# 5. Ensure no debug code
git diff | grep -i "console.log\|fmt.Println\|debugger"

# 6. Review your changes
git diff main

# 7. Push
git push origin feature/123-your-feature
```

## Testing Strategy

### Test-Driven Development (TDD)

**Recommended Approach**:
1. Write failing test first
2. Implement minimal code to pass
3. Refactor while keeping tests green
4. Repeat

**Example Workflow**:
```go
// 1. Write failing test
func TestCalculateNakshatra(t *testing.T) {
    nakshatra := CalculateNakshatra(0.0)
    assert.Equal(t, 1, nakshatra) // Fails - function doesn't exist
}

// 2. Implement minimal solution
func CalculateNakshatra(moonLong float64) int {
    return int(moonLong / 13.333) + 1
}

// 3. Test passes, now refactor if needed
func CalculateNakshatra(moonLong float64) int {
    const NakshatraSpan = 13.333333
    return int(moonLong / NakshatraSpan) + 1
}

// 4. Add more test cases
func TestCalculateNakshatra(t *testing.T) {
    tests := []struct {
        name     string
        moonLong float64
        want     int
    }{
        {"Ashwini start", 0.0, 1},
        {"Bharani start", 13.333, 2},
        {"Revati", 350.0, 27},
    }
    // ...
}
```

### Testing Levels

**Unit Tests** (90% coverage target):
- Test individual functions
- Mock dependencies
- Fast execution
- Run on every commit

**Integration Tests**:
- Test component interactions
- Use real dependencies when possible
- Slower execution
- Run before PR

**E2E Tests** (Frontend):
- Test complete user workflows
- Use Playwright
- Slowest execution
- Run in CI pipeline

## Continuous Integration

### CI Pipeline

**Triggered On**:
- Push to any branch
- Pull request creation
- Pull request update

**Pipeline Steps**:
1. Checkout code
2. Setup environment (Go, Node)
3. Install dependencies
4. Run linters
5. Run unit tests
6. Check code coverage (≥90%)
7. Run integration tests
8. Build artifacts
9. Report results

**Pipeline Configuration**: `.github/workflows/`

### Pre-commit Hooks (Recommended)

Install local pre-commit hooks:

```bash
# .git/hooks/pre-commit
#!/bin/bash

echo "Running pre-commit checks..."

# Run Go tests
echo "Running Go tests..."
make test || exit 1

# Run frontend tests
echo "Running frontend tests..."
cd ui && npm test -- --run || exit 1

# Check coverage
echo "Checking coverage..."
cd .. && make test-coverage || exit 1

echo "All checks passed!"
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

## Issue Management

### Creating Issues

**Bug Report Template**:
```markdown
## Description
Clear description of the bug

## Steps to Reproduce
1. Step 1
2. Step 2
3. Step 3

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS:
- Browser (if applicable):
- Version:

## Screenshots
[If applicable]

## Additional Context
Any other relevant information
```

**Feature Request Template**:
```markdown
## Problem/Motivation
Why is this feature needed?

## Proposed Solution
How should this work?

## Alternatives Considered
Other approaches you've thought about

## Implementation Ideas
Technical approach if you have ideas

## Additional Context
Any other relevant information
```

### Issue Labels

- `bug`: Something isn't working
- `enhancement`: New feature or request
- `documentation`: Documentation improvements
- `good first issue`: Good for newcomers
- `help wanted`: Extra attention needed
- `priority: high`: High priority
- `priority: medium`: Medium priority
- `priority: low`: Low priority
- `status: blocked`: Blocked by other work
- `status: in-progress`: Currently being worked on

## Release Process

### Versioning

Follow Semantic Versioning (SemVer):
```
MAJOR.MINOR.PATCH

Example: 1.2.3

MAJOR: Breaking changes
MINOR: New features (backward compatible)
PATCH: Bug fixes (backward compatible)
```

### Release Checklist

1. Update version numbers
2. Update CHANGELOG.md
3. Run full test suite
4. Create release branch
5. Tag release
6. Build artifacts
7. Deploy to staging
8. Test staging
9. Deploy to production
10. Announce release

## Communication

### Channels

- **GitHub Issues**: Bug reports, feature requests
- **Pull Requests**: Code discussions
- **Documentation**: `docs/` folder
- **README**: Project overview and setup

### Best Practices

- Be respectful and professional
- Provide context in discussions
- Use clear, concise language
- Include code examples when relevant
- Reference issues/PRs with #number
- Tag people when input needed (@username)

## Resources

- [Git Best Practices](https://git-scm.com/book/en/v2)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Code Review Guidelines](https://google.github.io/eng-practices/review/)
- [Semantic Versioning](https://semver.org/)

## Summary

**Key Principles**:
1. Branch for every task with issue numbers
2. Write tests achieving 90%+ coverage
3. Create focused, well-documented PRs
4. Review code thoroughly
5. Keep main branch clean and deployable
6. Communicate effectively
7. Respect the process and team

Following these guidelines ensures a smooth, collaborative development experience and maintains high code quality throughout the project.
