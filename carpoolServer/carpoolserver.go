package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"

	//"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//"golang.org/x/text/date"
)
//struct for the passenger object
type Passenger struct {
	PassengerID int    `json:"PassengerID"`
	FirstName   string `json:"First Name"`
	LastName    string `json:"Last Name"`
	PhoneNo     string `json:"Phone No"`
	Email       string `json:"Email"`
	CarOwner    int    `json:"Car Owner"`
}

//struct for the car owner object
type CarOwner struct {
	PassengerID  int    `json:"PassengerID"`
	CarPlateNo   string `json:"Car Plate No"`
	LicenseNo    string `json:"License No"`
	MaxPassenger int    `json:"Max Passenger"`
	CarColor     string `json:"Car Color"`
	YearRel      int    `json:"Year Release"`
}

//struct for the trip object
type Trip struct {
	TripRef         int    `json:"Trip Reference"`
	PickupAddr      string `json:"Pickup Addr"`
	AlterPickupAddr string `json:"Alter Pick Addr"`
	StartTrip       string `json:"Start Trip Date/Time"`
	DestinationAddr string `json:"Destination Addr"`
	MaxPassenger    int    `json:"Max Passenger"`
	PassengerID     int    `json:"PassengerID"`
}

//struct for the trip-passenger object (WIP)
type TripPassenger struct {
	PassengerID     int    `json:"PassengerID"`
	FirstName       string `json:"First Name"`
	LastName        string `json:"lastName"`
	PhoneNo         string `json:"Phone No"`
	Email           string `json:"Email"`
	PickupAddr      string `json:"Pickup Addr"`
	AlterPickupAddr string `json:"Alter Pick Addr"`
	StartTrip       string `json:"Start Trip Date/Time"`
	DestinationAddr string `json:"Destination Addr"`
	MaxPassenger    int    `json:"Max Passenger"`
	TripRef         int    `json:"Trip Reference"`
}

var (
	db  *sql.DB
	err error
)

//main server functions that handles all the microservices
func main() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")

	if err != nil {
		panic(err.Error())
	}
	// defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/carpool/passenger/{passengerID}", passengerdata).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	router.HandleFunc("/api/v1/carpool/passenger", allpassenger)
	router.HandleFunc("/api/v1/carpool/car", allcarowner)
	router.HandleFunc("/api/v1/carpool/car/{carPlateNo}", carowner)
	router.HandleFunc("/api/v1/carpool/trip", alltrip)
	router.HandleFunc("/api/v1/carpool/trip/{tripRef}", tripinfo)

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

//function that handles the passenger data, with methods to add, create, update and delete records.
func passengerdata(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Update passenger data by passenger ID
	if r.Method == "POST" {
		println("POST Passenger")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data Passenger
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println(params["passengerID"])
				fmt.Println(reflect.TypeOf(params["passengerID"]))
				if _, ok := isExist(params["passengerID"]); !ok {
					fmt.Println(data)
					//courses[params["courseid"]] = data
					insertPassenger(params["passengerID"], data)

					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusConflict)
					fmt.Fprintf(w, "passenger ID exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" {
		// Update passenger data by passenger ID to replace the entire resource with a new representation
		println("PUT Passenger")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data Passenger

			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isExist(params["passengerID"]); ok {
					fmt.Println(data)
					//courses[params["courseid"]] = data
					updatePassenger(params["passengerID"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "passenger ID does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PATCH" {
		println("PATCH Passenger")
		// Patch passenger by passenger ID to apply partial updates to a resource

		if body, err := io.ReadAll(r.Body); err == nil {
			var data map[string]interface{}

			if err := json.Unmarshal(body, &data); err == nil {
				fmt.Println("Json")
				if orig, ok := isExist(params["passengerID"]); ok {
					fmt.Println(data)

					for k, v := range data {
						switch k {
						case "First Name":
							orig.FirstName = v.(string)
						case "Last Name":
							orig.LastName = v.(string)
						case "Phone No":
							orig.PhoneNo = v.(string)
						case "Email":
							orig.Email = v.(string)
						}
					}
					//courses[params["courseid"]] = orig
					updatePassenger(params["passengerId"], orig)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "passenger ID does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "DELETE" {
		if _, err := io.ReadAll(r.Body); err == nil {
			//var data Passenger
			println("DELETE Passenger func")
			//if err := json.Unmarshal(body, &data); err == nil {
			println(params["passengerID"])
			if _, ok := isExist(params["passengerID"]); ok {
				println("DELETE Passenger")
				// Delete passenger by ID

				fmt.Fprintf(w, params["passengerID"]+" Deleted")
				//delete(courses, params["courseid"])
				delPassenger(params["passengerID"])
				// w.WriteHeader(http.StatusAccepted)

			} else {
				// input error
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Invalid passenger ID")
			}
			//}
		} else {
			// input error
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Invalid passenger ID")
		}
	}
}

// Get all passenger information from database
func allpassenger(w http.ResponseWriter, r *http.Request) {

	println("Query")
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")
	results, err := db.Query("select * from passenger")
	if err != nil {
		panic(err.Error())
	}

	var count int = 0

	//passenger := make(map[int]Passenger)
	var passengers map[int]Passenger = map[int]Passenger{}
	for results.Next() {
		var c Passenger

		_ = results.Scan(&c.PassengerID, &c.FirstName, &c.LastName, &c.PhoneNo, &c.Email, &c.CarOwner)

		count += 1
		passengers[count] = Passenger{PassengerID: c.PassengerID, FirstName: c.FirstName, LastName: c.LastName, PhoneNo: c.PhoneNo, Email: c.Email, CarOwner: c.CarOwner}
	}

	fmt.Println(count, " Records found")
	passengerWrapper := struct {
		Passengers map[int]Passenger `json:"Passenger"`
	}{passengers}

	jsonBytes, err := json.Marshal(passengerWrapper)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))

	json.NewEncoder(w).Encode(passengerWrapper)

}

func getPassenger() map[int]Passenger {
	results, err := db.Query("select * from passenger")
	if err != nil {
		panic(err.Error())
	}

	var passengers map[int]Passenger = map[int]Passenger{}

	for results.Next() {
		var c Passenger

		err = results.Scan(&c.PassengerID, &c.FirstName, &c.LastName, &c.PhoneNo, &c.Email, &c.CarOwner)
		if err != nil {
			panic(err.Error())
		}

		passengers[c.PassengerID] = c
		println(passengers[c.PassengerID].FirstName + " " + passengers[c.PassengerID].LastName)
	}

	return passengers
}

//function checks if the passenger already exists - used in CREATE and UPDATE functions
func isExist(id string) (Passenger, bool) {
	var c Passenger
	println("isExist")

	//db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")

	var pid, _ = strconv.Atoi(id)
	println(pid)
	result := db.QueryRow("select * from passenger where passengerID=?", pid)
	fmt.Println("QueryRow")
	err := result.Scan(&id, &c.FirstName, &c.LastName, &c.PhoneNo, &c.Email, &c.CarOwner)
	//err := result.Scan(&c.PassengerID)
	println("id:", c.PassengerID)
	//println("First Name: ", c.FirstName)
	if err == sql.ErrNoRows {
		fmt.Println("Found!")
		return c, false

	}
	fmt.Println("Not Found!")
	return c, true
}

//function to delete a passenger from the server
func delPassenger(id string) (int64, error) {
	fmt.Println("DeletePassenger func")
	var pid, _ = strconv.Atoi(id)
	fmt.Println("pid= ", pid)
	result, err := db.Exec("delete from passenger where passengerID=?", pid)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//function to insert a passenger into the server database
func insertPassenger(id string, c Passenger) {
	println("Insert Passenger")
	fmt.Println(" c.PassengerID= ", c.PassengerID)
	fmt.Println(" c.FirstName= ", c.FirstName)
	fmt.Println(" c.LastName= ", c.LastName)
	fmt.Println(" c.PhoneNo= ", c.PhoneNo)
	fmt.Println(" c.Email= ", c.Email)
	fmt.Println(" c.CarOwner= ", c.CarOwner)

	_, err := db.Exec("insert into passenger values(?,?,?,?,?,?)", c.PassengerID, c.FirstName, c.LastName, c.PhoneNo, c.Email, c.CarOwner)
	if err != nil {
		panic(err.Error())
	}
}

//function to update a passenger record in the server database
func updatePassenger(id string, c Passenger) {
	var pid, _ = strconv.Atoi(id)
	fmt.Println("Exec", pid)
	_, err := db.Exec("update passenger set passengerID=?, firstName=?, lastName=?, phoneNo=?, email=?, carOwner=?  where passengerID=?", pid, c.FirstName, c.LastName, c.PhoneNo, c.Email, c.CarOwner, pid)
	if err != nil {
		panic(err.Error())
	}
}

//function to send a passenger query to the server
func queryPassenger(query string) (map[int]Passenger, bool) {
	//results, err := db.Query(curl http://localhost:5000/api/v1/carpool/passenger?q="select"+"*"+"from"+"passenger"+"where"+"firstName="+"\"Goh\""m passenger where lower(firstName) like lower(?) or lower(lastName) Like lower(?)", "%"+query+"%",  "%"+query+"%")
	// curl http://localhost:5000/api/v1/carpool/passenger?q="Select"+"*"+"from"+"passenger"
	// curl http://localhost:5000/api/v1/carpool/passenger?q="Select%20*%20from%20passenger"
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")
	results, err := db.Query(query)
	println("queryPassenger routine: ", &results)

	if err != nil {
		panic(err.Error())
	}

	var passengers map[int]Passenger = map[int]Passenger{}

	for results.Next() {
		var c Passenger
		//var id string
		err = results.Scan(&c.PassengerID, &c.FirstName, &c.LastName, &c.PhoneNo, &c.Email, &c.CarOwner)
		if err != nil {
			panic(err.Error())
		}
		println("PassengerID: " + strconv.Itoa(c.PassengerID))
		passengers[c.PassengerID] = c
	}

	if len(passengers) == 0 {
		return passengers, false
	}
	return passengers, true
}

//function to find eligible passengers that fit the query
func findEligiblePassenger(name string) (map[int]Passenger, bool) {
	results, err := db.Query("select * from passenger where firstName like ?", name)
	if err != nil {
		panic(err.Error())
	}

	var passenger map[int]Passenger = map[int]Passenger{}

	for results.Next() {
		var c Passenger
		//var id string
		err = results.Scan(&c.PassengerID, &c.FirstName, &c.LastName, &c.PhoneNo, &c.Email, &c.CarOwner)
		if err != nil {
			panic(err.Error())
		}

		passenger[c.PassengerID] = c
	}

	if len(passenger) == 0 {
		return passenger, false
	}
	return passenger, true
}

//function prints all car owners
func allcarowner(w http.ResponseWriter, r *http.Request) {
	println("Allcarowner")
	// curl http://localhost:5000/api/v1/carpool/car?q="Select%20*%20from%20carowner"
	// query := r.URL.Query()

	println("Query")
	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")
	results, err := db.Query("select * from carOwner")
	if err != nil {
		panic(err.Error())
	}

	var count int = 0
	var carowners map[int]CarOwner = map[int]CarOwner{}
	//var courses[] Course
	for results.Next() {
		var c CarOwner
		//var id string
		_ = results.Scan(&c.PassengerID, &c.LicenseNo, &c.CarPlateNo, &c.MaxPassenger, &c.CarColor, &c.YearRel)
		count += 1
		carowners[count] = c
		println("count= ", count)
		println("c.= ", c.PassengerID)
		println(carowners[count].PassengerID)
		println(carowners[count].CarPlateNo)
		println(carowners[count].LicenseNo)
		println(carowners[count].MaxPassenger)
		println(carowners[count].CarColor)
		println(carowners[count].YearRel)
	}
	fmt.Println(count, " Records found")
	carWrapper := struct {
		Carowners map[int]CarOwner `json:"CarOwner"`
	}{carowners}

	jsonBytes, err := json.Marshal(carWrapper)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))

	json.NewEncoder(w).Encode(carWrapper)

	//defer db.Close()

}

//function handles methods for car owners, including ADDING, UPDATING and DELETING records.
func carowner(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// curl http://localhost:5000/api/v1/carpool/car/S8239162Z
	// Update CarOwner data by CarPlateNo
	if r.Method == "POST" {
		println("POST Car Owner")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data CarOwner
			fmt.Println("Car Owner!")
			fmt.Println("Passenger ID: ", params["passengerID"])
			fmt.Println("CarPlateNo: ", params["carPlateNo"])
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isCarExist(params["carPlateNo"]); !ok {
					fmt.Println("Car not Exist", data)
					//courses[params["courseid"]] = data
					insertCarOwner(params["CarPlateNo"], data)

					w.WriteHeader(http.StatusAccepted)
				} else {
					fmt.Println("Car  Exist", data)
					w.WriteHeader(http.StatusConflict)

				}
			} else {
				fmt.Println("err: ", err)
			}
		}
	} else if r.Method == "PUT" {
		// Update Carowner data by CarPlateNo to replace the entire resource with a new representation
		println("PUT Car Owner")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data CarOwner
			fmt.Println("Passenger ID: ", params["passengerID"])
			fmt.Println("CarPlateNo: ", params["carPlateNo"])
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isCarExist(params["carPlateNo"]); ok {
					fmt.Println("Car Plate: ", data.CarPlateNo)
					//courses[params["courseid"]] = data
					updateCarOwner(params["carPlateNo"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "CarPlateNo does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PATCH" {
		// Patch CarOwner by CarPlateNo to apply partial updates to a resource
		println("PATCH Car Owner")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data map[string]interface{}

			if err := json.Unmarshal(body, &data); err == nil {
				if orig, ok := isCarExist(params["carPlateNo"]); ok {
					fmt.Println(data)

					for k, v := range data {
						switch k {
						case "License No":
							orig.LicenseNo = v.(string)
						case "Car Plate No":
							orig.CarPlateNo = v.(string)
						case "OwnerID":
							orig.MaxPassenger = v.(int)
						case "Passenger ID":
							orig.PassengerID = v.(int)
						case "Max Passenger":
							orig.MaxPassenger = v.(int)
						case "Car Color":
							orig.CarColor = v.(string)
						}
					}
					//courses[params["courseid"]] = orig
					updateCarOwner(params["carPlateNo"], orig)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "CarPlateNo does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if val, ok := isCarExist(params["carPlateNo"]); ok {
		// Delete passenger by ID
		println("Get Car Owner")
		if r.Method == "DELETE" {
			fmt.Fprintf(w, params["carPlateNo"]+" Deleted")
			//delete(courses, params["courseid"])
			delCarOwner(params["carPlateNo"])
		} else {
			json.NewEncoder(w).Encode(val)
		}
	} else {
		// input error
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid carplateno")
	}
}

//function checks if the car owner record exists - used for CREATING and UPDATING
func isCarExist(id string) (CarOwner, bool) {
	var c CarOwner
	fmt.Println("isCarExist")
	fmt.Println("Car Plate No: ", id)
	result := db.QueryRow("select * from CarOwner where carPlateNo=?", id)
	err := result.Scan(&c.PassengerID, &c.CarPlateNo, &c.LicenseNo, &c.MaxPassenger, &c.CarColor, &c.YearRel)
	if err == sql.ErrNoRows {
		println("Error!")
		return c, false
	}

	return c, true
}

//function deletes the specified car owner record from the database.
func delCarOwner(id string) (int64, error) {
	result, err := db.Exec("delete from carowner where carPlateNo=?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//function inserts values to the query for adding another car owner entry
func insertCarOwner(id string, c CarOwner) {
	println("InsertCarOwner")
	_, err := db.Exec("insert into carowner values(?,?,?,?,?,?)", c.PassengerID, c.CarPlateNo, c.LicenseNo, c.MaxPassenger, c.CarColor, c.YearRel)
	if err != nil {
		panic(err.Error())
	}
}

//function updates a specific car owner entry
func updateCarOwner(id string, c CarOwner) {
	_, err := db.Exec("update carowner set passengerID=?, carPlateNo=?, licenseNo=?, maxPassenger=?, carColor=?, yearRel=? where carPlateNo=?", c.PassengerID, c.CarPlateNo, c.LicenseNo, c.MaxPassenger, c.CarColor, c.YearRel, c.CarPlateNo)
	if err != nil {
		panic(err.Error())
	}
}

// SELECT * FROM trip inner join passenger on trip.passengerID = passenger.passengerID;
//SELECT * FROM trip inner join passenger on trip.passengerID = passenger.passengerID Join Carowner on passenger.passengerID = carowner.passengerID;
// Get all trip information

/*
func alltrip(w http.ResponseWriter, r *http.Request) {
	println("Alltrip")
	// query := r.URL.Query()
	// curl http://localhost:5000/api/v1/carpool/trip
	println("Query")


	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")
	// results, err := db.Query("SELECT * FROM trip inner join passenger on trip.passengerID = passenger.passengerID Join Carowner on passenger.passengerID = carowner.passengerID;")
	results, err := db.Query("SELECT * FROM trip inner join passenger on trip.passengerID = passenger.passengerID")
	if err != nil {
		panic(err.Error())
	}

	var count int = 0;
	var passengers []Passenger
	var trips []Trip

	// var carowners map[int]CarOwner = map[int]CarOwner{}
	//var courses[] Course
	for results.Next() {
		var t Trip
		var p Passenger
		//var id string
		_ = results.Scan(&p.PassengerID, &p.FirstName, &p.LastName, &p.PhoneNo, &p.Email, &t.PassengerID, &t.PickupAddr, &t.AlterPickupAddr, &t.DestinationAddr, &t.MaxPassenger, &t.StartTrip, &t.TripRef)
		count += 1
		fmt.Println(p.FirstName)
		passengers = append(passengers, p)
		trips = append(trips, t)

		println("count= ", count)
		println("p.passengers[count-1].PassengerID= ", passengers[count-1].PassengerID)
		println("p.passengers[count-1].firstname= ", passengers[count-1].FirstName)
		println("p.passengers[count-1].PassengerID= ", passengers[count-1].LastName)
		println("p.passengers[count-1].PassengerID= ", passengers[count-1].PhoneNo)
		println("p.passengers[count-1].PassengerID= ", passengers[count-1].Email)
		println("trips[count-1].PassengerID= ", trips[count-1].PassengerID)
		println("trips[count-1].PickupAddr= ", trips[count-1].PickupAddr)
		println("trips[count-1].AlterPickupAddr= ", trips[count-1].AlterPickupAddr)
		println("trips[count-1].DestinationAddr= ", trips[count-1].DestinationAddr)
		println("trips[count-1].MaxPassenger= ", trips[count-1].MaxPassenger)
		fmt.Println("trips[count-1].StartTrip= ", trips[count-1].StartTrip)
		println("trips[count-1].TripRef= ", trips[count-1].TripRef)
	}
	fmt.Println(count, " Records found")
	passengerWrapper := struct {
		Passengers map[int]Passenger `json:"Passenger"`
	}{getPassenger()}

	jsonBytes, err := json.Marshal(passengerWrapper)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))

	defer db.Close()

}

*/

/*
func alltrip(w http.ResponseWriter, r *http.Request) {
	println("Alltrip")
	//query := r.URL.Query()
	// curl http://localhost:5000/api/v1/carpool/trip
	println("Query")

	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")
	// results, err := db.Query("SELECT * FROM trip inner join passenger on trip.passengerID = passenger.passengerID Join Carowner on passenger.passengerID = carowner.passengerID;")
	results, err := db.Query("SELECT * FROM trip inner join passenger on trip.passengerID = passenger.passengerID")
	if err != nil {
		panic(err.Error())
	}

	var count int = 0
	var trippassengers []TripPassenger

	for results.Next() {
		var tp TripPassenger
		//var id string
		//_ = results.Scan(&tp.PassengerID, &tp.FirstName, &tp.LastName, &tp.PhoneNo, &tp.Email, &tp.PickupAddr, &tp.AlterPickupAddr, &tp.DestinationAddr, &tp.StartTrip, &tp.MaxPassenger, &tp.TripRef)
		_ = results.Scan(&tp.TripRef, &tp.PickupAddr, &tp.PassengerID)

		println("tp.TripRef= ",tp.TripRef)
		println("tp.TripRef= ",tp.PickupAddr)
		println("tp.PassengerID= ",tp.PassengerID)
		println("tp.FirstName= ", tp.FirstName)
		println("tp.LastName= ", tp.LastName)


		trippassengers = append(trippassengers, tp)

		println("count= ", count)
		println("tp.passengers[count].PassengerID= ", trippassengers[count].PassengerID)
		println("tp.passengers[count].firstname= ", trippassengers[count].FirstName)
		println("tp.passengers[count].PassengerID= ", trippassengers[count].LastName)
		println("tp.passengers[count].PassengerID= ", trippassengers[count].PhoneNo)
		println("tp.passengers[count-1].PassengerID= ", trippassengers[count].Email)
		println("trippassengers[count-1].AlterPickupAddr= ", trippassengers[count].AlterPickupAddr)
		println("trippassengers[count-1].DestinationAddr= ", trippassengers[count].DestinationAddr)
		println("trippassengers[count-1].PickupAddr= ", trippassengers[count].PickupAddr)
		fmt.Println("trippassengers[count-1].StartTrip= ", trippassengers[count].StartTrip)
		println("trippassengers[count-1].MaxPassenger= ", trippassengers[count].MaxPassenger)
		println("trippassengers[count-1].TripRef= ", trippassengers[count].TripRef)
		count += 1
	}
	fmt.Println(count, " Records found")
	// carWrapper := struct {
	// 	Carowners map[int]CarOwner `json:"CarOwner"`
	// }{carowners}

	jsonBytes, err := json.Marshal(trippassengers)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))

	//json.NewEncoder(w).Encode(carWrapper)

	defer db.Close()

}

*/

//function displays all trip instances in the database
func alltrip(w http.ResponseWriter, r *http.Request) {
	println("Alltrip")
	// query := r.URL.Query()
	// curl http://localhost:5000/api/v1/carpool/trip
	println("Query")

	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/carpool")
	// results, err := db.Query("SELECT * FROM trip inner join passenger on trip.passengerID = passenger.passengerID Join Carowner on passenger.passengerID = carowner.passengerID;")
	results, err := db.Query("SELECT * FROM trip")
	if err != nil {
		panic(err.Error())
	}

	var count int = 0
	//var passengers []Passenger
	// var trips []Trip
	var trips map[int]Trip = map[int]Trip{}
	// var carowners map[int]CarOwner = map[int]CarOwner{}
	//var courses[] Course
	for results.Next() {
		var t Trip

		//var p Passenger
		//var id string
		//_ = results.Scan(&t.TripRef, &t.PickupAddr, &t.AlterPickupAddr, &t.StartTrip, &t.DestinationAddr, &t.MaxPassenger, &t.PassengerID)
		_ = results.Scan(&t.TripRef, &t.PickupAddr, &t.AlterPickupAddr, &t.StartTrip, &t.DestinationAddr, &t.MaxPassenger, &t.PassengerID)

		count += 1
		trips[count] = t
		fmt.Println(t.PassengerID)
		fmt.Println(trips[count].PickupAddr)
		fmt.Println(trips[count].StartTrip)
		fmt.Println(trips[count].MaxPassenger)
		//passengers = append(passengers, p)

		// trips[count] = Trip{TripRef: t.TripRef, PickupAddr: t.PickupAddr, AlterPickupAddr: t.AlterPickupAddr, StartTrip: t.StartTrip, DestinationAddr: t.DestinationAddr, MaxPassenger: t.MaxPassenger, PassengerID: t.PassengerID}
	}

	fmt.Println(count, " Records found")
	tripWrapper := struct {
		Trips map[int]Trip `json:"Trip"`
	}{trips}

	jsonBytes, err := json.Marshal(tripWrapper)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))

	json.NewEncoder(w).Encode((tripWrapper))
	//defer db.Close()

}

//Function handles all trip related methods, like CREATING, UPDATING and DELETING.
func tripinfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// curl http://localhost:5000/api/v1/carpool/trip/S8239162Z
	// Update Trip data by TripRefNo
	if r.Method == "POST" {
		println("POST trip")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data Trip
			fmt.Println("Trip!")
			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isTripExist(params["tripRef"]); !ok {
					fmt.Println("Trip does not exist", data)
					//courses[params["courseid"]] = data
					insertTrip(params["tripRef"], data)

					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusConflict)
					fmt.Fprintf(w, "Trip already exists")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PUT" {
		// Update Trip data by TripRef to replace the entire resource with a new representation
		println("PUT Trip")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data Trip

			if err := json.Unmarshal(body, &data); err == nil {
				if _, ok := isTripExist(params["tripRef"]); ok {
					fmt.Println(data)
					//courses[params["courseid"]] = data
					updateTrip(params["tripRef"], data)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Trip Referenece does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if r.Method == "PATCH" {
		// Patch Trip by TripRef to apply partial updates to a resource
		println("PATCH Trip")
		if body, err := io.ReadAll(r.Body); err == nil {
			var data map[string]interface{}

			if err := json.Unmarshal(body, &data); err == nil {
				if orig, ok := isTripExist(params["tripRef"]); ok {
					fmt.Println(data)

					for k, v := range data {
						switch k {
						case "Pickup Address":
							orig.PickupAddr = v.(string)
						case "Alter Pickup Address":
							orig.AlterPickupAddr = v.(string)
						case "Start Trip":
							orig.StartTrip = v.(string)
						case "Destination Address":
							orig.DestinationAddr = v.(string)
						case "Max Passenger":
							orig.MaxPassenger = v.(int)
						case "Passenger ID":
							orig.PassengerID = v.(int)
						}
					}
					//courses[params["courseid"]] = orig
					updateTrip(params["tripRef"], orig)
					w.WriteHeader(http.StatusAccepted)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "Trip reference does not exist")
				}
			} else {
				fmt.Println(err)
			}
		}
	} else if val, ok := isTripExist(params["tripRef"]); ok {
		// Delete passenger by ID
		println("Get Trip")
		if r.Method == "DELETE" {
			fmt.Fprintf(w, params["tripRef"]+" Deleted")
			//delete(courses, params["courseid"])
			delTrip(params["tripRef"])
		} else {
			json.NewEncoder(w).Encode(val)
		}
	} else {
		// input error
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid trip ID")
	}
}

//function checks if a specified trip instance exists in the database.
func isTripExist(id string) (Trip, bool) {
	var c Trip

	var pid, _ = strconv.Atoi(id)
	result := db.QueryRow("select * from Trip where tripRef=?", pid)
	err := result.Scan(&c.TripRef, &c.PickupAddr, &c.AlterPickupAddr, &c.StartTrip, &c.DestinationAddr, &c.MaxPassenger, &c.PassengerID)
	if err == sql.ErrNoRows {
		return c, false
	}

	return c, true
}

//function deletes a specified trip instance from the database.
func delTrip(id string) (int64, error) {
	fmt.Println("DeleteTrip func")
	var pid, _ = strconv.Atoi(id)
	fmt.Println("pid= ", pid)
	result, err := db.Exec("delete from trip where tripRef=?", pid)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//function inserts trip into the database
func insertTrip(id string, c Trip) {
	/*
		println("Insert Trip")
		fmt.Println(" c.TriRef= ", c.TripRef)
		fmt.Println(" c.pickupAddr= ", c.PickupAddr)
		fmt.Println(" c.alterPickupAddr= ", c.AlterPickupAddr)
		fmt.Println(" c.startTrip= ", c.StartTrip)
		fmt.Println(" c.destinationAddr= ", c.DestinationAddr)
		fmt.Println(" c.maxPassenger= ", c.MaxPassenger)
		fmt.Println(" c.passengerID= ", c.PassengerID)
	*/

	_, err := db.Exec("insert into trip values(?,?,?,?,?,?,?)", c.TripRef, c.PickupAddr, c.AlterPickupAddr, c.StartTrip, c.DestinationAddr, c.MaxPassenger, c.PassengerID)
	if err != nil {
		panic(err.Error())
	}
}

//function updates a specific trip instance
func updateTrip(id string, c Trip) {
	var pid, _ = strconv.Atoi(id)
	_, err := db.Exec("update trip set tripRef=?, pickupAddr=?, alterPickupAddr=?, startTrip=?, destinationAddr=?, maxPassenger=?  where tripRef=?", c.TripRef, c.PickupAddr, c.AlterPickupAddr, c.StartTrip, c.DestinationAddr, c.MaxPassenger, pid)
	if err != nil {
		panic(err.Error())
	}
}
