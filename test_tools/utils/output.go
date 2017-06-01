package utils
import "fmt"

func green(msg string) {
	fmt.Printf("%c[1;0;32m%s%c[0m\n", 0x1B, msg, 0x1B)
}

func red(msg string) {
	fmt.Printf("%c[1;1;31m%s%c[0m\n", 0x1B, msg, 0x1B)
}

func Ok(msg string) {
	msg = "Ok: " + msg
	green(msg)
}

func Failed(msg string) {
	msg = "Failed: " + msg
	red(msg)
}
