package lb

import (
	"fmt"
	"testing"
)

func TestListNic(t *testing.T) {
	a := NewDpvsAgentComm("")
	b, _ := a.ListNicStats()
	for _, v := range b {
		fmt.Printf("%v\n", v)
	}
}
func TestConn(t *testing.T) {
	a := NewDpvsAgentComm("")
	b, _ := a.ListVirtualServices()
	for _, v := range b.Items {
		fmt.Println(*v.AF)
	}

}
