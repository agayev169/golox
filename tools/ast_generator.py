import os
import sys
from typing import List, Tuple


def define_ast(
    out_dir: str, filename: str, base_name: str, subclasses: List[Tuple[str, List[Tuple[str, str]]]]
) -> None:
    try:
        os.mkdir(out_dir)
    except FileExistsError:
        pass

    path = os.path.join(out_dir, filename)

    visitor_name = f"{base_name}Visitor"

    with open(path, "w") as f:
        # package
        f.writelines(
            [
                "package golox\n",
                "\n",
                "import \"fmt\"\n",
                "\n"
            ]
        )

        # interface
        f.writelines(
            [
                f"type {base_name} interface {{\n",
                f"   Accept(v {visitor_name}) (interface{{}}, *LoxError)\n",
                "}\n",
                "\n",
            ]
        )

        visitor_methods = []

        # subclasses
        for subclass in subclasses:
            name, params = subclass

            lines = (
                [
                    f"// ================ {name} ================\n",
                    "\n",
                    f"type {name} struct {{\n",
                ]
                + ["    " +
                    f"{publicize(param[0])} {param[1]}\n" for param in params]
                + ["}\n", "\n"]
            )

            varName = name[0].lower()

            if varName == "v":
                varName = name[:2].lower()

            lines = lines + [
                f"func ({varName} *{name}) Accept(v {visitor_name}) (interface{{}}, *LoxError) {{\n",
                f"    return v.Accept{name}{base_name}({varName})\n",
                "}\n",
            ]

            lines += ["\n"]

            lines = lines + [
                f'func ({varName} *{name}) String() string {{\n'
                f'  return fmt.Sprintf("({name}): {{{"; ".join(f"{p[0]}:  %v" for p in params)}}}", {",".join(map(lambda p: f"{varName}.{publicize(p[0])}", params))})\n'
                '}\n'
            ]

            lines += ["\n"]

            f.writelines(lines)

            visitor_methods += [
                f"    Accept{name}{base_name}(*{name}) (interface{{}}, *LoxError)\n"]

        # Visitor
        f.writelines(
            [
                f"// ================ {visitor_name} ================\n",
                "\n",
                f"type {visitor_name} interface {{\n",
            ]
            + visitor_methods
            + ["}"]
        )

    os.system(f"gofmt -w {path}")


def publicize(s: str) -> str:
    return s[0].upper() + s[1:]


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <output directory>")
        exit(64)

    out_dir = sys.argv[1]

    define_ast(
        out_dir,
        "expr.go",
        "Expr",
        [
            ("Assign", [("name", "Token"), ("value", "Expr")]),
            ("Binary", [("left", "Expr"),
             ("operator", "Token"), ("right", "Expr")]),
            ("Grouping", [("expr", "Expr")]),
            ("Literal", [("value", "interface{}")]),
            ("Unary", [("operator", "Token"), ("right", "Expr")]),
            ("Call", [("callee", "Expr"), ("paren", "Token"), ("args", "[]Expr")]),
            ("Get", [("obj", "Expr"), ("name", "Token")]),
            ("Set", [("obj", "Expr"), ("name", "Token"), ("value", "Expr")]),
            ("Variable", [("name", "Token")]),
            ("Logical", [("left", "Expr"),
                         ("operator", "Token"), ("right", "Expr")]),
        ],
    )

    define_ast(
        out_dir,
        "stmt.go",
        "Stmt",
        [
            ("Block", [("stmts", "[]Stmt")]),
            ("Expression", [("expr", "Expr")]),
            ("Print", [("expr", "Expr")]),
            ("Var", [("name", "Token"), ("initializer", "Expr")]),
            ("Class", [("name", "Token"), ("methods", "[]Func")]),
            ("Func", [("name", "Token"), ("params", "[]Token"),
                      ("body", "[]Stmt")]),
            ("If", [("condition", "Expr"), ("body", "Stmt"),
                    ("elseBody", "Stmt")]),
            ("While", [("condition", "Expr"), ("body", "Stmt")]),
            ("Return", [("keyword", "Token"), ("value", "Expr")])
        ],
    )
