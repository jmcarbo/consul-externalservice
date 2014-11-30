package consul_externalservice

import (
	. "github.com/franela/goblin"
	"testing"
)

func TestLock(t *testing.T) {

	g := Goblin(t)

	g.Describe("lock", func() {
    g.Before(func(){
      initConsul()
    })
    g.After(func(){
      stopConsul()
    })

		g.It("create", func() {
			client := Connect()
			lock := NewLock(client, "testlock1")
			g.Assert(lock != nil).IsTrue()
			//g.Assert(lock.IsLeader()).IsTrue()
		})
		g.It("can be locked", func() {
			client := Connect()
			lock := NewLock(client, "testlock2")
			err := lock.Lock(nil)
			g.Assert(err == nil).IsTrue()
			g.Assert(lock.IsLeader()).IsTrue()
			lock.Destroy()
		})

		g.It("can be locked and unlocked", func() {
			client := Connect()
			lock := NewLock(client, "testlock3")
			err := lock.Lock(nil)
			g.Assert(err == nil).IsTrue()
			g.Assert(lock.IsLeader()).IsTrue()
			err = lock.Unlock()
			g.Assert(err == nil).IsTrue()
			g.Assert(lock.IsUnlocked()).IsTrue()
			g.Assert(lock.IsLeader()).IsFalse()
			lock.Destroy()
		})

		g.It("can't lock a locked key", func() {
			client := Connect()
			lock1 := NewLock(client, "testlockA")
			lock1.Lock(nil)
			lock2 := NewLock(client, "testlockA")
			err := lock2.Lock(nil)
			g.Assert(err != nil).IsTrue()
			g.Assert(lock2.IsLeader()).IsFalse()
			lock1.Destroy()
			lock2.Destroy()
		})
		g.It("can't be unlocked when not locked", func() {
			client := Connect()
			lock := NewLock(client, "testlock")
			err := lock.Unlock()
			g.Assert(err == nil).IsTrue()
			g.Assert(lock.IsUnlocked()).IsTrue()
			lock.Destroy()
		})
	})
}
