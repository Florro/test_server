package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
    "crypto/tls"
    "crypto/x509"
    "log"
    "flag"
)

func postFile(client *http.Client, filename string, targetUrl string) error {
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


    resp, err := client.Post(targetUrl, contentType, bodyBuf)
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

func getTLSclient(path string) *http.Client{

    caCert, err := ioutil.ReadFile(path + "/ca.crt")
    if err != nil {
        log.Fatal(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    tlsConfig := &tls.Config{
        RootCAs:    caCertPool,
    }
    tlsConfig.BuildNameToCertificate()
    transport := &http.Transport{TLSClientConfig: tlsConfig}
    client :=&http.Client{Transport: transport}

    return client
}

func main() {

    var tls = flag.Bool("tls", false, "bool for tls usage")
    flag.Parse()

	filename := os.Getenv("HOME") + "/Bilder/affe.jpg"

    var client *http.Client
    var target_url string
    if *tls == false {
        client = &http.Client{}
        target_url = "http://localhost:8080/todos"
    } else {
        client = getTLSclient("../ssl/localhost")
        target_url = "https://localhost:8080/todos"
    }
	postFile(client, filename, target_url)
}
