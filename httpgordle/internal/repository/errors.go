package repository

type apiError string

func (e apiError) Error() string {
	return string(e)
}
