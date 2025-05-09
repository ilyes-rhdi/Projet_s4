package models

type Role string

const (
    StaffAdmin      Role = "staff_admin"
    Professeur Role = "professeur"
    ChefDep Role = "chef_dep"
)

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Nom      string `gorm:"not null"`
    Prenom   string `gorm:"not null"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Role     Role   `gorm:"type:enum('staff_admin','professeur','chef_dep');default:'professeur';not null"`
}