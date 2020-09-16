package page

// Page ...
// pageSize表示每页要显示的条数，pageNumber表示页码
// offset = (pageNumber-1) * PageSize
// select * from table limit (pageNumber-1)*pageSize, pageSize
func Page(pageSize, pageNumber int64) (limit, offset int64) {
	if pageSize > 0 {
		limit = pageSize
	} else {
		limit = 10
	}

	if pageNumber > 0 {
		offset = (pageNumber - 1) * pageSize
	} else {
		offset = -1
	}
	return limit, offset
}
