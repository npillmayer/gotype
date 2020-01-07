package layout

//   Box type     | Height               Width                   Margin        Margin
//                |                                              (left/right)  (top/bottom)
//   -------------+------------------------------------------------------------------------
//   Inline       | N/A                  N/A                     auto->0       N/A
//   Block        | auto->content-based  auto->constraint-based  auto->center  auto->0
//   Float        | auto->content-based  auto->shrink-to-fit     auto->0       auto->0
//   Inline-block | auto->content-based  auto->shrink-to-fit     auto->0       auto->0
//   Absolute     | special              special                 special       special
