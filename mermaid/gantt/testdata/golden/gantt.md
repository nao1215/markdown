gantt
    title Every Task Kind
    dateFormat YYYY-MM-DD
    axisFormat %m-%d
    tickInterval 1week
    todayMarker stroke-width:4px
    excludes weekends
    excludes 2024-01-01
    section Plain tasks
    Task :2024-01-01, 2d
    Task with id :task-id, 2024-01-03, 2d
    Task after :after task-id, 1d
    Task after with id :after-id, after task-id, 1d
    section Marked tasks
    Critical :crit, 2024-01-06, 1d
    Critical with id :crit, crit-id, 2024-01-07, 1d
    Active :active, 2024-01-08, 1d
    Active with id :active, active-id, 2024-01-09, 1d
    Done :done, 2024-01-10, 1d
    Done with id :done, done-id, 2024-01-11, 1d
    Critical active :crit, active, 2024-01-12, 1d
    Critical active with id :crit, active, crit-active-id, 2024-01-13, 1d
    Critical done :crit, done, 2024-01-14, 1d
    Critical done with id :crit, done, crit-done-id, 2024-01-15, 1d

    section Milestones
    Milestone :milestone, 2024-01-16, 0d
    Milestone with id :milestone, milestone-id, 2024-01-17, 0d
    Critical milestone :crit, milestone, 2024-01-18, 0d
    Critical milestone with id :crit, milestone, crit-milestone-id, 2024-01-19, 0d