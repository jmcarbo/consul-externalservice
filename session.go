package consul_externalservice

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/armon/consul-api"
	"github.com/nu7hatch/gouuid"
	"time"
)

type Session struct {
	sessionID   string
	sessionName string
	client      *consulapi.Client
	stopCh      chan struct{}
	doneCh      chan struct{}
}

func NewSession(client *consulapi.Client, name string) *Session {
	if name == "" {
		uid, _ := uuid.NewV4()
		name = uid.String()
	}

	sess := &Session{sessionName: name, client: client}
	sess.doneCh = make(chan struct{})
	sess.stopCh = make(chan struct{})
	serviceCheck := consulapi.AgentServiceCheck{TTL: "10s"}
	err := client.Agent().CheckRegister(&consulapi.AgentCheckRegistration{sess.checkID(), sess.checkID(), "",
		serviceCheck})
	if err != nil {
		log.Error("Cannot register check")
		return nil
	}
	err = sess.maintainCheck()
	if err != nil {
		log.Error("Cannot maintain check")
		return nil
	}

	node, _ := client.Agent().NodeName()
	sess.sessionID, _, err = client.Session().Create(&consulapi.SessionEntry{ID: sess.sessionName, Name: sess.sessionName, Node: node,
		Checks: []string{sess.checkID(), "serfHealth"}}, nil)
	return sess
}

func (sess *Session) checkID() string {
	return fmt.Sprintf("session:%s", sess.sessionName)
}

func (sess *Session) IsHealthy() bool {
	se, _, err := sess.client.Session().Info(sess.sessionID, nil)
	if err != nil {
		return false
	}
	if se == nil {
		log.Error("Cannot find session ", sess.sessionName)
		return false
	}
	checks, err := sess.client.Agent().Checks()
	if err != nil {
		log.Error("Cannot find checks ", err)
		return false
	}
	for _, s := range se.Checks {
		if checks[s] != nil && checks[s].Status != "passing" {
			return false
		}
	}
	return true
}

func (sess *Session) maintainCheck() error {
	// Try to update the TTL immediately
	agent := sess.client.Agent()
	checkID := sess.checkID()
	if err := agent.PassTTL(checkID, ""); err != nil {
		log.Error(err)
		return err
	}

	// Run the update in the background
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := agent.PassTTL(checkID, ""); err != nil {
					log.Printf("[ERR] Failed to update check TTL: %v", err)
				}
			case <-sess.doneCh:
				return
			}
		}
	}()

	return nil
}

func (sess *Session) Destroy() error {
	close(sess.doneCh)
	err := sess.client.Agent().CheckDeregister(sess.checkID())
	if err != nil {
		return err
	}
	_, err = sess.client.Session().Destroy(sess.sessionID, nil)
	if err != nil {
		return err
	}
	return nil
}

func (sess *Session) IsDestroyed() bool {
	se, _, err := sess.client.Session().Info(sess.sessionID, nil)
	if err != nil {
		return false
	}
	if se != nil {
		return false
	}
	return true
}
