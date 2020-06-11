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

package casbin

import (
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

const (
	notImplemented = "not implemented"
)

// addPolicy adds a rule to the current policy.
func (e *Enforcer) addPolicy(sec string, ptype string, rule []string) (bool, error) {
	persistFn := func() error {
		if e.adapter != nil && e.autoSave {
			if err := e.adapter.AddPolicy(sec, ptype, rule); err != nil {
				if err.Error() != notImplemented {
					return err
				}
			}
		}
		return nil
	}

	ruleAdded, err := e.model.AddPolicy(persistFn, sec, ptype, rule)
	if err != nil {
		return false, err
	}
	if !ruleAdded {
		return ruleAdded, nil
	}

	if sec == "g" {
		err := e.BuildIncrementalRoleLinks(model.PolicyAdd, ptype, [][]string{rule})
		if err != nil {
			return ruleAdded, err
		}
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForAddPolicy(rule...)
		} else {
			err = e.watcher.Update()
		}
		return ruleAdded, err
	}

	return ruleAdded, nil
}

// addPolicies adds rules to the current policy.
func (e *Enforcer) addPolicies(sec string, ptype string, rules [][]string) (bool, error) {
	persistFn := func() error {
		if e.adapter != nil && e.autoSave {
			if err := e.adapter.(persist.BatchAdapter).AddPolicies(sec, ptype, rules); err != nil {
				if err.Error() != notImplemented {
					return err
				}
			}
		}
		return nil
	}

	rulesAdded, err := e.model.AddPolicies(persistFn, sec, ptype, rules)
	if err != nil {
		return false, err
	}
	if !rulesAdded {
		return rulesAdded, nil
	}

	if sec == "g" {
		err := e.BuildIncrementalRoleLinks(model.PolicyAdd, ptype, rules)
		if err != nil {
			return rulesAdded, err
		}
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		err := e.watcher.Update()
		if err != nil {
			return rulesAdded, err
		}
	}

	return rulesAdded, nil
}

// removePolicy removes a rule from the current policy.
func (e *Enforcer) removePolicy(sec string, ptype string, rule []string) (bool, error) {
	persistFn := func() error {
		if e.adapter != nil && e.autoSave {
			if err := e.adapter.RemovePolicy(sec, ptype, rule); err != nil {
				if err.Error() != notImplemented {
					return err
				}
			}
		}
		return nil
	}

	ruleRemoved, err := e.model.RemovePolicy(persistFn, sec, ptype, rule)
	if err != nil {
		return false, err
	}
	if !ruleRemoved {
		return ruleRemoved, nil
	}

	if sec == "g" {
		err := e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, [][]string{rule})
		if err != nil {
			return ruleRemoved, err
		}
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForRemovePolicy(rule...)
		} else {
			err = e.watcher.Update()
		}
		return ruleRemoved, err

	}

	return ruleRemoved, nil
}

// removePolicies removes rules from the current policy.
func (e *Enforcer) removePolicies(sec string, ptype string, rules [][]string) (bool, error) {
	persistFn := func() error {
		if e.adapter != nil && e.autoSave {
			if err := e.adapter.(persist.BatchAdapter).RemovePolicies(sec, ptype, rules); err != nil {
				if err.Error() != notImplemented {
					return err
				}
			}
		}
		return nil
	}

	rulesRemoved, err := e.model.RemovePolicies(persistFn, sec, ptype, rules)
	if err != nil {
		return false, err
	}
	if !rulesRemoved {
		return rulesRemoved, nil
	}

	if sec == "g" {
		err := e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, rules)
		if err != nil {
			return rulesRemoved, err
		}
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		err := e.watcher.Update()
		if err != nil {
			return rulesRemoved, err
		}
	}

	return rulesRemoved, nil
}

// removeFilteredPolicy removes rules based on field filters from the current policy.
func (e *Enforcer) removeFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (bool, error) {
	persistFn := func() error {
		if e.adapter != nil && e.autoSave {
			if err := e.adapter.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...); err != nil {
				if err.Error() != notImplemented {
					return err
				}
			}
		}
		return nil
	}

	ruleRemoved, effects, err := e.model.RemoveFilteredPolicy(persistFn, sec, ptype, fieldIndex, fieldValues...)
	if err != nil {
		return false, err
	}
	if !ruleRemoved {
		return ruleRemoved, nil
	}

	if sec == "g" {
		err := e.BuildIncrementalRoleLinks(model.PolicyRemove, ptype, effects)
		if err != nil {
			return ruleRemoved, err
		}
	}

	if e.watcher != nil && e.autoNotifyWatcher {
		var err error
		if watcher, ok := e.watcher.(persist.WatcherEx); ok {
			err = watcher.UpdateForRemoveFilteredPolicy(fieldIndex, fieldValues...)
		} else {
			err = e.watcher.Update()
		}
		return ruleRemoved, err
	}

	return ruleRemoved, nil
}
