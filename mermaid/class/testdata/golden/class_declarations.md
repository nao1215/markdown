---
title: "Declarations"
---
classDiagram
    direction LR
    %% Every declaration form of the package
    class Empty
    class Account {
        +string id
        +string owner
        -int balance
        #time.Time openedAt
        ~bool frozen
        +Transfer(to Account, amount int) error
        +Deposit(amount int) error
        -validate() error
        #audit() void
        ~reset() void
    }
    class Ledger["Append only ledger"]
    class Repository {
        <<interface>>
    }
    class Statement {
        +string period
        +Render() string
    }
    class Persistable {
        <<Interface>>
    }
    Ledger : +Append(entry Entry) error
    class Ledger {
        <<service>>
    }

    note "A note that belongs to the diagram"
    note for Account "A note that belongs to one class"