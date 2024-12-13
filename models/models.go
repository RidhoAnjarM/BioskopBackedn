package models

import "time"

type User struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Username string    `gorm:"unique" json:"username"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Bookings []Booking `gorm:"foreignKey:UserID"` 
}

type Movies struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Judul        string    `json:"judul"`
	Deskripsi    string    `json:"deskripsi"`
	TanggalRilis string    `json:"tanggal_rilis"`
	Photo        string    `json:"photo"`
	Harga        int       `json:"harga"`
	Trailer      string    `json:"trailer_url"`
	JumlahTiket  int       `json:"jumlah_tiket"`
	Genre        string    `json:"genre"`
	Status       string    `json:"status" gorm:"-"`
	Bookings     []Booking `gorm:"foreignKey:MovieID"` 
}

type Booking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`        
	User        User      `json:"-"`
	MovieID     uint      `json:"movie_id"`        
	Movie       Movies    `gorm:"foreignKey:MovieID"`
	Pembayaran  string    `json:"pembayaran"`
	Tiket       int       `json:"tiket"`
	MovieTitle  string    `json:"movie_title"`
	CreatedAt   time.Time `json:"created_at"`
}
