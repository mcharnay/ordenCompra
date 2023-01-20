package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/leekchan/accounting"
)

type Detail struct {
	Item        string
	Description string
	Quantity    string
	Price       string
	Tax         string
	Total       string
}

type PageData struct {
	Logo         template.URL
	OCNum        string
	Day          string
	Month        string
	Year         string
	Name         string
	Rut          string
	CommBusiness string
	Address      string
	Town         string
	Email        string
	Details      []Detail
	NetValue     string
	Tax          string
	Total        string
}

func main() {
	generateHtmlPage(false)
}

func generateHtmlPage(generarBase64 bool) {
	data := PageData{
		Logo:         template.URL(generateBase64Img("./template/logo.png")),
		OCNum:        "301",
		Day:          "17",
		Month:        "01",
		Year:         "2023",
		Name:         "Test nombre",
		Rut:          "12345678-9",
		CommBusiness: "Test giro",
		Address:      "Test direccion",
		Town:         "Santiago",
		Email:        "correo@prueba.com",
		Details: []Detail{
			{Item: "Item1", Description: "test descripcion 1 Mio", Quantity: "1", Price: "300", Tax: "0", Total: "300"},
			{Item: "Item2", Description: "test descripcion 2 Mio", Quantity: "2", Price: "200", Tax: "0", Total: "700"},
			{Item: "Item3", Description: "test descripcion 3 Mio", Quantity: "3", Price: "100", Tax: "0", Total: "1000"},
			{Item: "Item4", Description: "test descripcion 4 Mio", Quantity: "4", Price: "600", Tax: "0", Total: "1500000"},
			{Item: "Item5", Description: "test descripcion 5 Mio", Quantity: "3", Price: "7000", Tax: "0", Total: "1500000"},
		},
		NetValue: "1000",
		Tax:      "190",
		Total:    "",
	}

	//acá se llaman los formateos
	data.Total = sumarTotalDetallesFormateado(data.Details)
	//data.NetValue = sumarTotalDetallesFormateado(data.Details)
	data.Tax = calcularIva(data.Details)
	data.NetValue = calcularNeto(data.Details)

	fmt.Println(data.Total)
	formatearTotales(data.Details)
	formatearPrecio(data.Details)
	transformarImpuest(data.Details)
	formatearFecha(data)
	//calcularIva(data.Details)

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

// retorna el total de los detalles sumados (acá se debe hacer el cálculo)
func sumarTotalDetallesFormateado(totalDetalle []Detail) string {
	//recorrer total detalle con un for y después transformarlo a int
	var acumulador int
	var intCantidad int
	var cantidadPorPrecio int
	var ivaDetalle float64 = 0.19
	var intPrecio int

	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}

	for indice, objDetalle := range totalDetalle {
		intCantidad, _ = strconv.Atoi(objDetalle.Quantity) //guarda la cantidad (int)
		intPrecio, _ = strconv.Atoi(objDetalle.Price)      //guarda el precio (int)

		cantidadPorPrecio = intCantidad * intPrecio //cantidad por precio (int)
		fmt.Println("---------> cantidad por precio: ", cantidadPorPrecio)

		result := ivaDetalle * float64(cantidadPorPrecio) //ivaDetalle float, por (float)cantidadPorPrecio = 0.19 del total.
		resultIvaInt := int(result)
		totalPrecioCadaUno := resultIvaInt + cantidadPorPrecio   //calcula el total precio cantidad e iva individual
		objDetalle.Total = strconv.Itoa(int(totalPrecioCadaUno)) //transformo totalcalculado a string para ponerlo en
		fmt.Println("----------> result (iva de cant por precio) :", result)
		fmt.Println("----------> total calculado: ", objDetalle.Total)
		totalDetalle[indice].Total = objDetalle.Total //asigna el valor calculado final al total del item separado

		objDetalle.Tax = strconv.Itoa(int(result)) //tax ahora tiene el 0.19 como String
		var totalCompleto = int(result) + cantidadPorPrecio
		fmt.Println("total completo int ------->", totalCompleto) //muestro como int

		//totalMasIVA = int(result)

		// intTotal = totalMasIVA
		// acumulador = intTotal + acumulador
		fmt.Println(acumulador)

		acumulador = acumulador + totalCompleto
		fmt.Println("acumulador2------->", acumulador) //suma el total como int

		//lo guarda en neto
		//data.NetValue = ac.FormatMoney(acumulador2)
		//PageData.NetValue = "asdas"

		objDetalle.Total = ac.FormatMoney(acumulador)

		fmt.Println("----> ;", objDetalle.Total)
		totalDetalle[indice].Tax = objDetalle.Tax

	}
	fmt.Println(acumulador)
	return ac.FormatMoney(acumulador)

	//el iba, es la suma de los impuestos.
	//el neto, es cantidad por precio.
	//recordar que esta funcion se llama desde el maion, y desde ahi se puede acceder a la data neto e iva.
}

// formatea los totales por separado
func formatearTotales(totalDetalle []Detail) []Detail {
	var intTotal int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}
	for index, objDetalle := range totalDetalle {
		intTotal, _ = strconv.Atoi(objDetalle.Total)
		totalDetalle[index].Total = ac.FormatMoney(intTotal)
	}
	return totalDetalle
}

// formatear precios
func formatearPrecio(preciolDetalle []Detail) []Detail {
	var intPrecio int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}
	for index, objDetalle := range preciolDetalle {
		intPrecio, _ = strconv.Atoi(objDetalle.Price)
		preciolDetalle[index].Price = ac.FormatMoney(intPrecio)
	}
	return preciolDetalle
}

// formatear impuestos
func transformarImpuest(impuestolDetalle []Detail) []Detail {
	var intImpuesto int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}
	for index, objDetalle := range impuestolDetalle {
		intImpuesto, _ = strconv.Atoi(objDetalle.Tax)
		impuestolDetalle[index].Tax = ac.FormatMoney(intImpuesto)
	}
	return impuestolDetalle
}

func formatearFecha(fecha PageData) string {
	var dia string
	var mes string
	var año string
	var fechaFormat string
	dia = fecha.Day
	mes = fecha.Month
	año = fecha.Year
	fechaFormat = dia + "/" + mes + "/" + año

	fmt.Println("fecha " + fechaFormat)
	return fechaFormat
}

func calcularIva(objDetalle []Detail) string {
	var iva, acumulador int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}

	for _, objDetalle := range objDetalle {

		iva, _ = strconv.Atoi(objDetalle.Tax)
		acumulador = acumulador + iva
		fmt.Println("iva de mi funcion ----->", acumulador)
	}

	return ac.FormatMoney(acumulador)
}

func calcularNeto(objDetalle []Detail) string {
	var precio, acumulador, cantidad int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}

	for _, objDetalle := range objDetalle {
		precio, _ = strconv.Atoi(objDetalle.Price)
		cantidad, _ = strconv.Atoi(objDetalle.Quantity)
		acumulador = acumulador + (precio * cantidad)
		fmt.Println("iva de mi funcion ----->", acumulador)
	}

	return ac.FormatMoney(acumulador)

}

//desde el front debería venir como numérico, la función debería trabajar con datos numericos
//debe recibir ints
//voy a trabajar en seperar la funcionalidades como tipo dato, y formatear los datos a pasarlo a la plantilla.
//función queformatea los datos, los debe transformar en string para darle formato.

//hacer un metodo

//puedo cambiar el tipo de dato para hacer mio el proyecto.
//hacer otro método que los formatea a string.

//instalar corse (mati dijo)
