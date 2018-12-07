package data

import (
	"fmt"
	"regexp"
	"testing"
)

func TestExport(t *testing.T) {
	//20181130145227+0800
	reg, err := regexp.Compile("^[a-fA-F0-9]{32}$")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(reg.Find([]byte("3C183A30CFFCDA1408DAF1C61D47B274")))
	fmt.Println(isValidMessageAudio("3C183A30CFFCDA1408DAF1C61D47B274"))
}

func isValidMessageAudio(msg string) bool {
	cnt := 0
	for i := 0; i < len(msg); i++ {
		switch msg[i] {
		case '0':
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		case '6':
		case '7':
		case '8':
		case '9':
		case 'a':
		case 'b':
		case 'c':
		case 'd':
		case 'e':
		case 'f':
		case 'A':
		case 'B':
		case 'C':
		case 'D':
		case 'E':
		case 'F':
			cnt++
			if 32 <= cnt {
				return true
			}
		case '/':
			if (i + 10) < len(msg) { // "/storage/"
				ch1 := msg[i+1]
				ch2 := msg[i+8]
				if '/' == ch2 && ('s' == ch1 || 'S' == ch1) {
					return true
				}
			}
		default:
			cnt = 0
		}
	}
	return false
}
