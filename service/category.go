package service

import (
	"context"
	"github.com/devfullcycle/go-grpc/internal/database"
	"github.com/devfullcycle/go-grpc/internal/pb"
	"io"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDb database.CategoryDb
}

func NewCategoryService(categoryDb *database.CategoryDb) *CategoryService {
	return &CategoryService{
		CategoryDb: *categoryDb,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := cs.CategoryDb.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{Id: category.ID,
			Name: category.Name,
			Description: category.Description,
	}

	return categoryResponse, nil
}

func (cs *CategoryService) ListCategories(_ context.Context, _ *pb.Blank) (*pb.CategoryList, error) {
	categoriesDb, err := cs.CategoryDb.FindAll()
	if err != nil {
		return nil, err
	}

	categories := []*pb.Category{}
	for _, category := range categoriesDb {
		categories = append(categories, &pb.Category{
			Name: category.Name,
			Id: category.ID,
			Description: category.Description,
		})
	}

	return &pb.CategoryList{
		Categories: categories,
	}, nil
}

func (cs *CategoryService) GetCategory(_ context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := cs.CategoryDb.FindByID(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (cs *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		categoryResult, err := cs.CategoryDb.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (cs *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		categoryResult, err := cs.CategoryDb.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		err = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
		if err != nil {
			return err
		}
	}
}
