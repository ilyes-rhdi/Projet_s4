package models

type Module struct {
	ID          uint     `gorm:"primaryKey"`
	Nom         string   `gorm:"not null"`
	Niveaux     []Niveau `gorm:"many2many:module_niveaux;"`
	Voeux       []Voeux  `gorm:"foreignKey:ModuleID"`
	VolumeCours int      `gorm:"not null;default:0"`
	VolumeTD    int      `gorm:"not null;default:0"`
	VolumeTP    int      `gorm:"not null;default:0"`
}
