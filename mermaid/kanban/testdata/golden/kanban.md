---
title: "Sprint Board"
config:
  kanban:
    ticketBaseUrl: 'https://example.com/tickets/'
---
kanban
    [Todo]
        [Plain task]
        task-1[Task with metadata]@{ ticket: 'MB-101', assigned: 'Alice', priority: 'Very High' }
    in-progress[In Progress]
        [High]@{ priority: 'High' }
        [Low]@{ priority: 'Low' }
        [Very low]@{ priority: 'Very Low' }

    [Done]
        [Task added by column name]@{ assigned: 'Bob' }