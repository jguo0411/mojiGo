package image

/*
#include "img.h"
*/
import "C"

func GoSum(a, b int) C.int {
	s := C.sum(C.int(a), C.int(b))
	return s
}