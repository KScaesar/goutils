package goutils

func NewPageOption(pageNumber int, pageSize int) PageOption {
	option := PageOption{PageNumber: pageNumber, PageSize: pageSize}
	return option.Init()
}

type PageOption struct {
	PageNumber int `json:"page_number" form:"pNumber"`
	PageSize   int `json:"page_size" form:"pSize"`
}

func (p PageOption) Init() PageOption {
	const (
		defaultPageNumber = 1
		defaultPageSize   = 10
	)

	switch {
	case p.PageNumber > 0 && p.PageSize > 0:
		return PageOption{
			PageNumber: p.PageNumber,
			PageSize:   p.PageSize,
		}

	case p.PageNumber <= 0 && p.PageSize > 0:
		return PageOption{
			PageNumber: defaultPageNumber,
			PageSize:   p.PageSize,
		}

	case p.PageNumber > 0 && p.PageSize <= 0:
		return PageOption{
			PageNumber: p.PageNumber,
			PageSize:   defaultPageSize,
		}

	default:
		return PageOption{
			PageNumber: defaultPageNumber,
			PageSize:   defaultPageSize,
		}
	}
}

func (p PageOption) OffsetOrSkip() int64 {
	return int64(p.PageNumber-1) * int64(p.PageSize)
}

func NewPageResponse(opt PageOption, totalCount int) PageResponse {
	opt = opt.Init()

	quotient := totalCount / opt.PageSize
	remainder := totalCount % opt.PageSize

	var totalPageNumber int
	switch {
	case quotient > 0 && remainder == 0:
		totalPageNumber = quotient

	case quotient > 0 && remainder != 0:
		totalPageNumber = quotient + 1

	case quotient <= 0:
		totalPageNumber = 1
	}

	var targetPageSize int
	switch {
	case totalPageNumber > opt.PageNumber:
		targetPageSize = opt.PageSize

	case totalPageNumber == opt.PageNumber && remainder == 0:
		targetPageSize = opt.PageSize

	case totalPageNumber == opt.PageNumber && remainder != 0:
		targetPageSize = remainder

	case totalPageNumber < opt.PageNumber:
		targetPageSize = 0
	}

	return PageResponse{
		TotalPageNumber: totalPageNumber,
		TotalCount:      totalCount,

		TargetPageNumber: opt.PageNumber,
		TargetPageSize:   targetPageSize,
	}
}

type PageResponse struct {
	TotalPageNumber int `json:"totalPage"`
	TotalCount      int `json:"totalCount"`

	TargetPageNumber int `json:"targetPageNumber"`
	TargetPageSize   int `json:"targetPageSize"`
}
