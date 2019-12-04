package main
import ("fmt"
	"time"
)

func main(){
	fmt.Println(40*time.Second)
	fmt.Println(time.Now())
	fmt.Println("=====")
	currentyear := time.Unix(1512068141, 0).Year()
	t := time.Date(currentyear, 1, 3, 0, 0, 0, 0, time.UTC)
	fmt.Println(t)
	fmt.Println(t.Unix())
}
