package dto

type UploadQueryDTO struct {
	BucketName string `query:"bucketName" validate:"required,min=3,max=63,hostname_rfc1123"`
	FolderPath string `query:"folderPath" validate:"required,max=255"`
	IsPublic   bool   `query:"isPublic"`
}
