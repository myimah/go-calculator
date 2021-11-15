package web

import (
	"emorisse.fr/go-calculator/pkg/operation"
	"emorisse.fr/go-calculator/pkg/parser"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "HTTP ", log.Flags())
}

type errorRes struct {
	Error string `json:"error"`
}

type resultRes struct {
	Result string `json:"result"`
}

type calculationReq struct {
	Calculation string `json:"computation"`
}

func StartServer(address, port string) {
	var staticFolder = http.FileServer(http.Dir("./static"))

	http.Handle("/", staticFolder)
	http.HandleFunc("/api/calculate", handleApiCalculate)
	var addr = fmt.Sprintf("%s:%s", address, port)

	logger.Printf("Starting server listening to port %s...\n", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		logger.Fatalln(err)
	}
}

func handleApiCalculate(res http.ResponseWriter, req *http.Request) {
	defer logger.Printf("%s %s\n", req.RemoteAddr, req.URL.Path)

	// TODO : add more security, and input checking

	res.Header().Set("Content-Type", "application/json")

	if req.Method == "POST" {
		res.WriteHeader(200)

		var ope, err = readCalculation(req.Body)

		if err != nil {
			sendJson(errorRes{Error: err.Error()}, res)
			return
		}

		var result = fmt.Sprintf("%s", ope.Eval().GetString())
		sendJson(resultRes{result}, res)

	} else {
		res.WriteHeader(400)
		sendJson(errorRes{"Bad Request"}, res)
	}
}

func readCalculation(reader io.Reader) (operation.Operation, error) {
	var content, err = ioutil.ReadAll(reader)
	var ope operation.Operation

	if err != nil {
		return nil, err
	}

	var req calculationReq

	err = json.Unmarshal(content, &req)

	if err != nil {
		return nil, err
	}

	ope, err = parser.Parse(req.Calculation)

	if err != nil {
		return nil, err
	}

	return ope, nil
}

func sendJson(elem interface{}, res http.ResponseWriter) {
	jsonElem, _ := json.Marshal(elem)
	_, _ = res.Write(jsonElem)
}