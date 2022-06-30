package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/simonnik/GB_Backend2_GO/internal/entity"
	"github.com/simonnik/GB_Backend2_GO/internal/logic/storage"
)

type Router struct {
	*http.ServeMux
	stor *storage.DB
}

// NewRouter creates a router with specified storage and handlers
func NewRouter(stor *storage.DB) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		stor:     stor,
	}

	r.Handle("/create-user", http.HandlerFunc(r.CreateUser))
	r.Handle("/create-group", http.HandlerFunc(r.CreateGroup))
	r.Handle("/add-to-group", http.HandlerFunc(r.AddToGroup))
	r.Handle("/remove-from-group", http.HandlerFunc(r.RemoveFromGroup))
	r.Handle("/search-user", http.HandlerFunc(r.SearchUser))
	r.Handle("/search-group", http.HandlerFunc(r.SearchGroup))

	return r
}

// CreateUser adds a new user, passed with post-request
func (rt *Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method mismatch", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	u := entity.User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	newUser, err := rt.stor.CreateUser(r.Context(), u)
	if err != nil {
		http.Error(w, "error creating new user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(
		entity.User{
			ID:    newUser.ID,
			Name:  newUser.Name,
			Email: newUser.Email,
		},
	)
}

// CreateGroup adds a new group, passed with post-request
func (rt *Router) CreateGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method mismatch", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	g := entity.Group{}
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	newGroup, err := rt.stor.CreateGroup(r.Context(), g)
	if err != nil {
		http.Error(w, "error creating new user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(
		entity.Group{
			ID:          newGroup.ID,
			Name:        newGroup.Name,
			Description: newGroup.Description,
		},
	)
}

// AddToGroup add user specified into the group, passed with get-request
// .../add-to-group?uid=...&gid=...
func (rt *Router) AddToGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method mismatch", http.StatusMethodNotAllowed)
		return
	}

	// check and parse parameters
	uidParameter := r.URL.Query().Get("uid")
	if uidParameter == "" {
		http.Error(w, "uid should be set", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(uidParameter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (uid == uuid.UUID{}) {
		http.Error(w, "uid is nothing but ZERO", http.StatusBadRequest)
		return
	}

	gidParameter := r.URL.Query().Get("gid")
	if gidParameter == "" {
		http.Error(w, "gid should be set", http.StatusBadRequest)
		return
	}
	gid, err := uuid.Parse(gidParameter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (gid == uuid.UUID{}) {
		http.Error(w, "gid is nothing but ZERO", http.StatusBadRequest)
		return
	}

	err = rt.stor.AddToGroup(r.Context(), uid, gid)
	if err != nil {
		http.Error(w, "error adding user to the group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveFromGroup removes user specified from the group, passed with get-request
// .../remove-from-group?uid=...&gid=...
func (rt *Router) RemoveFromGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method mismatch", http.StatusMethodNotAllowed)
		return
	}

	// check and parse parameters
	uidParameter := r.URL.Query().Get("uid")
	if uidParameter == "" {
		http.Error(w, "uid should be set", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(uidParameter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (uid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	gidParameter := r.URL.Query().Get("gid")
	if gidParameter == "" {
		http.Error(w, "gid should be set", http.StatusBadRequest)
		return
	}
	gid, err := uuid.Parse(gidParameter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (gid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = rt.stor.RemoveFromGroup(r.Context(), uid, gid)
	if err != nil {
		http.Error(w, "error removing user from the group", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// SearchUser searches users by name or by group specified, passed with get-request
// .../search-user?uname=...&gid1=...&gid2=...&gid3=...
func (rt *Router) SearchUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method mismatch", http.StatusMethodNotAllowed)
		return
	}

	// check and parse parameters
	name := r.URL.Query().Get("uname")

	gidParams := make([]string, 1)
	gidParameter := r.URL.Query().Get("gid1")
	if gidParameter != "" {
		gidParams = append(gidParams, gidParameter)
	}

	gidParameter = r.URL.Query().Get("gid2")
	if gidParameter != "" {
		gidParams = append(gidParams, gidParameter)
	}

	gidParameter = r.URL.Query().Get("gid3")
	if gidParameter != "" {
		gidParams = append(gidParams, gidParameter)
	}

	gids := make([]uuid.UUID, 0)
	for _, gidParameter = range gidParams {
		gid, err := uuid.Parse(gidParameter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if (gid == uuid.UUID{}) {
			http.Error(w, "one of passed GIDs has nothing but zeroes", http.StatusBadRequest)
			return
		}
		gids = append(gids, gid)
	}

	users, err := rt.stor.SearchUser(r.Context(), name, gids)
	if err != nil {
		http.Error(w, "error searching users by name and group IDs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	for _, u := range users {
		_ = json.NewEncoder(w).Encode(
			entity.User{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
			},
		)
	}
}

// SearchGroup searches groups by name or by users specified, passed with get-request
// .../search-group?gname=...&uid1=...&uid2=...&uid3=...
func (rt *Router) SearchGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method mismatch", http.StatusMethodNotAllowed)
		return
	}

	// check and parse parameters
	name := r.URL.Query().Get("gname")

	uidParams := make([]string, 1)
	uidParameter := r.URL.Query().Get("uid1")
	if uidParameter != "" {
		uidParams = append(uidParams, uidParameter)
	}

	uidParameter = r.URL.Query().Get("uid2")
	if uidParameter != "" {
		uidParams = append(uidParams, uidParameter)
	}

	uidParameter = r.URL.Query().Get("uid3")
	if uidParameter != "" {
		uidParams = append(uidParams, uidParameter)
	}

	uids := make([]uuid.UUID, 0)
	for _, uidParameter = range uidParams {
		uid, err := uuid.Parse(uidParameter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if (uid == uuid.UUID{}) {
			http.Error(w, "one of passed GIDs has nothing but zeroes", http.StatusBadRequest)
			return
		}
		uids = append(uids, uid)
	}

	groups, err := rt.stor.SearchGroup(r.Context(), name, uids)
	if err != nil {
		http.Error(w, "error searching groups by name and user IDs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	for _, g := range groups {
		_ = json.NewEncoder(w).Encode(
			entity.Group{
				ID:          g.ID,
				Name:        g.Name,
				Type:        g.Type,
				Description: g.Description,
			},
		)
	}
}
