package util

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jung-kurt/gofpdf"
	"github.com/singhraushan/go-jwt-auth-role-access/entities"
)

// HEADER for report table
var HEADER = [...]string{"Transaction Time", "From Acc", "To Acc", "Available Balance", "Transaction Amount", "Transaction Type"}

//CreateToken will create bearer token
func CreateToken(username, role string) (string, error) {
	log.Println("util--CreateToken start.")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   username,
		"role":       role,
		"expiryTime": time.Now().Add(time.Minute * 5).Format(time.RFC822),
	})
	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		return "", err
	}
	log.Println("util--CreateToken end.")
	return tokenString, nil
}

//GeneratePdfReport will generate pdf report
func GeneratePdfReport(transactions *[]entities.Transaction) {
	log.Println("util--GeneratePdfReport start.")

	pdf := newReport()
	pdf = header(pdf, HEADER[:])
	pdf = table(pdf, transactions)

	if pdf.Err() {
		log.Fatalf("Failed creating PDF report: %s\n", pdf.Error())
	}

	err := savePDF(pdf)
	if err != nil {
		log.Fatalf("Cannot save PDF: %s|n", err)
	}
	log.Println("util--GeneratePdfReport end.")
}

func newReport() *gofpdf.Fpdf {
	pdf := gofpdf.New("L", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 28)
	pdf.Cell(40, 10, "Transaction statement")
	pdf.Ln(12)
	pdf.SetFont("Times", "", 20)
	pdf.Cell(40, 10, time.Now().Format("Mon Jan 2, 2006"))
	pdf.Ln(20)
	return pdf
}

func header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 12)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		pdf.CellFormat(40, 7, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}

func savePDF(pdf *gofpdf.Fpdf) error {
	return pdf.OutputFileAndClose("report.pdf")
}

func table(pdf *gofpdf.Fpdf, tbl *[]entities.Transaction) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 10)
	pdf.SetFillColor(255, 255, 255)
	//align := []string{"L", "C", "L", "R", "R", "R"}
	for _, row := range *tbl {
		pdf.CellFormat(40, 7, row.CreatedAt.String()[:19], "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%d", row.FromAccountNumber), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%d", row.ToAccountNumber), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, row.AvailableAmount.String(), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, row.TranferAmount.String(), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, row.Type, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	return pdf
}
