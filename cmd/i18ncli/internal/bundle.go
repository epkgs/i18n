package internal

// Bundle 存储bundle使用信息
type Bundle struct {
	Name  string              // bundle名称
	Trans map[string]struct{} // 翻译键集合，使用map[string]struct{}提高效率
	Vars  map[string]*VarInfo // 变量信息映射，PackagePath:VarName -> VarInfo
}

// VarInfo 存储变量的详细信息
type VarInfo struct {
	Pkg      string // 包的完整路径
	Name     string // 变量名
	FilePath string // 文件路径
}

func NewBundle(name string) *Bundle {
	return &Bundle{
		Name:  name,
		Trans: make(map[string]struct{}),
		Vars:  make(map[string]*VarInfo),
	}
}

func (b *Bundle) AddVarDefine(varName, pkg, filePath string) {
	b.Vars[pkg+"."+varName] = &VarInfo{
		Pkg:      pkg,
		Name:     varName,
		FilePath: filePath,
	}
}

func (b *Bundle) AddTrans(key string) {
	b.Trans[key] = struct{}{}
}
