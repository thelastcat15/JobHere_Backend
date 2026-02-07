package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Enums for status fields
type ReportState string
type BookingStatus string

const (
	ReportPending  ReportState = "PENDING"
	ReportApproved ReportState = "APPROVED"
	ReportRejected ReportState = "REJECTED"

	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingCompleted BookingStatus = "COMPLETED"
	BookingCancelled BookingStatus = "CANCELLED"
)

// ============================================
// Auth Model
// ============================================
type Auth struct {
	UID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"uid"`
	DisplayName string    `gorm:"type:varchar(255);not null" json:"display_name" validate:"required,min=2,max=255"`
	Email       string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email" validate:"required,email"`
	Phone       string    `gorm:"type:varchar(20);uniqueIndex;default:null" json:"phone" validate:"omitempty,e164"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Profile       *Profile       `gorm:"foreignKey:UID" json:"profile,omitempty"`
	Bookings      []Booking      `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
	Reports       []Report       `gorm:"foreignKey:UserID" json:"reports,omitempty"`
	CodeRedeems   []CodeRedeem   `gorm:"foreignKey:UserID" json:"code_redeems,omitempty"`
	Rewards       []Reward       `gorm:"foreignKey:UserID" json:"rewards,omitempty"`
	RewardRedeems []RewardRedeem `gorm:"foreignKey:UserID" json:"reward_redeems,omitempty"`
}

// ============================================
// Profile Model
// ============================================
type Profile struct {
	UID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"uid"`
	Point     int       `gorm:"type:integer;not null;default:0;check:point>=0" json:"point"`
	Phone     string    `gorm:"type:varchar(20);default:null" json:"phone" validate:"omitempty,e164"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	Auth *Auth `gorm:"foreignKey:UID;references:UID" json:"auth,omitempty"`
}

// ============================================
// Place Model
// ============================================
type Place struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type" validate:"required,max=50"`
	Contact     string    `gorm:"type:varchar(100);default:null" json:"contact" validate:"omitempty,max=100"`
	Address     string    `gorm:"type:text;not null" json:"address" validate:"required,min=5,max=500"`
	Description string    `gorm:"type:text;default:null" json:"description"`
	CoordinateX string    `gorm:"type:varchar(50);not null" json:"coordinate_x" validate:"required,latitude"`
	CoordinateY string    `gorm:"type:varchar(50);not null" json:"coordinate_y" validate:"required,longitude"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ParkingZones []ParkingZone `gorm:"foreignKey:PlaceID;references:ID" json:"parking_zones,omitempty"`
	Bookings     []Booking     `gorm:"foreignKey:PlaceID;references:ID" json:"bookings,omitempty"`
	Reports      []Report      `gorm:"foreignKey:PlaceID;references:ID" json:"reports,omitempty"`
	Images       []PlaceImage  `gorm:"foreignKey:PlaceID;references:ID" json:"images,omitempty"`
}

// ============================================
// ParkingZone Model
// ============================================
type ParkingZone struct {
	ID        uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlaceID   uuid.UUID          `gorm:"type:uuid;not null;index" json:"place_id"`
	HourRate  datatypes.JSONType `gorm:"type:numeric(10,2);not null;check:hour_rate>0" json:"hour_rate" validate:"required,min=0"`
	Name      string             `gorm:"type:varchar(100);not null" json:"name" validate:"required,max=100"`
	CreatedAt time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time          `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Place        *Place        `gorm:"foreignKey:PlaceID;references:ID" json:"place,omitempty"`
	ParkingSlots []ParkingSlot `gorm:"foreignKey:ZoneID;references:ID" json:"parking_slots,omitempty"`
	Bookings     []Booking     `gorm:"foreignKey:ZoneID;references:ID" json:"bookings,omitempty"`
	Reports      []Report      `gorm:"foreignKey:ZoneID;references:ID" json:"reports,omitempty"`
}

// ============================================
// ParkingSlot Model
// ============================================
type ParkingSlot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ZoneID    uuid.UUID `gorm:"type:uuid;not null;index" json:"zone_id"`
	Name      string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_zone_slot" json:"name" validate:"required,max=50"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;references:ID" json:"parking_zone,omitempty"`
	Sensors     []Sensor     `gorm:"foreignKey:SlotID;references:ID" json:"sensors,omitempty"`
	Bookings    []Booking    `gorm:"foreignKey:SlotID;references:ID" json:"bookings,omitempty"`
	Reports     []Report     `gorm:"foreignKey:SlotID;references:ID" json:"reports,omitempty"`
}

// ============================================
// Sensor Model
// ============================================
type Sensor struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,max=100"`
	SlotID    uuid.UUID `gorm:"type:uuid;not null;index" json:"slot_id"`
	URL       string    `gorm:"type:varchar(500);not null;uniqueIndex" json:"url" validate:"required,url,max=500"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	ParkingSlot *ParkingSlot `gorm:"foreignKey:SlotID;references:SlotID" json:"parking_slot,omitempty"`
}

// ============================================
// PlaceImage Model
// ============================================
type PlaceImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlaceID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_place_path" json:"place_id"`
	Path      string    `gorm:"type:varchar(500);not null;uniqueIndex:idx_place_path" json:"path" validate:"required,url,max=500"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	Place *Place `gorm:"foreignKey:PlaceID;references:PlaceID" json:"place,omitempty"`
}

// ============================================
// Booking Model
// ============================================
type Booking struct {
	ID              uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	PlaceID         uuid.UUID     `gorm:"type:uuid;not null;index" json:"place_id"`
	Status          BookingStatus `gorm:"type:varchar(50);not null;default:'CONFIRMED'" json:"status" validate:"required,oneof=CONFIRMED COMPLETED CANCELLED"`
	BookedTimeStart time.Time     `gorm:"not null" json:"booked_time_start" validate:"required"`
	BookedTimeEnd   time.Time     `gorm:"not null" json:"booked_time_end" validate:"required,gtfield=BookedTimeStart"`
	ZoneID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"zone_id"`
	SlotID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"slot_id"`
	CreatedAt       time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User        *Auth        `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
	Place       *Place       `gorm:"foreignKey:PlaceID;references:PlaceID" json:"place,omitempty"`
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;references:ZoneID" json:"parking_zone,omitempty"`
	ParkingSlot *ParkingSlot `gorm:"foreignKey:SlotID;references:SlotID" json:"parking_slot,omitempty"`
}

// ============================================
// Report Model
// ============================================
type Report struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	Image     string      `gorm:"type:varchar(500);default:null" json:"image" validate:"omitempty,url,max=500"`
	Content   string      `gorm:"type:text;not null" json:"content" validate:"required,min=10,max=2000"`
	Timestamp time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP" json:"timestamp"`
	State     ReportState `gorm:"type:varchar(50);not null;default:'PENDING'" json:"state" validate:"required,oneof=PENDING APPROVED REJECTED"`
	PlaceID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"place_id"`
	ZoneID    uuid.UUID   `gorm:"type:uuid;index;default:null" json:"zone_id"`
	SlotID    uuid.UUID   `gorm:"type:uuid;index;default:null" json:"slot_id"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User        *Auth        `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
	Place       *Place       `gorm:"foreignKey:PlaceID;references:ID" json:"place,omitempty"`
	ParkingZone *ParkingZone `gorm:"foreignKey:ZoneID;references:ID" json:"parking_zone,omitempty"`
	ParkingSlot *ParkingSlot `gorm:"foreignKey:SlotID;references:ID" json:"parking_slot,omitempty"`
}

// ============================================
// CodeRedeem Model
// ============================================
type CodeRedeem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      uuid.UUID `gorm:"type:uuid;not null" json:"type"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	User *Auth `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
}

// ============================================
// Reward Model
// ============================================
type Reward struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      uuid.UUID `gorm:"type:uuid;not null" json:"type"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	User          *Auth          `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
	RewardRedeems []RewardRedeem `gorm:"foreignKey:RewardID" json:"reward_redeems,omitempty"`
}

// ============================================
// RewardRedeem Model
// ============================================
type RewardRedeem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RewardID  uuid.UUID `gorm:"type:uuid;not null;index" json:"reward_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Reward *Reward `gorm:"foreignKey:RewardID;references:ID" json:"reward,omitempty"`
	User   *Auth   `gorm:"foreignKey:UserID;references:UID" json:"user,omitempty"`
}

// ============================================
// Database Migration Function
// ============================================
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Auth{},
		&Profile{},
		&Place{},
		&ParkingZone{},
		&ParkingSlot{},
		&Sensor{},
		&PlaceImage{},
		&Booking{},
		&Report{},
		&CodeRedeem{},
		&Reward{},
		&RewardRedeem{},
	)
}
