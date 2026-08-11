stateDiagram-v2
    Active : The order is active
    state Active {
        Reserved : Stock is reserved
        Reserved --> Packed
        Packed --> Handed over : to the carrier
        [*] --> Reserved
        Handed over --> [*]
    }
    Active --> Closed