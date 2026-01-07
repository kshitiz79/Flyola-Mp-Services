# Database Connectivity - Quick Reference Card

## 🚀 Quick Start (3 Steps)

```bash
# 1. Copy and configure .env
cp .env.example .env
# Edit .env with your database credentials

# 2. Test database connection
make test-db

# 3. Run the application
make run
```

## 📝 Configuration (.env file)

```bash
# Server
PORT=8080
ENVIRONMENT=development
GIN_MODE=debug

# Database (Option 1 - Recommended)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=flyola_services

# Database (Option 2 - Alternative)
# DATABASE_URL=root:password@tcp(localhost:3306)/flyola_services?charset=utf8mb4&parseTime=True&loc=Local
```

## 🛠️ Makefile Commands

```bash
# Database Commands
make test-db      # Test database connectivity
make db-create    # Create database
make db-drop      # Drop database (with confirmation)

# Application Commands
make build        # Build the application
make run          # Run the application
make dev          # Run with hot reload (requires air)

# Development Commands
make test         # Run tests
make fmt          # Format code
make tidy         # Tidy dependencies
```

## ✅ Database Connection Test

```bash
make test-db
```

**Expected Output:**
```
✅ Successfully loaded .env file
🌍 Environment: development
🚀 Server Port: 8080
🗄️  Database: root@localhost:3306/flyola
✅ Database connection successful!
📊 Connection Pool Statistics:
   Open Connections: 1
   Max Open Connections: 100
🗄️  MySQL Version: 9.5.0
✅ All database connectivity tests passed!
```

## 🔧 Configuration Priority

1. **DATABASE_URL** (if set) - Takes precedence
2. **Individual Parameters** - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
3. **Default Values** - Fallback if nothing is set

## 📊 Connection Pool Settings

| Setting | Value | Description |
|---------|-------|-------------|
| Max Idle Connections | 10 | Maximum idle connections in pool |
| Max Open Connections | 100 | Maximum total connections |
| Connection Lifetime | 1 hour | Auto-refresh after this time |

## 🐛 Troubleshooting

### Connection Refused
```bash
# Check if MySQL is running
brew services list
brew services start mysql
```

### Access Denied
```bash
# Verify credentials in .env file
# Check MySQL user permissions
mysql -u root -p
```

### Database Not Found
```bash
# Create the database
make db-create
# Or manually:
mysql -u root -p -e "CREATE DATABASE flyola_services;"
```

## 📁 Key Files

| File | Purpose |
|------|---------|
| `.env` | Environment configuration |
| `internal/config/config.go` | Configuration management |
| `internal/database/database.go` | Database initialization |
| `cmd/server/main.go` | Application entry point |
| `cmd/test-db/main.go` | Database test utility |

## 🔗 How It Works

```
.env file
   ↓
godotenv.Load()
   ↓
config.Load()
   ↓
config.GetDatabaseDSN()
   ↓
database.Initialize()
   ↓
MySQL Connection ✅
```

## 📚 Documentation

- **Setup Guide**: `docs/DATABASE_SETUP.md`
- **Architecture**: `docs/DATABASE_ARCHITECTURE.md`
- **Summary**: `docs/DATABASE_CONNECTIVITY_SUMMARY.md`

## 💡 Pro Tips

1. **Never commit `.env`** - It's gitignored for security
2. **Use `.env.example`** - Keep it updated as a template
3. **Test before deploying** - Run `make test-db` first
4. **Monitor connections** - Check pool stats in production
5. **Use environment-specific configs** - Different .env for dev/staging/prod

## 🔐 Security Best Practices

- ✅ Keep `.env` in `.gitignore`
- ✅ Use strong database passwords
- ✅ Limit database user permissions
- ✅ Never log passwords (already masked in code)
- ✅ Use environment variables in production

## 🎯 Common Tasks

### First Time Setup
```bash
cp .env.example .env
# Edit .env with your credentials
make db-create
make test-db
make run
```

### Daily Development
```bash
make dev  # Hot reload enabled
```

### Before Deployment
```bash
make test
make test-db
make build
```

### Production Deployment
```bash
# Set environment variables on server
# Don't use .env file in production
export DB_HOST=production-db-host
export DB_USER=production-user
# ... etc
./bin/flyola-services
```

## 📞 Need Help?

1. Check `docs/DATABASE_SETUP.md` for detailed troubleshooting
2. Run `make test-db` to diagnose connection issues
3. Verify `.env` configuration matches your MySQL setup
4. Check MySQL logs for detailed error messages

---

**Version**: 1.0  
**Last Updated**: 2026-01-01  
**Status**: ✅ Tested and Working
