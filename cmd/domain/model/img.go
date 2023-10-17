package model

type Img struct {
	id          string
	url         string
	isThumbnail bool
}

func NewImg(id, url string, isThumbnail bool) *Img {
	return &Img{
		id:          id,
		url:         url,
		isThumbnail: isThumbnail,
	}
}

func (i *Img) Id() string {
	return i.id
}

func (i *Img) URL() string {
	return i.url
}

func (i *Img) IsThumbnail() bool {
	return i.isThumbnail
}

func (i *Img) SetId(id string) {
	i.id = id
}

func (i *Img) SetURL(url string) {
	i.url = url
}

func (i *Img) SetIsThumbnail(isThumbnail bool) {
	i.isThumbnail = isThumbnail
}
