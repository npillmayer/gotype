package math

// Starting Point:
// https://www.overleaf.com/learn/latex/Articles/OpenType-based_math_typesetting:_An_introduction_to_the_STIX2_OpenType_fonts
//   darin z.B. Link auf https://www.tug.org/TUGboat/tb27-1/tb86jackowski.pdf

// Sehr gut:
// OpenType Math: // https://docs.microsoft.com/en-us/typography/opentype/spec/math
//
// http://helm.cs.unibo.it/mml-widget/gtkmathview.html#SEC17
// mathJax
// STIX Fonts: https://www.stixfonts.org/

// HTML / MathML is zentrales Format für gotype.
// Appendix G TeXbook ist interessant, aber Unterschiede:
// in gotype ist Parsen durch MathML schon erledigt. 2 Durchläufe sind wahrscheinlich
// nicht notwendig, sondern "nur" ein Bottom-up Durchlauf durch den MathML Baum.
// Eine Herausforderung stellt das Einlesen und Interpretieren von OpenType Math Fonts
// dar. Es gibt Go Libraries für SNFT Fonts, aber bisher nicht für die MATH tables
// von OpenType. In LuaTeX gibt es ein Modul, allerdings in Lua.
// Workaround noch unklar.

// Ein Khipu gibt es bei Math nicht. Internes Format ist in Baum.
// Teile/Blätter des Baums können ein Khipu sein.
