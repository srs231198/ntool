package main

import (
    "fmt"
    "time"
    "math"
    "sort"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/thatisuday/commando"
)

type Link_Object struct {
    Name string `json:"name"`
    Url string `json:"url"`
}

type json_Response struct {
    Links_List []Link_Object `json:"links"`
}

//function to calculate mean of response times
func calc_mean(resp_times []float64) (float64) {
    var sum float64
    var num_terms float64 = float64(len(resp_times))

    for resp_time := range resp_times{
        sum += resp_times[resp_time]
    }

    return sum / num_terms
}
//function to calculate the median of response times
func calc_median(resp_times []float64) (float64) {
    sort.Float64s(resp_times)

    if len(resp_times) % 2 == 0 {
        //if even number of elements then find the two middle most elements and return the average
        middle1 := resp_times[int((len(resp_times) + 1) / 2)]
        middle2 := resp_times[int((len(resp_times) - 1) / 2)]

        return (middle1 + middle2) / 2

    } else {
        //if odd number of elements, find middle element and return it
        middle := int(len(resp_times) / 2)

        return resp_times[middle]
    }
}

//function to implement the additional logic requested
func url_requests(url string, num_requests int) {

    var fastest_time float64 = math.Inf(1)
    var slowest_time float64 = math.Inf(-1)
    var response_time []float64
    var num_error int
    var err_resp []string
    var mean float64
    var median float64
    var smallest_response float64 = math.Inf(1)
    var largest_response float64 = math.Inf(-1)

    
    for req := 0; req < num_requests; req++ {
        start := time.Now()

        result, err := http.Get(url)

        if err != nil {
            err_resp = append(err_resp, err.Error())
            num_error++
        } else {
            defer result.Body.Close()

            body, err := ioutil.ReadAll(result.Body)

            if err != nil {
                fmt.Printf("Error")
            } else {
                size_response := float64(len(body))

                if size_response < smallest_response {
                    smallest_response = size_response
                }

                if size_response > largest_response {
                    largest_response = size_response
                }



            }

            elapsed := float64(time.Since(start).Seconds() * 1000)

            if elapsed < fastest_time {
                fastest_time = elapsed
            }

            if elapsed > slowest_time {
                slowest_time = elapsed
            }

            response_time = append(response_time, elapsed)
        }
    }

    mean = calc_mean(response_time)
    median = calc_median(response_time)
    success_percent := ((num_requests - num_error)/(num_requests)) * 100

    fmt.Printf("Number of request => %v\n", num_requests)
    fmt.Printf("Fastest Time => %v ms\n", fastest_time)
    fmt.Printf("Slowest Time => %v ms\n", slowest_time)
    fmt.Printf("Mean response Time => %v ms\n", mean)
    fmt.Printf("Median response Time => %v ms\n", median)
    fmt.Printf("Percentage of requests that succeeded => %v\n", success_percent)
    
    if num_error > 0 {
        fmt.Printf("Error Codes returned => \n")
    }
    for i := range err_resp {
        fmt.Println(err_resp[i])
    }

    fmt.Printf("Smallest Response size in bytes => %v\n", smallest_response)
    fmt.Printf("Largest Response size in bytes => %v\n", largest_response)

}

func main() {

    description := "Network CLI tool to measure stats of JSON response from a specified server \n"
    description += "Example command - \n"
    description += "go run main.go https://www.cloudflare.com --profile=10 (if using go runtime directly)\n"
    description += "ntool https://www.cloudflare.com --profile=10 (If using \"go get\" to install command from github)"
    //configure commando
    commando.
		SetExecutableName("ntool").
		SetVersion("1.0.0").
        SetDescription(description)
    
    //configure the root command
    commando.
            Register(nil).
            AddArgument("URL", "The full url of the website (eg - https://www.cloudflare.com) ", ""). // required
            AddFlag("profile", "number of requests", commando.Int, 1).   // optional
            SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
                
                // get the URL input
                url := ""
                url += args["URL"].Value

                if url == "https://new_worker.srs1029.workers.dev/links" {
                    resp, err := http.Get(url)
                    
                    if err != nil {
                        fmt.Printf("Error")
                    }
                    defer resp.Body.Close()

                    body, err := ioutil.ReadAll(resp.Body)

                    if err != nil {
                        fmt.Printf("Error")
                    }

                    var c json_Response
                    err = json.Unmarshal(body, &c)

                    if err != nil {
                        panic(err.Error())
                    }
                    

                    for t := range c.Links_List {
                        fmt.Println(c.Links_List[t])
                    }
                }

                number_Requests, _ := flags["profile"].GetInt()
                url_requests(url, number_Requests)
                
            })
    
    commando.Parse(nil)
}