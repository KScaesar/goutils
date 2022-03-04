package goutils

func NewPageOption(pageNumber int64, pageSize int64) PageOption {
	option := PageOption{PageNumber: pageNumber, PageSize: pageSize}
	return option.Init()
}

type PageOption struct {
	PageNumber int64 `json:"pNumber" form:"pNumber"`
	PageSize   int64 `json:"pSize" form:"pSize"`
}

func (p PageOption) Init() PageOption {
	const (
		defaultPageNumber = 1
		defaultPageSize   = 10

		maxPageSize = 2e3
	)

	pNumber := p.PageNumber
	pSize := p.PageSize

	if p.PageNumber <= 0 {
		pNumber = defaultPageNumber
	}

	switch {
	case p.PageSize <= 0:
		pSize = defaultPageSize

	case p.PageSize > maxPageSize:
		pSize = maxPageSize
	}

	return PageOption{
		PageNumber: pNumber,
		PageSize:   pSize,
	}
}

func (p PageOption) OffsetOrSkip() int64 {
	return (p.PageNumber - 1) * p.PageSize
}

func NewPageResponse(opt PageOption, totalCount int64) *PageResponse {
	opt = opt.Init()

	quotient := totalCount / opt.PageSize
	remainder := totalCount % opt.PageSize

	var totalPageNumber int64
	switch {
	case quotient > 0 && remainder == 0:
		totalPageNumber = quotient

	case quotient > 0 && remainder != 0:
		totalPageNumber = quotient + 1

	case quotient <= 0:
		totalPageNumber = 1
	}

	var targetPageSize int64
	switch {
	case totalPageNumber > opt.PageNumber:
		targetPageSize = opt.PageSize

	case totalPageNumber == opt.PageNumber && totalCount != 0 && remainder == 0:
		targetPageSize = opt.PageSize

	case totalPageNumber == opt.PageNumber && totalCount != 0 && remainder != 0:
		targetPageSize = remainder

	case totalPageNumber == opt.PageNumber && totalCount == 0:
		targetPageSize = 0

	case totalPageNumber < opt.PageNumber:
		targetPageSize = 0
	}

	return &PageResponse{
		TotalPageNumber: totalPageNumber,
		TotalCount:      totalCount,

		TargetPageNumber: opt.PageNumber,
		TargetPageSize:   targetPageSize,
	}
}

type PageResponse struct {
	TotalPageNumber int64 `json:"totalPage"`
	TotalCount      int64 `json:"totalCount"`

	TargetPageNumber int64 `json:"targetPageNumber"`
	TargetPageSize   int64 `json:"targetPageSize"`
}
