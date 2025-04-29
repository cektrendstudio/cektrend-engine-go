package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/logger"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"
	"github.com/cektrendstudio/cektrend-engine-go/service"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type googleDriveRepository struct {
	jwt *jwt.Config
}

func NewGoogleDriveRepository() (repo service.GoogleDriveRepository) {
	clientSecret := utstring.Env("CLIENT_SECRET_GOOGLE_API", "")
	if clientSecret == "" {
		logger.Err("ENV CLIENT_SECRET_GOOGLE_API not found")
		return
	}
	emailDelegate := utstring.Env("EMAIL_DELEGATION_GOOGLE_API", "")
	if emailDelegate == "" {
		logger.Err("ENV EMAIL_DELEGATION_GOOGLE_API not found")
		return
	}
	jsonData := []byte(clientSecret)

	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}

	if err := json.Unmarshal(jsonData, &s); err != nil {
		logger.Err("Error while unmarshal client secret")
		logger.Err(err.Error())
		return
	}

	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Subject:    emailDelegate,
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}

	return &googleDriveRepository{
		jwt: config,
	}
}

func (g *googleDriveRepository) GetFile(ctx context.Context, fileID string) (resp *http.Response, serr serror.SError) {
	srv, err := drive.NewService(ctx, option.WithHTTPClient(g.jwt.Client(ctx)))
	if err != nil {
		serr = serror.NewFromErrorc(err, "Failed to init service drive")
		return
	}

	meta, err := srv.Files.Get(fileID).Fields("size").Do()
	if err != nil {
		serr = serror.NewFromErrorc(err, "Failed to get metadata from drive")
		return
	}

	// Check content size
	if meta.Size > 10*1024*1024 {
		serr = serror.New("Google Drive file size exceeds 10MB limit")
		return
	}

	resp, err = srv.Files.Get(fileID).Download()
	if err != nil {
		serr = serror.NewFromErrorc(err, "Failed to get data from drive")
		return
	}

	return
}

func (g *googleDriveRepository) GetDriveFileExtension(ctx context.Context, fileURL string) (string, serror.SError) {
	srv, err := drive.NewService(ctx, option.WithHTTPClient(g.jwt.Client(ctx)))
	if err != nil {
		serr := serror.NewFromErrorc(err, "Failed to init service drive")
		return "", serr
	}

	re := regexp.MustCompile(`(?:drive\.google\.com\/.*\/d\/|drive\.google\.com\/file\/d\/|id=|open\?id=)([a-zA-Z0-9_-]{33,})`)
	matches := re.FindStringSubmatch(fileURL)
	if len(matches) < 2 {
		serr := serror.New("Failed to get google drive file id")
		return "", serr
	}

	fileId := matches[1]
	meta, err := srv.Files.Get(fileId).Fields("mimeType").Do()
	if err != nil {
		serr := serror.NewFromErrorc(err, "Failed to get metadata mimeType from drive")
		return "", serr
	}

	mimeType := meta.MimeType
	extensionMap := map[string]string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   ".docx",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
		"application/pdf": ".pdf",
		"image/png":       ".png",
		"image/jpeg":      ".jpg",
	}

	extension, ok := extensionMap[mimeType]
	if !ok {
		return "", nil
	}

	return extension, nil
}
