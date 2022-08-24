package log

type Config struct {
	Level      string
	Format     string
	Prefix     string
	Filename   string
	Maxsize    int
	Maxbackups int
	Maxage     int
	Compress   bool
}
