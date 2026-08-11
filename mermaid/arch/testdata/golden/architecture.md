architecture-beta
    group api(cloud)[API]
    group storage(database)[Storage] in api
    service gateway(internet)[Gateway]
    service db(database)[Database] in storage
    service disk(disk)[Disk] in storage
    service worker(server)[Worker] in api
    junction hub
    junction inner in storage

    gateway:R >-- L:hub
    hub:B --< T:worker
    worker{group}:R --> L:db{group}