class _Parser:
    """Recursive descent: one method per grammar rule, each consuming
    exactly its production and leaving self.i just past it. Every
    failure names the byte position."""

    def __init__(self, s):
        self.s = s
        self.i = 0

    def skip_space(self):
        while self.i < len(self.s) and self.s[self.i] in " \t\n\r":
            self.i += 1

    def parse_value(self):
        if self.i >= len(self.s):
            raise ValueError(f"unexpected end of input at position {self.i}")
        c = self.s[self.i]
        if c == "{":
            return self.parse_object()
        if c == "[":
            return self.parse_array()
        if c == '"':
            return self.parse_string()
        if c == "-" or c.isdigit():
            return self.parse_number()
        return self.parse_literal()

    def parse_object(self):
        obj = {}
        self.i += 1
        self.skip_space()
        if self.i < len(self.s) and self.s[self.i] == "}":
            self.i += 1
            return obj
        while True:
            self.skip_space()
            if self.i >= len(self.s) or self.s[self.i] != '"':
                raise ValueError(f"expected object key at position {self.i}")
            key = self.parse_string()
            self.skip_space()
            if self.i >= len(self.s) or self.s[self.i] != ":":
                raise ValueError(f"expected ':' at position {self.i}")
            self.i += 1
            self.skip_space()
            obj[key] = self.parse_value()
            self.skip_space()
            if self.i >= len(self.s):
                raise ValueError(f"unterminated object at position {self.i}")
            if self.s[self.i] == ",":
                self.i += 1
                continue
            if self.s[self.i] == "}":
                self.i += 1
                return obj
            raise ValueError(f"expected ',' or '}}' at position {self.i}")

    def parse_array(self):
        arr = []
        self.i += 1
        self.skip_space()
        if self.i < len(self.s) and self.s[self.i] == "]":
            self.i += 1
            return arr
        while True:
            self.skip_space()
            arr.append(self.parse_value())
            self.skip_space()
            if self.i >= len(self.s):
                raise ValueError(f"unterminated array at position {self.i}")
            if self.s[self.i] == ",":
                self.i += 1
                continue
            if self.s[self.i] == "]":
                self.i += 1
                return arr
            raise ValueError(f"expected ',' or ']' at position {self.i}")

    def parse_string(self):
        start = self.i
        self.i += 1
        out = []
        while self.i < len(self.s):
            c = self.s[self.i]
            if c == '"':
                self.i += 1
                return "".join(out)
            if c == "\\":
                if self.i + 1 >= len(self.s):
                    break
                nxt = self.s[self.i + 1]
                if nxt == '"':
                    out.append('"')
                elif nxt == "\\":
                    out.append("\\")
                else:
                    raise ValueError(f"unsupported escape at position {self.i}")
                self.i += 2
                continue
            out.append(c)
            self.i += 1
        raise ValueError(f"unterminated string starting at position {start}")

    def parse_number(self):
        start = self.i
        if self.s[self.i] == "-":
            self.i += 1
        digits = 0
        while self.i < len(self.s) and self.s[self.i].isdigit():
            self.i += 1
            digits += 1
        if digits == 0:
            raise ValueError(f"malformed number at position {start}")
        return int(self.s[start:self.i])

    def parse_literal(self):
        for lit, v in (("true", True), ("false", False), ("null", None)):
            if self.s.startswith(lit, self.i):
                self.i += len(lit)
                return v
        raise ValueError(f"unexpected value at position {self.i}")


def parse(input):
    p = _Parser(input)
    p.skip_space()
    v = p.parse_value()
    p.skip_space()
    if p.i != len(p.s):
        raise ValueError(f"trailing garbage at position {p.i}")
    return v
