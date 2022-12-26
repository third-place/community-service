package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
	"github.com/third-place/community-service/internal/util"
	"sort"
)

type ShareService struct {
	shareRepository *repository.ShareRepository
	postRepository  *repository.PostRepository
	userRepository  *repository.UserRepository
	likeRepository  *repository.LikeRepository
}

func CreateShareService() *ShareService {
	conn := db.CreateDefaultConnection()
	return &ShareService{
		repository.CreateShareRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateUserRepository(conn),
		repository.CreateLikeRepository(conn),
	}
}

func CreateTestShareService() *ShareService {
	conn := util.SetupTestDatabase()
	return &ShareService{
		repository.CreateShareRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateUserRepository(conn),
		repository.CreateLikeRepository(conn),
	}
}

func (s *ShareService) CreateShare(share *model.NewShare) (*model.Share, error) {
	user, err := s.userRepository.FindOneInGoodStandingByUuid(uuid.MustParse(share.User.Uuid))
	if err != nil {
		return nil, err
	}
	post, _ := s.postRepository.FindOneByUuid(uuid.MustParse(share.Post.Uuid))
	shareEntity := entity.CreateShare(user, post, share)
	s.shareRepository.Save(shareEntity)
	return mapper.GetShareModelFromEntity(shareEntity), nil
}

func (s *ShareService) GetShare(shareUuid uuid.UUID) (*model.Share, error) {
	share, err := s.shareRepository.FindOneByUuid(shareUuid)
	if err != nil {
		return nil, err
	}
	return mapper.GetShareModelFromEntity(share), nil
}

func (s *ShareService) GetShares(username *string, limit int) ([]*model.Share, error) {
	var selfShares []*entity.Post
	var followingShares []*entity.Post
	var publicShares []*entity.Post
	remaining := limit
	var user *entity.User
	if username != nil {
		user, _ = s.userRepository.FindOneByUsername(*username)
		selfShares = s.shareRepository.FindByUser(user, limit)
		remaining = remaining - len(selfShares)
	}
	if remaining > 0 && username != nil {
		followingShares = s.shareRepository.FindByUserFollows(*username, remaining)
		remaining -= len(followingShares)
	}
	if remaining > 0 {
		publicShares = s.shareRepository.FindAll(remaining)
	}
	allShares := append(selfShares, followingShares...)
	allShares = append(allShares, publicShares...)
	sort.SliceStable(allShares, func(i, j int) bool {
		return allShares[i].CreatedAt.After(allShares[j].CreatedAt)
	})
	fullList := removeDuplicatePosts(allShares)
	if user != nil {
		return s.populateModelsWithLikes(fullList, user), nil
	}
	return mapper.GetShareModelsFromEntities(allShares), nil
}

func (s *ShareService) populateModelsWithLikes(posts []*entity.Post, viewer *entity.User) []*model.Share {
	postIds := s.getPostIDs(posts)
	postLikes := s.likeRepository.FindLikesForPosts(postIds, viewer)
	likedPosts := make(map[uint]bool)
	for _, postLike := range postLikes {
		likedPosts[postLike.PostID] = true
	}
	fullListModels := mapper.GetShareModelsFromEntities(posts)
	for i, item := range posts {
		if likedPosts[item.ID] {
			fullListModels[i].SelfLiked = true
		}
	}
	return fullListModels
}

func (s *ShareService) getPostIDs(posts []*entity.Post) []uint {
	postIDs := make([]uint, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}
	return postIDs
}
