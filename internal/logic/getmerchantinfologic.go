package logic

import (
	"context"

	"github.com/starslipay/paycomm/xerror"
	"github.com/starslipay/user_mgr/internal/svc"
	"github.com/starslipay/user_mgr/internal/xerr"
	"github.com/starslipay/user_mgr/user_mgr_pb"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMerchantInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantInfoLogic {
	return &GetMerchantInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 模拟商户信息 mock
func (l *GetMerchantInfoLogic) GetMerchantInfo(in *user_mgr_pb.GetMerchantInfoReq) (*user_mgr_pb.GetMerchantInfoRsp, error) {
	if in.MerchantId == "" {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "merchant_id is empty")
	}

	if "2000000000" == in.MerchantId {
		return &user_mgr_pb.GetMerchantInfoRsp{
			MerchantId:   in.MerchantId,
			MerchantUid:  2000000000,
			MerchantName: "starslipay虚拟商城商户",
		}, nil
	} else if "3000000000" == in.MerchantId {
		return &user_mgr_pb.GetMerchantInfoRsp{
			MerchantId:   in.MerchantId,
			MerchantUid:  3000000000,
			MerchantName: "淘宝商户",
		}, nil
	} else if "4000000000" == in.MerchantId {
		return &user_mgr_pb.GetMerchantInfoRsp{
			MerchantId:   in.MerchantId,
			MerchantUid:  4000000000,
			MerchantName: "京东商户",
		}, nil
	} else if "5000000000" == in.MerchantId {
		return &user_mgr_pb.GetMerchantInfoRsp{
			MerchantId:   in.MerchantId,
			MerchantUid:  5000000000,
			MerchantName: "拼多多商户",
		}, nil
	} else if "6000000000" == in.MerchantId {
		return &user_mgr_pb.GetMerchantInfoRsp{
			MerchantId:   in.MerchantId,
			MerchantUid:  6000000000,
			MerchantName: "抖音商户",
		}, nil
	} else {
		return nil, xerror.NewBizError(codes.Internal, xerr.ErrCodeParam, "merchant_id is not found")
	}
}
