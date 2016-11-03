package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
)

// import (
// 	"fmt"
// 	"net/http"
// 	"io/ioutil"
// )

// func main() {
// 	resp, err := http.Get("http://localhost:8080/todos")
// 	if err != nil {
// 		fmt.Println("ERROR")
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(body)
// 	fmt.Println("############################")
// }

func postFile(filename string, targetUrl string) error {
    bodyBuf := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuf)

    // this step is very important
    fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
    if err != nil {
        fmt.Println("error writing to buffer")
        return err
    }

    // open file handle
    fh, err := os.Open(filename)
    if err != nil {
        fmt.Println("error opening file")
        return err
    }

    //iocopy
    _, err = io.Copy(fileWriter, fh)
    if err != nil {
        return err
    }

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

	fmt.Println("SENDSENDSEND")
	fmt.Println(contentType)
	fmt.Println("--------")

    resp, err := http.Post(targetUrl, contentType, bodyBuf)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    fmt.Println(resp.Status)
	fmt.Println("RESPONSE BODY")
    fmt.Println(string(resp_body))
	fmt.Println("RESPONSE BODY")
    return nil
}

func main() {
	target_url := "http://localhost:8080/todos"
	filename := "/home/florian/Downloads/affe.jpg"
	postFile(filename, target_url)
}
