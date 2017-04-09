package service

import "github.com/goushuyun/weixin-golang/pb"

func integreteInfo(to, from *pb.Book) *pb.Book {
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
	if book.Publisher != "" && book.Author != "" && book.Isbn != "" && book.Title != "" && book.Image != "" {
		return true
	} else {
		return false
	}
}
