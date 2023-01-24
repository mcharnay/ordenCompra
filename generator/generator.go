package generator

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"generate/pdf/oc/models"
	"generate/pdf/oc/pdfutil"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GenerateHtmlPage(generarBase64 bool) {
	data := models.PageData{
		Logo:         template.URL(generateBase64Img("./template/logo.png")),
		OCNum:        "301",
		Date:         "23/01/2023",
		Name:         "Test nombre",
		Rut:          "12345678-9",
		CommBusiness: "Test giro",
		Address:      "Test direccion",
		Town:         "Santiago",
		Email:        "correo@prueba.com",
		Details: []models.Detail{
			{Item: "Item1", Description: "test descripcion 1 Mio", Quantity: "1", Price: "300", Tax: "0", Total: "0"},
			{Item: "Item2", Description: "test descripcion 2 Mio", Quantity: "2", Price: "200", Tax: "0", Total: "0"},
			{Item: "Item3", Description: "test descripcion 3 Mio", Quantity: "3", Price: "100", Tax: "0", Total: "0"},
			{Item: "Item4", Description: "test descripcion 4 Mio", Quantity: "4", Price: "600", Tax: "0", Total: "0"},
			{Item: "Item5", Description: "test descripcion 5 Mio", Quantity: "3", Price: "7000", Tax: "0", Total: "0"},
			{Item: "Item6", Description: "test descripcion 6 Mio", Quantity: "3", Price: "500", Tax: "0", Total: "0"},
			{Item: "Item6", Description: "test descripcion 7 Mio", Quantity: "2", Price: "3500", Tax: "0", Total: "0"},
		},
		NetValue: "0",
		Tax:      "0",
		Total:    "",
	}

	//acá se llaman los formateos
	data.Total = pdfutil.AddTotalAmount(data.Details)
	//data.NetValue = sumarTotalDetallesFormateado(data.Details)
	data.Tax = pdfutil.CalculateIva(data.Details)
	data.NetValue = pdfutil.CalculateNetValue(data.Details)

	fmt.Println(data.Total)
	pdfutil.FormatEachDetail(data.Details)
	pdfutil.FormatPrice(data.Details)
	pdfutil.FormatTax(data.Details)
	//formatPdfDate(data)
	data.Date = pdfutil.FormatPdfDate(data)

	/////////////////////////////////////////////////////////////////////////////////////////////

	fmt.Println("leyendo template")
	tmpl := template.Must(template.ParseFiles("./template/template.html"))

	fmt.Println("guardando buffer con template modificado")
	var body bytes.Buffer
	tmpl.Execute(&body, data)

	fmt.Println("buscando path de wkhtmltopdf")
	path, err := exec.LookPath("wkhtmltopdf") //se debe instalar wkhtmltopdf y dejarlo en el path para que esta linea funcione
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("path encontrado: " + path)

	wkhtmltopdf.SetPath(path)
	fmt.Println("Inicializando pdf generator")
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inicializando page reader")
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))
	page.EnableLocalFileAccess.Set(true)
	fmt.Println("añadiendo page a pdf")
	pdfg.AddPage(page)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)

	fmt.Println("creando pdf")
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	os.RemoveAll("./result/")
	_ = os.Mkdir("./result/", os.ModePerm)

	fmt.Println("guardando archivo html con template modificado")
	f := createFile("html")
	writeString, err := f.WriteString(body.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("escritos %d bytes\n", writeString)
	defer f.Close()

	if generarBase64 {
		saveBase64OnFile(pdfg.Bytes())
	} else {
		pdfg.WriteFile("./result/pdfResult.pdf")
	}
}

func generateBase64Img(file string) string {
	fBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	mimeType := http.DetectContentType(fBytes)

	var base64Encoding string
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(fBytes)
	return base64Encoding
}

func saveBase64OnFile(file []byte) {
	base64String := base64.StdEncoding.EncodeToString(file)
	f := createFile("txt")

	writeString, err := f.WriteString(base64String)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("escritos %d bytes\n", writeString)

	defer f.Close()
}

func createFile(ext string) *os.File {
	f, err := os.Create("./result/result." + ext)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
