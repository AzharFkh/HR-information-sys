package seed

import (
	"hris-tools/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func helperHashed(password string) (string, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func SetupAndSeedDatabase(db *gorm.DB) {

	// roles
	roles := []models.Role{
		{ID: 1, Name: "admin"},
		{ID: 2, Name: "hr"},
		{ID: 3, Name: "employee"},
	}

	for _, role := range roles {
		db.Where(models.Role{ID: role.ID}).FirstOrCreate(&role)
	}

	password, _ := helperHashed("1234")

	defaultUsers := []models.User{
		{
			Name:     "budi_admin",
			Password: password,
			RoleID:   1, // ID Role Admin
		},
		{
			Name:     "andi_hr",
			Password: password,
			RoleID:   2, // ID Role Hr
		},
		{
			Name:     "bagus_employee",
			Password: password,
			RoleID:   3, // ID Role Hr
		},
	}

	for _, user := range defaultUsers {
		db.Where(models.User{Name: user.Name}).FirstOrCreate(&user)

		balance := models.LeaveBalance{
			UserID: user.ID,
			Quota:  12,
			Used:   0,
		}

		db.Where(models.LeaveBalance{UserID: user.ID}).FirstOrCreate(&balance)
	}

	//bisa ditambah jenis cuti
	leaveTypes := []models.LeaveType{
		{ID: 1, Name: "cuti tahunan", IsDeductQuota: true},
		{ID: 2, Name: "sakit", IsDeductQuota: false},
	}

	for _, leaveType := range leaveTypes {
		db.Where(models.LeaveType{Name: leaveType.Name}).FirstOrCreate(&leaveType)
	}
}
