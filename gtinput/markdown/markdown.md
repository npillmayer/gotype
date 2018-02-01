Parsing Markdown
================

Ziele bei "normalen" MD parsern:

* Speed (online-Darstellung im Web)
* HTML-Inlining
* "einfache" Programmierung
* vollständige Dokumente

Parser werden top-down programmiert, so dass die Implementierung einigermaßen einfach ist (PEG).

**Stattdessen**:

* Fragmente parsen
* Dokumente enthalten _fast sicher_ Fehler
* so fehlertolerant wie möglich parsen
* mit Zusatz-Features entstehen mit hoher Wahrscheinlichkeit Mehrdeutigkeiten
* erweiterbar mit weiteren Parsern, die aggregiert werden

**Folgerungen**:

* Bottom-up Parsing
* context-sensitiv? (BCS-Grammar?)
* LR(0)-Menge berechnen (diese gibt für jedes Handle an, ob es korrekt top-down herleitbar ist
* LR(1)-Lookahead berechnen, aber evtl. (nur) SLR
* LA ist evtl in den Produktionen mit angegeben
* strikte Verkürzung (bei den Reduktionen)
* Umgang mit eps-Produktionen? Sind eps-Prod nützlich, um vergessenen Input zu interpolieren? z.B. fehlende Leerzeile?
* high-Level Struktur durch 1. Pass Block-Parsing, dann diese korrigieren, dann inline-Parsing

[Earley-Parser](https://joshuagrams.github.io/pep/) Algorithmus zur Erstellung in Meta-Notation und Impl. (Link) in JavaScript

[Set of Integers](https://github.com/karlseguin/intset)
(für die Mengenoperationen bei der DFA-Bildung/LR(0)-Set)

[Lexer in Go](https://github.com/timtadh/lexmachine#narrative-documentation)

**nachdenken**:

* LR(0) DFA Generator selbst schreiben
* DFS als lesbare Funktionen (siehe Pike on Template lexing)
