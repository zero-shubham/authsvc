
„
service.protogrpc"E

AppRequest
org_id (	RorgId 
app_group_id (	R
appGroupId"V
AppResponse
id (	Rid
org_id (	RorgId 
app_group_id (	R
appGroupId"d
AppGroupRequest
id (	Rid
org_id (	RorgId
name (	Rname
scopes (	Rscopes"e
AppGroupResponse
id (	Rid
org_id (	RorgId
name (	Rname
scopes (	Rscopes"x
UserRequest
email (	Remail
password (	Rpassword
org_id (	RorgId 
app_group_id (	R
appGroupId"m
UserResponse
id (	Rid
email (	Remail
org_id (	RorgId 
app_group_id (	R
appGroupId"D
UserTokenRequest
email (	Remail
password (	Rpassword"W
TokenResponse!
access_token (	RaccessToken#
refresh_token (	RrefreshToken",
ServiceTokenRequest
app_id (	RappId"u
ExchangeTokenRequest!
access_token (	RaccessToken#
refresh_token (	RrefreshToken
app_id (	RappId2™
AuthService8
	UserToken.grpc.UserTokenRequest.grpc.TokenResponse>
ServiceToken.grpc.ServiceTokenRequest.grpc.TokenResponse@
ExchangeToken.grpc.ExchangeTokenRequest.grpc.TokenResponse3

CreateUser.grpc.UserRequest.grpc.UserResponse3

UpdateUser.grpc.UserRequest.grpc.UserResponse?
CreateAppGroup.grpc.AppGroupRequest.grpc.AppGroupResponse?
UpdateAppGroup.grpc.AppGroupRequest.grpc.AppGroupResponse0
	CreateApp.grpc.AppRequest.grpc.AppResponse0
	UpdateApp.grpc.AppRequest.grpc.AppResponseB0Z.github.com/zero-shubham/authsvc/transport/grpcbproto3