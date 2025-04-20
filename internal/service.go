package internal

import (
	"context"
	"errors"

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

func (s *Server) UserToken(ctx context.Context, in *authrpc.UserTokenRequest) (*authrpc.TokenResponse, error) {

	orm := db.New(s.conn)

	user, err := orm.GetUserAndAppGrpByEmail(ctx, in.Email)
	if err != nil {
		log.Err(err).Msg("failed to retrieve user record by email")
		return &authrpc.TokenResponse{}, err
	}

	if !CheckPasswordHash(in.Password, user.Password) {
		log.Err(err).Msg("failed to retrieve user record by email")
		return &authrpc.TokenResponse{}, err
	}

	token, err := CreateToken(user.ID.String(), user.Scopes, user.ID.String())
	if err != nil {
		log.Err(err).Msg("internal server error")
		return &authrpc.TokenResponse{}, err
	}

	resp := authrpc.TokenResponse{
		AccessToken: token,
	}

	return &resp, err
}

func (s *Server) UpdateAppGroup(ctx context.Context, in *authrpc.AppGroupRequest) (*authrpc.AppGroupResponse, error) {
	orm := db.New(s.conn)

	if in.Id == "" {
		log.Info().Msg("Id not passed for UpdateAppGroup")
		return &authrpc.AppGroupResponse{}, errors.New("required Id missing")
	}

	id, err := uuid.Parse(in.Id)
	if err != nil {
		log.Err(err).Msg("failed to parse Id")
		return &authrpc.AppGroupResponse{}, err
	}

	// First get the existing app group to ensure it exists
	appGroup, err := orm.GetAppGroupByID(ctx, id)
	if err != nil {
		log.Err(err).Msg("failed to retrieve app group")
		return &authrpc.AppGroupResponse{}, err
	}

	params := db.UpdateAppGroupParams{
		ID:     appGroup.ID,
		Name:   in.Name,
		Scopes: in.Scopes,
	}

	if in.OrgId != "" {
		oid, err := uuid.Parse(in.OrgId)
		if err != nil {
			log.Err(err).Msg("failed to parse OrgIdId")
			return &authrpc.AppGroupResponse{}, err
		}
		params.OrgID = oid
	}

	// Update the app group with new values
	updatedAppGroup, err := orm.UpdateAppGroup(ctx, params)
	if err != nil {
		log.Err(err).Msg("failed to update app group")
		return &authrpc.AppGroupResponse{}, err
	}

	return &authrpc.AppGroupResponse{
		Id:     updatedAppGroup.ID.String(),
		OrgId:  updatedAppGroup.OrgID.String(),
		Name:   updatedAppGroup.Name,
		Scopes: updatedAppGroup.Scopes,
	}, nil
}

func (s *Server) ExchangeToken(ctx context.Context, in *authrpc.ExchangeTokenRequest) (*authrpc.TokenResponse, error) {
	if in.AppId == "" {
		return nil, errors.New("missing required AppId")
	}

	if in.AccessToken == "" {
		return nil, errors.New("missing required AccessToken")
	}

	if in.AppId == "" {
		return nil, errors.New("missing required AppId")
	}

	appID, err := uuid.Parse(in.AppId)
	if err != nil {
		return nil, errors.New("error parsing AppId")
	}

	orm := db.New(s.conn)
	appGrp, err := orm.GetAppGrpByAppID(ctx, appID)
	if err != nil {
		return nil, errors.New("error while retrieving app")
	}

	claims, err := ParseToken(in.AccessToken)
	if err != nil {
		log.Err(err).Msg("error while parsing token")
		return nil, errors.New("error while parsing token")
	}

	token, err := CreateToken(in.AppId, appGrp.Scopes, claims.Subject)
	if err != nil {
		return nil, errors.New("error while creating token")
	}

	return &authrpc.TokenResponse{
		AccessToken: token,
	}, nil
}
