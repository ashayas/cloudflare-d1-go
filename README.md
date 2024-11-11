# Cloudflare D1 Go Client ☁️ 
- This is a lightweight Go client for the Cloudflare D1 database
- D1 is a cool serverless, zero-config, transactional SQL database built by [Cloudflare](https://www.cloudflare.com/) built for the edge and cost-effective

## Installation 📦

```bash
go get github.com/ashayas/cloudflare-d1-go
```

## Usage 💻

### Initialize the client 🔑

```go
client := cloudflare_d1_go.NewClient("account_id", "api_token")
```

### Connect to a database 📁

```go
client.ConnectDB("database_id")
```

### Query the database 🔍

```go
// Execute a SQL query with optional parameters
// query: SQL query string
// params: Array of parameter values to bind to the query (use ? placeholders in query)
client.Query("SELECT * FROM users WHERE age > ?", []string{"18"})
```

Example with parameters:
```go
// Find users in a specific city
client.Query("SELECT * FROM users WHERE city = ?", []string{"San Francisco"})

// Find products in a price range
client.Query("SELECT * FROM products WHERE price >= ? AND price <= ?", []string{"10.00", "50.00"})
```

### Create a table 📄

```go
client.CreateTable("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)")
```

### Remove a table 🗑️

```go
client.RemoveTable("users")
```

### Method 2
- Specify the database ID in the client. 
- Useful if you have multiple databases and want to switch between them

```go
client := cloudflare_d1_go.NewClient("account_id", "api_token")
client.QueryDB(databaseID, "SELECT * FROM users", nil)
```

### List Of Methods

#### Database Management
- `NewClient(accountID, apiToken string) *Client` - Creates a new D1 client
- `ListDB() (*APIResponse, error)` - Lists all databases in the account
- `CreateDB(name string) (*APIResponse, error)` - Creates a new database
- `DeleteDB(databaseID string) (*APIResponse, error)` - Deletes a database
- `ConnectDB(name string) error` - Connects to a database by name for subsequent operations

#### Table Operations
- `CreateTable(createQuery string) (*APIResponse, error)` - Creates a table in the connected database
- `RemoveTable(tableName string) (*APIResponse, error)` - Removes a table from the connected database
- `CreateTableWithID(databaseID, createQuery string) (*APIResponse, error)` - Creates a table in a specific database
- `RemoveTableWithID(databaseID, tableName string) (*APIResponse, error)` - Removes a table from a specific database

#### Query Execution
- `Query(query string, params []string) (*APIResponse, error)` - Executes a query on the connected database
- `QueryDB(databaseID string, query string, params []string) (*APIResponse, error)` - Executes a query on a specific database

## TODO
- Better error handling 🛡️
- More comprehensive test coverage 🧪
- Improvements on persistence of database connections 🔄 
- Integration with GORM 🦕 

## Testing 
- Run `go test` to run the tests
- You can use the `-v` flag to see more verbose test results

## Contributing 🤝
Contributions are welcome! Please feel free to submit a Pull Request.

## License 📄
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support 💪
If you encounter any issues or have questions, please file an issue on GitHub.