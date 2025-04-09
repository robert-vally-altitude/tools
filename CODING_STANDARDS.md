# Coding Standards

## Comments

### Avoid Redundant Comments

Comments should add value by explaining *why* something is done, not *what* is being done. The code itself should be clear enough to show what it's doing.

❌ Bad:
```go
// Execute the query
rows, err := db.Query(sqlQuery)

// Define the DSN
dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, server, dbname)
```

✅ Good:
```go
// Using 3306 as the default MySQL port since it's the most common configuration
dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, server, dbname)

// Scan into interface{} array to handle dynamic column types
values := make([]interface{}, numColumns)
```

### When to Use Comments

Comments are valuable when they:
1. Explain complex business logic or non-obvious requirements
2. Document why a particular approach was chosen over alternatives
3. Warn about edge cases or potential pitfalls
4. Provide context that can't be inferred from the code
5. Document public APIs or packages

### Comment Format

- Keep comments concise and to the point
- Use complete sentences for longer explanations
- Place comments on the line before the code they describe
- Use proper grammar and punctuation
- For Go code, follow the official Go commenting conventions for documentation 