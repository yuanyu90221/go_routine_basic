# 學習golang channel跨process間的溝通

## reference doc
[go routine with channel](https://blog.wu-boy.com/2020/01/when-to-use-go-channel-and-goroutine/)
## 1 communicate with share memory
舉例： 產生一個 int 的slice
產生10個goroutine去填入資料
由於10個goroutine各自去執行
因此有可能同時存取到slice
所以有機會會導致存取到同個位置 前面寫入的值被後來執行的值複寫
原始程式碼
```golang===
func addByShareMemory(n int) []int {
	var ints []int
	var wg sync.WaitGroup
	var mux sync.Mutex

	wg.Add(n) //add n counter
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done() // counter--
			ints = append(ints, i)
		}(i)
	}

	wg.Wait()// make the goroutine wait for last goroutine execution
	return ints
}

func main() {
	foo := addByShareMemory(10)
	fmt.Println(len(foo))
	fmt.Println(foo)
}
```
那要如何保證 每次都只有一個goroutine去存取到slice呢？

作法1: 把runtime.GOMAXPROCS(1) 也就是限制同時只有一個goroutine能執行

不太好的作法 不實際
```golang===
function init(){
    runtime.GOMAXPROCS(1) 
}
```

作法2:使用mutex來lock 住存取slice的區塊 讓每次在存取的時候都限制必須拿到lock來才能存取 直到寫完值才能release lock讓其他routine拿

```golang===
func addByShareMemory(n int) []int {
	var ints []int
	var wg sync.WaitGroup
	var mux sync.Mutex

	wg.Add(n) //add n counter
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done() // counter--
			mux.Lock()
			ints = append(ints, i)
			mux.Unlock()
		}(i)
	}

	wg.Wait()
	return ints
}

func main() {
	foo := addByShareMemory(10)
	fmt.Println(len(foo))
	fmt.Println(foo)
}
```
## 2 share by communicating

作法3 建立一個只允許寫入的buffered channel, 讓10個goroutine去寫入channel

然後把channel值寫入ints

然後最後等到channel 的buffer size 到達10個之後在關閉channel 

```golang===
// share memory by communicate
func addByShareCommunicate(n int) []int {
	var ints []int
	channel := make(chan int, n)// build a channel

	for i := 0; i < n; i++ {
		go func(channel chan<- int, order int) { // write into channel value
			channel <- order
		}(channel, i)
	}

	for i := range channel {
		ints = append(ints, i)

		if len(ints) == n {
			break
		}
	}
	close(channel)

	return ints
}
```