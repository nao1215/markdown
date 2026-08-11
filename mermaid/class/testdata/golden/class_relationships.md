classDiagram
    Source <|-- Target
    Source --|> Target
    Source *-- Target
    Source --* Target
    Source o-- Target
    Source --o Target
    Source --> Target
    Source <-- Target
    Source -- Target
    Source ..> Target
    Source <.. Target
    Source ..|> Target
    Source <|.. Target
    Source .. Target

    Source --> Target : labelled
    Source "1" --> "many" Target
    Source "0..1" --> "0..*" Target : labelled

    Whole *-- Part
    Whole *-- Part : contains
    Whole "1" *-- "1..*" Part
    Whole "1" *-- "many" Part : contains

    Left --> Right
    Left --> Right : uses
    Left "1" --> "1" Right
    Left "many" --> "many" Right : uses

    Persistable ()-- Account
    Account --() Persistable