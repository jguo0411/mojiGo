package image

/*
#include "img.h"
*/
import "C"

func GoSum(a, b int) C.int {
	s := C.sum(C.int(a), C.int(b))
	return s
}

//func toCPoints(points []image.Point) C.struct_Points {
//	cPointSlice := make([]C.struct_Point, len(points))
//	for i, point := range points {
//		cPointSlice[i] = C.struct_Point{
//			x: C.int(point.X),
//			y: C.int(point.Y),
//		}
//	}
//
//	return C.struct_Points{
//		points: (*C.Point)(&cPointSlice[0]),
//		length: C.int(len(points)),
//	}
//}
