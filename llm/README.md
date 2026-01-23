# LLM Context Documentation

This folder contains critical documentation to help AI/LLM agents understand the Panchangam project's coding standards, architecture, and testing requirements. These documents ensure consistency and quality when AI agents assist with development tasks.

## Purpose

This documentation helps LLM agents:
- Understand the project structure and architecture
- Follow established coding standards and conventions
- Maintain 90% code coverage requirement
- Work with the specific tech stack (Go backend, React/TypeScript frontend)
- Understand the domain-specific concepts of Panchangam astronomy

## Documentation Structure

- **[coding-standards.md](./coding-standards.md)** - Coding conventions and best practices for Go and TypeScript
- **[testing-guidelines.md](./testing-guidelines.md)** - Testing requirements, frameworks, and coverage standards
- **[project-architecture.md](./project-architecture.md)** - System architecture, component structure, and design patterns
- **[domain-context.md](./domain-context.md)** - Panchangam-specific domain knowledge and astronomical concepts
- **[development-workflow.md](./development-workflow.md)** - Git workflow, branching strategy, and PR requirements

## Quick Reference

### Tech Stack
- **Backend**: Go 1.23.0, gRPC, Redis, OpenTelemetry
- **Frontend**: React 18, TypeScript, Vite, TailwindCSS
- **Testing**: Go testing + testify, Vitest + Testing Library
- **Build**: Makefile

### Critical Requirements
- **Code Coverage**: Minimum 90% for all PRs
- **Branching**: Always branch out, include issue numbers
- **Testing**: Both unit and integration tests required
- **Documentation**: Keep API docs and user guides updated

## For LLM Agents

When working on this project:
1. **Always read relevant documentation** from this folder before starting tasks
2. **Follow coding standards** strictly for consistency
3. **Write tests first** or alongside implementation
4. **Verify coverage** meets 90% threshold before submitting
5. **Include issue numbers** in branch names and PR descriptions
6. **Update documentation** when adding new features or changing architecture

## Maintenance

This documentation should be updated when:
- Major architectural changes occur
- New coding standards are adopted
- Testing requirements change
- New technologies are added to the stack
- Domain knowledge needs clarification
