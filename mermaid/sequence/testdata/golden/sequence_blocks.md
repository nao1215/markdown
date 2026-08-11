sequenceDiagram
    create participant Worker
    create actor Auditor
    loop every minute
    Client->>Worker: poll
    end
    alt the queue is empty
    Worker-->>Client: nothing to do
    else the queue has work
    Worker-->>Client: one job
    end
    opt the caller asked for details
    Worker-->>Client: job details
    end
    par fan out
    Worker->>Auditor: record start
    and and
    Worker->>Auditor: record end
    end
    critical acquire the lock
    Worker->>Auditor: lock
    option the lock is taken
    Worker->>Auditor: wait
    end
    break the job failed
    Worker-->>Client: error
    end
    destroy Auditor
    destroy Worker