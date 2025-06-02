package entity

type Documento struct {
	Id            int64
	Doc           []byte
	FuncionarioId int64
}

func newDocumento(doc []byte, FuncionarioId int64) *Documento {
	d := new(Documento)
	d.Doc = doc
	d.FuncionarioId = FuncionarioId
	return d
}
