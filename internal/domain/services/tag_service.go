package services

import (
	"product-service/internal/domain/models"
	"product-service/internal/domain/repositories"
	"product-service/pkg/logger"

	// protobuf
	pbTag "product-service/proto/tag"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TagService struct {
	tagRepo *repositories.TagRepository
}

func NewTagService() *TagService {
	return &TagService{
		tagRepo: &repositories.TagRepository{},
	}
}

func (c *TagService) GetTags(name string) (*pbTag.GetTagsResponse, error) {
	var pbTags []*pbTag.Tag
	tags, err := c.tagRepo.GetTags(name)
	if err != nil {
		logger.Errorf("Failed to get tags: %s", err)
		return nil, status.Errorf(codes.Internal, "Failed to get tags")
	}

	err = copier.Copy(&pbTags, &tags)
	if err != nil {
		logger.Errorf("Failed to copy tags: %s", err)
		return nil, status.Errorf(codes.Internal, "Failed to get tags")
	}
	res := &pbTag.GetTagsResponse{
		Tags: pbTags,
	}

	return res, nil
}

func (c *TagService) GetTagById(query *pbTag.TagId) (*pbTag.Tag, error) {
	var tagData pbTag.Tag
	tag, err := c.tagRepo.GetTagById(query.Id)
	if err != nil {
		logger.Errorf("Failed to get a tag by id: %s", err)
		return nil, status.Errorf(codes.NotFound, "Tag not found")
	}

	err = copier.Copy(&tagData, tag)
	if err != nil {
		logger.Errorf("Failed to copy tag: %s", err)
		return nil, status.Errorf(codes.Internal, "Failed to get tag")
	}

	return &tagData, nil
}

func (c *TagService) SaveTag(tagReq *pbTag.SaveTagRequest) (*pbTag.Tag, error) {
	tag := &models.Tag{
		Name: tagReq.Name,
	}

	err := c.tagRepo.Save(tag)
	if err != nil {
		logger.Errorf("Failed to save tag: %s", err)
		return nil, status.Errorf(codes.Internal, "Failed to save tag")
	}

	// convert models.Tag to pbTag.Tag
	tagData := &pbTag.Tag{}
	if err := copier.Copy(tagData, tag); err != nil {
		logger.Errorf("Failed to copy tag: %s", err)
		return nil, status.Errorf(codes.Internal, "Failed to save tag")
	}

	return tagData, nil
}

func (c *TagService) UpdateTag(request *pbTag.UpdateTagRequest) (*pbTag.Tag, error) {
	tag, err := c.tagRepo.GetTagById(request.Id)
	if err != nil {
		logger.Errorf("Failed to get a tag by id: %s", err)
		return nil, status.Errorf(codes.NotFound, "Tag not found")
	}

	tag.Name = request.TagReq.Name

	err = c.tagRepo.Update(tag)
	if err != nil {
		logger.Errorf("Failed to update tag: %s", err)
		return nil, status.Errorf(codes.Internal, "Failed to update tag")
	}

	var tagData pbTag.Tag
	if err = copier.Copy(&tagData, tag); err != nil {
		logger.Errorf("Failed to copy tag: %s", err)
		return nil, status.Errorf(codes.Internal, "Failed to update tag")
	}

	return &tagData, nil
}

func (c *TagService) DeleteTag(request *pbTag.TagId) error {
	tag, err := c.tagRepo.GetTagById(request.Id)
	if err != nil {
		logger.Errorf("Failed to get a tag by id: %s", err)
		return status.Errorf(codes.NotFound, "Tag not found")
	}

	err = c.tagRepo.Delete(tag, false)
	if err != nil {
		logger.Errorf("Failed to delete tag: %s", err)
		return status.Errorf(codes.Internal, "Failed to delete tag")
	}
	return nil
}
