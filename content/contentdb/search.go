package contentdb

import (
	"time"
)

type SearchParameter struct {
	Clause   string
	Argument interface{}
}

type SearchOption func() SearchParameter

func SearchVideoByName(name string) SearchOption {
	return func() SearchParameter {
		return SearchParameter{
			Clause:   "name LIKE ?",
			Argument: name,
		}
	}
}

func SearchVideoByLicense(lic string) SearchOption {
	return func() SearchParameter {
		return SearchParameter{
			Clause:   "license LIKE ?",
			Argument: lic,
		}
	}
}

func SearchVideoByEncoding(enc string) SearchOption {
	return func() SearchParameter {
		return SearchParameter{
			Clause:   "encoding = ?",
			Argument: enc,
		}
	}
}

func SearchVideoByDuration(dur time.Duration, shorter bool) SearchOption {
	return func() SearchParameter {
		if shorter {
			return SearchParameter{
				Clause:   "duration_seconds <= ?",
				Argument: dur.Milliseconds() / 1000,
			}
		}
		return SearchParameter{
			Clause:   "duration_seconds >= ?",
			Argument: dur.Milliseconds() / 1000,
		}
	}
}

func SearchVideoByUploadedBeforeOrAfterDate(t time.Time, before bool) SearchOption {
	return func() SearchParameter {
		if before {
			return SearchParameter{
				Clause:   "uploaded <= ?",
				Argument: t,
			}
		}
		return SearchParameter{
			Clause:   "uploaded >= ?",
			Argument: t,
		}
	}
}

func SearchVideoByResolution(res string) SearchOption {
	return func() SearchParameter {
		return SearchParameter{
			Clause:   "resolution = ?",
			Argument: res,
		}
	}
}

func SearchVideoByCategory(cat string) SearchOption {
	return func() SearchParameter {
		return SearchParameter{
			Clause:   "category = ?",
			Argument: cat,
		}
	}
}

func SearchVideoByAttribution(attr string) SearchOption {
	return func() SearchParameter {
		return SearchParameter{
			Clause:   "attribution LIKE ?",
			Argument: attr,
		}
	}
}

func SearchVideoByFileSize(limit uint64, lower bool) SearchOption {
	return func() SearchParameter {
		if lower {
			return SearchParameter{
				Clause:   "file_size_bytes <= ?",
				Argument: limit,
			}
		}
		return SearchParameter{
			Clause:   "file_size_bytes >= ?",
			Argument: limit,
		}
	}
}
