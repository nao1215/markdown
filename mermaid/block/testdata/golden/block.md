---
title: "Every Token"
---
block
    columns 3
    plain labeled["With a label"] wide:2
    space space:2 raw:3
    explicit<["&nbsp;"]>(right) labeled<["to the right"]>(right) secondary<["&nbsp;"]>(right, down)
    right<["&nbsp;"]>(right) left<["&nbsp;"]>(left) up<["&nbsp;"]>(up)
    down<["&nbsp;"]>(down) x<["&nbsp;"]>(x) y<["&nbsp;"]>(y)

    block:grouped:2
        columns 2
        inner1 inner2
        style inner1 fill:#f9f
    end
    style plain fill:#eee
    plain --> labeled
    labeled -- "connects to" --> wide
    style plain fill:#fff
    classDef highlight fill:#ff0
    class plain highlight