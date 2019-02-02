/*
 * Copyright (c) 2019 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package system

type SSequence baseContext

const (
	AccountIdSequence = "noah.account_ids"
)

func (sequences *SSequence) GetNextValueForSequence(sequenceName string) (*uint64, error) {
	return sequences.db.NextSequenceValueById(sequenceName)
}

func (sequences *SSequence) GetSequenceIndex(sequenceName string) (uint64, error) {
	return sequences.db.SequenceIndexById(sequenceName)
}

func (sequences *SSequence) NewAccountID() (uint64, uint64, error) {
	accountId, err := sequences.db.NextSequenceValueById(AccountIdSequence)
	if err != nil {
		return 0, 0, err
	}
	accountIdIndex, err := sequences.db.SequenceIndexById(AccountIdSequence)
	if err != nil {
		return 0, 0, err
	}
	return *accountId, accountIdIndex, nil
}
