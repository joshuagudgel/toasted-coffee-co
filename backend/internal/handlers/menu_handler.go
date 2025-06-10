package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// MenuHandler handles HTTP requests for menu items
type MenuHandler struct {
	repo database.MenuRepository
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler(repo database.MenuRepository) *MenuHandler {
	return &MenuHandler{
		repo: repo,
	}
}

// GetAll returns all menu items
func (h *MenuHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.repo.GetAll(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve menu items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// GetByType returns menu items of a specific type
func (h *MenuHandler) GetByType(w http.ResponseWriter, r *http.Request) {
	itemType := chi.URLParam(r, "type")
	if itemType != string(models.CoffeeFlavor) && itemType != string(models.MilkOption) {
		http.Error(w, "Invalid item type", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	items, err := h.repo.GetByType(ctx, models.ItemType(itemType))
	if err != nil {
		http.Error(w, "Failed to retrieve menu items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Create handles POST /menu requests to add a new menu item
func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the menu item
	if menuItem.Value == "" || menuItem.Label == "" {
		http.Error(w, "Value and label are required", http.StatusBadRequest)
		return
	}

	if menuItem.Type != models.CoffeeFlavor && menuItem.Type != models.MilkOption {
		http.Error(w, "Type must be either coffee_flavor or milk_option", http.StatusBadRequest)
		return
	}

	// Create the menu item
	id, err := h.repo.Create(r.Context(), &menuItem)
	if err != nil {
		http.Error(w, "Failed to create menu item", http.StatusInternalServerError)
		return
	}

	// Set the ID on the returned item
	menuItem.ID = id

	// Return the created item
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(menuItem)
}

// Update handles PUT /menu/{id} requests to update an existing menu item
func (h *MenuHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	// Decode request body
	var menuItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the menu item
	if menuItem.Value == "" || menuItem.Label == "" {
		http.Error(w, "Value and label are required", http.StatusBadRequest)
		return
	}

	if menuItem.Type != models.CoffeeFlavor && menuItem.Type != models.MilkOption {
		http.Error(w, "Type must be either coffee_flavor or milk_option", http.StatusBadRequest)
		return
	}

	// Update the menu item
	if err := h.repo.Update(r.Context(), id, &menuItem); err != nil {
		http.Error(w, "Failed to update menu item", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"message": "Menu item updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete handles DELETE /menu/{id} requests to remove a menu item
func (h *MenuHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse ID from URL
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid menu item ID", http.StatusBadRequest)
		return
	}

	// Delete the menu item
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete menu item", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Menu item deleted successfully",
	})
}
