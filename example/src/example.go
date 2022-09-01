package src

type Example struct {
	Var1 int    `db:"var1"`
	Var2 string `db:"var2"`
	Var3 bool   `db:"var3"`
	Var4 bool   `db:"var4"`
}

type Example2 struct {
	Var5 int    `db:"var5"`
	Var6 string `db:"var6"`
	Var7 bool   `db:"var7"`
	Var8 bool   `db:"var8"`
}

type JoinExample struct {
	Var1  int     `db:"var1"`
	Var9  int     `db:"var9"`
	Var10 string  `db:"var10"`
	Var11 bool    `db:"var11"`
	Var12 bool    `db:"var12"`
	Var13 Example `db:"var13"`
}
