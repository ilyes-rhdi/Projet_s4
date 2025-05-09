package models

type Module struct {
    ID       uint    `gorm:"primaryKey"`
    Nom      string  `gorm:"not null"`
    VolumeTD int
    VolumeTP int
    VolumeCours int
    Niveaux  []Niveau `gorm:"many2many:module_niveaux;"`
    Voeux    []Voeux `gorm:"foreignKey:ModuleID"`
}