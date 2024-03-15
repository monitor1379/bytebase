package pg

// Framework code is generated by the generator.

import (
	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/plugin/advisor"
	"github.com/bytebase/bytebase/backend/plugin/parser/sql/ast"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

var (
	_ advisor.Advisor = (*StatementDisallowMixDdlDmlAdvisor)(nil)
)

func init() {
	advisor.Register(storepb.Engine_POSTGRES, advisor.PostgreSQLStatementDisallowMixDDLDML, &StatementDisallowMixDdlDmlAdvisor{})
}

// StatementDisallowMixDdlDmlAdvisor is the advisor checking for disallow mix DDL and DML.
type StatementDisallowMixDdlDmlAdvisor struct {
}

// Check checks for disallow mix DDL and DML.
func (*StatementDisallowMixDdlDmlAdvisor) Check(ctx advisor.Context, _ string) ([]advisor.Advice, error) {
	stmtList, ok := ctx.AST.([]ast.Node)
	if !ok {
		return nil, errors.Errorf("failed to convert to Node")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	title := string(ctx.Rule.Type)

	var hasDDL, hasDML bool
	for _, stmt := range stmtList {
		if _, ok := stmt.(ast.DDLNode); ok {
			hasDDL = true
		}
		if _, ok := stmt.(ast.DMLNode); ok {
			hasDML = true
		}
		if hasDDL && hasDML {
			break
		}
	}

	if hasDDL && hasDML {
		return []advisor.Advice{
			{
				Status:  level,
				Title:   title,
				Content: "Mixing DDL with DML is not allowed",
				Code:    advisor.StatementDisallowMixDDLDML,
			},
		}, nil
	}

	return []advisor.Advice{
		{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		},
	}, nil
}
