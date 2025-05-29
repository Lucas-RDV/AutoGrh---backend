package entity

type Documento struct {
	Id  string
	Doc []byte
}

func newDocumento(doc []byte) *Documento {
	d := new(Documento)
	d.Doc = doc
	return d
}
