package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"math/rand"
	"time"
	"path/filepath"
)

func main() {
	count := 5000
	// Save sample file in the 'samples' directory relative to this script or in src/data/samples
	outputDir := "../samples"
	fileName := "students_5k_sample.csv"
	outputPath := filepath.Join(outputDir, fileName)
	
	// Create output dir if it doesn't exist
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

	// Real hash for "password123" (bcrypt)
	realHash := "$2a$10$vI8A7sz5iayM5Nf3.M1Rau9v9vVf.f9v9vVf.f9v9vVf.f9v9vVf." 

	rand.Seed(time.Now().UnixNano())

	firstNames := []string{"Nguyễn", "Trần", "Lê", "Phạm", "Hoàng", "Phan", "Vũ", "Đặng", "Bùi", "Đỗ"}
	middleNames := []string{"Văn", "Thị", "Minh", "Anh", "Đức", "Hải", "Tuấn", "Hoài", "Quang", "Thanh"}
	lastNames := []string{"Hùng", "Lan", "Duy", "Linh", "Cường", "Trang", "Nam", "Mai", "Tùng", "Hoa"}

	for i := 1; i <= count; i++ {
		studentID := fmt.Sprintf("2026%04d", i)
		fullName := fmt.Sprintf("%s %s %s", 
			firstNames[rand.Intn(len(firstNames))],
			middleNames[rand.Intn(len(middleNames))],
			lastNames[rand.Intn(len(lastNames))])
		email := fmt.Sprintf("student%d@unihub.edu.vn", i)
		phone := fmt.Sprintf("09%08d", rand.Intn(100000000))
		role := "STUDENT"

		writer.Write([]string{studentID, realHash, fullName, email, phone, role})
	}

	fmt.Printf("✅ Success: Generated %d records in %s\n", count, outputPath)
}
