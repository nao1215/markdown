sequenceDiagram
    autonumber
    box Client & Server
    participant Client
    participant Server
    end
    actor Operator

    Client->>Server: GET /users
    Client->>Server: GET /users/1
    Server-->>Client: 200 OK
    Server-->>Client: 200 OK
    Client->)Server: publish event
    Client->)Server: publish event
    Server--)Client: ack
    Server--)Client: ack
    Client-xServer: malformed body
    Client-xServer: malformed body
    Server--xClient: 500 Internal Server Error
    Server--xClient: 500 Internal Server Error

    activate Server
    Client->>+Server: activate on request
    Client->>+Server: activate on request
    Server-->>-Client: deactivate on response
    Server-->>-Client: deactivate on response
    Client->>+Server: async activate
    Client->>+Server: async activate
    Server-->>-Client: async deactivate
    Server-->>-Client: async deactivate
    deactivate Server

    note over Server: a note over the participant
    note right of Server: a note to the right
    note left of Client: a note to the left