package comment

import (
	"go/ast"
	"go/token"
	"strings"
)

// Maps is slice of ast.CommentMap.
type Maps []ast.CommentMap

// New creates new a CommentMap slice from specified files.
func New(fset *token.FileSet, files []*ast.File) Maps {
	maps := make(Maps, len(files))
	for i := range files {
		maps[i] = ast.NewCommentMap(fset, files[i], files[i].Comments)
	}
	return maps
}

// Comments returns correspond a CommentGroup slice to specified AST node.
func (maps Maps) Comments(n ast.Node) []*ast.CommentGroup {
	for i := range maps {
		if maps[i][n] != nil {
			return maps[i][n]
		}
	}
	return nil
}

// Annotated checks either specified AST node is annotated or not.
func (maps Maps) Annotated(n ast.Node, annotation string) bool {
	for _, cg := range maps.Comments(n) {
		if strings.HasPrefix(strings.TrimSpace(cg.Text()), annotation) {
			return true
		}
	}
	return false
}

// Ignore checks either specified AST node is ignored by the check.
// It follows staticcheck style as the below.
//   //lint:ignore Check1[,Check2,...,CheckN] reason
func (maps Maps) Ignore(n ast.Node, check string) bool {
	for _, cg := range maps.Comments(n) {
		if hasIgnoreCheck(cg.Text(), check) {
			return true
		}
	}
	return false
}

// IgnorePos checks either specified postion of AST node is ignored by the check.
// It follows staticcheck style as the below.
//   //lint:ignore Check1[,Check2,...,CheckN] reason
func (maps Maps) IgnorePos(pos token.Pos, check string) bool {

	for n, cg := range maps {
		if n.Pos() != pos {
			continue
		}

		if hasIgnoreCheck(cg.Text(), check) {
			return true
		}
	}

	return false
}

func hasIgnoreCheck(s, check string) bool {
	txt := strings.Split(s, " ")
	if len(txt) < 3 && txt[0] != "lint:ignore" {
		continue
	}
	checks := strings.Split(txt[1], ",")
	for i := range checks {
		if check == checks[i] {
			return true
		}
	}
	return false
}
