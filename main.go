// Fetch prints the content fount at a URL
package main

import(
	"fmt"
	"io"
	"time"
	"net/http"
	"os"
)
func WithoutConcurrency(arguments []string){
	start:=time.Now()
	for _,url := range arguments{
		resp , err := http.Get(url)
		if err!=nil{
			fmt.Fprintf(os.Stderr,"fetch : %v\n",err)
			os.Exit(1)
		}
		b,err2 := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err2 != nil{
			fmt.Fprintf(os.Stderr,"fetch : reading %s : %v\n",url,err2)
		}
		fmt.Printf("%s",b)
	}
	fmt.Printf("%fs elapsed\n",time.Since(start).Seconds())
}

func fetch(url string, ch chan <- string){
	start:=time.Now()
	resp,err := http.Get(url)
	if err!=nil{
		ch <- fmt.Sprint(err)
		return
	}
	html , err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err!=nil{
		ch <- fmt.Sprintf("while reading %s : %v",url,err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%fs \t %s \t %s ",secs,url,html)
}
func WithConcurrency(arguments []string){
	start := time.Now()
	ch := make(chan string)
	for _ , url := range arguments{
		go fetch(url,ch)
	}
	for range arguments {
		fmt.Println(<-ch) // wait for each fetch to send a result
	}
	fmt.Printf("%fs elapsed\n",time.Since(start).Seconds())
}
func main(){
	WithoutConcurrency(os.Args[1:])
	WithConcurrency(os.Args[1:])
}
