package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/DarkReduX/productservice/internal/service"
	"github.com/DarkReduX/productservice/model"
	"github.com/DarkReduX/productservice/protobuf"
)

type Product struct {
	svc *service.Product
	protobuf.UnimplementedProductServiceServer
}

func NewProduct(svc *service.Product) *Product {
	return &Product{svc: svc}
}

func (p *Product) Create(ctx context.Context, pbReq *protobuf.CreateRequest) (*protobuf.Product, error) {
	userID, err := uuid.Parse(pbReq.GetUserId())
	if err != nil {
		slog.Error("Failed to parse user id: ", slog.String("error", err.Error()), slog.Any("id", pbReq.GetUserId()))
		return nil, status.Error(codes.InvalidArgument, ErrMsgInvalidID)
	}

	req := &model.Product{
		UserID:      userID,
		Name:        pbReq.GetName(),
		Description: pbReq.GetDescription(),
		Price:       pbReq.GetPrice(),
	}

	product, err := p.svc.Create(ctx, req)
	if err != nil {
		slog.Error("Failed to create product: ", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, "failed to create product")
	}

	resp := &protobuf.Product{
		Id:          product.ID.String(),
		UserId:      product.UserID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return resp, nil
}

func (p *Product) Get(ctx context.Context, req *protobuf.GetRequest) (*protobuf.Product, error) {
	pID, err := uuid.Parse(req.GetId())
	if err != nil {
		slog.Error("Failed to parse product id: ", slog.String("error", err.Error()), slog.String("id", req.GetId()))
		return nil, status.Error(codes.InvalidArgument, ErrMsgInvalidID)
	}

	product, err := p.svc.Get(ctx, pID)
	if err != nil {
		slog.Error("Failed to get product: ", slog.String("error", err.Error()))
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "not found product")
		}
		return nil, status.Error(codes.Internal, "not found product")
	}

	resp := &protobuf.Product{
		Id:          product.ID.String(),
		UserId:      product.UserID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return resp, nil
}

func (p *Product) List(ctx context.Context, _ *emptypb.Empty) (*protobuf.ListResponse, error) {
	products, err := p.svc.List(ctx)
	if err != nil {
		slog.Error("Failed to list products: ", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "")
	}

	pbResp := &protobuf.ListResponse{
		Products: make([]*protobuf.Product, 0, len(products)),
	}

	for _, v := range products {
		pbResp.Products = append(pbResp.Products, &protobuf.Product{
			Id:          v.ID.String(),
			UserId:      v.UserID.String(),
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
		})
	}

	return pbResp, nil
}

func (p *Product) ListByUser(ctx context.Context, request *protobuf.ListByUserRequest) (*protobuf.ListResponse, error) {
	userID, err := uuid.Parse(request.GetUserId())
	if err != nil {
		slog.Error("failed to parse user id", slog.String("error", err.Error()), slog.String("id", request.GetUserId()))
		return nil, status.Error(codes.InvalidArgument, ErrMsgInvalidID)
	}

	products, err := p.svc.ListByUser(ctx, userID)
	if err != nil {
		slog.Error("failed to list products by user", slog.String("error", err.Error()))
		return nil, status.Error(codes.NotFound, "not found product")
	}

	pbResp := &protobuf.ListResponse{Products: make([]*protobuf.Product, 0, len(products))}

	for _, v := range products {
		pbResp.Products = append(pbResp.Products, &protobuf.Product{
			Id:          v.ID.String(),
			UserId:      v.UserID.String(),
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
		})
	}

	return pbResp, nil
}

func (p *Product) Update(ctx context.Context, pbProduct *protobuf.Product) (*protobuf.Product, error) {
	productID, err := uuid.Parse(pbProduct.GetId())
	if err != nil {
		slog.Error("failed to parse product id", slog.String("error", err.Error()), slog.String("id", pbProduct.GetId()))
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("product %v", ErrMsgInvalidID))
	}

	userID, err := uuid.Parse(pbProduct.GetUserId())
	if err != nil {
		slog.Error("failed to parse user id", slog.String("error", err.Error()), slog.String("id", pbProduct.GetUserId()))
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("user  %v", ErrMsgInvalidID))
	}

	product := &model.Product{
		ID:          productID,
		UserID:      userID,
		Name:        pbProduct.GetName(),
		Description: pbProduct.GetDescription(),
		Price:       pbProduct.GetPrice(),
	}

	if err = p.svc.Update(ctx, product); err != nil {
		slog.Error("failed to update product", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to update product")
	}

	return pbProduct, nil
}

func (p *Product) Delete(ctx context.Context, request *protobuf.DeleteRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		slog.Error("failed to parse product id", slog.String("error", err.Error()), slog.String("id", request.GetId()))
		return nil, status.Error(codes.InvalidArgument, ErrMsgInvalidID)
	}

	if err = p.svc.Delete(ctx, id); err != nil {
		slog.Error("failed to delete product", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to delete product")
	}

	return new(emptypb.Empty), nil
}
