package firestore

import (
	"Assignment_2_CST_412/structs"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/fxtlabs/date"
	_ "github.com/fxtlabs/date"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

//Variable for interacting with the Firestore DB.
var ctx context.Context
var client *firestore.Client

// InitDB Initiate the database.
func InitDB(credPath string) {

	// Get the context.
	ctx = context.Background()

	// Get the credential file from the given path and set up a firebase app.
	sa := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println(err)
	}

	// Instantiate client
	client, err = app.Firestore(ctx)

	// Checks for error when instantiating the client.
	if err != nil {
		log.Println(err)
	}
}

//CloseDB close the DB.
func CloseDB() {
	// Close down client
	defer func() {
		err := client.Close()
		if err != nil {
			log.Println(err)
		}
	}()
}

// FindAsset Method that iterates through the DB using the value of the search field to look for a matching object.
func FindAsset(w http.ResponseWriter, search string, r *http.Request) {

	//Check if any object in the DB is retired and should be deleted.
	CheckRetired(w, r)

	// Iterator to loop through all entries in collection "assets"
	w.Header().Add("content-type", "application/json")
	iter := client.Collection("assets").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}

		// A map with string keys. Each key is one field. "m" lets us get data from the DB object.
		m := doc.Data()

		//Simply stop the search if the input is nothing.
		if search == "" {
			return
		}

		//If any of the data of the object matches the search input, display this object.
		if m["asset_id"] == search || m["manufacturer"] == search || m["manufacturer_address"] == search || m["manufacturer_phone"] == search ||
			m["manufacturer_web"] == search || m["model"] == search || m["date_purchased"] == search || m["purchase_price"] == search ||
			m["warranty_date"] == search || m["retired_date"] == search || m["description"] == search {

			//In order to iterate over the comments in the HTML template we need to set it up like this.
			var commentsFromDB []structs.Comment

			//We gather all the comments from the object into "s" and then we iterate over add all these comments to the array.
			s := reflect.ValueOf(m["comments"])
			for i := 0; i < s.Len(); i++ {
				if m["comments"] != "" {
					str := fmt.Sprintf("%v", s.Index(i))
					commentsFromDB = append(commentsFromDB, structs.Comment{Comment: str})
				}
			}

			//All the data gathered from the DB object put into a struct or an Object.
			data := structs.Asset{
				AssetNumber:         fmt.Sprint(m["asset_id"]),
				Manufacturer:        fmt.Sprint(m["manufacturer"]),
				ManufacturerAddress: fmt.Sprint(m["manufacturer_address"]),
				ManufacturerPhone:   fmt.Sprint(m["manufacturer_phone"]),
				ManufacturerWeb:     fmt.Sprint(m["manufacturer_web"]),
				Model:               fmt.Sprint(m["model"]),
				DatePurchased:       fmt.Sprint(m["date_purchased"]),
				PurchasePrice:       fmt.Sprint(m["purchase_price"]),
				WarrantyDate:        fmt.Sprint(m["warranty_date"]),
				RetiredDate:         fmt.Sprint(m["retired_date"]),
				Description:         fmt.Sprint(m["description"]),
				Comments:            commentsFromDB,
			}

			//Execute the search result template with the data picked out of the DB.
			t := template.Must(template.ParseFiles("html_files\\layout_2.html"))
			err := t.Execute(w, data)
			if err != nil {
				return
			}
		}
	}
}

//AddAsset method to add a new asset to the DB.
func AddAsset(w http.ResponseWriter, r *http.Request) {

	CheckRetired(w, r)

	w.Header().Add("content-type", "application/json")
	iter := client.Collection("assets").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}

		// A map with string keys. Each key is one field.
		m := doc.Data()

		if m["asset_id"] == r.FormValue("asset_id") {
			data := structs.Asset{
				AssetNumber:         r.FormValue("invalid"),
				Manufacturer:        r.FormValue("asset_manufacturer"),
				ManufacturerAddress: r.FormValue("asset_manufacturer_address"),
				ManufacturerPhone:   r.FormValue("asset_manufacturer_phone"),
				ManufacturerWeb:     r.FormValue("asset_manufacturer_website"),
				Model:               r.FormValue("asset_model"),
				DatePurchased:       r.FormValue("asset_date_purchased"),
				PurchasePrice:       r.FormValue("asset_purchase_price"),
				WarrantyDate:        r.FormValue("asset_warranty_date"),
				RetiredDate:         r.FormValue("asset_date_retired"),
				Description:         r.FormValue("asset_description"),
				Comments:            []structs.Comment{},
			}

			//Execute a new template in order to refresh the page.
			t := template.Must(template.ParseFiles("html_files\\layout_4.html"))
			err := t.Execute(w, data)
			if err != nil {
				return
			}
			return
		}
	}

	//In order to pass data to Firestore we need to organize it in json form and use an interface.
	var newAsset = map[string]interface{}{
		"asset_id":             r.FormValue("asset_id"),
		"comments":             []string{r.FormValue("asset_new_comment")},
		"date_purchased":       r.FormValue("asset_date_purchased"),
		"description":          r.FormValue("asset_description"),
		"manufacturer":         r.FormValue("asset_manufacturer"),
		"manufacturer_address": r.FormValue("asset_manufacturer_address"),
		"manufacturer_phone":   r.FormValue("asset_manufacturer_phone"),
		"manufacturer_web":     r.FormValue("asset_manufacturer_website"),
		"model":                r.FormValue("asset_model"),
		"purchase_price":       r.FormValue("asset_purchase_price"),
		"retired_date":         "",
		"warranty_date":        r.FormValue("asset_warranty_date"),
	}

	//Add the new asset to the DB.
	id, err, _ := client.Collection("assets").Add(ctx, &newAsset)
	fmt.Println(id)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

}

//UpdateAsset will update the current asset with any new data.
func UpdateAsset(w http.ResponseWriter, r *http.Request) {

	CheckRetired(w, r)

	w.Header().Add("content-type", "application/json")
	iter := client.Collection("assets").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}

		// A map with string keys. Each key is one field.
		m := doc.Data()

		if m["asset_id"] == r.FormValue("asset_id") {

			var commentsFromDB []string

			s := reflect.ValueOf(m["comments"])

			for i := 0; i < s.Len(); i++ {
				if m["comments"] != "" {
					str := fmt.Sprintf("%v", s.Index(i))
					commentsFromDB = append(commentsFromDB, str)
				}
			}

			//If the form now has data in the "New Comment" field, we add this to the array containing the old comments.
			if r.FormValue("asset_new_comment") != "" {
				commentsFromDB = append(commentsFromDB, r.FormValue("asset_new_comment"))
			}

			//Pass the new data to the DB.
			var updatedAsset = map[string]interface{}{
				"asset_id":             r.FormValue("asset_id"),
				"comments":             commentsFromDB,
				"date_purchased":       r.FormValue("asset_date_purchased"),
				"description":          r.FormValue("asset_description"),
				"manufacturer":         r.FormValue("asset_manufacturer"),
				"manufacturer_address": r.FormValue("asset_manufacturer_address"),
				"manufacturer_phone":   r.FormValue("asset_manufacturer_phone"),
				"manufacturer_web":     r.FormValue("asset_manufacturer_website"),
				"model":                r.FormValue("asset_model"),
				"purchase_price":       r.FormValue("asset_purchase_price"),
				"retired_date":         r.FormValue("asset_date_retired"),
				"warranty_date":        r.FormValue("asset_warranty_date"),
			}

			//The .Set method does not add a new document if the current document already exists it will simply overwrite it with the new data.
			//doc.Ref.ID will give us the ID of the current document we are updating.
			_, err = client.Collection("assets").Doc(doc.Ref.ID).Set(ctx, updatedAsset)
			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", err)
			}
		}
	}
}

//CheckRetired will check if any objects in the DB has passed their retired date. If so they are to be deleted.
func CheckRetired(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("content-type", "application/json")
	iter := client.Collection("assets").Documents(ctx)

	//We loop through all the objects in the DB.
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}

		m := doc.Data()

		//If the current objects(document) does not have a retired date, nothing will happen.
		if m["retired_date"] == "" {
			return
		} else {
			//Since we're in the else section it means that the current object DOES have a retired date, so we must check it.

			//We fetch the date from the DB and parse it into a date type.
			//We have to save the date as a string in the DB because the firestore timestamp format does not work with the .After() method.
			timeFromDBString, _ := date.ParseISO(fmt.Sprint(m["retired_date"]))

			//We fetch today date. We use the Date type in order to use the .After() method later on.
			currentTimeString := date.Today()

			//If the current date is AFTER the retire date, the object will be deleted.
			//The .After() comes from the same package as the dates we use so that everything works as it should.
			if currentTimeString.After(timeFromDBString) {
				_, err := client.Collection("assets").Doc(doc.Ref.ID).Delete(ctx)
				if err != nil {
					log.Printf("An error has occurred: %s", err)
				}
				return
			}
		}
	}
}
