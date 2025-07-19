# Panchangam Project Design Documentation

This directory contains comprehensive design documentation for the Panchangam project, including architecture specifications, system design decisions, and implementation guidelines.

## 📁 Documentation Structure

### **Architecture Documents**
- **[high-level-architecture.md](./high-level-architecture.md)** - Complete system architecture overview
  - Executive summary and project scope
  - System architecture diagrams and patterns
  - Component architecture and interfaces
  - Service architecture and future expansion plans
  - Data flow and processing pipelines
  - Technology stack and integration patterns
  - Implementation roadmap with 4-phase approach
  - Quality attributes and architectural principles

## 🎯 Design Principles

### **Core Architectural Principles**
1. **Hexagonal Architecture**: Clean separation of concerns with well-defined ports and adapters
2. **Domain-Driven Design**: Astronomical calculations as the core domain with rich business logic
3. **Microservices Architecture**: Independently deployable services with clear boundaries
4. **Context Awareness**: Request context propagation for distributed tracing and logging
5. **Observability First**: Comprehensive monitoring, tracing, and logging from the ground up

### **Quality Attributes**
- **Accuracy**: ±1 minute precision for astronomical calculations
- **Performance**: <100ms response time for standard calculations
- **Scalability**: Horizontal scaling with load balancing capabilities
- **Reliability**: 99.9% uptime with graceful degradation
- **Extensibility**: Plugin architecture for regional variations and custom calculations

## 🚀 Implementation Phases

### **Phase 1: Foundation** (Months 1-3)
- **Status**: ✅ Partially Complete
- **Focus**: Astronomical calculation foundation and observability infrastructure
- **Key Components**: Sunrise/sunset calculations, OpenTelemetry tracing, service logging

### **Phase 2: Core Panchangam** (Months 4-6)
- **Status**: 🔄 Planning
- **Focus**: Implementation of the five fundamental panchangam elements
- **Key Components**: Tithi, Nakshatra, Yoga, Karana, Vara calculations

### **Phase 3: Regional & Events** (Months 7-9)
- **Status**: 📋 Planned
- **Focus**: Regional variations and practical panchangam features
- **Key Components**: Regional systems, event generation, festival calculations

### **Phase 4: Advanced Features** (Months 10-12)
- **Status**: 📋 Planned
- **Focus**: Advanced astronomical features and user experience
- **Key Components**: Planetary integration, UI/UX, validation systems

## 🏗️ Architecture Patterns

### **Hexagonal Architecture Implementation**
```
Domain Layer (Core) → Application Layer → Infrastructure Layer
     ↑                      ↑                    ↑
Astronomical Logic    Service Orchestration    External APIs
```

### **Microservices Decomposition**
- **Panchangam Service**: Main API endpoint (✅ Implemented)
- **Astronomical Service**: Dedicated calculations (📋 Planned)
- **Regional Service**: Regional variations (📋 Planned)
- **Event Service**: Event generation (📋 Planned)
- **Validation Service**: Cross-validation (📋 Planned)

## 📊 Technology Decisions

### **Current Technology Stack**
- **Language**: Go 1.23 (chosen for performance and concurrency)
- **API Protocol**: gRPC with Protocol Buffers (type-safe, efficient)
- **Observability**: OpenTelemetry standard (vendor-neutral, comprehensive)
- **Authentication**: Custom AAA interceptors (flexible, secure)
- **Testing**: Comprehensive test suite with >95% coverage

### **Planned Technologies**
- **Ephemeris**: JPL DE440, Swiss Ephemeris (astronomical accuracy)
- **Storage**: Redis (cache), PostgreSQL (configuration)
- **Monitoring**: Prometheus (metrics), Grafana (dashboards)
- **Deployment**: Docker containers, Kubernetes orchestration

## 📖 Design Document Guidelines

### **Documentation Standards**
- **Format**: Markdown with consistent structure
- **Diagrams**: ASCII art for architecture diagrams
- **Code Examples**: Go interfaces and implementation patterns
- **Version Control**: Document version, creation date, last updated
- **Review Process**: Regular review after each phase completion

### **Review and Update Process**
1. **Phase Completion Reviews**: Update architecture after each phase
2. **Quarterly Reviews**: Assess architectural decisions and technical debt
3. **Annual Reviews**: Comprehensive architecture evaluation and planning
4. **Change Requests**: Document significant architectural changes

## 🔗 Related Documentation

- **[Implementation Plan](../../implementation_plan.md)** - Detailed task breakdown
- **[GitHub Issues](https://github.com/naren-m/panchangam/issues)** - Active development tracking
- **[GitHub Milestones](https://github.com/naren-m/panchangam/milestones)** - Phase-based project planning
- **[CLAUDE.md](../../CLAUDE.md)** - Development workflow guidelines

## 📞 Contact & Contributions

For questions about the architecture or suggestions for improvements:
- **GitHub Issues**: Technical questions and feature requests
- **Pull Requests**: Documentation improvements and architecture updates
- **Discussions**: Architecture discussions and design decisions

---

**Last Updated**: 2025-07-18  
**Document Version**: 1.0  
**Maintainer**: Panchangam Development Team