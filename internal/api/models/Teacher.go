package models

type Teacher struct {
    ID            uint        `gorm:"primaryKey"`
    UserID        uint        `gorm:"not null;unique"`
    YearEntrance  string      `gorm:"unique;not null"`  // Année d'entrée
    Grade         string      `gorm:"not null"`
    User          User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Specialities  []Speciality `gorm:"many2many:teacher_specialities;"`
    Voeux         []Voeux     `gorm:"foreignKey:TeacherID"`
    Name          string      `gorm:"not null"`         // Teacher's name
    SpecialtyID   uint        `gorm:"not null"`         // ID of the teacher's specialty
    MaxHours      int         `gorm:"not null"`         // Maximum hours the teacher can work
    CurrentHours  int         `gorm:"default:0"`        // Current hours assigned to the teacher
}