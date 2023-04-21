package prompt

type Progress struct {
	frame   int
	animate []string
}

func NewProgress(animate []string) *Progress {
	return &Progress{
		animate: animate,
		frame:   -1,
	}
}

func (p *Progress) Next() string {
	p.frame = (p.frame + 1) % len(p.animate)
	return p.animate[p.frame]
}
