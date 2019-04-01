package main

var StoredProcConfig = []CrudStoredProcConfig{}

// Table based end points
var TableConfig = []CrudConfig{
	{
		URIPath:        "/api/v1/documents",
		AuthKey:        true,
		JWTKey:         false,
		MethodsAllowed: []string{"GET", "POST", "PUT", "DELETE"},
		TableName:      "documents",
		InsertCols: []string{
			"id", "document_hash", "email", "real_name", "phone_number", "address_usps", "document_file_name",
			"file_name", "orig_file_name", "txid", "note"},
		InsertPkCol: "id",
		UpdateCols: []string{
			"id", "document_hash", "email", "real_name", "phone_number", "address_usps", "document_file_name",
			"file_name", "orig_file_name", "txid", "note"},
		UpdatePkCol: "id",
		WhereCols: []string{
			"id", "document_hash", "email", "real_name", "phone_number", "address_usps", "document_file_name",
			"file_name", "orig_file_name", "txid", "note"},
	},
}
