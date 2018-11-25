package data

import (
	"fmt"
	"testing"
)

func TestExport(t *testing.T) {
	err := exportUser("select * from pure_user", "test")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

}
