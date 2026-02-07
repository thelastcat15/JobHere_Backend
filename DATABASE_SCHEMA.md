# JobHere Database Schema Documentation

## Overview
This document outlines the database schema for the JobHere parking management system. The database is designed using PostgreSQL/MongoDB compatible structure with proper relationships and constraints.

---

## Entity Relationship Diagram
![ER Diagram](./diagram/JodHere_ERdiagram_2.drawio)

---

## Tables and Fields

### 1. **Auth**
User authentication and basic profile information.

| Field | Type | Description |
|-------|------|-------------|
| UID | UUID | Primary Key - Unique user identifier |
| Display name | String | User's display name |
| Email | String | User's email address |
| Phone | String | User's phone number |

**Purpose:** Stores core authentication data for users.

---

### 2. **Profile**
Extended user profile information including loyalty points.

| Field | Type | Description |
|-------|------|-------------|
| UID | UUID | Primary Key / Foreign Key (Auth) |
| Point | Number | Loyalty points balance |
| Phone | String | Contact phone number |

**Purpose:** Extends Auth with additional user profile data and rewards tracking.

---

### 3. **Report**
User-submitted reports about parking locations or issues.

| Field | Type | Description |
|-------|------|-------------|
| reportId | UUID | Primary Key |
| userId | UUID | Foreign Key (Auth) |
| image | String | Image URL/path |
| content | String | Report description |
| timestamp | Timestamp | Report submission time |
| state | ReportState | Report status (pending, approved, rejected) |
| placeId | UUID | Foreign Key (Place) |
| zoneId | UUID | Foreign Key (ParkingZone) |
| slotId | UUID | Foreign Key (ParkingSlot) |

**Purpose:** Tracks user-generated reports about parking locations or issues.

---

### 4. **Booking**
Parking slot reservations and bookings.

| Field | Type | Description |
|-------|------|-------------|
| bookId | UUID | Primary Key |
| userId | UUID | Foreign Key (Auth) |
| placeId | UUID | Foreign Key (Place) |
| status | BookingStatus | Booking status (confirmed, completed, cancelled) |
| bookedTimeStart | Timestamp | Reservation start time |
| bookedTimeEnd | Timestamp | Reservation end time |
| zoneId | UUID | Foreign Key (ParkingZone) |
| slotId | UUID | Foreign Key (ParkingSlot) |

**Purpose:** Manages parking reservations and tracks booking history.

---

### 5. **Place**
Parking location/facility information.

| Field | Type | Description |
|-------|------|-------------|
| placeId | UUID | Primary Key |
| type | String | Facility type (parking lot, garage, etc.) |
| contact | String | Contact information |
| address | String | Physical address |
| description | String | Place description |
| coordinateX | String | Latitude coordinate |
| coordinateY | String | Longitude coordinate |

**Purpose:** Stores parking facility/location details and metadata.

---

### 6. **ParkingZone**
Zones within a parking facility.

| Field | Type | Description |
|-------|------|-------------|
| zoneId | UUID | Primary Key |
| placeId | UUID | Foreign Key (Place) |
| hourRate | Number | Hourly parking rate |
| name | String | Zone name/identifier |

**Purpose:** Divides parking places into zones with different rates.

---

### 7. **ParkingSlot**
Individual parking spaces.

| Field | Type | Description |
|-------|------|-------------|
| slotId | UUID | Primary Key |
| zoneId | UUID | Foreign Key (ParkingZone) |
| name | String | Slot identifier (e.g., "A-01") |

**Purpose:** Represents individual parking spaces within zones.

---

### 8. **Sensor**
IoT sensors for slot occupancy detection.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary Key |
| name | String | Sensor identifier |
| slotId | UUID | Foreign Key (ParkingSlot) |
| url | String | Sensor API/connection URL |

**Purpose:** Tracks IoT sensors monitoring parking slot occupancy.

---

### 9. **PlaceImage**
Images associated with parking places.

| Field | Type | Description |
|-------|------|-------------|
| placeId | UUID | Foreign Key (Place) |
| path | String | Image file path/URL |

**Purpose:** Stores multiple images for each parking location.

---

### 10. **CodeRedeem**
Redemption codes for discounts or rewards.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary Key |
| userId | UUID | Foreign Key (Auth) |
| type | UUID | Redemption code type identifier |

**Purpose:** Manages discount/reward code usage per user.

---

### 11. **Reward**
Reward types and definitions.

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary Key |
| userId | UUID | Foreign Key (Auth) |
| type | UUID | Reward type identifier |

**Purpose:** Stores user reward definitions and allocations.

---

### 12. **RewardRedeem**
Redemption transactions for rewards.

| Field | Type | Description |
|-------|------|-------------|
| rewardId | UUID | Primary Key / Foreign Key (Reward) |
| userId | UUID | Foreign Key (Auth) |

**Purpose:** Tracks reward redemption transactions.

---

## Column Constraints

### Auth Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| UID | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| Display name | VARCHAR | 255 | NOT NULL | - | NO |
| Email | VARCHAR | 255 | NOT NULL, UNIQUE | - | NO |
| Phone | VARCHAR | 20 | UNIQUE | NULL | YES |

### Profile Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| UID | UUID | - | PRIMARY KEY, FK(Auth) | - | NO |
| Point | INTEGER | - | NOT NULL, CHECK(Point >= 0) | 0 | NO |
| Phone | VARCHAR | 20 | - | NULL | YES |

### Place Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| placeId | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| type | VARCHAR | 50 | NOT NULL | - | NO |
| contact | VARCHAR | 100 | - | NULL | YES |
| address | TEXT | - | NOT NULL | - | NO |
| description | TEXT | - | - | NULL | YES |
| coordinateX | VARCHAR | 50 | NOT NULL | - | NO |
| coordinateY | VARCHAR | 50 | NOT NULL | - | NO |

### ParkingZone Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| zoneId | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| placeId | UUID | - | NOT NULL, FK(Place) | - | NO |
| hourRate | DECIMAL | (10,2) | NOT NULL, CHECK(hourRate > 0) | - | NO |
| name | VARCHAR | 100 | NOT NULL | - | NO |

### ParkingSlot Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| slotId | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| zoneId | UUID | - | NOT NULL, FK(ParkingZone) | - | NO |
| name | VARCHAR | 50 | NOT NULL, UNIQUE(zoneId, name) | - | NO |

### Sensor Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| id | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| name | VARCHAR | 100 | NOT NULL | - | NO |
| slotId | UUID | - | NOT NULL, FK(ParkingSlot) | - | NO |
| url | VARCHAR | 500 | NOT NULL, UNIQUE | - | NO |

### PlaceImage Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| placeId | UUID | - | NOT NULL, FK(Place) | - | NO |
| path | VARCHAR | 500 | NOT NULL, UNIQUE(placeId, path) | - | NO |

### Booking Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| bookId | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| userId | UUID | - | NOT NULL, FK(Auth) | - | NO |
| placeId | UUID | - | NOT NULL, FK(Place) | - | NO |
| status | VARCHAR | 50 | NOT NULL | 'CONFIRMED' | NO |
| bookedTimeStart | TIMESTAMP | - | NOT NULL | - | NO |
| bookedTimeEnd | TIMESTAMP | - | NOT NULL, CHECK(bookedTimeEnd > bookedTimeStart) | - | NO |
| zoneId | UUID | - | NOT NULL, FK(ParkingZone) | - | NO |
| slotId | UUID | - | NOT NULL, FK(ParkingSlot) | - | NO |

### Report Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| reportId | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| userId | UUID | - | NOT NULL, FK(Auth) | - | NO |
| image | VARCHAR | 500 | - | NULL | YES |
| content | TEXT | - | NOT NULL | - | NO |
| timestamp | TIMESTAMP | - | NOT NULL | CURRENT_TIMESTAMP | NO |
| state | VARCHAR | 50 | NOT NULL | 'PENDING' | NO |
| placeId | UUID | - | NOT NULL, FK(Place) | - | NO |
| zoneId | UUID | - | - | NULL | YES |
| slotId | UUID | - | - | NULL | YES |

### CodeRedeem Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| id | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| userId | UUID | - | NOT NULL, FK(Auth) | - | NO |
| type | UUID | - | NOT NULL | - | NO |

### Reward Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| id | UUID | - | PRIMARY KEY, UNIQUE | uuid_generate_v4() | NO |
| userId | UUID | - | NOT NULL, FK(Auth) | - | NO |
| type | UUID | - | NOT NULL | - | NO |

### RewardRedeem Table
| Field | Type | Length | Constraint | Default | Nullable |
|-------|------|--------|-----------|---------|----------|
| rewardId | UUID | - | PRIMARY KEY, FK(Reward) | - | NO |
| userId | UUID | - | NOT NULL, FK(Auth) | - | NO |

---

## Validation Rules

### Auth Table
| Field | Validation Rule | Error Message | Example |
|-------|-----------------|---------------|---------|
| Display name | Min 2, Max 255 chars | "Name must be between 2-255 characters" | "John Doe" |
| Email | RFC 5322 format, unique | "Invalid email format" | "user@example.com" |
| Phone | E.164 format (optional) | "Phone must be in +1234567890 format" | "+66812345678" |

### Place Table
| Field | Validation Rule | Error Message | Example |
|-------|-----------------|---------------|---------|
| type | Max 50 chars, predefined list | "Invalid place type" | "parking_lot", "garage", "street" |
| address | Min 5, Max 500 chars | "Address must be between 5-500 characters" | "123 Main St, City, Country" |
| contact | Valid phone/email | "Invalid contact format" | "+66812345678" or "contact@place.com" |
| coordinateX | Valid latitude (-90 to 90) | "Latitude must be between -90 and 90" | "13.7563" |
| coordinateY | Valid longitude (-180 to 180) | "Longitude must be between -180 and 180" | "100.5018" |

### ParkingZone Table
| Field | Validation Rule | Error Message | Example |
|-------|-----------------|---------------|---------|
| hourRate | Decimal > 0, 2 decimals | "Rate must be greater than 0" | "50.00" |
| name | Max 100 chars, not empty | "Zone name required" | "Zone A", "Level 2" |

### ParkingSlot Table
| Field | Validation Rule | Error Message | Example |
|-------|-----------------|---------------|---------|
| name | Max 50 chars, alphanumeric | "Slot name invalid" | "A-01", "B-102" |

### Sensor Table
| Field | Validation Rule | Error Message | Example |
|-------|-----------------|---------------|---------|
| name | Max 100 chars, not empty | "Sensor name required" | "Sensor-A01" |
| url | Valid URL format | "Invalid URL format" | "https://api.sensor.io/device/123" |

### Booking Table
| Field | Validation Rule | Error Message | Example |
|-------|-----------------|---------------|---------|
| status | Enum: CONFIRMED, COMPLETED, CANCELLED | "Invalid booking status" | "CONFIRMED" |
| bookedTimeStart | ISO 8601, future date | "Start time must be in future" | "2026-02-10T14:00:00Z" |
| bookedTimeEnd | Must be > bookedTimeStart | "End time must be after start time" | "2026-02-10T16:00:00Z" |

### Report Table
| Field | Validation Rule | Error Message | Example |
|-------|-----------------|---------------|---------|
| content | Min 10, Max 2000 chars | "Content must be between 10-2000 characters" | "Description of the issue..." |
| state | Enum: PENDING, APPROVED, REJECTED | "Invalid report state" | "PENDING" |
| image | Valid URL/image path (optional) | "Invalid image URL" | "https://cdn.example.com/img.jpg" |

---

## Key Relationships

### One-to-Many Relationships:
- **Place** → **ParkingZone** (One place has many zones)
- **ParkingZone** → **ParkingSlot** (One zone has many slots)
- **ParkingSlot** → **Sensor** (One slot can have multiple sensors)
- **Place** → **PlaceImage** (One place has multiple images)
- **Auth/User** → **Booking** (One user has many bookings)
- **Auth/User** → **Report** (One user can submit many reports)
- **Auth/User** → **Profile** (One user has one profile)
- **Auth/User** → **CodeRedeem** (One user can redeem many codes)
- **Auth/User** → **Reward** (One user can have many rewards)

### Foreign Key Constraints:
- `Booking.userId` → `Auth.UID`
- `Booking.placeId` → `Place.placeId`
- `Booking.zoneId` → `ParkingZone.zoneId`
- `Booking.slotId` → `ParkingSlot.slotId`
- `Report.userId` → `Auth.UID`
- `Report.placeId` → `Place.placeId`
- `Report.zoneId` → `ParkingZone.zoneId`
- `Report.slotId` → `ParkingSlot.slotId`
- `ParkingZone.placeId` → `Place.placeId`
- `ParkingSlot.zoneId` → `ParkingZone.zoneId`
- `Sensor.slotId` → `ParkingSlot.slotId`
- `PlaceImage.placeId` → `Place.placeId`
- `Profile.UID` → `Auth.UID`

---

## Enums

### ReportState
- `PENDING` - Awaiting review
- `APPROVED` - Verified and approved
- `REJECTED` - Rejected report

### BookingStatus
- `CONFIRMED` - Booking confirmed
- `COMPLETED` - Parking session completed
- `CANCELLED` - Booking cancelled

---

## Notes

1. **ID Scheme:** All primary keys use UUID for distributed system compatibility
2. **Timestamps:** Use ISO 8601 format for timestamp fields
3. **Coordinates:** Store latitude/longitude as strings for flexibility with different coordinate systems
4. **Rates:** Store monetary values as numbers (integer for cents or decimal)
5. **Status Tracking:** Use enums for predictable, bounded values

---

## Migration Status
- [ ] Create Auth table
- [ ] Create Profile table
- [ ] Create Place table
- [ ] Create ParkingZone table
- [ ] Create ParkingSlot table
- [ ] Create Sensor table
- [ ] Create PlaceImage table
- [ ] Create Booking table
- [ ] Create Report table
- [ ] Create CodeRedeem table
- [ ] Create Reward table
- [ ] Create RewardRedeem table
- [ ] Add all foreign key constraints
- [ ] Add indexes for frequently queried fields

---

## Related Files
- ER Diagram: `./diagram/JodHere_ERdiagram_2.drawio`
- Previous version: `./diagram/JodHere_ERdiagram.drawio`
