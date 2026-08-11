---
title: "Order Lifecycle"
---
stateDiagram-v2
    direction LR
    Draft : The order is being written
    Placed : The order has been placed
    Shipped : The order is on its way
    [*] --> Draft
    [*] --> Placed : restored from a backup
    Draft --> Placed
    Placed --> Shipped : after payment
    Shipped --> [*]
    Placed --> [*] : cancelled by the customer

    note left of Draft : editable
    note right of Shipped : immutable
    note left of Placed
        waiting for payment
        then for the warehouse
    end note
    note right of Draft
        no payment yet
        no reservation yet
    end note

    state split <<fork>>
    state merge <<join>>
    state decide <<choice>>
    ---