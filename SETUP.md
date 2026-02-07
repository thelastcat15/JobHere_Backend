# JobHere Backend Setup Guide

## 🚀 Quick Start

### 1. Prerequisites
- Go 1.20+ installed
- Supabase project created
- PostgreSQL database from Supabase

### 2. Environment Setup

Copy `.env.example` to `.env` and fill in your Supabase credentials:

```bash
cp .env.example .env
```

**Fill in these values from your Supabase Dashboard:**

```env
# Get from: Project Settings > API
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key

# Get from: Project Settings > Database
DB_HOST=db.your-project.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run Server

```bash
go run main.go
```

The server will start on `http://localhost:3000`

### 5. Test Connection

```bash
curl http://localhost:3000/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "JobHere API is running"
}
```

---

## 📁 Project Structure

```
JobHere_Backend/
├── config/
│   ├── database.go       # Supabase DB connection
│   └── env.go            # Environment variables
├── models/
│   └── models.go         # GORM models (12 tables)
├── main.go               # Entry point
├── .env                  # Environment variables (local)
├── .env.example          # Example env variables
├── go.mod                # Go module
└── go.sum                # Dependencies lock file
```

---

## 🔗 Supabase Connection Details

The connection uses PostgreSQL over SSL/TLS to your Supabase database.

**Features:**
- ✅ Automatic migrations on startup
- ✅ Connection pooling
- ✅ Error handling
- ✅ Graceful shutdown

---

## 📝 Database Tables Created

When the server starts, GORM will automatically create these tables:

1. **auth** - User authentication
2. **profile** - User profiles & loyalty points
3. **place** - Parking locations
4. **parking_zone** - Zones within places
5. **parking_slot** - Individual parking spaces
6. **sensor** - IoT sensors for slot detection
7. **place_image** - Images for places
8. **booking** - Parking reservations
9. **report** - User reports
10. **code_redeem** - Redemption codes
11. **reward** - Reward definitions
12. **reward_redeem** - Reward redemptions

---

## 🔑 Key Endpoints

### Health Check
```
GET /health
```

### API Base
```
GET /api/v1/
```

---

## 🛠️ Development

### Add New Route Example

Edit `main.go`:
```go
// In the api group
api.Get("/places", func(c *fiber.Ctx) error {
    var places []models.Place
    db := config.GetDB()
    db.Find(&places)
    return c.JSON(places)
})
```

### Database Queries

Use the global `DB` instance:
```go
import "jobhere.backend/config"

db := config.GetDB()

// Create
user := models.Auth{...}
db.Create(&user)

// Read
db.First(&user, "uid = ?", id)

// Update
db.Model(&user).Update("email", "new@example.com")

// Delete
db.Delete(&user)
```

---

## 📚 Dependencies

- **Fiber** - Web framework
- **GORM** - ORM for Go
- **PostgreSQL Driver** - Database driver
- **UUID** - Unique identifiers
- **godotenv** - Environment variables

---

## ⚠️ Important Notes

1. **Never commit `.env`** - Always use `.env.example` for reference
2. **Always use environment variables** in production
3. **Enable Row Level Security (RLS)** in Supabase for security
4. **Backup database** before running migrations

---

## 🐛 Troubleshooting

### Connection Failed
- Check your Supabase credentials in `.env`
- Ensure database is accessible (IP whitelisting)
- Verify network connectivity

### Migration Failed
- Check database permissions
- Ensure DB_USER has table creation rights
- Check PostgreSQL version compatibility (GORM requires 9.1+)

### Port Already in Use
Change `SERVER_PORT` in `.env` or modify `main.go`:
```go
app.Listen(":8080") // Different port
```

---

## 📖 Resources

- [Supabase Docs](https://supabase.com/docs)
- [GORM Documentation](https://gorm.io)
- [Fiber Documentation](https://docs.gofiber.io)
- [PostgreSQL Driver](https://github.com/jackc/pgx)
