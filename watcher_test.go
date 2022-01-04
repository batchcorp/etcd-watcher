// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdwatcher

import (
	"log"
	"testing"

	"github.com/batchcorp/casbin/v2"
)

func updateCallback(rev string) {
	log.Println("New revision detected:", rev)
}

func TestWatcher(t *testing.T) {
	// updater represents the Casbin enforcer instance that changes the policy in DB.
	// Use the endpoints of etcd cluster as parameter.
	updater, _ := NewWatcher([]string{"http://127.0.0.1:2379"}, "/casbin", nil)

	// listener represents any other Casbin enforcer instance that watches the change of policy in DB.
	listener, _ := NewWatcher([]string{"http://127.0.0.1:2379"}, "/casbin", nil)
	// listener should set a callback that gets called when policy changes.
	_ = listener.SetUpdateCallback(updateCallback)

	// updater changes the policy, and sends the notifications.
	err := updater.Update()
	if err != nil {
		panic(err)
	}

	// Now the listener's callback updateCallback() should be called,
	// because it receives the notification of policy update.
	// You should see "[New revision detected: X]" in the log.
}

func TestWithEnforcer(t *testing.T) {
	// Initialize the watcher.
	// Use the endpoints of etcd cluster as parameter.
	w, _ := NewWatcher([]string{"http://127.0.0.1:2379"}, "/casbin", nil)

	// Initialize the enforcer.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// Set the watcher for the enforcer.
	_ = e.SetWatcher(w)

	// By default, the watcher's callback is automatically set to the
	// enforcer's LoadPolicy() in the SetWatcher() call.
	// We can change it by explicitly setting a callback.
	_ = w.SetUpdateCallback(updateCallback)

	// Update the policy to test the effect.
	// You should see "[New revision detected: X]" in the log.
	err := e.SavePolicy()
	if err != nil {
		t.Fail()
	}
}

// Uncomment if you want to perform TLS tests
//var (
//	authConfig = &AuthConfig{
//		UseTLS:     true,
//		CACert:     `fill this in`,
//		ClientKey:  `fill this in`,
//		ClientCert: `fill this in`,
//	}
//)
//
//func TestWatcherTLS(t *testing.T) {
//	// updater represents the Casbin enforcer instance that changes the policy in DB.
//	// Use the endpoints of etcd cluster as parameter.
//	updater, _ := NewWatcher([]string{"https://etcd-1.sfo3.dev.batch.sh:2379"}, "/casbin-tls", authConfig)
//
//	// listener represents any other Casbin enforcer instance that watches the change of policy in DB.
//	listener, _ := NewWatcher([]string{"https://etcd-1.sfo3.dev.batch.sh:2379"}, "/casbin-tls", authConfig)
//	// listener should set a callback that gets called when policy changes.
//	_ = listener.SetUpdateCallback(updateCallback)
//
//	// updater changes the policy, and sends the notifications.
//	err := updater.Update()
//	if err != nil {
//		panic(err)
//	}
//
//	// Now the listener's callback updateCallback() should be called,
//	// because it receives the notification of policy update.
//	// You should see "[New revision detected: X]" in the log.
//}
//
//func TestWithEnforcerTLS(t *testing.T) {
//	// Initialize the watcher.
//	// Use the endpoints of etcd cluster as parameter.
//	w, _ := NewWatcher([]string{"https://etcd-1.sfo3.dev.batch.sh:2379"}, "/casbin-tls", authConfig)
//
//	// Initialize the enforcer.
//	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")
//
//	// Set the watcher for the enforcer.
//	_ = e.SetWatcher(w)
//
//	// By default, the watcher's callback is automatically set to the
//	// enforcer's LoadPolicy() in the SetWatcher() call.
//	// We can change it by explicitly setting a callback.
//	_ = w.SetUpdateCallback(updateCallback)
//
//	// Update the policy to test the effect.
//	// You should see "[New revision detected: X]" in the log.
//	err := e.SavePolicy()
//	if err != nil {
//		t.Fail()
//	}
//}
