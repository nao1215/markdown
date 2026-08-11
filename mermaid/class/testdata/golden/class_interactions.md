classDiagram
    class Account
    link Account "https://example.com/account" "Open the docs"
    callback Account "showDetails" "Show the details"
    click Account call openDetails() "Open the details"
    click Account href "https://example.com/account" "Open in a new tab"

    style Account fill:#f9f,stroke:#333
    classDef highlight fill:#ff0
    cssClass "Account" highlight;
    class Account:::highlight