file-system を使用して、ファイルを直接読み書きして作業を進めてください。
回答は日本語で行ってください。
"\\wsl.localhost\Ubuntu-22.04\home\moai\instant-backlog"

- この配下で実装している Go プロジェクトの概要と構成を理解してください。
- 以下のような problem があるので、すべて解消します。

```
[{
	"resource": "/home/moai/instant-backlog/internal/commands/rename.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.ParseEpicFile",
	"source": "compiler",
	"startLineNumber": 44,
	"startColumn": 23,
	"endLineNumber": 44,
	"endColumn": 36
},{
	"resource": "/home/moai/instant-backlog/internal/commands/rename.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.ParseIssueFile",
	"source": "compiler",
	"startLineNumber": 88,
	"startColumn": 24,
	"endLineNumber": 88,
	"endColumn": 38
},{
	"resource": "/home/moai/instant-backlog/internal/fileops/reader.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.ParseIssueFile",
	"source": "compiler",
	"startLineNumber": 25,
	"startColumn": 24,
	"endLineNumber": 25,
	"endColumn": 38
},{
	"resource": "/home/moai/instant-backlog/internal/fileops/reader.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.ParseEpicFile",
	"source": "compiler",
	"startLineNumber": 51,
	"startColumn": 23,
	"endLineNumber": 51,
	"endColumn": 36
},{
	"resource": "/home/moai/instant-backlog/internal/fileops/writer.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.GenerateMarkdown",
	"source": "compiler",
	"startLineNumber": 15,
	"startColumn": 27,
	"endLineNumber": 15,
	"endColumn": 43
},{
	"resource": "/home/moai/instant-backlog/internal/fileops/writer.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.GenerateMarkdown",
	"source": "compiler",
	"startLineNumber": 31,
	"startColumn": 27,
	"endLineNumber": 31,
	"endColumn": 43
},{
	"resource": "/home/moai/instant-backlog/internal/parser/markdown.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"severity": 8,
	"message": "raw string literal not terminated",
	"source": "syntax",
	"startLineNumber": 17,
	"startColumn": 43,
	"endLineNumber": 17,
	"endColumn": 43
},{
	"resource": "/home/moai/instant-backlog/test/integration_test.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.ParseIssueFile",
	"source": "compiler",
	"startLineNumber": 118,
	"startColumn": 23,
	"endLineNumber": 118,
	"endColumn": 37
},{
	"resource": "/home/moai/instant-backlog/test/integration_test.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.ParseEpicFile",
	"source": "compiler",
	"startLineNumber": 157,
	"startColumn": 22,
	"endLineNumber": 157,
	"endColumn": 35
},{
	"resource": "/home/moai/instant-backlog/test/integration_test.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.GenerateMarkdown",
	"source": "compiler",
	"startLineNumber": 246,
	"startColumn": 22,
	"endLineNumber": 246,
	"endColumn": 38
},{
	"resource": "/home/moai/instant-backlog/test/integration_test.go",
	"owner": "_generated_diagnostic_collection_name_#1",
	"code": {
		"value": "UndeclaredImportedName",
		"target": {
			"$mid": 1,
			"path": "/golang.org/x/tools/internal/typesinternal",
			"scheme": "https",
			"authority": "pkg.go.dev",
			"fragment": "UndeclaredImportedName"
		}
	},
	"severity": 8,
	"message": "undefined: parser.GenerateMarkdown",
	"source": "compiler",
	"startLineNumber": 260,
	"startColumn": 23,
	"endLineNumber": 260,
	"endColumn": 39
}]
```

- まず、コードを読み、それぞれの問題の原因を特定してください。
- その後、修正を行ってください。
