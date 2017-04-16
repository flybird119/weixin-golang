package db

import (
	"database/sql"
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/wothing/log"
)

func GetBookInfo(book *pb.Book) error {
	query := "select title, price, author, publisher, pubdate, subtitle, image, summary, author_intro, isbn from books where id = $1"
	err := DB.QueryRow(query, book.Id).Scan(&book.Title, &book.Price, &book.Author, &book.Publisher, &book.Pubdate, &book.Subtitle, &book.Image, &book.Summary, &book.AuthorIntro, &book.Isbn)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func GetBookInfoByISBN(book *pb.Book) error {
	query := "select id, title, price, author, publisher, pubdate, subtitle, image, summary, author_intro from books where isbn = $1 %s order by level DESC limit 1"

	// first, get his book
	condition := fmt.Sprintf("and store_id = '%s'", book.StoreId)
	err := DB.QueryRow(fmt.Sprintf(query, condition), book.Isbn).Scan(&book.Id, &book.Title, &book.Price, &book.Author, &book.Publisher, &book.Pubdate, &book.Subtitle, &book.Image, &book.Summary, &book.AuthorIntro)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	log.Debugf("select id, title, price, author, publisher, pubdate, subtitle, image, summary, author_intro from books where isbn = '%s' %s order by level DESC limit 1", book.Isbn, condition)

	// get book info by isbn, return level is most high
	err = DB.QueryRow(fmt.Sprintf(query, ""), book.Isbn).Scan(&book.Id, &book.Title, &book.Price, &book.Author, &book.Publisher, &book.Pubdate, &book.Subtitle, &book.Image, &book.Summary, &book.AuthorIntro)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func SaveBook(book *pb.Book) error {
	//save book info, which level plus one and return new ID
	query := "insert into books(store_id, title, isbn, price, author, publisher, pubdate, subtitle, image, summary, author_intro, level) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) returning id"

	log.Debugf("insert into books(store_id, title, isbn, price, author, publisher, pubdate, subtitle, image, summary, author_intro, level) values('%s', '%s', '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d) returning id", book.StoreId, book.Title, book.Isbn, book.Price, book.Author, book.Publisher, book.Pubdate, book.Subtitle, book.Image, book.AuthorIntro, book.Level+1)

	err := DB.QueryRow(query, book.StoreId, book.Title, book.Isbn, book.Price, book.Author, book.Publisher, book.Pubdate, book.Subtitle, book.Image, book.Summary, book.AuthorIntro, book.Level+1).Scan(&book.Id)

	if err != nil {
		log.Error(err)
		return nil
	}

	return nil
}
