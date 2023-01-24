package pdfutil

import (
	"generate/pdf/oc/models"
	"strconv"

	"github.com/leekchan/accounting"
)

// sum all totals finals mounts, all > ((quantity * price) * tax) + (quantity * price) = total,  and give format. $x.xxx
func AddTotalAmount(totalDetail []models.Detail) string {
	var accumulator, intAmount, quantityForPrice, intPrice int
	var taxDetail float64 = 0.19
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}

	for index, objDetalle := range totalDetail {
		intAmount, _ = strconv.Atoi(objDetalle.Quantity) //save Quantity (int)
		intPrice, _ = strconv.Atoi(objDetalle.Price)     //save price (int)

		quantityForPrice = intAmount * intPrice //quantity * price (int)

		result := taxDetail * float64(quantityForPrice) //taxDetail float, for (float)quantityForPrice = 0.19 of Totalamount
		resultIvaInt := int(result)
		detailTotal := resultIvaInt + quantityForPrice
		objDetalle.Total = strconv.Itoa(int(detailTotal))
		totalDetail[index].Total = objDetalle.Total

		objDetalle.Tax = strconv.Itoa(int(result)) //tax now have  0.19 like String
		var totalComplete = int(result) + quantityForPrice

		accumulator = accumulator + totalComplete

		objDetalle.Total = ac.FormatMoney(accumulator)

		totalDetail[index].Tax = objDetalle.Tax

	}
	return ac.FormatMoney(accumulator)
}

// Give each total price format $x.xxx
func FormatEachDetail(totalDetail []models.Detail) []models.Detail {
	var intTotal int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}
	for index, objDetail := range totalDetail {
		intTotal, _ = strconv.Atoi(objDetail.Total)
		totalDetail[index].Total = ac.FormatMoney(intTotal)
	}
	return totalDetail
}

// Give price format $x.xxx
func FormatPrice(detailPriceObj []models.Detail) []models.Detail {
	var intPrice int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}
	for index, detailObj := range detailPriceObj {
		intPrice, _ = strconv.Atoi(detailObj.Price)
		detailPriceObj[index].Price = ac.FormatMoney(intPrice)
	}
	return detailPriceObj
}

// Give tax format $x.xxx
func FormatTax(detailTaxObj []models.Detail) []models.Detail {
	var intTax int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}
	for index, objDetalle := range detailTaxObj {
		intTax, _ = strconv.Atoi(objDetalle.Tax)
		detailTaxObj[index].Tax = ac.FormatMoney(intTax)
	}
	return detailTaxObj
}

// func formatearFechaDate(objPageData PageData) time.Time {

// 	fechaActual := objPageData.Date.Format("02/01/2006") //transforma fecha date a string
// 	fmt.Println("Fechaaaaaa String---------------->: ", fechaActual)
// 	fechaDate, err := time.Parse("02/01/2006", fechaActual)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("Fechaaaaaa Date---------------->: ", fechaDate)
// 	return fechaDate
// }

// func formatearFechaString(objPageData PageData) string {
// 	fechaActual := objPageData.Date.Format("02/01/2006") //transforma fecha date a string
// 	fmt.Println("Fechaaaaaa String HOLAAA---------------->: ", fechaActual)
// 	objPageData.Date
// 	return fechaActual
// }

func FormatPdfDate(objPageData models.PageData) string {
	dateString := objPageData.Date
	day := []rune(dateString)
	daySubstring := string(day[0:2])

	month := []rune(dateString)
	monthSubstring := string(month[3:5])

	year := []rune(dateString)
	yearSubstring := string(year[6:10])

	dateStringFormat := daySubstring + " / " + monthSubstring + " / " + yearSubstring
	return dateStringFormat
}

func CalculateIva(objDetail []models.Detail) string {
	var tax, accumulator int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}

	for _, objDetail := range objDetail {
		tax, _ = strconv.Atoi(objDetail.Tax)
		accumulator = accumulator + tax
	}

	return ac.FormatMoney(accumulator)
}

func CalculateNetValue(objDetail []models.Detail) string {
	var price, accumulator, amount int
	var ac = accounting.Accounting{Symbol: "$", Precision: 0, Thousand: ".", Decimal: ","}

	for _, objDetail := range objDetail {
		price, _ = strconv.Atoi(objDetail.Price)
		amount, _ = strconv.Atoi(objDetail.Quantity)
		accumulator = accumulator + (price * amount)
	}

	return ac.FormatMoney(accumulator)

}

//have to install corse
