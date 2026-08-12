wardley-beta
    title Checkout, "as it stands" #1
    anchor Customer [0.95, 0.95]
    component Checkout (web) [0.6, 0.8]
    component Payment (refunds) [0.75, 0.5]

    component Card network [0.95, 0.2]
    Customer -> Checkout (web)
    Checkout (web) -> Payment (refunds)
    Payment (refunds) -> Card network
    evolve Payment (refunds) 0.9