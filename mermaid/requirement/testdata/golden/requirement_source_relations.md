requirementDiagram
    requirement root {
        id: "1"
        text: "the root requirement"
        risk: Low
        verifymethod: Test
    }
    root - contains -> child
    root - copies -> copy
    root - derives -> derived
    root - satisfies -> satisfied
    root - verifies -> verified
    root - refines -> refined
    root - traces -> traced
    root - contains -> explicit