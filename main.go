package main

import (
	"Assignment_2_CST_412/firestore"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

//Main method that runs the app.
func main() {

	//Initiate the DB with the private key.
	firestore.InitDB("firebase_key.json")
	defer firestore.CloseDB()

	//We use whatever port the OS provides us.
	port := os.Getenv("PORT")
	fmt.Println(":" + port)

	//Handlefunc tells the app that when anyone enters the default page the "mainPage" method will run.
	http.HandleFunc("/", mainPage)

	log.Println("Running on port :" + port)

	//Listen on chosen port and handle error if present.
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err)
	}
}

//Main page that runs each time someone accesses the webapp.
func mainPage(w http.ResponseWriter, r *http.Request) {

	//Parse the initial HTML template.
	tmpl := template.Must(template.ParseFiles("html_files\\layout_1.html"))

	//Execute the first template with "nil" data (null). This is what will make the html file display in the browser.
	err := tmpl.Execute(w, nil)
	if err != nil {
		return
	}
	if err != nil {
		return
	}

	//Check if the "Add New" button is pressed. If so initiate the layout_3 template and execute it.
	if r.FormValue("newAsset") == "newAsset" {
		tmpl = template.Must(template.ParseFiles("html_files\\layout_3.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			return
		}
		return
	}

	//If the Add button is pressed the AddAsset method will run, adding the new asset to the DB.
	if r.FormValue("createAsset") == "createAsset" {
		firestore.AddAsset(w, r)
		return
	}

	//If the update button is pressed, the UpdateAsset method will run, updating the specific asset with any new data.
	if r.FormValue("updateAsset") == "updateAsset" {
		firestore.UpdateAsset(w, r)
		return
	}

	//If the request method simply is a POST method the FindAsset method will run. This is because we can't get the value of the Search button.
	if r.Method == "POST" {
		firestore.FindAsset(w, r.FormValue("asset"), r)
		return
	}
}
