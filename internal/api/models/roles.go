package models

var RolePermissions = map[Role]map[string]bool{
    StaffAdmin: {
        "view_all":        true,
        "leave_comments":  true,
        "export_organigram": true,
    },
    Professeur: {
        "edit_profile":    true,
        "update_wishlist": true,
        "view_assignments": true,
    },
    ChefDep: {
        "edit_organigram": true,
        "crud_teachers":   true,
        "crud_modules":    true,
        "crud_specialties": true,
        "view_comments":   true,
        "respond_comments": true,
        "export_organigram": true,
    },
}