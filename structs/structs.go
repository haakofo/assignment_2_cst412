package structs

//These are the structs for the "objects" we use in the app. They work kinda like Java constructors.

// Asset This struct represents an Asset.
type Asset struct {
	AssetNumber         string
	Manufacturer        string
	ManufacturerAddress string
	ManufacturerPhone   string
	ManufacturerWeb     string
	Model               string
	DatePurchased       string
	PurchasePrice       string
	WarrantyDate        string
	RetiredDate         string
	Description         string
	Comments            []Comment //We have an array of comments setup like this in order to iterate over them in the HTML template.
}

// Comment Struct represting a comment.
type Comment struct {
	Comment string
}
