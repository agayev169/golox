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

    with open(path, "w") as f:
        # package
        f.writelines(
            [
                "package golox\n",
                "\n",
            ]
        )

        # interface
        f.writelines(
            [
                f"type {base_name} interface {{\n",
                "   Accept(v *Visitor)\n",
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
                + ["    " + " ".join(param) + "\n" for param in params]
                + ["}\n", "\n"]
            )

            lines = lines + [
                f"func ({name[0].lower()} *{name}) Accept(v Visitor) {{\n",
                f"    v.Accept{name}{base_name}({name[0].lower()})\n",
                "}\n",
            ]

            lines += ["\n"]
            f.writelines(lines)

            visitor_methods += [f"    Accept{name}{base_name}(*{name})\n"]

        # Visitor
        f.writelines(
            [
                "// ================ Visitor ================\n",
                "\n",
                "type Visitor interface {\n",
            ]
            + visitor_methods
            + ["}"]
        )
    
    os.system(f"gofmt -w {path}")


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
            ("Binary", [("left", "Expr"), ("operator", "Token"), ("right", "Expr")]),
            ("Grouping", [("expr", "Expr")]),
            ("Literal", [("value", "interface{}")]),
            ("Unary", [("operator", "Token"), ("right", "Expr")]),
        ],
    )
