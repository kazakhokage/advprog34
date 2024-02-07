package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	os.Setenv("DATABASE_URL", "user=postgres dbname=postgres1 password=Maryam06. sslmode=disable")
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{}, &Good{}) // Create tables User and Good
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")

	db.Create(&User{Username: "Esen", Email: "esen.com", Password: "Batur"})
	db.Unscoped().Delete(&User{}, "id = ?", 3)

	insertGoods()

	user, err := getUserByID(4)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("User by ID:", user)
	}

	allUsers, err := getAllUsers()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("All users:", allUsers)
	}

	// Define routes
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/goods", handleGoodsPage)

	// Start the HTTP server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type JsonRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Message  string `json:"message"`
}

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Role     string
	Token    string // Add a new field for session token
}

type Good struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func insertGoods() {
	goods := []Good{
		{Name: "Smartphone 1", Description: "Description for Smartphone 1", Price: 499.99},
		{Name: "Smartphone 16", Description: "Description for Smartphone 1", Price: 100.99},
		{Name: "Smartphone 2", Description: "Description for Smartphone 2", Price: 599.99},
		{Name: "Smartphone 3", Description: "Description for Smartphone 3", Price: 699.99},
		{Name: "Smartphone 4", Description: "Description for Smartphone 4", Price: 799.99},
		{Name: "Smartphone 5", Description: "Description for Smartphone 5", Price: 899.99},
		{Name: "Smartphone 6", Description: "Description for Smartphone 6", Price: 999.99},
		{Name: "Smartphone 7", Description: "Description for Smartphone 7", Price: 1099.99},
		{Name: "Smartphone 8", Description: "Description for Smartphone 8", Price: 1199.99},
		{Name: "Smartphone 9", Description: "Description for Smartphone 9", Price: 1299.99},
		{Name: "Smartphone 10", Description: "Description for Smartphone 10", Price: 1399.99},
		{Name: "Smartphone 11", Description: "Description for Smartphone 11", Price: 1499.99},
		{Name: "Smartphone 12", Description: "Description for Smartphone 12", Price: 1599.99},
		{Name: "Smartphone 13", Description: "Description for Smartphone 13", Price: 1699.99},
		{Name: "Smartphone 14", Description: "Description for Smartphone 14", Price: 1799.99},
		{Name: "Smartphone 15", Description: "Description for Smartphone 15", Price: 1899.99},
		{Name: "Smartphone 17", Description: "Description for Smartphone 17", Price: 1999.99},
		{Name: "Smartphone 18", Description: "Description for Smartphone 18", Price: 2099.99},
		{Name: "Smartphone 19", Description: "Description for Smartphone 19", Price: 2199.99},
		{Name: "Smartphone 20", Description: "Description for Smartphone 20", Price: 2299.99},
		{Name: "Smartphone 21", Description: "Description for Smartphone 21", Price: 2399.99},
		{Name: "Smartphone 22", Description: "Description for Smartphone 22", Price: 2499.99},
	}

	for _, good := range goods {
		result := db.Create(&good)
		if result.Error != nil {
			panic(result.Error)
		}
	}

	fmt.Println("Inserted 20 goods into the database")
}

func loginUser(username, password string) (string, error) {
	var user User
	result := db.Where("username = ? AND password = ?", username, password).First(&user)
	if result.Error != nil {
		return "", fmt.Errorf("Invalid credentials")
	}

	// Set a session token for the user
	token := generateSessionToken()
	user.Token = token
	db.Save(&user)

	return token, nil
}

func generateSessionToken() string {
	// Implement your logic to generate a session token
	// For simplicity, let's generate a random string for now
	return "random_session_token"
}

func logoutUser(token string) error {
	var user User
	result := db.Where("token = ?", token).First(&user)
	if result.Error != nil {
		return fmt.Errorf("User not found with the given token")
	}

	// Clear the session token for the user
	user.Token = ""
	db.Save(&user)

	return nil
}

func handleMainPage(w http.ResponseWriter, r *http.Request, token string) {
	goods, err := getAllGoods()
	if err != nil {
		http.Error(w, "Failed to fetch goods", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		Goods  []Good `json:"goods"`
	}{
		Status: "success",
		Goods:  goods,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

func handleDashboardPage(w http.ResponseWriter, r *http.Request, token string) {
	// Implement your logic for the dashboard page
	// For now, let's return a simple message
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := JsonResponse{
		Status:  "success",
		Message: "Welcome to the Dashboard!",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// ... (unchanged code for handleRequest)

	switch r.Method {
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodGet:
		handleGetRequest(w, r)
	case http.MethodOptions:
		handleOptionsRequest(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	var requestData JsonRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, `{"status": "400", "message": "Invalid JSON message"}`, http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case "/login":
		handleLoginRequest(w, requestData)
	case "/logout":
		handleLogoutRequest(w, requestData)
	default:
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	name := queryParams.Get("name")
	age := queryParams.Get("age")

	if name == "" || age == "" {
		http.Error(w, `{"status": "400", "message": "Both name and age parameters are required in the GET request"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf("Received GET request with name: %s, age: %s\n", name, age)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/main":
		handleMainPage(w, r, "")
	case "/dashboard":
		handleDashboardPage(w, r, "")
	default:
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}

func handleLoginRequest(w http.ResponseWriter, requestData JsonRequest) {
	if requestData.Message == "" && (requestData.Username == "" || requestData.Password == "") {
		http.Error(w, `{"status": "400", "message": "Either 'message' or both 'username' and 'password' fields must be provided"}`, http.StatusBadRequest)
		return
	}

	if requestData.Message != "" {
		fmt.Println("Received POST message:", requestData.Message)
	} else {
		token, err := loginUser(requestData.Username, requestData.Password)
		if err != nil {
			http.Error(w, `{"status": "401", "message": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		handleDashboardPage(w, nil, token)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := JsonResponse{
		Status:  "success",
		Message: "Data successfully received",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

func handleLogoutRequest(w http.ResponseWriter, requestData JsonRequest) {
	if requestData.Message == "" || requestData.Username != "" || requestData.Password != "" {
		http.Error(w, `{"status": "400", "message": "Invalid parameters for logout request"}`, http.StatusBadRequest)
		return
	}

	logoutUser(requestData.Message)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := JsonResponse{
		Status:  "success",
		Message: "Logout successful",
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

func handleOptionsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func isAuthorized(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		return false
	}

	var user User
	result := db.Where("token = ?", token).First(&user)
	return result.Error == nil
}

func getAllGoods() ([]Good, error) {
	var goods []Good
	result := db.Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}
	return goods, nil
}

func saveRegistrationData(username, email, password, role string) error {
	user := User{Username: username, Email: email, Password: password, Role: role}
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func getUserByID(userID uint) (User, error) {
	var user User
	result := db.Select("id, username, email, password").First(&user, userID)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func getAllUsers() ([]User, error) {
	var users []User
	result := db.Select("id, username, email, password").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func handleGoodsPage(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	sortField := r.URL.Query().Get("sort")
	filterName := r.URL.Query().Get("filter_name")
	pageStr := r.URL.Query().Get("page")
	itemsPerPageStr := r.URL.Query().Get("per_page")

	// Set default values for pagination
	page := 1
	itemsPerPage := 10

	// Parse page and itemsPerPage from query parameters
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}
	if itemsPerPageStr != "" {
		itemsPerPageInt, err := strconv.Atoi(itemsPerPageStr)
		if err == nil && itemsPerPageInt > 0 {
			itemsPerPage = itemsPerPageInt
		}
	}

	// Calculate offset for pagination
	offset := (page - 1) * itemsPerPage

	// Retrieve goods with filtering and pagination
	goods, totalItems, err := getAllGoodsWithFilteringAndPagination(filterName, sortField, itemsPerPage, offset)
	if err != nil {
		http.Error(w, "Failed to fetch goods", http.StatusInternalServerError)
		return
	}

	// Parse HTML template
	tmpl, err := template.New("goods").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Goods Page</title>
			<style>
				/* Add your styles here */
				body {
					font-family: Arial, sans-serif;
					margin: 20px;
				}

				h1 {
					color: #333;
				}

				ul {
					list-style-type: none;
					padding: 0;
				}

				li {
					border: 1px solid #ddd;
					padding: 10px;
					margin-bottom: 10px;
					border-radius: 5px;
				}

				.pagination {
					display: inline-block;
				}

				.pagination a {
					color: black;
					float: left;
					padding: 8px 16px;
					text-decoration: none;
					transition: background-color .3s;
					border: 1px solid #ddd;
					margin: 0 4px;
				}

				.pagination a.active {
					background-color: #4CAF50;
					color: white;
					border: 1px solid #4CAF50;
				}

				.pagination a:hover:not(.active) {background-color: #ddd;}
			</style>
		</head>
		<body>
			<h1>Goods List</h1>
			<form action="/goods" method="get">
				<label for="sort">Sort:</label>
				<select id="sort" name="sort">
					<option value="name">By Name</option>
					<option value="description">By Description</option>
					<option value="price">By Price</option>
				</select>
				<label for="filter_name">Filter by Name:</label>
				<input type="text" id="filter_name" name="filter_name">
				<button type="submit">Apply</button>
			</form>
			<ul>
				{{range .Goods}}
					<li>
						<strong>Name:</strong> {{.Name}}<br>
						<strong>Description:</strong> {{.Description}}<br>
						<strong>Price:</strong> ${{.Price}}
					</li>
					<br>
				{{end}}
			</ul>
			<div class="pagination">
				{{if gt .TotalPages 1}}
					{{if gt .Page 1}}
						<a href="/goods?page=1&per_page={{.ItemsPerPage}}&sort={{.SortField}}&filter_name={{.FilterName}}">First</a>
						<a href="/goods?page={{.PrevPage}}&per_page={{.ItemsPerPage}}&sort={{.SortField}}&filter_name={{.FilterName}}">Previous</a>
					{{end}}
					{{range $i := .PageNumbers}}
						{{if eq $i $.Page}}
							<a class="active" href="#">{{$i}}</a>
						{{else}}
							<a href="/goods?page={{$i}}&per_page={{$.ItemsPerPage}}&sort={{$.SortField}}&filter_name={{$.FilterName}}">{{$i}}</a>
						{{end}}
					{{end}}
					{{if lt .Page .TotalPages}}
						<a href="/goods?page={{.NextPage}}&per_page={{.ItemsPerPage}}&sort={{.SortField}}&filter_name={{.FilterName}}">Next</a>
						<a href="/goods?page={{.TotalPages}}&per_page={{.ItemsPerPage}}&sort={{.SortField}}&filter_name={{.FilterName}}">Last</a>
					{{end}}
				{{end}}
			</div>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Calculate total number of pages
	totalPages := totalItems / itemsPerPage
	if totalItems%itemsPerPage != 0 {
		totalPages++
	}

	// Generate page numbers for pagination links
	var pageNumbers []int
	for i := 1; i <= totalPages; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	// Execute the template with the goods data
	err = tmpl.Execute(w, struct {
		Goods        []Good
		Page         int
		TotalPages   int
		PageNumbers  []int
		ItemsPerPage int
		SortField    string
		FilterName   string
	}{
		Goods:        goods,
		Page:         page,
		TotalPages:   totalPages,
		PageNumbers:  pageNumbers,
		ItemsPerPage: itemsPerPage,
		SortField:    sortField,
		FilterName:   filterName,
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getAllGoodsWithFilteringAndPagination(filterName, sortField string, itemsPerPage, offset int) ([]Good, int, error) {
	var goods []Good
	dbQuery := db

	// Apply filtering if filterName is provided
	if filterName != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+filterName+"%")
	}

	// Count total items before pagination
	var totalItems int64
	result := dbQuery.Model(&Good{}).Count(&totalItems)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Apply sorting and pagination
	result = dbQuery.Order(sortField).Limit(itemsPerPage).Offset(offset).Find(&goods)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return goods, int(totalItems), nil
}
