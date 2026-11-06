package filestream

import "maps"

// Progress is what the server has acknowledged for a run's filestream.
type Progress struct {
	// AckedOffsets maps each filestream file name to the offset just past
	// the last line the server acknowledged.
	AckedOffsets map[string]int

	// CompleteAcked is whether the server acknowledged the request that
	// marked the run complete.
	CompleteAcked bool
}

// addAcked records that the server accepted the request data.
func (p *Progress) addAcked(data *FileStreamRequestJSON) {
	if data == nil {
		return
	}

	if data.Complete != nil && *data.Complete {
		p.CompleteAcked = true
	}

	for name, chunk := range data.Files {
		if len(chunk.Content) == 0 {
			continue
		}
		if p.AckedOffsets == nil {
			p.AckedOffsets = make(map[string]int)
		}
		p.AckedOffsets[name] = max(
			p.AckedOffsets[name],
			chunk.Offset+len(chunk.Content),
		)
	}
}

// clone returns a copy that does not share the offsets map.
func (p Progress) clone() Progress {
	p.AckedOffsets = maps.Clone(p.AckedOffsets)
	return p
}
