package main

import (
    "fmt"
    "crypto/md5"
    "io"
    "net/http"
    "time"
    "encoding/json"
    "encoding/base64"
    "image"
    "image/jpeg"
    "github.com/nfnt/resize"
    "strconv"
    "os"
    "bytes"
    "html/template"

    "github.com/gorilla/mux"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// pb "google.golang.org/grpc/examples/helloworld/helloworld"
	pb "github.com/florro/test_server/proto/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)
type Todo struct {
    Name      string        `json:"name"`
    Completed bool          `json:"completed"`
    Due       time.Time     `json:"due"`
}

type Todos []Todo

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
    router.HandleFunc("/todos", TodoIndex)
    router.HandleFunc("/img", IMGIndex)
    router.HandleFunc("/todos/{todoID}", TodoShow)
    router.HandleFunc("/upload",upload)

    log.Fatal(http.ListenAndServe(":8080", router))

}


func getimg(path string) image.Image {
    fimg, _ := os.Open(path)
    defer fimg.Close()
    img, _, _ := image.Decode(fimg)

    return img

}


var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

var UploadTemplate string = `<html>
<head>
    <title>Upload file</title>
<form enctype="multipart/form-data" action="http://127.0.0.1:8080/upload" method="post">
      <input type="file" name="uploadfile" />
      <input type="hidden" name="token" value="{{.}}"/>
      <input type="submit" value="upload" />
</head>
<body>
</form>
</body>
</html>`

func hellogrpc(name string) string {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)

    // Contact the server and print out its response.
    // name := defaultName
    if len(os.Args) > 1 {
        name = os.Args[1]
    }
    r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    // log.Printf("Greeting: %s", r.Message)
    return r.Message

}

func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method)
    if r.Method == "GET" {
        fmt.Println("%d", 32 << 20)
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        if t, err := template.New("UploadTemplate").Parse(UploadTemplate); err!= nil {
            log.Println("fuck")
        } else {
            t.Execute(w, token)
        }
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }
        img, err := jpeg.Decode(file)
        img = resize.Resize(160, 0, img, resize.Lanczos3)
        // defer file.Close()
        // fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        // defer f.Close()
        buffer := new(bytes.Buffer)
        if err := jpeg.Encode(buffer, img, nil); err != nil {
            log.Println("unable to encode image.")
        }
        // // var buffer []byte
        // buffer := new(bytes.Buffer)
        // _, err = buffer.ReadFrom(file)
        // if err != nil {
        //     fmt.Println(err)
        //     return
        // }

        // respnse := hellogrpc("you")

        // log.Printf("Greeting: %s", respnse)



        // str := base64.StdEncoding.EncodeToString(buffer.Bytes())
        // if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
        //     log.Println("unable to parse image template.")
        // } else {
        //     data := map[string]interface{}{"Image": str}
        //     if err = tmpl.Execute(w, data); err != nil {
        //         log.Println("unable to execute template.")
        //     }
        // }
        w.Header().Set("Content-Type", "image/jpeg")
        w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
        if _, err := w.Write(buffer.Bytes()); err != nil {
            log.Println("unable to write image.")
        }
        if _, err := io.Copy(f, buffer); err != nil {
            log.Fatal(err)
        }
        f.Close()
        file.Close()
    }
}


// func upload(w http.ResponseWriter, r *http.Request){

//     if tmpl, err := template.New("image").Parse(UploadTemplate); err != nil {
//         log.Println("unable to parse image template.")
//     } else {
//         // data := map[string]interface{}{"Image": str}
//         if err = tmpl.Execute(w, ""); err != nil {
//             log.Println("unable to execute template.")
//         }
//     }
// }

func IMGIndex(w http.ResponseWriter, r *http.Request){

    buffer := new(bytes.Buffer)
    img := getimg("/home/jean/Bilder/Wallpapers/affe.JPG")
    if err := jpeg.Encode(buffer, img, nil); err!=nil{
        log.Println("unable to encode image.")
    }
    str := base64.StdEncoding.EncodeToString(buffer.Bytes())

    if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
        log.Println("unable to parse image template.")
    } else {
        data := map[string]interface{}{"Image": str}
        if err = tmpl.Execute(w, data); err != nil {
            log.Println("unable to execute template.")
        }
    }
}

func Index(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintf(w, "Welcome")
    buffer := new(bytes.Buffer)
    img := getimg("/home/jean/Bilder/Wallpapers/affe.JPG")
    if err := jpeg.Encode(buffer, img, nil); err!=nil{
        log.Println("unable to encode image.")
    }
    w.Header().Set("Content-Type", "image/jpeg")
    w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
    if _, err := w.Write(buffer.Bytes()); err != nil {
        log.Println("unable to write image.")
    }
    // buffer := new(bytes.Buffer)
    // if err:=jpeg.Encode(buffer, *img, nil); err!=nil
}


func TodoShow(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    fmt.Fprintf(w, "Todo show:", todoId)
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
    todos := Todos{
        Todo{Name: "Write presentation"},
        Todo{Name: "Host meetup"},
    }

    json.NewEncoder(w).Encode(todos)
}
