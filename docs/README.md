# Panchangam Project Documentation

## Overview

This documentation provides comprehensive coverage of the Panchangam project's algorithms, APIs, regional adaptations, and validation frameworks. The documentation is organized to support developers, astronomers, cultural experts, and system integrators.

## Documentation Structure

### 📊 Algorithms Documentation
**Location**: `docs/algorithms/`

#### [Tithi Calculation Algorithm](algorithms/TITHI_CALCULATION.md)
- **Mathematical Foundation**: Core formula and astronomical principles
- **Implementation Details**: Go code structure and calculation process
- **Tithi Classification**: Five-category system (Nanda, Bhadra, Jaya, Rikta, Purna)
- **Validation Methods**: Accuracy standards and cross-verification
- **Performance Optimization**: Caching and computational efficiency

**Key Topics Covered**:
- Angular difference calculations (Moon longitude - Sun longitude) ÷ 12°
- Variable tithi durations (19h 59m to 26h 47m)
- Paksha classification (Shukla/Krishna)
- Swiss Ephemeris integration
- OpenTelemetry observability

### 🔌 API Documentation
**Location**: `docs/api/`

#### [Ephemeris Integration API](api/EPHEMERIS_API.md)
- **Provider Architecture**: Swiss Ephemeris and JPL integration
- **Data Structures**: Planetary positions, solar/lunar positions
- **Manager Interface**: Multi-provider orchestration with fallback
- **Utility Functions**: Time conversion and astronomical calculations
- **Performance Features**: Caching, batch operations, health monitoring

**Key Components**:
- `EphemerisProvider` interface with multiple implementations
- `PlanetaryPositions`, `SolarPosition`, `LunarPosition` data structures
- Cross-provider validation and error handling
- Real-time health monitoring and observability integration

### 🌍 Regional Documentation
**Location**: `docs/regional/`

#### [Regional Calculations and Cultural Adaptations](regional/REGIONAL_CALCULATIONS.md)
- **Plugin Architecture**: Modular regional logic and cultural customization
- **Regional Classifications**: North vs South India calendar systems
- **Calculation Methods**: Drik Ganita vs Vakya Ganita approaches
- **Localization Framework**: Multi-language support with cultural sensitivity
- **Festival Systems**: Regional event calculations and muhurta timing

**Regional Implementations**:
- **Tamil Nadu**: Amanta system, Naazhikai time units, Tamil festivals
- **Kerala**: Malayalam era, solar calendar, Drik preference
- **Bengal**: Dual system support, Durga Puja calculations
- **Gujarat/Maharashtra**: Vikrama Samvat, business muhurtas

### ✅ Validation Documentation
**Location**: `docs/validation/`

#### [Validation Framework](validation/VALIDATION_FRAMEWORK.md)
- **Multi-Layer Strategy**: 8-tier validation system
- **Mathematical Accuracy**: Astronomical formula verification
- **Reference Source Validation**: Cross-checking with established panchangams
- **Regional Compliance**: Cultural and traditional accuracy
- **Performance Testing**: Load testing and response time monitoring

**Validation Components**:
- Cross-provider ephemeris verification
- Historical consistency validation
- Edge case and boundary condition testing
- Continuous integration and quality gates

### 📐 Design Documentation
**Location**: `docs/design/`

#### [High-Level Architecture](design/high-level-architecture.md)
- System architecture overview
- Component interactions and data flow
- Scalability and performance considerations

#### [Implementation Phase Diagrams](design/implementation-phase-diagrams.md)
- Development workflow and milestone planning
- Component integration strategies

## Quick Reference Guides

### 🚀 Getting Started

#### For Developers
1. **Setup**: Clone repository and install dependencies
2. **Core Concepts**: Read [Tithi Calculation](algorithms/TITHI_CALCULATION.md) for mathematical foundation
3. **API Integration**: Use [Ephemeris API](api/EPHEMERIS_API.md) for astronomical calculations
4. **Testing**: Follow [Validation Framework](validation/VALIDATION_FRAMEWORK.md) for quality assurance

#### For Cultural Experts
1. **Regional Systems**: Review [Regional Calculations](regional/REGIONAL_CALCULATIONS.md)
2. **Localization**: Understand multi-language and cultural adaptation frameworks
3. **Festival Calculations**: Explore regional event and muhurta systems
4. **Validation**: Contribute to cultural accuracy verification

#### For System Integrators
1. **API Reference**: Start with [Ephemeris API](api/EPHEMERIS_API.md)
2. **Performance**: Review performance optimization strategies
3. **Regional Support**: Configure regional plugins for target markets
4. **Monitoring**: Implement observability and health checking

### 📊 Key Metrics and Standards

#### Accuracy Benchmarks
- **Tithi Calculations**: ±1 minute precision
- **Ephemeris Data**: ±0.001 arcsecond accuracy
- **Sunrise/Sunset**: Sub-minute accuracy globally
- **Festival Dates**: 90%+ agreement with established sources

#### Performance Targets
- **API Response Time**: <5 seconds for single date calculations
- **Batch Processing**: <30 seconds for monthly data
- **Memory Usage**: <500MB under normal load
- **Availability**: 99.9% uptime for production systems

#### Quality Gates
- **Mathematical Validation**: 95%+ pass rate required
- **Regional Compliance**: 80%+ cultural accuracy
- **Performance Tests**: 90%+ under load conditions
- **Integration Tests**: 100% critical path coverage

## Technology Stack

### Core Technologies
- **Language**: Go 1.21+
- **Astronomical Engine**: Swiss Ephemeris, JPL integration
- **Observability**: OpenTelemetry tracing and metrics
- **Testing**: Go testing framework with validation suites
- **Documentation**: Markdown with technical specifications

### External Dependencies
- Swiss Ephemeris library for planetary calculations
- JPL planetary ephemeris for cross-validation
- Time zone databases for global location support
- Regional panchangam sources for validation

## Contributing Guidelines

### Documentation Standards
1. **Technical Accuracy**: All algorithms must be mathematically verified
2. **Cultural Sensitivity**: Regional content reviewed by cultural experts
3. **Code Examples**: Include working Go code snippets
4. **Validation**: Document all accuracy claims with test results
5. **Updates**: Maintain version compatibility and migration guides

### Review Process
1. **Technical Review**: Algorithm accuracy and implementation quality
2. **Cultural Review**: Regional accuracy and cultural appropriateness  
3. **Performance Review**: System impact and optimization opportunities
4. **Documentation Review**: Clarity, completeness, and maintenance

## Support and Resources

### Issue Tracking
- **GitHub Issues**: Technical bugs and feature requests
- **Cultural Feedback**: Regional accuracy and cultural concerns
- **Performance Issues**: System performance and scalability
- **Documentation**: Clarity and completeness improvements

### Community Resources
- **Traditional Panchangam Makers**: Cultural accuracy validation
- **Academic Institutions**: Astronomical algorithm verification
- **Regional Communities**: Cultural adaptation and localization
- **Developer Community**: Technical implementation and optimization

### Reference Materials
- Ancient texts: Surya Siddhanta, Siddhanta Shiromani
- Modern sources: Swiss Ephemeris documentation, JPL planetary data
- Cultural sources: Regional panchangam publications, festival calendars
- Academic papers: Hindu calendar research, astronomical algorithms

## Roadmap and Future Development

### Planned Enhancements
1. **Extended Regional Support**: Additional cultural variations
2. **AI-Enhanced Localization**: Improved translation and cultural adaptation
3. **Real-time Validation**: Continuous accuracy monitoring
4. **Performance Optimization**: Advanced caching and computation
5. **Mobile Integration**: Lightweight implementations for mobile platforms

### Research Areas
- Historical calendar system evolution
- Regional variation documentation
- Performance optimization techniques
- Cultural sensitivity frameworks
- Astronomical accuracy improvements

---

**Last Updated**: July 2025  
**Documentation Version**: 1.0.0  
**Project Version**: Compatible with Panchangam v1.x  
**Maintainers**: Panchangam Development Team

For questions or contributions, please refer to the project's GitHub repository and follow the established contribution guidelines.