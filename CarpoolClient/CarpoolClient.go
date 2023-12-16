package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	//"strings"
	//"time"
)
//Struct for Passenger
type Passenger struct {
	PassengerID int    `json:"PassengerID"`
	FirstName   string `json:"First Name"`
	LastName    string `json:"Last Name"`
	PhoneNo     string `json:"Phone No"`
	Email       string `json:"Email"`
	CarOwner    int    `json:"Car Owner"`
}
//Struct for CarOwner
type CarOwner struct {
	PassengerID  int    `json:"PassengerID"`
	CarPlateNo   string `json:"Car Plate No"`
	LicenseNo    string `json:"License No"`
	MaxPassenger int    `json:"Max Passenger"`
	CarColor     string `json:"Car Color"`
	YearRel      int    `json:"Year Release"`
}
//Struct for Trip
type Trip struct {
	TripRef         int    `json:"Trip Reference"`
	PickupAddr      string `json:"Pickup Addr"`
	AlterPickupAddr string `json:"Alter Pick Addr"`
	StartTrip       string `json:"Start Trip Date/Time"`
	DestinationAddr string `json:"Destination Addr"`
	MaxPassenger    int    `json:"Max Passenger"`
	PassengerID     int    `json:"PassengerID"`
}

/*
	type Passengers struct {
		Passengers map[string]Passenger `json:"Passenger"`
	}
*/
type Passengers struct {
	Passengers map[string]Passenger `json:"Passenger"`
}

type CarOwners struct {
	CarOwners map[string]CarOwner `json:"Carowner"`
}

type Trips struct {
	Trips map[string]Trip `json:"Trip"`
}
//menu display for client service
func main() {

	var keyin int
	var quitprog = 1

	scanner := bufio.NewScanner(os.Stdin)

	for quitprog == 1 {
		mainmenu()
		// fmt.Scan(&keyin)
		scanner.Scan()
		keyin, _ = strconv.Atoi(scanner.Text())

		switch keyin {
		case 1: // Passenger
			passengermenu()
		case 2: // Car
			carmenu()
		case 3: // Trip
			tripmenu()
		case 9: // Return
			//  Quit
			quitprog = 0
		default:
			fmt.Println("### Invalid Input ###")
		}

	}
	os.Exit(0)
}
// function to list all passengers in the database
func listAllPassenger() {
	//client := &http.Client{}
	resp, err := http.Get("http://localhost:5000/api/v1/carpool/passenger")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//var res Passengers
	var passengers Passengers

	err1 := json.Unmarshal(body, &passengers)
	if err != nil {
		fmt.Println(err1)
	}

	for k, v := range passengers.Passengers {
		fmt.Println("(", k, ") First Name: ", v.FirstName)
		fmt.Println("       Last Name: ", v.LastName)
		fmt.Println("       Phone: ", v.PhoneNo)
		fmt.Println("       Email: ", v.Email)
		fmt.Println()
	}

}

// function to create a new passenger
func createPassenger() {
	var passenger Passenger
	var passengerID string
	var updateCar bool = false
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Passenger ID to be created: ")
	scanner.Scan()
	passengerID = scanner.Text()
	passenger.PassengerID, _ = strconv.Atoi(passengerID)

	fmt.Print("First name: ")
	scanner.Scan()
	passenger.FirstName = scanner.Text()
	fmt.Print("Last Name: ")
	scanner.Scan()
	passenger.LastName = scanner.Text()
	fmt.Print("Phone no: ")
	scanner.Scan()
	passenger.PhoneNo = scanner.Text()
	fmt.Print("Email: ")
	scanner.Scan()
	passenger.Email = scanner.Text()
	fmt.Print("Car owner (Y/N): ")
	scanner.Scan()
	var yesNo string = scanner.Text()
	if yesNo == "y" {
		passenger.CarOwner = 1
		fmt.Print("Update owner Info(Y/N): ")
		scanner.Scan()
		var toUpdate string = scanner.Text()
		if toUpdate == "y" {
			updateCar = true
		} else {
			updateCar = false
		}
	} else {
		passenger.CarOwner = 0
	}
	postBody, _ := json.Marshal(passenger)
	resBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/carpool/passenger/"+passengerID, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Passenger", passengerID, "created successfully")
			} else if res.StatusCode == 409 {
				fmt.Println("Error - course", passengerID, "exists")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
	if updateCar == true {

	}
}

/*
func UpdatedCarOwner() {
// Update CarOWner
fmt.Print("LicenseNo: ")
scanner.Scan()
carowner.LicenseNo = scanner.Text()
fmt.Print("Car Plate No: ")
scanner.Scan()
carowner.CarPlateNo = scanner.Text()
fmt.Print("Max passenger: ")
scanner.Scan()
carowner.Maxpassenger, _ = strconv.Atoi(scanner.Text())
fmt.Print("Car color: ")
scanner.Scan()
carowner.Carcolor = scanner.Text()
fmt.Print("Year Release: ")
scanner.Scan()
carowner.YearRel, _ = strconv.Atoi(scanner.Text())
carowner.PassengerID = passenger.PassengerID

carpostBody, _ := json.Marshal(carowner)
carresBody := bytes.NewBuffer(carpostBody)

client := &http.Client{}

	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/carpool/car/"+passengerID, carresBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("CarOwner", passengerID, "created Car Owner successfully")
			}
		}
	}

}
*/

//function to update a passenger record
func updatePassenger() {
	var passenger Passenger
	var passengerID string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Passenger ID to be created: ")
	scanner.Scan()
	passengerID = scanner.Text()
	passenger.PassengerID, _ = strconv.Atoi(passengerID)
	fmt.Print("First name: ")
	scanner.Scan()
	passenger.FirstName = scanner.Text()
	fmt.Print("Last Name: ")
	scanner.Scan()
	passenger.LastName = scanner.Text()
	fmt.Print("Phone no: ")
	scanner.Scan()
	passenger.PhoneNo = scanner.Text()
	fmt.Print("Email: ")
	scanner.Scan()
	passenger.Email = scanner.Text()

	postBody, _ := json.Marshal(passenger)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/carpool/passenger/"+passengerID, bytes.NewBuffer(postBody)); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Passenger ", passengerID, "updated successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Passenger ", passengerID, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

//function to delete a passenger record
func deletePassenger() {
	var passengerID string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the ID of the Passenger to be deleted: ")
	scanner.Scan()
	passengerID = scanner.Text()

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5000/api/v1/carpool/passenger/"+passengerID, nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 200 {
				fmt.Println("Passenger", passengerID, "deleted successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - passenger", passengerID, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

//main display for the client service
func mainmenu() {
	println("==========")
	println("Carpool Server System")
	println("  1.  Passenger")
	println("  2.  Car")
	println("  3.  Trip")
	println("  9.  Quit")
	print("Enter an option: ")

}

//menu display for the passenger service
func passengermenu() {
	var keyin int
	scanner := bufio.NewScanner(os.Stdin)
loop:
	for {
		println("=================================")
		println("Carpool Server System")
		println("======== Passenger ==============")
		println("  1.  List all passenger")
		println("  2.  Create new passenger")
		println("  3.  Update passenger profile")
		println("  4.  Delete passenger")
		println("  9.  Back to Mainmenu")
		print("Enter an option: ")

		scanner.Scan()
		keyin, _ = strconv.Atoi(scanner.Text())
		switch keyin {
		case 1:
			listAllPassenger()
		case 2:
			createPassenger()
		case 3:
			updatePassenger()
		case 4:
			deletePassenger()
		case 9:
			//  Quit
			break loop
		default:
			fmt.Println("### Invalid Input ###")
		}
	}
}

// menu display for the car owner service
func carmenu() {
	var keyin int
	scanner := bufio.NewScanner(os.Stdin)
loop:
	for {
		println("=================================")
		println("Carpool Server System")
		println("======= Car Detail ==============")
		println("  1.  List all car")
		println("  2.  Create new car")
		println("  3.  Update car Info")
		println("  4.  Delete car")
		println("  9.  Back to Mainmenu")
		print("Enter an option: ")

		scanner.Scan()
		keyin, _ = strconv.Atoi(scanner.Text())
		switch keyin {
		case 1:
			listAllCar()
		case 2:
			createCar()
		case 3:
			updateCar()
		case 4:
			deleteCar()
		case 9:
			//  Quit
			break loop
		default:
			fmt.Println("### Invalid Input ###")
		}
	}
}

//menu display for the trip service
func tripmenu() {
	var keyin int
	scanner := bufio.NewScanner(os.Stdin)
loop:
	for {
		println("=================================")
		println("Carpool Server System")
		println("====== Trip Detail ==============")
		println("  1.  List all Trip")
		println("  2.  Create new Trip")
		println("  3.  Update Trip")
		println("  4.  Cancel Trip")
		println("  9.  Back to Mainmenu")
		print("Enter an option: ")
		scanner.Scan()
		keyin, _ = strconv.Atoi(scanner.Text())
		switch keyin {
		case 1:
			listAllTrip()
		case 2:
			createTrip()
		case 3:
			updateTrip()
		case 4:
			deleteTrip()
		case 9:
			//  Quit
			break loop
		default:
			fmt.Println("### Invalid Input ###")
		}
	}

}

//function to list all car owners
func listAllCar() {
	//client := &http.Client{}
	resp, err := http.Get("http://localhost:5000/api/v1/carpool/car")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//var res Passengers
	var carowners CarOwners

	err1 := json.Unmarshal(body, &carowners)
	if err != nil {
		fmt.Println(err1)
	}

	fmt.Println("============ Car Information =========")
	for k, v := range carowners.CarOwners {
		fmt.Println("(", k, ") Owner ID:", v.PassengerID)
		fmt.Println("      Car PlateNo:", v.CarPlateNo)
		fmt.Println("      License No:", v.LicenseNo)
		fmt.Println("      Max Capacity:", v.MaxPassenger)
		fmt.Println("      Car color:", v.CarColor)
		fmt.Println("      Year Release:", v.YearRel)
		fmt.Println()
	}

}

//function to add a new car owner record
func createCar() {
	var carowner CarOwner
	var passengerID string
	var updateCar bool = false
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Passenger ID to be created: ")
	scanner.Scan()
	passengerID = scanner.Text()
	carowner.PassengerID, _ = strconv.Atoi(passengerID)

	fmt.Print("CarPlateNo: ")
	scanner.Scan()
	var carPlateNo = scanner.Text()
	carowner.CarPlateNo = carPlateNo
	fmt.Print("LicenseNo: ")
	scanner.Scan()
	carowner.LicenseNo = scanner.Text()
	fmt.Print("Maxpassenger: ")
	scanner.Scan()
	carowner.MaxPassenger, _ = strconv.Atoi(scanner.Text())
	fmt.Print("Carcolor: ")
	scanner.Scan()
	carowner.CarColor = scanner.Text()
	fmt.Print("YearRel: ")
	scanner.Scan()
	carowner.YearRel, _ = strconv.Atoi(scanner.Text())

	postBody, _ := json.Marshal(carowner)
	resBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/carpool/car/"+carPlateNo, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Car for Passenger", carPlateNo, "created successfully")
			} else if res.StatusCode == 409 {
				fmt.Println("Error - Car ", carPlateNo, "exists")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
	if updateCar == true {

	}
}

//function to update a car owner record
func updateCar() {
	var carowner CarOwner
	var carplateno string
	var passengerID string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Passenger ID to be Updated: ")
	scanner.Scan()
	passengerID = scanner.Text()
	carowner.PassengerID, _ = strconv.Atoi(passengerID)
	fmt.Print("Car Plate No: to be updated: ")
	scanner.Scan()
	carplateno = scanner.Text()
	carowner.CarPlateNo = carplateno
	fmt.Print("License No: ")
	scanner.Scan()
	carowner.LicenseNo = scanner.Text()
	fmt.Print("Max Capacity: ")
	scanner.Scan()
	carowner.MaxPassenger, _ = strconv.Atoi(scanner.Text())
	fmt.Print("YearRel: ")
	scanner.Scan()
	carowner.YearRel, _ = strconv.Atoi(scanner.Text())
	fmt.Print("Car Color: ")
	scanner.Scan()
	carowner.CarColor = scanner.Text()

	postBody, _ := json.Marshal(carowner)
	resBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/carpool/car/"+carplateno, resBody); err == nil {
		fmt.Println("req")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Car Plate No.: ", carplateno, "updated successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Car Plate No.: ", carplateno, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

//function to delete a car owner record
func deleteCar() {
	var carplateno string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the Car Plate No.to be deleted: ")
	scanner.Scan()
	carplateno = scanner.Text()

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5000/api/v1/carpool/car/"+carplateno, nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 200 {
				fmt.Println("Car Plate No.", carplateno, "deleted successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Car Plate No.: ", carplateno, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

//function to list all trip records
func listAllTrip() {
	//client := &http.Client{}
	resp, err := http.Get("http://localhost:5000/api/v1/carpool/trip")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//var res Passengers
	var trips Trips

	err1 := json.Unmarshal(body, &trips)
	if err != nil {
		fmt.Println(err1)
	}
	fmt.Println("=========Trip Information:=========")
	for k, v := range trips.Trips {
		fmt.Println("(", k, ") Passenger ID", v.PassengerID)
		fmt.Println("      Trip Ref no:", v.TripRef)
		fmt.Println("      Pickup Addr:", v.PickupAddr)
		fmt.Println("      Alter Pickup Addr:", v.AlterPickupAddr)
		fmt.Println("      Start Trip:", v.StartTrip)
		fmt.Println("      Destination Addr:", v.DestinationAddr)
		fmt.Println("      Number of Passenger:", v.MaxPassenger)
		fmt.Println()
	}
}

//function to create a trip record
func createTrip() {
	var trip Trip
	var passengerID string
	var updateTrip bool = false
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Passenger ID to be created: ")
	scanner.Scan()
	passengerID = scanner.Text()
	trip.PassengerID, _ = strconv.Atoi(passengerID)
	/*
		fmt.Print("Enter the Trip ID to be created: ")
		scanner.Scan()
		tripID = scanner.Text()
		trip.TripRef, _ = strconv.Atoi(tripID)
	*/
	fmt.Print("Start Trip Date/Time: (yyyy-mm-dd hh:mm:ss) ")
	scanner.Scan()
	trip.StartTrip = scanner.Text()
	fmt.Print("Pickup Addr: ")
	scanner.Scan()
	trip.PickupAddr = scanner.Text()
	fmt.Print("Alt Pickup Addr: ")
	scanner.Scan()
	trip.AlterPickupAddr = scanner.Text()
	fmt.Print("Maxpassenger: ")
	scanner.Scan()
	trip.MaxPassenger, _ = strconv.Atoi(scanner.Text())
	fmt.Print("Destination Addr: ")
	scanner.Scan()
	trip.DestinationAddr = scanner.Text()

	postBody, _ := json.Marshal(trip)
	resBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/api/v1/carpool/trip/"+passengerID, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Passenger ID", passengerID, "created successfully")
			} else if res.StatusCode == 409 {
				fmt.Println("Error - passengerID", passengerID, "Post Error")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
	if updateTrip == true {

	}
}

//function to update a trip record
func updateTrip() {
	var trip Trip
	var tripref string
	var passengerID string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the Trip Ref to be update: ")
	scanner.Scan()
	tripref = scanner.Text()
	trip.TripRef, _ = strconv.Atoi(tripref)
	fmt.Print("Pickup Addr: ")
	scanner.Scan()
	trip.PickupAddr = scanner.Text()
	fmt.Print("Alt Pickup Addr: ")
	scanner.Scan()
	trip.AlterPickupAddr = scanner.Text()
	fmt.Print("Start Trip Date/Time: (yyyy-mm-dd hh:mm:ss) ")
	scanner.Scan()
	trip.StartTrip = scanner.Text()
	fmt.Print("Destination Addr: ")
	scanner.Scan()
	trip.DestinationAddr = scanner.Text()
	fmt.Print("Maxpassenger: ")
	scanner.Scan()
	trip.MaxPassenger, _ = strconv.Atoi(scanner.Text())
	fmt.Print("Passenger ID: ")
	scanner.Scan()
	passengerID = scanner.Text()
	trip.PassengerID, _ = strconv.Atoi(passengerID)

	postBody, _ := json.Marshal(trip)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, "http://localhost:5000/api/v1/carpool/trip/"+tripref, bytes.NewBuffer(postBody)); err == nil {
		fmt.Println("tripref")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				fmt.Println("Trip No.: ", tripref, "updated successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Trip No.: ", tripref, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}

//function to delete a trip record
func deleteTrip() {
	var tripref string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the Trip Ref to be deleted: ")
	scanner.Scan()
	tripref = scanner.Text()

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5000/api/v1/carpool/trip/"+tripref, nil); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 200 {
				fmt.Println("Trip Reference No.", tripref, "deleted successfully")
			} else if res.StatusCode == 404 {
				fmt.Println("Error - Trip Reference No.: ", tripref, "does not exist")
			}
		} else {
			fmt.Println(2, err)
		}
	} else {
		fmt.Println(3, err)
	}
}
