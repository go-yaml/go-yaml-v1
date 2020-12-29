package yaml_test

// This is for testing multi-document parsing and emitting with NewDecoder and NewEncoder with the Node interface.

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
	. "gopkg.in/check.v1"
)


func walkTree(indent int, node *yaml.Node) {
	fmt.Printf("%s{%d %d %#v:%#v anchor:%#v head:%#v line:%#v foot:%#v %d:%d}\n", strings.Repeat("  ", indent), node.Kind, node.Style, node.Tag, node.Value, node.Anchor, node.HeadComment, node.LineComment, node.FootComment, node.Line, node.Column)
	for _, item := range node.Content {
		walkTree(indent + 1, item)
	}
	if node.Alias != nil {
		walkTree(indent + 1, node.Alias)
	}
}


func testMDCycle(c *C, input string, expectedOutput string) {
	fmt.Printf("New multi-doc input:\n")
	var nodes []*yaml.Node
	d := yaml.NewDecoder(bytes.NewReader([]byte(input)))
	for true {
		node := &yaml.Node{}
		err := d.Decode(node)
		if err == io.EOF {
			break
		}
		c.Assert(err, IsNil)
		nodes = append(nodes, node)
		walkTree(0, node)
	}
	var b bytes.Buffer
	e := yaml.NewEncoder(io.Writer(&b))
	e.SetIndent(4)
	for _, node := range nodes {
		err := e.Encode(node)
		c.Assert(err, IsNil)
	}
	e.Close()
	out := b.Bytes()
	c.Assert(string(out), DeepEquals, expectedOutput)
	if len(expectedOutput) == 0 {
		c.Assert(out, DeepEquals, []byte(nil))
	} else {
		c.Assert(out, DeepEquals, []byte(expectedOutput))
	}
}


func testMDIdempotent(c *C, data string) {
	testMDCycle(c, data, data)
}


func (s *S) TestCommentDocSkip(c *C) {
	testMDIdempotent(c, `key: value

# foo
---
key: value
`)
	testMDIdempotent(c, `# foo
---
key: value
`)
	testMDIdempotent(c, `# foo
---
# bar
key: value
`)
}
