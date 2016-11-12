package handlers

//TODO correct logging
import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "io/ioutil"
	"bytes"
	"html/template"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"

	tf "github.com/florro/test_server/proto/tensorflow"
	pb "github.com/florro/test_server/proto/tensorflow_serving"
	"github.com/florro/test_server/templates"
	db "gitlab.com/gabbagandalf/protos/db_service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write presentation"},
		Todo{Name: "Host meetup"},
	}

	// Send back content type, tells client to expect json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Set status code to OK
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {

	// 	var todo Todo
	// 	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if err := r.Body.Close(); err != nil {
	//         panic(err)
	//     }

	// fmt.Println(r.Method)
	// fmt.Println(r.URL)
	// fmt.Println(r.Header)
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	img = resize.Resize(128, 128, img, resize.Lanczos3)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100})
	if err != nil {
		panic(err)
	}

	var jpg_name []string
	jpg_name = strings.Split(handler.Filename, "/")

	fmt.Println(jpg_name[len(jpg_name)-1])
	var name string = jpg_name[len(jpg_name)-1]

	f, err := os.OpenFile(os.Getenv("HOME")+"/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// var buffer2 = new(bytes.Buffer)
	// *buffer2 = *buffer

	//fmt.Println(r.Body)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to display image")
		panic(err)
	}

	if _, err := io.Copy(f, buffer); err != nil {
		log.Println("unable to save image")
		panic(err)
	}

	// if err := json.Unmarshal(body, &todo); err != nil {
	//     w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//     w.WriteHeader(422) // unprocessable entity
	//     if err := json.NewEncoder(w).Encode(err); err != nil {
	//         panic(err)
	//     }
	// }

	// t := RepoCreateTodo(todo)
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusCreated)
	// if err := json.NewEncoder(w).Encode(t); err != nil {
	//     panic(err)
	// }
}

func makePredRequestfromImgBuffer(buf *bytes.Buffer) *pb.PredictRequest {
	var tmp = [][]byte{buf.Bytes()}
	// var buffer_test = [][]byte{[]byte("teststhi"), []byte("testagaine")}

	Shape := &tf.TensorShapeProto{
		Dim: []*tf.TensorShapeProto_Dim{
			&tf.TensorShapeProto_Dim{Name: "x", Size: 1},
		},
	}

	tensor := &tf.TensorProto{
		StringVal:   tmp,
		TensorShape: Shape,
		Dtype:       tf.DataType_DT_STRING,
	}

	tmpmap := make(map[string]*tf.TensorProto)
	tmpmap["images"] = tensor

	req := &pb.PredictRequest{
		ModelSpec: &pb.ModelSpec{Name: "inception"},
		Inputs:    tmpmap,
	}
	return req
}

func SendImg(w http.ResponseWriter, r *http.Request) {

	imgpath := os.Getenv("HOME") + "/Bilder/affe.jpg"
	file, err := os.Open(imgpath)
	if err != nil {
		log.Fatal(err)
	}

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, file); err != nil {
		log.Fatal(err)
	}
	req := makePredRequestfromImgBuffer(buffer)

	fmt.Println("##########")

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// fmt.Println(req)
	var client = pb.NewPredictionServiceClient(conn)
	feature, err := client.Predict(context.Background(), req)
	if err != nil {
		fmt.Println(grpc.ErrorDesc(err))
		log.Fatal(err)
	}

	fmt.Println(string(feature.Outputs["classes"].StringVal[0]))
	for i, xx := range feature.Outputs["classes"].StringVal {
		fmt.Printf("Num: %v, Class: %s\n", i, xx)
	}
	for i, xx := range feature.Outputs["scores"].FloatVal {
		fmt.Printf("Num: %v, Class: %v\n", i, xx)
	}

}

func SendImgwithTemplate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		//TODO ?
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		if t, err := template.New("UploadTemplate").Parse(templates.UploadTemplate); err != nil {
			log.Fatal(err)
		} else {
			t.Execute(w, token)
		}
	} else {
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		//modify img
		// img, err := jpeg.Decode(file)
		// img = resize.Resize(160, 0, img, resize.Lanczos3)
		// defer file.Close()
		// fmt.Fprintf(w, "%v", handler.Header)
		// f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		// if err != nil {
		//     fmt.Println(err)
		//     return
		// }
		// defer f.Close()
		// buffer := new(bytes.Buffer)
		// if err := jpeg.Encode(buffer, img, nil); err != nil {
		//     log.Println("unable to encode image.")
		// }

		buffer := new(bytes.Buffer)
		_, err = buffer.ReadFrom(file)
		if err != nil {
			log.Fatal(err)
		}

		//Create Request
		req := makePredRequestfromImgBuffer(buffer)

		//Create Stub
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// fmt.Println(req)
		var client = pb.NewPredictionServiceClient(conn)
		feature, err := client.Predict(context.Background(), req)
		if err != nil {
			fmt.Println(grpc.ErrorDesc(err))
			log.Fatal(err)
		}

		fmt.Println(string(feature.Outputs["classes"].StringVal[0]))
		for i, xx := range feature.Outputs["classes"].StringVal {
			fmt.Printf("Num: %v, Class: %s\n", i, xx)
		}
		for i, xx := range feature.Outputs["scores"].FloatVal {
			fmt.Printf("Num: %v, Class: %v\n", i, xx)
		}

		//Show Image
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
		}

		//Save file
		// if _, err := io.Copy(f, buffer); err != nil {
		//     log.Fatal(err)
		// }
	}
}

func SimilaritywithTemplate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		//TODO ?
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		if t, err := template.New("UploadTemplate").Parse(templates.UploadTemplate2); err != nil {
			log.Fatal(err)
		} else {
			t.Execute(w, token)
		}
	} else {
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		//modify img
		// img, err := jpeg.Decode(file)
		// img = resize.Resize(160, 0, img, resize.Lanczos3)
		// defer file.Close()
		// fmt.Fprintf(w, "%v", handler.Header)
		// f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		// if err != nil {
		//     fmt.Println(err)
		//     return
		// }
		// defer f.Close()
		// buffer := new(bytes.Buffer)
		// if err := jpeg.Encode(buffer, img, nil); err != nil {
		//     log.Println("unable to encode image.")
		// }

		buffer := new(bytes.Buffer)
		_, err = buffer.ReadFrom(file)
		if err != nil {
			log.Fatal(err)
		}

		//Create Request
		req := makePredRequestfromImgBuffer(buffer)

		//Create Stub
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// fmt.Println(req)
		var client = pb.NewPredictionServiceClient(conn)
		feature, err := client.Predict(context.Background(), req)
		if err != nil {
			fmt.Println(grpc.ErrorDesc(err))
			log.Fatal(err)
		}
		// fmt.Println(feature.Outputs["bools"].BoolVal)

		conn2, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		// feature.Outputs["bools"]
		var client2 = db.NewDBServerClient(conn2)
		info, err := client2.Info(context.Background(), &db.DBRequest{})
		fmt.Println(info)
		// resp, err := client2.Create(context.Background(), &db.DBRequest{DbName: "test_char"})
		// fmt.Println(resp)
		// var ran int = rand.Int()
        var ran int = rand.Int()
		resp, err := client2.Put(context.Background(),
			&db.PutRequest{DbName: "test_char",
				Key:  "char_a" + strconv.Itoa(ran),
				Data: feature.Outputs["bools"].BoolVal})
		fmt.Println("Put Respnse: ", resp, "err: ", err)

		resp2, err := client2.Search(context.Background(),
			&db.SearchRequest{DbName: "test_char",
				Number: 10,
				Data:   feature.Outputs["bools"].BoolVal})
		fmt.Println("Search Response: ", resp2, "err: ", err)
		// fmt.Println(string(feature.Outputs["classes"].StringVal[0]))
		// for i, xx := range(feature.Outputs["classes"].StringVal) {
		//     fmt.Printf("Num: %v, Class: %s\n", i, xx)
		// }
		// for i, xx := range(feature.Outputs["scores"].FloatVal) {
		//     fmt.Printf("Num: %v, Class: %v\n", i, xx)
		// }

		//Show Image
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
		}

		//Save file
		// if _, err := io.Copy(f, buffer); err != nil {
		//     log.Fatal(err)
		// }
	}
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})
