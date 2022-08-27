package index

type Emplacement struct {
	Partition string
	Name      string
	Offset    int64
	Size      int64
}

type ContentEmplacement struct {
	Emplacements []*Emplacement
}
