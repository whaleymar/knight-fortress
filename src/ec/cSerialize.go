package ec

type CSerialize struct {
	FileName string
}

func (comp *CSerialize) getType() ComponentType {
	return CMP_SERIALIZE
}
