---
title: "System Requirements"
---
requirementDiagram
    direction LR
    requirement "plain requirement" {
        id: "1"
        text: "the system shall do the thing"
        risk: Low
        verifymethod: Test
    }
    requirement "full requirement" {
        id: "2"
        text: "the system shall do the other thing"
        risk: High
        verifymethod: Test
    }
    requirement "typed requirement" {
        id: "3"
        text: "stated with an explicit type"
        risk: Medium
        verifymethod: Test
    }
    functionalRequirement functional {
        id: "4"
        text: "a functional requirement"
        risk: Low
        verifymethod: Test
    }
    interfaceRequirement interface {
        id: "5"
        text: "an interface requirement"
        risk: Medium
        verifymethod: Analysis
    }
    performanceRequirement performance {
        id: "6"
        text: "a performance requirement"
        risk: Medium
        verifymethod: Inspection
    }
    physicalRequirement physical {
        id: "7"
        text: "a physical requirement"
        risk: Medium
        verifymethod: Demonstration
    }
    designConstraint "design constraint" {
        id: "8"
        text: "a design constraint"
        risk: High
        verifymethod: Inspection
    }
    requirement "classified requirement":::important:::reviewed {
        id: "9"
        text: "a requirement carrying classes"
        risk: Low
        verifymethod: Analysis
    }

    element "test suite" {
        type: "simulation"
        docRef: "./tests"
    }
    element "classified element":::important {
    }

    "plain requirement" - contains -> functional
    "plain requirement" - contains -> interface
    functional - copies -> interface
    functional - derives -> performance
    "test suite" - satisfies -> functional
    "test suite" - verifies -> interface
    performance - refines -> physical
    physical - traces -> "design constraint"

    style functional fill:#f9f
    classDef important fill:#ffa
    classDef reviewed stroke:#0a0
    classDef legacy stroke-dasharray: 5 5
    class interface important
    performance:::important:::reviewed