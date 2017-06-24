package service

import "github.com/goushuyun/weixin-golang/pb"

func integreteInfo(to, from *pb.Book) *pb.Book {
	if to == nil && from != nil {
		return from
	} else if to == nil && from == nil {
		return nil
	}

	if to.Title == "" {
		to.Title = from.Title
	}

	if to.Price == 0 {
		to.Price = from.Price
	}

	if to.Author == "" {
		to.Author = from.Author
	}

	if to.Publisher == "" {
		to.Publisher = from.Publisher
	}

	if to.Image == "" {
		to.Image = from.Image
	}

	if to.Pubdate == "" {
		to.Pubdate = from.Pubdate
	}

	if to.Subtitle == "" {
		to.Subtitle = from.Subtitle
	}

	return to
}

func bookInfoIsOk(book *pb.Book) bool {
	if book == nil {
		return false
	}

	if book.Price != 0 && book.Isbn != "" && book.Title != "" {
		return true
	} else {
		return false
	}
}
