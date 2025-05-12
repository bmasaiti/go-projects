package concurrency


type WebsiteChecker func(string) bool 
type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string ) map[string]bool{
	results:= make(map[string]bool)
	resultChannel := make(chan result)

	// for _, url:= range urls{
	// 	go func() {   //go routine anytime you put go before a function definition it because a go routine
	// 		results[url] = wc(url)
	// 	}() // this  () means we execute as soon as the function is declared
	// 	//anonymous functions also maintain lexical scope , ie they can access stuff 
	// 	//here each run of the for loop starts a new go routine whose results gets added to map

	// }

	// time.Sleep(20*time.Second)
	// return results

	for _, url := range urls{
		go func(){
			resultChannel <- result{url, wc(url)} //sending results of goroutine call to a channel
		}()
	}
	for i:= 0; i<len(urls); i++{
		r:= <-resultChannel  //receiving values from channel 
		results[r.string]=r.bool //now we won't get race condition , map is upddated one at a time
	}
	return results
}