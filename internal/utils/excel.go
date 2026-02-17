package utils

import (
	"io"

	"github.com/xuri/excelize/v2"
)

// ApplicantRow represents a row from Excel
type ApplicantRow struct {
	FullName                  string
	CondominiumRegistrationID string
}

// ParseApplicantExcel parses an Excel file and returns applicant rows
// func ParseApplicantExcel(file io.Reader) ([]ApplicantRow, error) {
// 	f, err := excelize.OpenReader(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer f.Close()

// 	sheet := f.GetSheetName(0)
// 	if sheet == "" {
// 		return nil, errors.New("no sheets found in excel")
// 	}

// 	rows, err := f.GetRows(sheet)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var applicants []ApplicantRow

// 	// Skip header row
// 	for i, row := range rows {
// 		if i == 0 {
// 			continue
// 		}
// 		if len(row) < 2 {
// 			continue
// 		}
// 		applicants = append(applicants, ApplicantRow{
// 			FullName:                  row[0],
// 			CondominiumRegistrationID: row[1],
// 		})
// 	}

// 	return applicants, nil
// }

func ParseApplicantExcel(file io.Reader) ([]ApplicantRow, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var result []ApplicantRow

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		if len(row) < 2 {
			continue
		}

		result = append(result, ApplicantRow{
			FullName:                  row[0],
			CondominiumRegistrationID: row[1],
		})
	}

	return result, nil
}
