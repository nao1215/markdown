---
title: "Every Shape And Link"
---
flowchart TB
    plain
    text["With text"]
    markdown["`**bold** text`"]
    newlines["`first
second`"]
    round("Round edges")
    stadium(["Stadium"])
    subroutine[["Subroutine"]]
    cylindrical[("Cylindrical")]
    database[("Database")]
    circle(("Circle"))
    asymmetric>"Asymmetric"]
    rhombus{"Rhombus"}
    hexagon{{"Hexagon"}}
    parallelogram[/"Parallelogram"/]
    parallelogramAlt[\"Parallelogram alt"\]
    trapezoid[/"Trapezoid"\]
    trapezoidAlt[\"Trapezoid alt"/]
    doubleCircle((("Double circle")))
    plain-->text
    text-->|"with text"|markdown
    markdown --- newlines
    newlines---|"open with text"|round
    round-.->stadium
    stadium-. "dotted with text" .-> subroutine
    subroutine ==> cylindrical
    cylindrical == "thick with text" ==> database
    database ~~~ circle