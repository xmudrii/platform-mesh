# ADR: Account Model Implementation in Platform Mesh using KCP

## Status: Proposed

## Deciders:
- TBD
- ...

## Date: TBD

## Technical Story:
Evaluate implementation options for the account model in Platform Mesh using KCP to create a flexible, scalable, and interoperable system for managing accounts, services and workload instances.

## Context and Problem Statement

We need to implement an account model for Platform Mesh using KCP. The account model should be simple, scalable, and not locked to regions. It should support service and workload management, distinguish between services and applications, and allow for the decoupling of orthogonal aspects such as quotas, service validation, and access control. The core question is how to implement this account model effectively using KCP's workspace concepts and the Kubernetes Resource Model (KRM).

## Decision Drivers

1. Need for a simple and scalable account model
2. Requirement to support hierarchical account structures
3. Desire to leverage KCP's workspace capabilities efficiently
4. Need for clear distinction between services and workloads
5. Requirement for extensibility to accommodate various provider-specific needs
6. Desire for orthogonal aspects to be decoupled from the core account model
7. Requirement for consistent and atomic account creation
8. Need for predictable resource hierarchy and relationship management
9. Requirement for cross-account operations and visibility
10. Desire for audit capability and compliance support

## Considered Options

### Option 1: Custom Resource Definition (CRD) for Account Model

This option involves creating a new CRD in KCP to define the account model, with accounts managed as custom resources.

Pros:
- Native Kubernetes approach, easily integrable with KCP
- Allows for declarative management of accounts
- Can be extended using additional fields or annotations
- Facilitates versioning and API evolution

Cons:
- May require frequent updates to the CRD as new requirements emerge
- Could become complex if trying to accommodate all provider-specific needs in one CRD
- Potential performance impact with a large number of custom resources

### Option 2: KCP Workspace as the Core Account Representation

This option uses KCP workspaces as the primary representation of accounts, with additional metadata stored in workspace annotations or labels.

Pros:
- Leverages KCP's existing hierarchical workspace model
- Provides built-in isolation and access control mechanisms
- Allows for easy implementation of hierarchical structures
- Facilitates management of resources within account context

Cons:
- Limited flexibility in storing complex account data
- May require additional controllers to manage account-specific operations
- Could lead to overloading of workspace concepts

### Option 3: KCP as a Service with Encapsulated Account Model

This option positions KCP as a service that can be consumed by other teams, with the account model built as a layer on top.

Pros:
- Clear separation of concerns between KCP and account management
- Flexibility to evolve account model independently
- Can leverage existing KCP features while maintaining abstraction

Cons:
- Additional complexity in managing service layer
- Requires dedicated team for service maintenance
- Potential performance overhead from additional abstraction layer

### Option 4: Initializer-Based Account-Workspace Binding

This option implements accounts as CRDs with a strict 1:1 mapping to KCP workspaces, using initializers for atomic creation and setup.

Pros:
- Atomic account creation through initializer pattern
- Clear, predictable 1:1 relationship with workspaces
- Built-in validation and dependency management
- Supports hierarchical account structures
- Facilitates staged resource creation
- Clear status tracking through conditions

Cons:
- More complex initialization flow
- Potential for stuck initializers requiring manual intervention
- Need for careful timeout and retry handling

## Decision Outcome (Proposed)

Recommend Option 4: Initializer-Based Account-Workspace Binding for the following reasons:

### Positive Consequences

1. Atomic Operations:
- Guaranteed consistent account setup
- Clear initialization status tracking
- Built-in failure handling

2. Clear Relationships:
- One-to-one account-workspace mapping
- Explicit parent-child relationships
- Clear resource ownership

3. Extensibility:
- Custom initializers for different providers
- Pluggable initialization steps
- Support for future requirements

4. Operational Benefits:
- Clear status tracking
- Built-in retry mechanisms
- Audit trail of account setup

### Negative Consequences

1. Complexity:
- More complex initialization flow
- Need for careful error handling
- Potential for initialization deadlocks

2. Operational Overhead:
- Need for initializer management
- Potential for stuck initializations
- More complex debugging

## Risk Mitigation

1. Initialization Risks:
- Implement timeout mechanisms for initializers
- Create clear initialization dependency graphs
- Monitor initialization progress
- Provide manual intervention capabilities

2. Resource Management:
- Implement comprehensive cleanup mechanisms
- Define clear resource ownership
- Track resource dependencies
- Provide rollback capabilities

3. Operational Risks:
- Create comprehensive monitoring strategy
- Define clear operational procedures
- Implement automated recovery where possible
- Provide clear debugging tools

## Action Items

1. Create detailed design document for account-workspace relationship
2. Define initialization workflow and dependencies
3. Design monitoring and debugging strategy
4. Create operational procedures documentation
5. Define recovery and rollback procedures
6. Create test plan for initialization scenarios
7. Design audit and compliance tracking mechanisms

## Related Documents

- KCP Workspace Documentation
- Platform Mesh Architecture Overview
- Service Provider Integration Guide

## Notes

This is a living document and should be updated as implementation progresses and new requirements or challenges are discovered.