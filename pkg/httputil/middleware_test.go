package httputil

import (
	"fmt"
	"net/http"
)

func mid1() Middleware {
	return func(tr http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			fmt.Println(1)
			defer func() {
				fmt.Println("mid1 defer")
			}()
			return tr.RoundTrip(req)
		})
	}
}

func mid2() Middleware {
	return func(tr http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			fmt.Println(2)
			defer func() {
				fmt.Println("mid2 defer")
			}()
			return tr.RoundTrip(req)
		})
	}
}

func tr() http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		fmt.Println("tr")
		w := &http.Response{
			StatusCode: 200,
		}

		return w, nil
	})
}

func Example_chain() {
	tr := chain(tr(), []Middleware{mid1(), mid2()}...)

	client := http.Client{
		Transport: tr,
	}

	req, _ := http.NewRequest(http.MethodGet, "https://demo.com/test", nil)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("status code error")
	}

	// Output:
	// 1
	// 2
	// tr
	// mid2 defer
	// mid1 defer
}
