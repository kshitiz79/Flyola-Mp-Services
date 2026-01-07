# Database Connectivity Guide

This guide explains how to configure and test database connectivity for the Flyola Services Backend.

## Configuration

The application supports two methods for database configuration:

### Method 1: Individual Parameters (Recommended)

Set individual database parameters in your `.env` file:

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=flyola_services
```

### Method 2: Complete Database URL

Alternatively, you can provide a complete database URL (this takes precedence over individual parameters):

```bash
DATABASE_URL=root:password@tcp(localhost:3306)/flyola_services?charset=utf8mb4&parseTime=True&loc=Local
```

## Setup Instructions

### 1. Copy Environment File

```bash
cp .env.example .env
```

### 2. Update Database Credentials

Edit the `.env` file and update the database configuration with your actual credentials:

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=MyNewPassword123!
DB_NAME=flyola_services
```

### 3. Create Database

Make sure the database exists. You can create it using MySQL CLI:

```bash
mysql -u root -p
```

Then run:

```sql
CREATE DATABASE IF NOT EXISTS flyola_services CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. Test Database Connection

Run the database connectivity test:

```bash
go run cmd/test-db/main.go
```

This will:
- ✅ Load configuration from `.env`
- ✅ Test database connection
- ✅ Display connection pool statistics
- ✅ Show MySQL version
- ✅ List available databases

Expected output:
```
🧪 Testing Database Connectivity...
==================================================

✅ Successfully loaded .env file
🌍 Environment: development
🚀 Server Port: 8080
🗄️  Database: root@localhost:3306/flyola_services

📋 Configuration loaded:
   Database Host: localhost
   Database Port: 3306
   Database Name: flyola_services
   Database User: root

🔗 Connection String: root:****@tcp(localhost:3306)/flyola_services?charset=utf8mb4&parseTime=True&loc=Local

🔌 Attempting to connect to database...
🔌 Database connection pool configured successfully

✅ Database connection successful!

📊 Connection Pool Statistics:
   Open Connections: 1
   In Use: 0
   Idle: 1
   Max Open Connections: 100

🗄️  MySQL Version: 8.0.x

📚 Available Databases:
   [ ] information_schema
   [ ] mysql
   [✓] flyola_services
   [ ] performance_schema

✅ All database connectivity tests passed!
```

## Connection Pool Settings

The application is configured with the following connection pool settings:

- **Max Idle Connections**: 10
- **Max Open Connections**: 100
- **Connection Max Lifetime**: 1 hour

These settings can be adjusted in `internal/database/database.go` if needed.

## Troubleshooting

### Connection Refused

If you see "connection refused" errors:

1. Ensure MySQL is running:
   ```bash
   # macOS
   brew services list
   brew services start mysql
   
   # Linux
   sudo systemctl status mysql
   sudo systemctl start mysql
   ```

2. Verify MySQL is listening on the correct port:
   ```bash
   netstat -an | grep 3306
   ```

### Access Denied

If you see "access denied" errors:

1. Verify your credentials are correct in `.env`
2. Check user permissions:
   ```sql
   SHOW GRANTS FOR 'root'@'localhost';
   ```

3. Create user if needed:
   ```sql
   CREATE USER 'root'@'localhost' IDENTIFIED BY 'your_password';
   GRANT ALL PRIVILEGES ON flyola_services.* TO 'root'@'localhost';
   FLUSH PRIVILEGES;
   ```

### Database Does Not Exist

If you see "unknown database" errors:

```sql
CREATE DATABASE flyola_services CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Running the Application

Once database connectivity is confirmed, start the server:

```bash
# Using go run
go run cmd/server/main.go

# Or build and run
make build
./bin/server

# Or use air for hot reload (if installed)
air
```

## Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host | localhost | Yes |
| `DB_PORT` | Database port | 3306 | Yes |
| `DB_USER` | Database user | root | Yes |
| `DB_PASSWORD` | Database password | password | Yes |
| `DB_NAME` | Database name | flyola_services | Yes |
| `DATABASE_URL` | Complete connection string | - | No* |

*If `DATABASE_URL` is set, it takes precedence over individual parameters.

## Additional Resources

- [GORM Documentation](https://gorm.io/docs/)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [Go MySQL Driver](https://github.com/go-sql-driver/mysql)
