erDiagram
    teachers ||--o{ students : "teaches"
    schools }|..|| teachers : "employs"
    students |o--|{ clubs : "joins"
    clubs {
        int id PK "Club ID"
    }
    rooms {

    }
    schools {
        int id PK "School ID"
    }
    students {
        int id PK,UK "Student ID"
        int teacher_id FK "Teacher ID"
    }
    teachers {
        int id PK,UK "Teacher ID"
        string name  "Teacher name"
    }
