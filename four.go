package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getFieldSet() (fields map[string]bool) {
	fields = make(map[string]bool)
	fields["byr"] = false // Birth Year
	fields["iyr"] = false // Issue Year
	fields["eyr"] = false // Expiration Year
	fields["hgt"] = false // Height
	fields["hcl"] = false // Hair Color
	fields["ecl"] = false // Eye Color
	fields["pid"] = false // Passport ID
	fields["cid"] = false // Country ID
	return
}

func isValidField(fieldName string, fieldValue string) bool {
	switch fieldName {
	case "byr":
		birthYear, err := strconv.Atoi(fieldValue)
		if err != nil || birthYear < 1920 || birthYear > 2002 {
			return false
		}
	case "iyr":
		issueYear, err := strconv.Atoi(fieldValue)
		if err != nil || issueYear < 2010 || issueYear > 2020 {
			return false
		}
	case "eyr":
		expirationYear, err := strconv.Atoi(fieldValue)
		if err != nil || expirationYear < 2020 || expirationYear > 2030 {
			return false
		}
	case "hgt":
		// a number followed by either cm or in
		if strings.HasSuffix(fieldValue, "cm") {
			centimeters, err := strconv.Atoi(strings.TrimSuffix(fieldValue, "cm"))
			if err != nil || centimeters < 150 || centimeters > 193 {
				return false
			}
		} else if strings.HasSuffix(fieldValue, "in") {
			inches, err := strconv.Atoi(strings.TrimSuffix(fieldValue, "in"))
			if err != nil || inches < 59 || inches > 76 {
				return false
			}
		} else {
			// Any other suffix is not allowed
			return false
		}
	case "hcl":
		if len(fieldValue) != 7 || string(fieldValue[0]) != "#" {
			return false
		}
		for i := 1; i < len(fieldValue); i++ {
			if !strings.Contains("0123456789abcdef", string(fieldValue[i])) {
				return false
			}
		}
	case "ecl":
		eyeColors := []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}
		isValid := false
		for _, color := range eyeColors {
			if color == fieldValue {
				isValid = true
				break
			}
		}
		return isValid
	case "pid":
		_, err := strconv.Atoi(fieldValue)
		if err != nil || len(fieldValue) != 9 {
			// Must be a number, represented by 9 digits (w leading zeros allowed)
			return false
		}
	}
	// If a field isn't present, it can be valid by default
	return true
}

func main() {
	file, err := os.Open("./four-input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	fields := getFieldSet()
	numValid := 0
	isValid := true
	for scanner.Scan() {
		// Example: hcl:#ae17e1 iyr:2013
		line := scanner.Text()
		if err != nil {
			log.Fatal(err)
		}
		if line == "" {
			// End of this passport
			// Check if all fields are present
			for name, isPresent := range fields {
				if !isPresent && name != "cid" {
					isValid = false
					break
				}
			}
			if isValid {
				numValid++
			}
			// Reset for next one
			fields = getFieldSet()
			isValid = true
		} else {
			entries := strings.Fields(line)
			for _, field := range entries {
				// name:value
				fieldInfo := strings.Split(field, ":")
				fields[fieldInfo[0]] = true
				if !isValidField(fieldInfo[0], fieldInfo[1]) {
					isValid = false
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of valid passports: ", numValid)
}
