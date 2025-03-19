package middleware

import (
	"context"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"vm/internal/interfaces/token"
)

// JWTInterceptor 是 gRPC Unary 拦截器，用于验证 JWT Token
func JWTInterceptor(rdb *redis.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 从上下文获取 Metadata
		// metadata.FromIncomingContext() 是 gRPC 提供的一个函数，用于从传入的 context.Context 对象 ctx 中提取元数据。元数据是一个键值对的集合，通常用于存储请求的一些额外信息，如请求头。
		// 该函数返回两个值：第一个值 md 是一个 metadata.MD 类型的对象，它实际上是一个 map[string][]string，用于存储提取到的元数据；第二个值 ok 是一个布尔类型，表示是否成功从上下文中提取到元数据。
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		
		// 获取 authorization 头
		authHeaders, ok := md["authorization"]
		if !ok || len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}
		
		// 提取 Token，格式为 "Bearer <token>"
		authHeader := authHeaders[0]
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
		}
		
		// 检查 Token 是否在黑名单中
		isBlacklisted, err := token.IsTokenBlacklisted(rdb, tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to check blacklist: %v", err)
		}
		if isBlacklisted {
			return nil, status.Errorf(codes.Unauthenticated, "token is blacklisted")
		}
		
		// 验证 Token
		claims, err := token.ValidateToken(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		
		// 将邮箱放入上下文
		newCtx := context.WithValue(ctx, "email", claims.Email)
		// newCtx = context.WithValue(newCtx, "userID", claims.ID)
		
		// 调用实际的处理函数
		return handler(newCtx, req)
	}
}
