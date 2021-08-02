// You can edit this code!
// Click here and start typing., go run main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	gofpdf "github.com/jung-kurt/gofpdf"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"

	"database/sql"

	_ "github.com/lib/pq"
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type ErrorModel struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
	Status  bool        `json:"status"`
}

type HackerRand struct {
	CodeEmp string `json:"codeEmp"`
}

type _DateReport struct {
	Fromdate string `json:"Fromdate"`
	Todate   string `json:"Todate"`
	Note     string `json:"Note"`
}

type errYsPH struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type ysPhuHuynh struct {
	Success bool          `json:"success"`
	Data    interface{}   `json:"data"`
	List    []interface{} `json:"list"`
	Error   errYsPH       `json:"error"`
}

type StatusCodes int

const (
	StatusCodeOK                       = 200
	StatusCodeBadRequest               = 400
	SearchRequestUNIVERSAL StatusCodes = 0 // UNIVERSAL
	SearchRequestWEB       StatusCodes = 1 // WEB
	SearchRequestIMAGES    StatusCodes = 2 // IMAGES
	SearchRequestLOCAL     StatusCodes = 3 // LOCAL
	SearchRequestNEWS      StatusCodes = 4 // NEWS
	SearchRequestPRODUCTS  StatusCodes = 5 // PRODUCTS
	SearchRequestVIDEO     StatusCodes = 6 // VIDEO
)

var Articles []Article
var ErrorContains ErrorModel

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func postAPI(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		Articles = []Article{
			Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
			Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
		}
		// json.NewEncoder(w).Encode(Articles)
		ErrorContains = ErrorModel{Code: 200, Status: true, Message: "", Data: Articles}

	case http.MethodHead:
		w.Header().Set("GiaTriNek", "catch me if you can")

	case http.MethodPut:
		_str := r.Header.Get("Authorize")
		ErrorContains = ErrorModel{Code: 200, Status: true, Message: "Put thành công", Data: _str}

	case http.MethodPatch:

		var _requestBody _DateReport

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ErrorContains = ErrorModel{Code: StatusCodeOK, Message: "API Failed", Status: false, Data: err}
		} else {
			json.Unmarshal(b, &_requestBody)
			ErrorContains = ErrorModel{Code: StatusCodeOK, Message: "Success", Status: true, Data: _requestBody}
		}

	default:
		ErrorContains = ErrorModel{Code: 0, Message: "Method not support"}

	}

	js, err := json.Marshal(ErrorContains)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func testAPIGin(c *gin.Context) {
	var _Err ErrorModel = ErrorModel{Code: StatusCodeOK, Data: "Hello Gin"}

	c.IndentedJSON(http.StatusOK, _Err)
}

func postAPIGin(c *gin.Context) {

	var _d string = c.GetHeader("Authorization")

	jsonData, err := c.GetRawData()

	if err != nil {
		var _Err ErrorModel = ErrorModel{Code: StatusCodeBadRequest, Status: false, Message: err.Error()}
		c.IndentedJSON(http.StatusOK, _Err)
	}

	var _requestBody interface{}

	json.Unmarshal(jsonData, &_requestBody)
	var _Err ErrorModel = ErrorModel{Code: StatusCodeOK, Data: _requestBody}

	c.Writer.Header().Add("ValueReadMe", _d)
	c.IndentedJSON(http.StatusOK, _Err)
}

func mainGin() {
	router := gin.Default()

	router.GET("/test", testAPIGin)
	router.POST("/test", postAPIGin)
	router.POST("/downloadFile", downloadFromGin)
	router.POST("/reados", readOSVariable)
	router.POST("/export-pdf", downloadFilePDF)
	router.GET("/call-another-api", callAPIToAnother_GetMethod)
	router.POST("/call-another-api", postToAnotherAPI_PostMethod)
	router.POST("/api-header-body", postToAnotherAPI_With_HeaderAndBody)

	// api := router.Group("/api")
	// {
	// 	api.GET("/test", func(ctx *gin.Context) {
	// 		ctx.JSON(200, gin.H{
	// 			"message": "test successful",
	// 		})
	// 	})
	// }

	router.Run(":6673")
}

func saveExcelFile() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet2021")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "000101"
	cell = row.AddCell()
	cell.Value = "Chinese"
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func downloadFromGin(c *gin.Context) {

	// file, err := os.Open("MyXLSXFile.xlsx") // For read access.
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var _nameFile string = "file_golang_2021.xlsx"

	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+_nameFile)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File("./MyXLSXFile.xlsx")
}

func readOSVariable(c *gin.Context) {
	var _var string = os.Getenv("GOPATH")
	fmt.Println(_var)
	c.IndentedJSON(http.StatusOK, _var)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/hellopost", postAPI)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func mainExportPDF() {
	err := GeneratePdf("hello.pdf")
	if err != nil {
		panic(err)
	}
}

// GeneratePdf generates our pdf by adding text and images to the page
// then saving it to a file (name specified in params).
func GeneratePdf(filename string) error {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	// pdf.CellFormat(190, 7, "Welcome to golangcode.com", "0", 0, "CM", false, 0, "")
	pdf.CellFormat(190, 7, "<i>Hello p</i>", "5", 0, "CM", false, 0, "")
	pdf.Cell(100, 0, "Pin")
	pdf.Cell(100, 0, "Pin 22")

	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	// pdf.ImageOptions(
	// 	"minus.jpg",
	// 	80, 20,
	// 	0, 0,
	// 	false,
	// 	gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
	// 	0,
	// 	"",
	// )

	return pdf.OutputFileAndClose(filename)
}

func loremList() []string {
	return []string{
		"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod " +
			"tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut " +
			"aliquip ex ea commodo consequat.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum " +
			"dolore eu fugiat nulla pariatur.",
		"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui " +
			"officia deserunt mollit anim id est laborum.",
	}
}

func mainTable_GoPDF() {
	const (
		colCount = 3
		colWd    = 60.0
		marginH  = 15.0
		lineHt   = 5.5
		cellGap  = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	var (
		cellList [colCount]cellType
		cell     cellType
	)

	pdf := gofpdf.New("P", "mm", "A4", "") // 210 x 297
	header := [colCount]string{"Column A", "Column B", "Column C"}
	alignList := [colCount]string{"L", "C", "R"}
	strList := loremList()
	pdf.SetMargins(marginH, 15, marginH)
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()

	// Headers
	pdf.SetTextColor(224, 224, 224)
	pdf.SetFillColor(64, 64, 64)
	for colJ := 0; colJ < colCount; colJ++ {
		pdf.CellFormat(colWd, 10, header[colJ], "1", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	count := 0
	for rowJ := 0; rowJ < 2; rowJ++ {
		maxHt := lineHt
		// Cell height calculation loop
		for colJ := 0; colJ < colCount; colJ++ {
			count++
			if count > len(strList) {
				count = 1
			}
			cell.str = strings.Join(strList[0:count], " ")
			cell.list = pdf.SplitLines([]byte(cell.str), colWd-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList[colJ] = cell
		}
		// Cell render loop
		x := marginH
		for colJ := 0; colJ < colCount; colJ++ {
			pdf.Rect(x, y, colWd, maxHt+cellGap+cellGap, "D")
			cell = cellList[colJ]
			cellY := y + cellGap + (maxHt-cell.ht)/2
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(colWd-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0,
					alignList[colJ], false, 0, "")
				cellY += lineHt
			}
			x += colWd
		}
		y += maxHt + cellGap + cellGap
	}

	pdf.OutputFileAndClose("PDFFile_Table.pdf")

	// fileStr := example.Filename("Fpdf_SplitLines_tables")
	// err := pdf.OutputFileAndClose(fileStr)
	// example.Summary(err, fileStr)
}

func downloadFilePDF(c *gin.Context) {
	mainExportPDF()

	c.Writer.Header().Add("Content-Disposition", "attachment; filename=pdfSample.pdf")
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File("./PDFFile_Table.pdf")
}

func callAPIToAnother_GetMethod(c *gin.Context) {

	resp, err := http.Get(`https://apiys.yschool.vn/api/dangkytaikhoan/DangKyOTP?_soDienThoai=0909000001`)
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)

	sb3 := sb + ""
	c.IndentedJSON(StatusCodeOK, sb3)

	// var _mYS ysPhuHuynh
	// json.Unmarshal(body, &_mYS)
	// c.IndentedJSON(StatusCodeOK, _mYS)
}

func postToAnotherAPI_PostMethod(c *gin.Context) {
	values := map[string]string{"PhoneNumber": "0223355617", "Password": "123123123"}
	jsonValue, _ := json.Marshal(values)
	respLogin, errLogin := http.Post(`https://apiys.yschool.vn/api/account/login_v2`, "application/json", bytes.NewBuffer(jsonValue))
	if errLogin != nil {
		log.Fatalln(errLogin)
		c.IndentedJSON(StatusCodeOK, errLogin.Error())
	}
	bodyLogin, _ := ioutil.ReadAll(respLogin.Body)
	sb2 := string(bodyLogin)
	c.IndentedJSON(StatusCodeOK, sb2)
}

func postToAnotherAPI_With_HeaderAndBody(c *gin.Context) {

	client := http.Client{}
	valuesCheckUpdate := map[string]string{"DevicesToken": "bf7ab492-ad95-402b-92dd-418860273dd5", "DevicesInfo": "android", "IdUser": "189127", "Status": "1", "IdApp": "mn_ph_m"}
	jsonValueCheckUpd, _ := json.Marshal(valuesCheckUpdate)

	reqYS, _ := http.NewRequest(http.MethodPost, `https://apiys.yschool.vn/api/thongbaophuhuynh/ThongBao_Insert_Update_V2`, bytes.NewBuffer(jsonValueCheckUpd))
	reqYS.Header.Add("Token", "365f9c64-d5d2-485f-9da9-b5d66ed37995")
	reqYS.Header.Add("Content-Type", "application/json")

	resTB, errTB := client.Do(reqYS)
	if errTB != nil {
		c.IndentedJSON(StatusCodeOK, errTB.Error())
	}
	bodyLogin, _ := ioutil.ReadAll(resTB.Body)
	sb2 := string(bodyLogin)
	c.IndentedJSON(StatusCodeOK, sb2)
}

func callImportModuleSQL() {

	connStr := "user=postgres dbname=SampleDesmo password=123123 host=localhost sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//region

	rows, err := db.Query("SELECT username FROM accounts;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			log.Fatal(err)
		}
		fmt.Println(title)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	//endregion

	fmt.Printf("\nSuccessfully connected to database!\n")

}

func main() {

	callImportModuleSQL()

	saveExcelFile()
	// handleRequests()
	mainGin()
	// handleRequests_Methods()
}
