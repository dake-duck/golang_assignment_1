package mysql

import (
	"DAKExDUCK/assignment_1/pkg/models"
	"database/sql"
	"errors"
	"reflect"
)

type NewsModel struct {
	DB *sql.DB
}

func PointersOf(v interface{}) interface{} {
	in := reflect.ValueOf(v)
	out := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(in.Type().Elem())), in.Len(), in.Len())
	for i := 0; i < in.Len(); i++ {
		out.Index(i).Set(in.Index(i).Addr())
	}
	return out.Interface()
}

func (model *NewsModel) get_tags(newsID int) ([]models.Tag, error) {
	query := `
		SELECT t.tag_id, t.tag_name_en, t.tag_name_ru
		FROM tags t
		INNER JOIN news_tags_pairs ntp ON t.tag_id = ntp.pair_tag_id
		WHERE ntp.pair_news_id = $1;
	`

	rows, err := model.DB.Query(query, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		tag := models.Tag{}
		err := rows.Scan(&tag.ID, &tag.NameEN, &tag.NameRU)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (model *NewsModel) Insert(moodle_id int, title string, content string, attachments []map[string]interface{}) (int, error) {
	query := `
	INSERT INTO
		news (news_moodle_id, news_title_html, news_body_html, news_attachments, news_created)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP())
	`

	result, err := model.DB.Exec(query, moodle_id, title, content, attachments)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (model *NewsModel) Get(id int) ([]models.News, error) {
	query := `
	SELECT
		news_id, news_moodle_id, news_title_html, news_body_html, news_attachments, news_created
	FROM
		public.news
	WHERE
		news_id = $1;`
	row := model.DB.QueryRow(query, id)
	news := []models.News{}
	n := models.News{}

	err := row.Scan(&n.ID, &n.MoodleID, &n.Title, &n.Body, &n.Attachments, &n.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	tags, err := model.get_tags(n.ID)
	if err != nil {
		return nil, err
	}
	n.Tags = tags

	news = append(news, n)
	return news, nil
}

func (model *NewsModel) GetByCategory(category int) ([]models.News, error) {
	query := `
        SELECT
            n.news_id, n.news_moodle_id, n.news_title_html, n.news_body_html,
            n.news_attachments, n.news_created
        FROM
            public.news n
        INNER JOIN
            public.news_tags_pairs ntp ON n.news_id = ntp.pair_news_id
        WHERE
            ntp.pair_tag_id = $1
		ORDER BY
            n.news_created DESC
    `
	row := model.DB.QueryRow(query, category)
	news := []models.News{}
	n := models.News{}

	err := row.Scan(&n.ID, &n.MoodleID, &n.Title, &n.Body, &n.Attachments, &n.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	tags, err := model.get_tags(n.ID)
	if err != nil {
		return nil, err
	}
	n.Tags = tags

	news = append(news, n)
	return news, nil
}

func (model *NewsModel) Latest() ([]models.News, error) {
	query := `
	SELECT
		news_id, news_moodle_id, news_title_html, news_body_html, news_attachments, news_created
	FROM
		public.news
	ORDER BY
		news_created DESC
	LIMIT 10;`
	rows, err := model.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	news := []models.News{}
	for rows.Next() {
		n := models.News{}
		err = rows.Scan(&n.ID, &n.MoodleID, &n.Title, &n.Body, &n.Attachments, &n.Created)
		if err != nil {
			return nil, err
		}

		tags, err := model.get_tags(n.ID)
		if err != nil {
			return nil, err
		}
		n.Tags = tags

		news = append(news, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}
