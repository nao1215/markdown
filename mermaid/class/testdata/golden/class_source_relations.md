classDiagram
    Account *-- Entry
    Account --> Ledger : writes to
    Account ..> Clock
    Account "1" --> "1" OneToOne
    Account "1" --> "many" OneToMany
    Account "many" --> "1" ManyToOne
    Account "many" --> "many" ManyToMany
    Account "0..1" --> "1..*" Explicit : explicit cardinality