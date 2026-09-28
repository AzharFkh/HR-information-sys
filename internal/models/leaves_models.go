package models

import (
	"time"

	"gorm.io/gorm"
)

// jenis cuti, bisa ditambahkan oleh HR terkait jenisnya
type LeaveType struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"type:varchar(100);unique;not null"` // Contoh: "Cuti Tahunan", "Sakit"
	IsDeductQuota bool   `gorm:"default:false"`                     // True jika cuti memotong jatah tahunan
}

// jatah cuti karyawan
type LeaveBalance struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null;index"`
	LeaveTypeID uint `gorm:"not null;index"`      //for next feature
	Quota       int  `gorm:"not null;default:12"` // Sisa jatah yang bisa diambil (atur sesuai perusahaan)
	Used        int  `gorm:"not null;default:0"`  // Opsional: Tracking jumlah yang sudah dipakai

	// Relasi untuk memudahkan query preload di GORM
	LeaveType LeaveType `gorm:"foreignKey:LeaveTypeID"`
}

type LeaveStatus string

const (
	StatusPending  LeaveStatus = "PENDING"
	StatusApproved LeaveStatus = "APPROVED"
	StatusRejected LeaveStatus = "REJECTED"
)

type LeaveRequest struct {
	ID           uint        `gorm:"primaryKey"`
	UserID       uint        `gorm:"not null;index"`     // Karyawan yang mengajukan
	LeaveTypeID  uint        `gorm:"not null"`           // Jenis cuti yang diambil
	StartDate    time.Time   `gorm:"type:date;not null"` // tanggal cuti mulai (menggunakan type date)
	EndDate      time.Time   `gorm:"type:date;not null"`
	Reason       string      `gorm:"type:text"` // Menggunakan TEXT agar cukup untuk ~200 kata
	Status       LeaveStatus `gorm:"type:varchar(20);default:'PENDING'"`
	ApprovedByID *uint       `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"` // GORM Soft Delete

	// Relasi
	LeaveType LeaveType `gorm:"foreignKey:LeaveTypeID"`
	User      User      `gorm:"foreignKey:UserID"`
}
