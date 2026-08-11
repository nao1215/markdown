---
title: "Release Flow"
---
gitGraph
    commit
    commit id: "init" tag: "v0.1.0"
    commit type: NORMAL
    commit id: "reverse" type: REVERSE
    commit id: "highlight" type: HIGHLIGHT
    branch develop
    branch feature order: 2
    checkout develop
    commit id: "dev-1"

    checkout main
    merge develop tag: "v1.0.0"
    cherry-pick id: "dev-1"
    cherry-pick id: "dev-1" parent: "init"
    reset id: "init"