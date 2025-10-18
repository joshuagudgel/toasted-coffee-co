package handlers

import (
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/services"
)

type Handlers struct {
	Auth    *AuthHandler
	Booking *BookingHandler
	Contact *ContactHandler
	Menu    *MenuHandler
	Package *PackageHandler
}

func NewHandlers(repos *database.Repositories, emailService *services.EmailService) *Handlers {
	return &Handlers{
		Auth:    NewAuthHandler(repos.User),
		Booking: NewBookingHandler(repos.Booking, emailService),
		Contact: NewContactHandler(emailService),
		Menu:    NewMenuHandler(repos.Menu),
		Package: NewPackageHandler(repos.Package),
	}
}
