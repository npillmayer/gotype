package bidi

/*
W1.
 ✓  AL  NSM NSM → AL  AL  AL  // done by scanner
 ✓  sos NSM     → sos R       // done by scanner
 -  LRI NSM     → LRI ON      // no overflow LRIs, etc.
 -  PDI NSM     → PDI ON      //

W2.
 ✓  AL … EN     → AL … AN     // done by the scanner

W3.


*/
