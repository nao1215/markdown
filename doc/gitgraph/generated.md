## Git Graph
```mermaid
---
title: Release Flow
---
gitGraph
    commit id: "init" tag: "v0.1.0"
    branch develop order: 2
    checkout develop
    commit type: HIGHLIGHT
    checkout main
    merge develop tag: "v1.0.0"
```