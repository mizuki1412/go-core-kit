package class

import "mime/multipart"

type File struct {
	File   multipart.File
	Header *multipart.FileHeader
}
