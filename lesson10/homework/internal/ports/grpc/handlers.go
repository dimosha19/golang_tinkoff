package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	myerrors "homework10/internal/errors"
	"strconv"
	"time"
)

func (s Service) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := s.a.CreateAd(request.Title, request.Text, request.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, myerrors.ErrBadRequest.Error())
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

func (s Service) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := s.a.UpdateAdStatus(request.AdId, request.UserId, request.Published)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, myerrors.ErrBadRequest.Error())
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

func (s Service) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	ad, err := s.a.UpdateAd(request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, myerrors.ErrBadRequest.Error())
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

func (s Service) ListAds(ctx context.Context, ads *GetAds) (*ListAdResponse, error) {
	pub := ads.Pub
	if pub != "true" && pub != "false" && pub != "all" {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprint("invalid pub filter"))
	}
	author, err := strconv.Atoi(ads.Author)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprint("invalid author filter"))
	}

	_, err = time.Parse("02-01-06", ads.Date)
	if ads.Date != "all" && err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprint("invalid date filter"))
	}

	ad, err := s.a.GetAds(ads.Pub, int64(author), ads.Date, ads.Title)
	if err != nil {
		if errors.Is(myerrors.ErrBadRequest, err) {
			return nil, status.Error(codes.InvalidArgument, myerrors.ErrBadRequest.Error())
		}
		return nil, status.Error(codes.Unknown, fmt.Sprint("unknown error"))
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

func (s Service) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	user, err := s.a.CreateUser(request.Nickname, request.Email)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, myerrors.ErrBadRequest.Error())

	}
	return &UserResponse{
		Id:    user.ID,
		Name:  user.Nickname,
		Email: user.Email,
	}, nil
}

func (s Service) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	user, err := s.a.GetUser(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, myerrors.ErrBadRequest.Error())
	}
	return &UserResponse{
		Id:    user.ID,
		Name:  user.Nickname,
		Email: user.Email,
	}, nil
}

func (s Service) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*Success, error) {
	err := s.a.DeleteUser(request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprint("internal error"))
	}
	return &Success{
		Success: "user was successfully deleted",
	}, nil
}

func (s Service) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*Success, error) {
	err := s.a.DeleteAd(request.AdId, request.AuthorId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprint("internal error"))
	}
	return &Success{
		Success: "ad was successfully deleted",
	}, nil
}
