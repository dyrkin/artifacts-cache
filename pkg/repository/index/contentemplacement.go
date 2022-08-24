package index

type Emplacement struct {
	Partition string
	Path      string
	Offset    int64
	Size      int64
}

type ContentEmplacement struct {
	Emplacements []*Emplacement
}
