package model

type Img struct {
	id  string
	url string
}

func NewImg(id, url string) *Img {
	return &Img{
		id:  id,
		url: url,
	}
}

func (i *Img) Id() string {
	return i.id
}

func (i *Img) URL() string {
	return i.url
}
