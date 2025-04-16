package internal

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	db "github.com/zero-shubham/authsvc/db/orm"
	authrpc "github.com/zero-shubham/authsvc/transport/grpc"
)

type Server struct {
	authrpc.AuthServiceServer
	conn db.DBTX
}

func NewServer(conn db.DBTX) *Server {
	return &Server{
		conn: conn,
	}
}

func (s *Server) UserToken(ctx context.Context, in *authrpc.UserTokenRequest) (*authrpc.TokenResponse, error) {
	return &authrpc.TokenResponse{
		AccessToken: "test-token",
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, in *authrpc.UserRequest) (*authrpc.UserResponse, error) {
	orm := db.New(s.conn)

	agId, err := uuid.Parse(in.AppGroupId)
	if err != nil {
		log.Err(err).Msg("failed to parse AppGroupId")
		return &authrpc.UserResponse{}, err
	}

	params := db.CreateUserParams{
		Email:      in.Email,
		AppGroupID: agId,
	}

	if in.OrgId != "" {
		oid, err := uuid.Parse(in.OrgId)
		if err != nil {
			log.Err(err).Msg("failed to parse OrdId")
			return &authrpc.UserResponse{}, err
		}
		params.OrgID = uuid.NullUUID{
			UUID:  oid,
			Valid: true,
		}
	}

	hpass, err := HashPassword(in.Password)
	if err != nil {
		log.Err(err).Msg("failed to parse password")
		return &authrpc.UserResponse{}, err
	}
	params.Password = hpass

	user, err := orm.CreateUser(ctx, params)
	if err != nil {
		log.Err(err).Msg("failed to create db record")
		return &authrpc.UserResponse{}, err
	}

	resp := authrpc.UserResponse{
		Id:         user.ID.String(),
		Email:      user.Email,
		AppGroupId: user.AppGroupID.String(),
	}
	if user.OrgID.Valid {
		resp.OrgId = user.OrgID.UUID.String()
	}

	return &resp, nil
}

func (s *Server) CreateApp(ctx context.Context, in *authrpc.AppRequest) (*authrpc.AppResponse, error) {
	orm := db.New(s.conn)

	agId, err := uuid.Parse(in.AppGroupId)
	if err != nil {
		log.Err(err).Msg("failed to parse AppGroupId")
		return &authrpc.AppResponse{}, err
	}

	params := db.CreateAppParams{
		AppGrpID: agId,
	}

	if in.OrgId != "" {
		oid, err := uuid.Parse(in.OrgId)
		if err != nil {
			log.Err(err).Msg("failed to parse OrgId")
			return &authrpc.AppResponse{}, err
		}
		params.OrgID = uuid.NullUUID{
			UUID:  oid,
			Valid: true,
		}
	}

	app, err := orm.CreateApp(ctx, params)
	if err != nil {
		log.Err(err).Msg("failed to create app record")
		return &authrpc.AppResponse{}, err
	}

	resp := authrpc.AppResponse{
		Id:         app.ID.String(),
		AppGroupId: app.AppGrpID.String(),
	}
	if app.OrgID.Valid {
		resp.OrgId = app.OrgID.UUID.String()
	}

	return &resp, nil
}

func (s *Server) CreateAppGroup(ctx context.Context, in *authrpc.AppGroupRequest) (*authrpc.AppGroupResponse, error) {
	orm := db.New(s.conn)

	oid, err := uuid.Parse(in.OrgId)
	if err != nil {
		log.Err(err).Msg("failed to parse OrgId")
		return &authrpc.AppGroupResponse{}, err
	}

	appGroup, err := orm.CreateAppGroup(ctx, db.CreateAppGroupParams{
		OrgID:  oid,
		Name:   in.Name,
		Scopes: in.Scopes,
	})
	if err != nil {
		log.Err(err).Msg("failed to create app group record")
		return &authrpc.AppGroupResponse{}, err
	}

	return &authrpc.AppGroupResponse{
		Id:     appGroup.ID.String(),
		OrgId:  appGroup.OrgID.String(),
		Name:   appGroup.Name,
		Scopes: appGroup.Scopes,
	}, nil
}
