/* AuthPlz Authentication and Authorization Microservice
 * Application context session storage and binding
 *
 * Copyright 2018 Ryan Kurte
 */

package appcontext

import (
	"log"

	"github.com/gocraft/web"
	"github.com/gorilla/sessions"
)

// BindInst Binds an object instance to a session key and writes to the browser session store
// TODO: Bindings should probably time out eventually
func (c *AuthPlzCtx) BindInst(rw web.ResponseWriter, req *web.Request, sessionKey, dataKey string, inst interface{}) error {
	session, err := c.Global.SessionStore.Get(req.Request, sessionKey)
	if err != nil {
		log.Printf("AuthPlzCtx.Bind error fetching session-key:%s (%s)", sessionKey, err)
		c.WriteInternalError(rw)
		return err
	}

	session.Values[dataKey] = inst
	session.Save(req.Request, rw)

	return nil
}

// GetInst Fetches an object instance by session key from the browser session store
func (c *AuthPlzCtx) GetInst(rw web.ResponseWriter, req *web.Request, sessionKey, dataKey string) (interface{}, error) {
	session, err := c.Global.SessionStore.Get(req.Request, sessionKey)
	if err != nil {
		log.Printf("AuthPlzCtx.GetInst error fetching session-key:%s (%s)", sessionKey, err)
		c.WriteInternalError(rw)
		return nil, err
	}

	if session.Values[dataKey] == nil {
		log.Printf("AuthPlzCtx.GetInst error no dataKey: %s found in session: %s", dataKey, sessionKey)
		c.WriteInternalError(rw)
		return nil, err
	}

	return session.Values[dataKey], nil
}

// GetNamedSession fetches a session by name
func (c *AuthPlzCtx) GetNamedSession(rw web.ResponseWriter, req *web.Request, sessionKey string) (*sessions.Session, error) {
	session, err := c.Global.SessionStore.Get(req.Request, sessionKey)
	if err != nil {
		log.Printf("AuthPlzCtx.GetInst error fetching session-key:%s (%s)", sessionKey, err)
		return nil, err
	}
	return session, nil
}
