# Research: Adding rundeck/rundeck Provider to Terraformer Community Providers

## Overview

This document outlines the requirements and steps for integrating the `rundeck/rundeck` provider into terraformer's community providers, based on analysis of existing terraformer patterns and the rundeck provider documentation.

## Rundeck Provider Information

### Provider Details
- **Registry**: `rundeck/rundeck`
- **Latest Version**: 0.5.2 (as of research date)
- **Source**: https://github.com/rundeck/terraform-provider-rundeck
- **Tier**: Partner
- **Downloads**: 492,201+

### Supported Resources
1. `rundeck_acl_policy` - ACL policies for authorization
2. `rundeck_job` - Job definitions and configurations
3. `rundeck_password` - Password key storage
4. `rundeck_private_key` - Private key storage
5. `rundeck_project` - Project management
6. `rundeck_public_key` - Public key storage

## Implementation Requirements

### 1. File Structure Requirements

Based on terraformer patterns, the following files need to be created:

```
providers/rundeck/
├── rundeck_provider.go     # Main provider implementation
├── rundeck_service.go      # Base service struct
├── project.go              # Project resource generator
├── job.go                  # Job resource generator
├── acl_policy.go          # ACL policy generator
├── password.go            # Password key storage generator
├── private_key.go         # Private key storage generator
├── public_key.go          # Public key storage generator
└── helpers.go             # Common utilities (optional)
```

### 2. CLI Integration Files

```
cmd/
├── provider_cmd_rundeck.go # CLI command implementation
└── root.go                 # Updated to register rundeck provider
```

### 3. Documentation

```
docs/rundeck.md            # Provider documentation
```

### 4. Testing Structure (Optional)

```
tests/rundeck/             # Integration tests
├── README.md              # Testing instructions
└── provider.tf           # Test configuration
```

## Technical Implementation Checklist

### Core Provider Implementation

- [ ] **Provider Structure** (`providers/rundeck/rundeck_provider.go`)
  - [ ] Implement `terraformutils.Provider` interface
  - [ ] Define provider configuration fields (URL, API key, username, password)
  - [ ] Implement `Init()`, `GetName()`, `GetConfig()`, `GetBasicConfig()` methods
  - [ ] Implement `InitService()` and `GetSupportedService()` methods
  - [ ] Define resource connections mapping

- [ ] **Service Base** (`providers/rundeck/rundeck_service.go`)
  - [ ] Create RundeckService struct extending `terraformutils.Service`

### Authentication & Configuration

Based on rundeck provider patterns, authentication typically involves:

- [ ] **Environment Variables**
  - [ ] `RUNDECK_URL` - Rundeck server URL
  - [ ] `RUNDECK_API_VERSION` - API version (optional)
  - [ ] `RUNDECK_USERNAME` - Username for authentication  
  - [ ] `RUNDECK_PASSWORD` - Password for authentication
  - [ ] `RUNDECK_TOKEN` - API token (alternative to username/password)
  - [ ] `RUNDECK_INSECURE_SSL` - Skip SSL verification (boolean)

- [ ] **Provider Configuration**
  - [ ] URL field for Rundeck server endpoint
  - [ ] Authentication credentials (token or username/password)
  - [ ] SSL/TLS configuration options
  - [ ] API version specification

### Resource Generators

For each supported resource type, implement:

- [ ] **Project Generator** (`providers/rundeck/project.go`)
  - [ ] `ProjectGenerator` struct with `InitResources()` method
  - [ ] API calls to list projects
  - [ ] Resource mapping to `rundeck_project` terraform resources

- [ ] **Job Generator** (`providers/rundeck/job.go`)
  - [ ] `JobGenerator` struct with `InitResources()` method
  - [ ] API calls to list jobs by project
  - [ ] Resource mapping to `rundeck_job` terraform resources

- [ ] **ACL Policy Generator** (`providers/rundeck/acl_policy.go`)
  - [ ] `AclPolicyGenerator` struct with `InitResources()` method
  - [ ] API calls to list ACL policies
  - [ ] Resource mapping to `rundeck_acl_policy` terraform resources

- [ ] **Key Storage Generators**
  - [ ] `PasswordGenerator` for `rundeck_password` resources
  - [ ] `PrivateKeyGenerator` for `rundeck_private_key` resources
  - [ ] `PublicKeyGenerator` for `rundeck_public_key` resources
  - [ ] API calls to enumerate key storage items

### CLI Integration

- [ ] **Command Implementation** (`cmd/provider_cmd_rundeck.go`)
  - [ ] Create `newCmdRundeckImporter()` function
  - [ ] Environment variable parsing for authentication
  - [ ] Command flags for filtering and options
  - [ ] Support for listing resources (`terraformer import rundeck list`)

- [ ] **Root Registration** (`cmd/root.go`)
  - [ ] Add `newCmdRundeckImporter` to `providerImporterSubcommands()`
  - [ ] Add `newRundeckProvider` to `providerGenerators()` map
  - [ ] Place in appropriate category (Community providers section)

## Authentication Investigation

### Required Research
1. **Rundeck API Documentation**: Study authentication methods
2. **Provider Source Code**: Examine authentication implementation
3. **API Endpoints**: Identify required API calls for each resource type
4. **Permissions**: Understand required permissions for read-only operations

### Expected Authentication Methods
- HTTP Basic Authentication (username/password)
- API Token authentication
- SSL/TLS configuration options

## Potential Blockers & Considerations

### Technical Blockers
1. **API Limitations**: Some resources might not be listable via API
2. **Authentication Complexity**: Multiple auth methods may complicate implementation
3. **API Rate Limiting**: May need to implement rate limiting or retry logic
4. **Resource Dependencies**: Some resources may depend on others (jobs depend on projects)

### Implementation Considerations
1. **Resource Filtering**: Jobs should be organized by project
2. **Large Environments**: May need pagination for large Rundeck installations
3. **Resource Naming**: Ensure generated resource names follow terraform conventions
4. **State Management**: Handle resources that might be modified outside terraform

### Missing Information
1. **Rundeck API Documentation**: Need detailed API endpoint documentation
2. **Provider Source Code**: Access to implementation details for proper mapping
3. **Testing Environment**: Need Rundeck instance for development and testing
4. **Resource Relationships**: Understanding dependencies between resource types

## Next Steps for Implementation

### Phase 1: Setup & Research
1. **Deep Dive into Rundeck Provider Source**
   - Clone `https://github.com/rundeck/terraform-provider-rundeck`
   - Study authentication implementation
   - Understand resource schemas and API mappings

2. **API Documentation Review**
   - Study Rundeck REST API documentation
   - Identify required endpoints for each resource type
   - Understand authentication mechanisms

### Phase 2: Core Implementation
1. **Provider Foundation**
   - Implement `rundeck_provider.go` with authentication
   - Create `rundeck_service.go` base structure
   - Add CLI integration in `cmd/provider_cmd_rundeck.go`

2. **Resource Generators**
   - Start with `project.go` as foundation (projects are typically the base entity)
   - Implement `job.go` with project-based filtering
   - Add key storage generators (`password.go`, `private_key.go`, `public_key.go`)
   - Implement `acl_policy.go`

### Phase 3: Testing & Documentation
1. **Testing Setup**
   - Create test environment setup instructions
   - Implement integration tests
   - Test with various Rundeck configurations

2. **Documentation**
   - Create `docs/rundeck.md` with usage examples
   - Update main README.md with rundeck provider
   - Document authentication and configuration options

### Phase 4: Integration & Review
1. **Code Review**
   - Ensure compliance with terraformer patterns
   - Performance optimization
   - Error handling improvements

2. **Community Integration**
   - Submit pull request
   - Address feedback and reviews
   - Update documentation based on feedback

## Resource Relationships

Based on rundeck concepts:
- **Projects** are the top-level containers
- **Jobs** belong to specific projects
- **ACL Policies** control access to projects and jobs
- **Key Storage** items are global but can be project-specific

## Estimated Implementation Effort

- **Phase 1 (Research)**: 1-2 days
- **Phase 2 (Core Implementation)**: 3-4 days
- **Phase 3 (Testing & Documentation)**: 2-3 days
- **Phase 4 (Review & Integration)**: 1-2 days

**Total Estimated Effort**: 7-11 days

## Implementation Checklist

### Pre-Implementation Research
- [ ] Study rundeck provider source code for authentication patterns
- [ ] Review Rundeck API documentation for required endpoints
- [ ] Understand resource schemas and terraform resource definitions
- [ ] Set up local Rundeck environment for testing

### Code Implementation

#### Core Provider Files
- [ ] `providers/rundeck/rundeck_provider.go`
  - [ ] RundeckProvider struct with required fields
  - [ ] Authentication configuration (url, token/credentials)
  - [ ] Init() method for argument parsing
  - [ ] GetName() returning "rundeck"
  - [ ] GetConfig() and GetBasicConfig() with cty.Value
  - [ ] InitService() method
  - [ ] GetSupportedService() mapping to generators
  - [ ] GetResourceConnections() for dependencies

- [ ] `providers/rundeck/rundeck_service.go`
  - [ ] RundeckService struct extending terraformutils.Service

#### Resource Generators
- [ ] `providers/rundeck/project.go`
  - [ ] ProjectGenerator struct
  - [ ] InitResources() method with API calls
  - [ ] createResources() helper method
  - [ ] Resource mapping to terraform state

- [ ] `providers/rundeck/job.go`
  - [ ] JobGenerator struct
  - [ ] Project-aware job listing
  - [ ] Job configuration mapping

- [ ] `providers/rundeck/acl_policy.go`
  - [ ] AclPolicyGenerator struct
  - [ ] Policy enumeration and mapping

- [ ] `providers/rundeck/password.go` (key storage)
- [ ] `providers/rundeck/private_key.go` (key storage)  
- [ ] `providers/rundeck/public_key.go` (key storage)

#### CLI Integration
- [ ] `cmd/provider_cmd_rundeck.go`
  - [ ] newCmdRundeckImporter() function
  - [ ] Environment variable handling
  - [ ] Command line flags
  - [ ] Error handling and validation
  - [ ] List command support

- [ ] `cmd/root.go` updates
  - [ ] Add to providerImporterSubcommands()
  - [ ] Add to providerGenerators() map
  - [ ] Import statement for rundeck provider

#### Documentation & Testing
- [ ] `docs/rundeck.md`
  - [ ] Installation instructions
  - [ ] Authentication setup
  - [ ] Usage examples
  - [ ] Resource listing

- [ ] `tests/rundeck/` (optional)
  - [ ] README.md with setup instructions
  - [ ] Test configuration files
  - [ ] Integration test scripts

### Validation & Testing
- [ ] Code compilation and build verification
- [ ] Unit tests for resource generators
- [ ] Integration testing with live Rundeck instance
- [ ] Documentation review and examples
- [ ] Performance testing with large Rundeck environments

### Code Review & Integration
- [ ] Code style compliance (gofmt, linting)
- [ ] Error handling completeness
- [ ] Documentation accuracy
- [ ] Community feedback incorporation
- [ ] Final pull request submission

## Success Criteria

The integration will be considered successful when:

1. **Functional Requirements**
   - [ ] All 6 rundeck resource types can be imported
   - [ ] Authentication works with both API tokens and credentials
   - [ ] Resources are properly mapped to terraform state
   - [ ] Generated HCL is valid and applies without errors

2. **Technical Requirements**
   - [ ] Code follows terraformer patterns and conventions
   - [ ] Proper error handling and logging
   - [ ] No breaking changes to existing functionality
   - [ ] Performance is acceptable for typical use cases

3. **Documentation Requirements**
   - [ ] Clear setup and usage instructions
   - [ ] Working examples for all authentication methods
   - [ ] Integration with main terraformer documentation

4. **Testing Requirements**
   - [ ] Integration tests pass
   - [ ] Code builds successfully across platforms
   - [ ] No regressions in existing providers

## Risk Assessment

### Low Risk
- Basic resource enumeration (projects, jobs)
- Standard authentication patterns
- CLI integration following existing patterns

### Medium Risk  
- Key storage resource handling (may require special permissions)
- Resource dependency management (jobs → projects)
- Large environment performance

### High Risk
- Complex ACL policy mapping
- Non-standard authentication requirements
- API rate limiting or availability issues

## References

1. [Terraform Registry: rundeck/rundeck](https://registry.terraform.io/providers/rundeck/rundeck/latest/docs)
2. [Terraformer README - High-Level steps to add new provider](../README.md)
3. [Existing Community Provider Examples](../providers/keycloak/)
4. [Rundeck Provider Source](https://github.com/rundeck/terraform-provider-rundeck)
5. [Rundeck API Documentation](https://docs.rundeck.com/docs/api/)
6. [Terraformer Contributing Guide](../CONTRIBUTING.md)