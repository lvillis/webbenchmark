package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/EDDYCJY/fake-useragent"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/valyala/fasthttp"
)

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

type P map[string]interface{}

func log(format string, p P) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	fmt.Println(strings.NewReplacer(args...).Replace(format))
	return strings.NewReplacer(args...).Replace(format)
}

func RandStringBytesMaskImpr(n int) string {
	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func generateRandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func readableBytes(bytes float64) (expression string) {
	if bytes == 0 {
		return "0B"
	}
	var i = math.Floor(math.Log(bytes) / math.Log(1024))
	var sizes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	return fmt.Sprintf("%.3f%s", bytes/math.Pow(1024, i), sizes[int(i)])
}

func showStat() {
	initialNetCounter, _ := net.IOCounters(true)
	for true {
		cpuStat, _ := cpu.Percent(time.Second, false)
		memStat, _ := mem.VirtualMemory()
		loadStat, _ := load.Avg()
		netCounter, _ := net.IOCounters(true)
		log("CPU: {cpuStat} Mem: {memStat} Load: {loadStat}", P{"cpuStat": cpuStat, "memStat": memStat.UsedPercent, "loadStat": loadStat.Load1})
		for i := 0; i < len(netCounter); i++ {
			if netCounter[i].BytesRecv == 0 && netCounter[i].BytesSent == 0 {
				continue
			}
			netName := netCounter[i].Name
			netRecv := readableBytes(float64(netCounter[i].BytesRecv - initialNetCounter[i].BytesRecv))
			netSent := readableBytes(float64(netCounter[i].BytesSent - initialNetCounter[i].BytesSent))
			log("Nic: {netName} ↓ {netRecv} | ↑ {netSent}", P{"netName": netName, "netRecv": netRecv, "netSent": netSent})
		}
		initialNetCounter = netCounter
		time.Sleep(1 * time.Millisecond)
	}
}

func benchmark(url string, method string, postData string, referer string, xForwardFor bool, wg *sync.WaitGroup) {
	for true {
		var request *fasthttp.Request
		var response *fasthttp.Response
		request = fasthttp.AcquireRequest()
		request.Header.SetMethod(method)
		request.SetRequestURI(url)
		if method == "POST" {
			request.SetBodyString(postData)
		}
		if len(referer) == 0 {
			referer = url
		}
		request.Header.Add("Cookie", RandStringBytesMaskImpr(12))
		request.Header.Add("User-Agent", browser.Random())
		request.Header.Add("Referer", referer)
		if xForwardFor {
			randomIp := generateRandomIPAddress()
			request.Header.Add("X-Forwarded-For", randomIp)
			request.Header.Add("X-Real-IP", randomIp)
		}
		response = fasthttp.AcquireResponse()
		if err := fasthttp.Do(request, response); err != nil {
			fasthttp.ReleaseRequest(request)
			fasthttp.ReleaseResponse(response)
			continue
		}
		_, err := io.Copy(ioutil.Discard, bytes.NewReader(response.Body()))
		if err != nil {
			fasthttp.ReleaseRequest(request)
			fasthttp.ReleaseResponse(response)
			continue
		}
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}
	wg.Done()
}

func main() {
	var thread, _ = strconv.Atoi(getEnv("THREAD", "16"))
	var url = getEnv("URL", "http://cachefly.cachefly.net/100mb.test")
	var method = getEnv("METHOD", "GET")
	var postData = getEnv("POST_DATA", "")
	var referer = getEnv("referer", "")
	var xForwardedFor, _ = strconv.ParseBool(getEnv("URL", "false"))
	log("THREAD: {thread}  URL: {url}", P{"thread": thread, "url": url})
	log("method: {method} postData: {postData} referer: {referer} xForwardedFor: {xForwardedFor}", P{"method": method, "postData": postData, "referer": referer, "xForwardedFor": xForwardedFor})

	go showStat()
	var waitGroup sync.WaitGroup
	for i := 0; i < thread; i++ {
		waitGroup.Add(1)
		go benchmark(url, method, postData, referer, xForwardedFor, &waitGroup)
	}
	waitGroup.Wait()
}
