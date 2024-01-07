package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioRepository interface {
	UploadServiceImage(periodID int, userID uint, imageBytes []byte, contentType string) (string, error)
	RemoveServiceImage(periodID int, userID uint) error
}

func (r *Repository) UploadServiceImage(periodID int, userID uint, imageBytes []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("period/%d/image", periodID)

	reader := io.NopCloser(bytes.NewReader(imageBytes))

	_, err := r.mc.PutObject(context.TODO(), "images-bucket", objectName, reader, int64(len(imageBytes)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", errors.New("ошибка при добавлении изображения в минио бакет")
	}

	// Формирование URL изображения
	imageURL := fmt.Sprintf("http://localhost:9000/images-bucket/%s", objectName)

	return imageURL, nil
}

func (r *Repository) RemoveServiceImage(periodID int, userID uint) error {
	objectName := fmt.Sprintf("period/%d/image", periodID)
	err := r.mc.RemoveObject(context.TODO(), "images-bucket", objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.New("не удалось удалить изображение из бакет")
	}

	if err := r.db.Table("baggages").
		Where("id_period = ?", periodID).
		Update("photo", nil).Error; err != nil {
		return errors.New("ошибка при обновлении URL изображения в базе данных")
	}

	return nil
}
