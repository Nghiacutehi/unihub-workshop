package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	outputDir := "../samples"
	fileName := "students_error_test.csv"
	outputPath := filepath.Join(outputDir, fileName)
	
	os.MkdirAll(outputDir, 0755)

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	writer.Write([]string{"student_id", "password", "full_name", "email", "phone", "role"})

	// Pre-hashed password
	realHash := "$2a$10$vI8A7sz5iayM5Nf3.M1Rau9v9vVf.f9v9vVf.f9v9vVf.f9v9vVf." 

	// 1. Valid records
	writer.Write([]string{"20269901", realHash, "Sinh viên Hợp lệ 1", "valid1@unihub.edu.vn", "0911111111", "STUDENT"})
	writer.Write([]string{"20269902", realHash, "Sinh viên Hợp lệ 2", "valid2@unihub.edu.vn", "0922222222", "STUDENT"})

	// 2. ERROR: Missing Student ID
	writer.Write([]string{"", realHash, "Lỗi Thiếu MSSV", "error1@unihub.edu.vn", "0933333333", "STUDENT"})

	// 3. ERROR: Duplicate ID (within file)
	writer.Write([]string{"20269905", realHash, "Sinh viên Gốc", "original@unihub.edu.vn", "0955555555", "STUDENT"})
	writer.Write([]string{"20269905", realHash, "Lỗi Trùng MSSV", "duplicate@unihub.edu.vn", "0966666666", "STUDENT"})

	// 4. ERROR: Invalid Email format
	writer.Write([]string{"20269907", realHash, "Lỗi Sai Email", "invalid-email-format", "0977777777", "STUDENT"})

	// 5. More valid records to make it realistic
	for i := 10; i <= 30; i++ {
		writer.Write([]string{
			fmt.Sprintf("2026%04d", i), 
			realHash, 
			fmt.Sprintf("Sinh viên Thử nghiệm %d", i), 
			fmt.Sprintf("test%d@unihub.edu.vn", i), 
			"0900000000", 
			"STUDENT",
		})
	}

	fmt.Printf("⚠️  Generated test file with errors at: %s\n", outputPath)
}
