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

type BSTKind uint16

const (
	KindSplayTree BSTKind = iota
	KindTreap
	KindRedBlackTree
	KindScapegoatTree
	KindBTree
	KindSkipList
)

// This interface describes the main methods of general BST.
type BinarySearchTree[K any, V any] interface {
	// The number of elements in the current tree.
	Len() int
	// Search for the element corresponding to the specified key.
	Search(K) (V, bool)
	// Insert or update elements.
	Insert(K, V)
	// Delete specified element.
	Delete(K) (V, bool)
	//Comparing elements.
	Compare(K, K) int
	//Return the minimum element.
	Min() (K, V, bool)
	// Return the maximum element.
	Max() (K, V, bool)
	// Create an Iterator
	Cursor() Cursor[K, V]
	// Clear the current BST, which means deleting all elements.
	Clean()
}

// Create a BST object based on its type.
func NewBinarySearchTree[K any, V any](comp CompareFunc[K], kind BSTKind) (BinarySearchTree[K, V], bool) {
	switch kind {
	case KindBTree:
		cfg := &BTreeConfig{MinDegree: DefaultBTreeMinDegree}
		return NewBTree[K, V](cfg, comp), true
	case KindTreap:
		return NewTreap[K, V](comp), true
	case KindRedBlackTree:
		return NewRedBlackTree[K, V](comp), true
	case KindScapegoatTree:
		return NewScapegoatTree[K, V](DefaultScapegoatTreeAlpha, comp), true
	case KindSplayTree:
		return NewSplayTree[K, V](comp), true
	case KindSkipList:
		cfg := &SkipListConfig{}
		return NewSkipList[K, V](cfg, comp), true
	default:
		return nil, false
	}
}
