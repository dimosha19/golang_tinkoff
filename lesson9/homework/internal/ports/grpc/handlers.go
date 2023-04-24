package grpc

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

func (G Service) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := G.a.CreateAd(request.Title, request.Text, request.UserId)
	if err != nil {
		return nil, err
	}
	return &AdResponse{
		Id:            ad.ID,
		Title:         ad.Title,
		Text:          ad.Text,
		AuthorId:      ad.AuthorID,
		Published:     ad.Published,
		PublishedTime: timestamppb.New(ad.PublishedTime),
		UpdatedTimem:  timestamppb.New(ad.UpdateTime),
	}, nil
}

func (G Service) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := G.a.UpdateAdStatus(request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, err
	}
	return &AdResponse{
		Id:            ad.ID,
		Title:         ad.Title,
		Text:          ad.Text,
		AuthorId:      ad.AuthorID,
		Published:     ad.Published,
		PublishedTime: timestamppb.New(ad.PublishedTime),
		UpdatedTimem:  timestamppb.New(ad.UpdateTime),
	}, nil
}

func (G Service) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	ad, err := G.a.UpdateAd(request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, err
	}
	return &AdResponse{
		Id:            ad.ID,
		Title:         ad.Title,
		Text:          ad.Text,
		AuthorId:      ad.AuthorID,
		Published:     ad.Published,
		PublishedTime: timestamppb.New(ad.PublishedTime),
		UpdatedTimem:  timestamppb.New(ad.UpdateTime),
	}, nil
}

func (G Service) ListAds(ctx context.Context, ads *GetAds) (*ListAdResponse, error) {
	pub := ads.Pub
	if pub != "true" && pub != "false" && pub != "all" {
		return nil, errors.New("invalid pub filter")
	}
	author, err := strconv.Atoi(ads.Author)
	if err != nil {
		return nil, errors.New("invalid author filter")
	}

	_, err = time.Parse("02-01-06", ads.Date)
	if ads.Date != "all" && err != nil {
		return nil, errors.New("invalid date filter")
	}

	ad, err := G.a.GetAds(ads.Pub, int64(author), ads.Date, ads.Title)
	if err != nil {
		return nil, err
	}
	var response ListAdResponse
	for i := range *ad {
		response.List = append(response.List, &AdResponse{
			Id:            (*ad)[i].ID,
			Title:         (*ad)[i].Title,
			Text:          (*ad)[i].Text,
			AuthorId:      (*ad)[i].AuthorID,
			Published:     (*ad)[i].Published,
			PublishedTime: timestamppb.New((*ad)[i].PublishedTime),
			UpdatedTimem:  timestamppb.New((*ad)[i].UpdateTime),
		})
	}
	return &response, nil

}

func (G Service) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	user, err := G.a.CreateUser(request.Nickname, request.Email)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		Id:    user.ID,
		Name:  user.Nickname,
		Email: user.Email,
	}, nil
}

func (G Service) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	user, err := G.a.GetUser(request.Id)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		Id:    user.ID,
		Name:  user.Nickname,
		Email: user.Email,
	}, nil
}

func (G Service) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*Success, error) {
	err := G.a.DeleteUser(request.Id)
	if err != nil {
		return nil, err
	}
	return &Success{
		Success: "user was successfully deleted",
	}, nil
}

func (G Service) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*Success, error) {
	err := G.a.DeleteAd(request.AdId, request.AuthorId)
	if err != nil {
		return nil, err
	}
	return &Success{
		Success: "ad was successfully deleted",
	}, nil
}
