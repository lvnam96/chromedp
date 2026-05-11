package chromedp

import (
	"testing"

	"github.com/chromedp/cdproto/cdp"
)

// FuzzAttributeModified fuzzes the attributeModified node operation with
// arbitrary name/value pairs to detect panics or unexpected behavior.
func FuzzAttributeModified(f *testing.F) {
	f.Add("class", "foo bar")
	f.Add("id", "main")
	f.Add("href", "https://example.com")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, name, value string) {
		n := &cdp.Node{}
		op := attributeModified(name, value)
		op(n)

		// Verify invariant: attributes slice has even length (name/value pairs).
		if len(n.Attributes)%2 != 0 {
			t.Fatalf("attributes length is odd: %d", len(n.Attributes))
		}

		// Apply again; the existing attribute should be updated, not appended.
		op(n)
		if len(n.Attributes) != 2 {
			t.Fatalf("expected 2 attributes after second apply, got %d", len(n.Attributes))
		}
	})
}

// FuzzAttributeRemoved fuzzes the attributeRemoved node operation to ensure it
// does not panic and keeps the attributes slice length even.
func FuzzAttributeRemoved(f *testing.F) {
	f.Add("class")
	f.Add("id")
	f.Add("")

	f.Fuzz(func(t *testing.T, name string) {
		n := &cdp.Node{Attributes: []string{"class", "foo", "id", "bar"}}
		op := attributeRemoved(name)
		op(n)

		if len(n.Attributes)%2 != 0 {
			t.Fatalf("attributes length is odd after removal: %d", len(n.Attributes))
		}
	})
}

// FuzzInsertNode fuzzes the insertNode helper with arbitrary node IDs.
func FuzzInsertNode(f *testing.F) {
	f.Add(int64(0), int64(1))
	f.Add(int64(5), int64(5))
	f.Add(int64(1), int64(0))

	f.Fuzz(func(t *testing.T, existingID, prevID int64) {
		existing := &cdp.Node{NodeID: cdp.NodeID(existingID)}
		child := &cdp.Node{NodeID: cdp.NodeID(99)}
		nodes := []*cdp.Node{existing}
		result := insertNode(nodes, cdp.NodeID(prevID), child)
		if len(result) != 2 {
			t.Fatalf("expected 2 nodes, got %d", len(result))
		}
	})
}

// FuzzRemoveNode fuzzes the removeNode helper with arbitrary node IDs.
func FuzzRemoveNode(f *testing.F) {
	f.Add(int64(1), int64(1))
	f.Add(int64(1), int64(2))
	f.Add(int64(0), int64(0))

	f.Fuzz(func(t *testing.T, existingID, removeID int64) {
		existing := &cdp.Node{NodeID: cdp.NodeID(existingID)}
		nodes := []*cdp.Node{existing}
		result := removeNode(nodes, cdp.NodeID(removeID))
		if existingID == removeID {
			if len(result) != 0 {
				t.Fatalf("expected empty slice after removing matching node, got %d", len(result))
			}
		} else {
			if len(result) != 1 {
				t.Fatalf("expected 1 node after removing non-matching node, got %d", len(result))
			}
		}
	})
}
