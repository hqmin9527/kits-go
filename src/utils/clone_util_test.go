package utils

import (
	"fmt"
	"testing"
)

func TestCloneSlicePointer(t *testing.T) {
	s := []string{"A"}
	sp1 := &s
	sp2 := CloneSlicePtr(sp1)
	fmt.Printf("sp2: %v\n", *sp2)
	// ["A"] // 克隆结果
	(*sp2)[0] = "B"
	fmt.Printf("sp1: %v\n", *sp1)
	// ["A"] // 不影响sp1
}

func TestClonePointerSlice(t *testing.T) {
	str := "A"
	ps1 := []*string{&str}
	ps2 := ClonePtrSlice(ps1)
	fmt.Println("ps2[0]:", *ps2[0]) // A
	*(ps2[0]) = "B"
	fmt.Println("ps1[0]:", *ps1[0]) // A

	type User struct {
		Name string
	}
	user := &User{Name: "A"}
	us1 := []*User{user}
	us2 := ClonePtrSlice(us1)
	fmt.Println("us2[0].Name:", us2[0].Name) // A
	us2[0].Name = "B"
	fmt.Println("us1[0].Name:", us1[0].Name) // A
}

func TestCloneSlice(t *testing.T) {
	ss1 := []string{"A"}
	ss2 := CloneSlice(ss1)
	fmt.Println("ss2[0]:", ss2[0]) // A
	ss2[0] = "B"
	fmt.Println("ss1[0]:", ss1[0]) // A
}

func TestClonePointer(t *testing.T) {
	str := "A"
	sp1 := &str
	sp2 := ClonePtr(sp1)
	fmt.Println("sp2:", *sp2) // A
	*sp2 = "B"
	fmt.Println("sp1:", *sp1) // A

	type User struct {
		Name string
	}
	up1 := &User{Name: "A"}
	up2 := ClonePtr(up1)
	fmt.Println("up2.Name:", up2.Name) // A
	up2.Name = "B"
	fmt.Println("up1.Name:", up1.Name) // A
}
