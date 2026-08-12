## Wardley map

```mermaid
wardley-beta
    title Checkout, as it stands
    anchor Customer [0.95, 0.95]
    component Checkout (web) [0.6, 0.8]
    component Payment service [0.75, 0.5]
    component Card network [0.95, 0.2]
    Customer -> Checkout (web)
    Checkout (web) -> Payment service
    Payment service -> Card network
    evolve Payment service 0.9
```
