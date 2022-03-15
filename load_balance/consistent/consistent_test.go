package consistent

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetMultipleQuick(t *testing.T) {

	basepath := "/censor"
	pathList := make([]string, 0, 100)
	for j := 0; j < 3; j++ {
		basepathTmp := basepath + strconv.FormatInt(int64(j), 10)
		for i := 0; i < 100; i++ {
			pathList = append(pathList, fmt.Sprintf(basepathTmp+"-"+strconv.FormatInt(int64(i), 10)))
		}
	}

	x := NewConsistentWithServKeys(pathList)

	for i := 0; i < 20; i++ {
		// v, err := x.Get(strconv.FormatInt(int64(i), 10))
		v, err := x.Get("3")
		if err != nil {
			fmt.Println("error is ", err)
			return
		}
		fmt.Println(v)
	}

}
