// square sum of integers
package main

import (
	"bufio"
	"bytes"
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var slice []int

func main() {
	var numOfTestCase int
	fmt.Printf("Enter number of test case: ")
	fmt.Scanln(&numOfTestCase)
	if numOfTestCase >= 1 && numOfTestCase <= 100 {
		checkNumOfTestCase(numOfTestCase, 1)
		fmt.Println("Output")
		printSum(0)
	} else {
		fmt.Println(" !! Number of test case must be within 1 <= N <= 100")
	}
}

func checkNumOfTestCase(numOfTestCase int, serial int) {
	if numOfTestCase-serial >= 0 {
		var itemSize int
		fmt.Printf("Enter number of integer, test case %d will consist: ", serial)
		fmt.Scanln(&itemSize)
		if itemSize > 0 && itemSize <= 100 {
			checkTestCaseInput(serial, itemSize)
			serial = serial + 1
			checkNumOfTestCase(numOfTestCase, serial)
		} else {
			fmt.Println(" !! Number of integers in a line must be within 0 < X <= 100")
		}
	}
}

func checkTestCaseInput(serial int, itemSize int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter space-separated integers for test case %d: ",
		serial)
	scanner.Scan()
	arr := strings.Fields(scanner.Text())
	if len(arr) == itemSize {
		calculateSum(0, arr, 0)
	} else {
		fmt.Printf(" !! Test case must be %d space-separated integers \n",
			itemSize)
		checkTestCaseInput(serial, itemSize)
	}
}

func calculateSum(serial int, arr []string, sum int) {
	if len(arr)-serial > 0 {
		dgt, err := strconv.Atoi(arr[serial])
		if err != nil {
			fmt.Printf("Error: string (%s) to integer convertion error: %v \n",
				arr[serial], err)
		}
		if dgt >= -100 && dgt <= 100 {
			if dgt > 0 {
				sum = sum + dgt*dgt
			}
			if len(arr)-serial == 1 {
				slice = append(slice, sum)
			}
			serial = serial + 1
			calculateSum(serial, arr, sum)
		} else {
			fmt.Println("Error: space seperated interger must be within -100 <= Yn <= 100")
		}
	}
}

func printSum(serial int) {
	if len(slice)-serial > 0 {
		fmt.Println(slice[serial])
		serial = serial + 1
		printSum(serial)
	}
}

func postRequest() {
	url := "https://api.challenge.hennge.com/challenges/003"
	var jsonStr = `{
		"github_url": "https://gist.github.com/ashish041/f4bd16affbc033b987dc07611a2b5980",
		"contact_email": "as.ku.041@gmail.com"
	}`
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	password := generatePassCode("as.ku.041@gmail.comHENNGECHALLENGE003")
	fmt.Printf("password: %s \n", password)
	req.SetBasicAuth("as.ku.041@gmail.com", password)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Response status: %s \n", response.Status)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	log.Printf("Respose Body: %s \n", string(body))
}

func generatePassCode(utf8string string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    10,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return ""
	}
	return passcode
}
