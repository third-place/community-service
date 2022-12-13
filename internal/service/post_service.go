package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/kafka"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
	"github.com/third-place/community-service/internal/util"
	"sort"
)

type PostService struct {
	userRepository   *repository.UserRepository
	postRepository   *repository.PostRepository
	followRepository *repository.FollowRepository
	imageRepository  *repository.ImageRepository
	likeRepository   *repository.LikeRepository
	kafkaWriter      kafka.Producer
	securityService  *SecurityService
}

func CreatePostService() *PostService {
	conn := db.CreateDefaultConnection()
	producer := kafka.CreateProducer()
	return &PostService{
		postRepository:   repository.CreatePostRepository(conn),
		userRepository:   repository.CreateUserRepository(conn),
		followRepository: repository.CreateFollowRepository(conn),
		imageRepository:  repository.CreateImageRepository(conn),
		likeRepository:   repository.CreateLikeRepository(conn),
		kafkaWriter:      producer,
		securityService:  CreateSecurityService(),
	}
}

func CreateTestPostService() *PostService {
	conn := util.SetupTestDatabase()
	producer := util.CreateTestProducer()
	return &PostService{
		postRepository:   repository.CreatePostRepository(conn),
		userRepository:   repository.CreateUserRepository(conn),
		followRepository: repository.CreateFollowRepository(conn),
		imageRepository:  repository.CreateImageRepository(conn),
		likeRepository:   repository.CreateLikeRepository(conn),
		kafkaWriter:      producer,
		securityService:  CreateSecurityService(),
	}
}

func (p *PostService) GetPost(session *model.Session, postUuid uuid.UUID) (*model.Post, error) {
	post, err := p.postRepository.FindOneByUuid(postUuid)
	if err != nil {
		return nil, err
	}
	if post.User == nil {
		return nil, errors.New(constants.ErrorMessageUserNotFound)
	}
	if !p.canSee(session, post) {
		return nil, errors.New("not accessible")
	}
	posts := make([]*entity.Post, 1)
	posts[0] = post
	postsWithShare := p.populateSharePosts(posts)
	return mapper.GetPostModelFromEntity(postsWithShare[0]), nil
}

func (p *PostService) CreatePost(session *model.Session, newPost *model.NewPost) (*model.Post, error) {
	user, err := p.userRepository.FindOneByUuid(uuid.MustParse(session.User.Uuid))
	if err != nil {
		return nil, err
	}
	post := entity.CreatePost(user, newPost)
	p.postRepository.Create(post)
	var imageEntities []*entity.Image
	for _, newImage := range newPost.Images {
		imageEntity := entity.CreateImage(user, post, &newImage)
		p.imageRepository.Create(imageEntity)
		imageEntities = append(imageEntities, imageEntity)
	}
	search, _ := p.postRepository.FindOneByUuid(*post.Uuid)
	postsWithShare := p.populateSharePosts([]*entity.Post{search})
	postModel := mapper.GetPostModelFromEntity(postsWithShare[0])
	err = p.publishPostToKafka(postModel)
	return postModel, err
}

func (p *PostService) UpdatePost(session *model.Session, postModel *model.Post) error {
	postEntity, err := p.postRepository.FindOneByUuid(uuid.MustParse(postModel.Uuid))
	if err != nil || !p.securityService.Owns(session, postEntity) {
		return errors.New("user cannot update post")
	}
	postEntity.Text = postModel.Text
	postEntity.Draft = postModel.Draft
	p.postRepository.Save(postEntity)
	err = p.publishPostToKafka(postModel)
	return nil
}

func (p *PostService) DeletePost(session *model.Session, postUuid uuid.UUID) error {
	post, err := p.postRepository.FindOneByUuid(postUuid)
	if err != nil {
		return err
	}
	if !p.securityService.Owns(session, post) {
		return errors.New("cannot delete post")
	}
	p.postRepository.Delete(post)
	postModel := mapper.GetPostModelFromEntity(post)
	err = p.publishPostToKafka(postModel)
	return nil
}

func (p *PostService) GetPostsForUser(username string, viewerUuid *uuid.UUID, limit int) ([]*model.Post, error) {
	user, err := p.userRepository.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}
	postEntities := p.populateSharePosts(p.postRepository.FindPublishedByUser(user, limit))
	var fullListModels []*model.Post
	var viewer *entity.User
	if viewerUuid != nil {
		viewer, _ = p.userRepository.FindOneByUuid(*viewerUuid)
	}
	if viewer != nil {
		fullListModels = p.populateModelsWithLikes(postEntities, viewer)
	} else {
		fullListModels = mapper.GetPostModelsFromEntities(postEntities)
	}
	return fullListModels, nil
}

func (p *PostService) GetPostsForUserFollows(usernameFollowing string, viewerUserUuid uuid.UUID, limit int) ([]*model.Post, error) {
	_, err := p.userRepository.FindOneByUsername(usernameFollowing)
	if err != nil {
		return nil, err
	}
	viewer, err := p.userRepository.FindOneByUuid(viewerUserUuid)
	if err != nil {
		return nil, err
	}
	posts := p.postRepository.FindByUserFollows(usernameFollowing, limit)
	postsWithShares := p.populateSharePosts(posts)
	postModels := p.populateModelsWithLikes(postsWithShares, viewer)
	return postModels, nil
}

func (p *PostService) GetAllPosts(limit int) []*model.Post {
	posts := p.postRepository.FindAll(limit)
	return mapper.GetPostModelsFromEntities(posts)
}

func (p *PostService) GetDraftPosts(sessionUuid uuid.UUID, limit int) []*model.Post {
	user, _ := p.userRepository.FindOneByUuid(sessionUuid)
	posts := p.postRepository.FindDraftsByUser(user, limit)
	return mapper.GetPostModelsFromEntities(posts)
}

func (p *PostService) GetPosts(username string, limit int) []*model.Post {
	user, _ := p.userRepository.FindOneByUsername(username)
	posts := p.postRepository.FindPublishedByUser(user, limit)
	return mapper.GetPostModelsFromEntities(posts)
}

func (p *PostService) GetPostsFirehose(username *string, limit int) ([]*model.Post, error) {
	var selfPosts []*entity.Post
	var followingPosts []*entity.Post
	var publicPosts []*entity.Post
	var user *entity.User
	if username != nil {
		followingPosts = p.postRepository.FindByUserFollows(*username, limit)
		limit -= len(followingPosts)
	}
	if limit > 0 && username != nil {
		user, _ = p.userRepository.FindOneByUsername(*username)
		selfPosts = p.postRepository.FindPublishedByUser(user, limit)
		limit -= len(selfPosts)
	}
	if limit > 0 {
		publicPosts = p.postRepository.FindAll(limit)
	}
	allPosts := append(selfPosts, followingPosts...)
	allPosts = append(allPosts, publicPosts...)
	sort.SliceStable(allPosts, func(i, j int) bool {
		return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)
	})
	fullList := p.populateSharePosts(removeDuplicatePosts(allPosts))
	if user != nil {
		return p.populateModelsWithLikes(fullList, user), nil
	}
	return mapper.GetPostModelsFromEntities(fullList), nil
}

func (p *PostService) GetLikedPosts(username string, limit int) ([]*model.Post, error) {
	user, err := p.userRepository.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}
	posts := p.postRepository.FindByLikes(user, limit)
	models := mapper.GetPostModelsFromEntities(p.populateSharePosts(posts))
	for _, m := range models {
		m.SelfLiked = true
	}
	return models, nil
}

func (p *PostService) populateSharePosts(posts []*entity.Post) []*entity.Post {
	postIDs := make([]uint, len(posts))
	postMap := make(map[uint]int)
	j := 0
	for i, post := range posts {
		if post.SharePostID != 0 {
			postIDs[j] = post.SharePostID
			postMap[post.SharePostID] = i
			j += 1
		}
	}
	shares := p.postRepository.FindByIDs(postIDs)
	for _, share := range shares {
		posts[postMap[share.ID]].SharePost = share
	}
	return posts
}

func (p *PostService) populateModelsWithLikes(posts []*entity.Post, viewer *entity.User) []*model.Post {
	postIds := p.getPostIDs(posts)
	postLikes := p.likeRepository.FindLikesForPosts(postIds, viewer)
	likedPosts := make(map[uint]bool)
	for _, postLike := range postLikes {
		likedPosts[postLike.PostID] = true
	}
	fullListModels := mapper.GetPostModelsFromEntities(posts)
	for i, item := range posts {
		if likedPosts[item.ID] {
			fullListModels[i].SelfLiked = true
		}
	}
	return fullListModels
}

func (p *PostService) getPostIDs(posts []*entity.Post) []uint {
	postIDs := make([]uint, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}
	return postIDs
}

func (p *PostService) publishPostToKafka(post *model.Post) error {
	topic := "posts"
	data, _ := json.Marshal(post)
	return p.kafkaWriter.Produce(kafka.CreateMessage(data, topic), nil)
}

func (p *PostService) canSee(session *model.Session, post *entity.Post) bool {
	if post.Visibility == model.PUBLIC {
		return true
	}
	if session == nil {
		return false
	}
	if (post.Draft || post.Visibility == model.PRIVATE) && !p.securityService.Owns(session, post) {
		return false
	}
	sessionUuid := uuid.MustParse(session.User.Uuid)
	follow := p.followRepository.FindByUserAndFollowing(*post.User.Uuid, sessionUuid)
	return follow != nil
}

func removeDuplicatePosts(posts []*entity.Post) []*entity.Post {
	var dedup []*entity.Post
	allKeys := make(map[uint]bool)
	for _, item := range posts {
		if value := allKeys[item.ID]; !value {
			allKeys[item.ID] = true
			dedup = append(dedup, item)
		}
	}
	return dedup
}
