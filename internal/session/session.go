package session

import (
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SessionStore interface {
	Delete(token string) (err error)
	Find(token string) (b []byte, found bool, err error)
	Commit(token string, b []byte, expiry time.Time) (err error)
}

type SessionManager struct {
	manager *scs.SessionManager
}

func NewSessionManager(store SessionStore) *SessionManager {
	sessionManager := SessionManager{
		manager: scs.New(),
	}

	sessionManager.manager.Lifetime = 24 * time.Hour
	sessionManager.manager.IdleTimeout = 2 * time.Hour
	sessionManager.manager.Cookie.Secure = true
	sessionManager.manager.Cookie.HttpOnly = true
	sessionManager.manager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.manager.Cookie.Name = "litespend_session"

	sessionManager.manager.Store = store

	return &sessionManager
}

func (s *SessionManager) GetManager() *scs.SessionManager {
	return s.manager
}

func NewSessionInMemoryStore() SessionStore {
	return memstore.New()
}

func (s *SessionManager) LoadAndSave(ginCtx *gin.Context) {
	respWriter := ginCtx.Writer
	req := ginCtx.Request

	var token string
	cookie, err := req.Cookie(s.manager.Cookie.Name)
	if err == nil {
		token = cookie.Value
	}

	ctx, err := s.manager.Load(req.Context(), token)
	if err != nil {
		s.manager.ErrorFunc(respWriter, req, err)
		return
	}

	sessionReq := req.WithContext(ctx)
	respWriter.Header().Add("Vary", "Cookie")

	ginCtx.Request = sessionReq
	ginCtx.Next()
}

// Put adds a key and corresponding value to the session data. Any existing
// value for the key will be replaced. The session data status will be set to
// Modified.
func (s *SessionManager) Put(ctx *gin.Context, key string, val interface{}) {
	s.manager.Put(ctx.Request.Context(), key, val)
	tok, exp, _ := s.manager.Commit(ctx.Request.Context())
	s.manager.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, tok, exp)
}

// Get returns the value for a given key from the session data. The return
// value has the type interface{} so will usually need to be type asserted
// before you can use it. For example:
//
//	foo, ok := session.Get(r, "foo").(string)
//	if !ok {
//		return errors.New("type assertion to string failed")
//	}
//
// Also see the GetString(), GetInt(), GetBytes() and other helper methods which
// wrap the type conversion for common types.
func (s *SessionManager) Get(ctx *gin.Context, key string) interface{} {
	val := s.manager.Get(ctx.Request.Context(), key)
	s.manager.Commit(ctx.Request.Context())
	return val
}

// Destroy deletes the session data from the session store and sets the session
// status to Destroyed. Any further operations in the same request cycle will
// result in a new session being created.
func (s *SessionManager) Destroy(ctx *gin.Context) error {
	err := s.manager.Destroy(ctx.Request.Context())
	if err != nil {
		return err
	}
	s.manager.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, "", time.Time{})
	return nil
}

// RenewToken updates the session data to have a new session token while
// retaining the current session data. The session lifetime is also reset and
// the session data status will be set to Modified.
//
// The old session token and accompanying data are deleted from the session store.
//
// To mitiste the risk of session fixation attacks, it's important that you call
// RenewToken before making any changes to privilege levels (e.g. login and
// logout operations). See https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md#renew-the-session-id-after-any-privilege-level-change
// for additional information.
func (s *SessionManager) RenewToken(ctx *gin.Context) error {
	err := s.manager.RenewToken(ctx.Request.Context())
	if err != nil {
		return err
	}
	tok, exp, _ := s.manager.Commit(ctx.Request.Context())
	s.manager.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, tok, exp)
	return nil
}

// RememberMe controls whether the session cookie is persistent (i.e  whether it
// is retained after a user closes their browser). RememberMe only has an effect
// if you have set SessionManager.Cookie.Persist = false (the default is true) and
// you are using the standard LoadAndSave() middleware.
func (s *SessionManager) RememberMe(ctx *gin.Context, val bool) {
	s.manager.RememberMe(ctx.Request.Context(), val)
	tok, exp, _ := s.manager.Commit(ctx.Request.Context())
	s.manager.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, tok, exp)
}

// GetString returns the string value for a given key from the session data.
// The zero value for a string ("") is returned if the key does not exist or the
// value could not be type asserted to a string.
func (s *SessionManager) GetString(ctx *gin.Context, key string) string {
	val := s.manager.GetString(ctx.Request.Context(), key)
	tok, exp, _ := s.manager.Commit(ctx.Request.Context())
	s.manager.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, tok, exp)
	return val
}
