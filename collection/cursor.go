/*
	Copyright (C) 2023 flxj(https://github.com/flxj)

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package collection

// Cursors are used to access ordered collections.
type Cursor[K, V any] interface {
	// If the collection object is in concurrent security mode,
	// the Open method needs to be called to attempt locking before using the cursor.
	// After use, the Close method must be called to release the lock.
	Open() error
	// release cursor resources.
	Close()
	// The Seek(key) method locates the cursor at the key.
	// If the key does not exist, it locates at the next key and returns it
	Seek(K) (K, V, bool)
	// The First method locates the cursor at the minimum element of the set.
	// If there is no minimum element (the set is empty), it returns false
	First() (K, V, bool)
	// The Last method locates the cursor at the maximum element of the set.
	// If there is no maximum element (the set is empty), it returns false.
	Last() (K, V, bool)
	// HasNext returns whether the next element exists relative to the current cursor position.
	HasNext() bool
	// Next() moves the cursor backwards and returns the element.
	// If the element does not exist, it returns a type zero value.
	Next() (K, V)
	// HasPrev() returns whether the previous element exists relative to the current cursor position.
	HasPrev() bool
	// Prev() moves the cursor forward and returns the element.
	// If the element does not exist, it returns a type value of zero.
	Prev() (K, V)
}
