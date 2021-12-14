package goutils

// FilterOption 簡易查詢用途
//
// example:
//
// filter := FilterOption{
//     {"user_id = ?", 123},
//     {"datetime in (?)", []string{"2020-10-17", "2021-09-16"}},
//     {"is_admin = ?", true},
//     {"blocked_at = ?", "null"},
// }
type FilterOption []struct {
	Key   string
	Value interface{}
}

func (f FilterOption) FilterOption() {}
