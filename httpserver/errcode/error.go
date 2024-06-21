package errcode

import (
	"strings"
)

const (
	SuccessCodeBusiness = 0
	SuccessMsgBusiness  = "success"

	// ErrorCodeInternal HTTP CODE
	ErrorCodeInternal = 500
	// ErrorMsgInternal 默认错误信息
	ErrorMsgInternal          = "系统错误"
	ErrorCodeInvalidParameter = 400
	ErrorCodeUnauthorized     = 401
	ErrorMsgUnauthorized      = "认证或授权失败"
	ErrorCodeNotfound         = 404
	ErrorMsgNotfound          = "无资源错误"

	// ErrorMsgPrefixInvalidParameter 错误信息前缀
	ErrorMsgPrefixInvalidParameter = "参数验证错误%v"

	ErrorMsgSuffixParamExists    = "%v已经存在"
	ErrorMsgSuffixParamNotExists = "%v不存在"

	// ErrorCodeBusiness Business Code
	ErrorCodeBusiness = 9999

	ErrorMsgMethodNotFound = "请求方法不允许"
	ErrorMsgHandleNotFound = "请求URL不存在"

	ErrorMsgLoadParam     = "读取请求参数失败"
	ErrorMsgValidateParam = "参数验证错误"
	ErrorMsgAtoiParam     = "参数转换失败"

	ErrorMsgUserNameOrPassword = "用户名尚未注册"
	ErrorMsgUserLogin          = "用户登录失败"

	ErrorMsgGetUserInfo  = "用户信息获取失败"
	ErrorMsgUserLoginOut = "退出登录失败"

	ErrorMsgCreateData     = "创建数据失败"
	ErrorMsgListData       = "获取数据失败"
	ErrorMsgUpdateData     = "修改数据失败"
	ErrorMsgDeleteData     = "删除数据失败"
	ErrorMsgCreateOrUpdate = "修改/创建数据失败"
	ErrorMsgBatchCreate    = "批量创建数据失败"
	ErrorMsgCancel         = "取消失败"

	ErrorMsgTrainTypeNameExists     = "车型名称重复"
	ErrorMsgRegionCarriageName      = "车厢名称超出车型车厢数"
	ErrorMsgRegionNameExists        = "区域名称重复"
	ErrorMsgRegionIDExists          = "区域ID重复"
	ErrorMsgTrainRegionNotMatch     = "绑定区域和车型绑定区域不匹配"
	ErrorMsgTrainFrameTypeNotMatch  = "绑定车型和配准帧绑定车型不匹配"
	ErrorMsgRegionName              = "非法区域名称"
	ErrorMsgFrameMerged             = "帧信息已完成合并，不允许编辑"
	ErrorMsgMergeFrame              = "合并车头帧失败"
	ErrorMsgFileRead                = "读取file失败"
	ErrorMsgFileEmpty               = "文件为空"
	ErrorFileSave                   = "文件存储失败"
	ErrorMsgOSMkdir                 = "创建目录失败"
	ErrorMsgRegionBindPointCLoud    = "该区域已经绑定其他点云信息"
	ErrorMsgRegionMatchTrainType    = "不允许不同车型下的数据进行配准"
	ErrorMsgRegionMatchRegion       = "配准帧和点云数据不属于同一区域"
	ErrorMsgTrainFrameMatch         = "车头帧未进行合并"
	ErrorMsgStopComputing           = "停车点计算失败"
	ErrorMsgStopComputeCancel       = "停车点停止计算失败"
	ErrorMsgStopComputeProgress     = "获取停车点计算进度失败"
	ErrorMsgDataUpload              = "数据上传失败,请重新上传"
	ErrorMsgGetSocketAddress        = "获取配准服务地址失败"
	ErrorMsgSocketConn              = "配准工具打开失败"
	ErrorMsgUploadToken             = "文件上传Token失效"
	ErrorMsgUploadSeek              = "偏移量处理失败"
	ErrorMsgUploadWrite             = "写入文件失败"
	ErrorMsgUploadStat              = "获得文件信息失败"
	ErrorMsgOccupiedModify          = "模型占用状态禁止修改"
	ErrorMsgRedisSetExpire          = "上传token生成失败"
	ErrorMsgNerfOccupiedDelete      = "nerf模型状态不支持删除"
	ErrorMsgNerfModelExists         = "nerf模型已经存在"
	ErrorMsgNerfDataGroupProcessing = "数据正在处理，请先停止处理"
	ErrorMsgNerfDataOccupied        = "数据组被占用，请先删除相关模型"

	ErrorMsgCameraNotModify            = "相机配置占用，仅焦距和光心支持修改"
	ErrorMsgCameraOccupied             = "相机配置占用"
	ErrorMsgPullProgress               = "进度同步失败"
	ErrorMsgGetProgress                = "进度获取失败"
	ErrorMsgNerfModelNotUploading      = "nerf模型未处于上传状态"
	ErrorMsgNerfModelCreate            = "nerf模型导入失败"
	ErrorMsgNerfModelCancel            = "nerf模型导入取消"
	ErrorMsgRegionNoRelatedFrame       = "区域和示教停车点不具有关联关系"
	ErrorMsgPhotoIsShield              = "重复屏蔽拍照点"
	ErrorMsgPhotoIsNoShield            = "拍照点未屏蔽"
	ErrorMsgTrainTypeNotMatch          = "绑定车型不匹配"
	ErrorMsgDefaultCorrectionNotModify = "默认纠偏轨迹不允许修改或删除"
	ErrorMsgNerfModelOccupied          = "nerf模型被占用"
	ErrorMsgPlanOccupied               = "示教方案占用"
	ErrorMsgNerfPhotoPointDelete       = "已存档实拍图不可删除"
	ErrorMsgPlanNerfViewerStart        = "虚拟示教启动失败"
	ErrorMsgPlanNerfViewerNoStart      = "虚拟示教未启动"

	ErrorMsgPlanNerfViewerStop = "虚拟示教关闭失败"

	ErrorMsgScaleNoRegionFrame = "暂无区域配准帧数据，请先上传配准帧数据"
	ErrorMsgGenScalePointCloud = "点云调整失败"
	ErrorMsgNerfModelStatus    = "模型状态不支持当前操作"

	ErrorMsgNerfDataStatus     = "数据组状态不支持当前操作"
	ErrorMsgNerfDataVideoImage = "视频抽帧图不允许删除"
	ErrorMsgNerfDataImage      = "nerf源数据不存在图片数据"
	ErrorMsgCarriageGroupName  = "车厢编组信息缺失"
	ErrorMsgCarriageInfo       = "车厢信息缺失"

	ErrorMsgBogieTypeRef          = "待删除转向架类型被占用，请先结束占用"
	ErrorMsgRegionVerifyState     = "验证区域前置信息不完整"
	ErrorMsgIssuedPlan            = "下发方案失败"
	ErrorMsgVerifyTask            = "存在机器人验证任务，请先结束任务"
	ErrorMsgApplyPlan             = "同步虚拟方案失败"
	ErrorMsgVerifyTaskStateRepOpt = "任务不允许重复操作"
	ErrorMsgVerifyTaskStateOpt    = "不支持的任务状态操作"

	ErrorMsgOccupyRobot = "机器人占用失败"

	ErrorMsgNerfStudioOccupied = "获取NerfStudio全局锁失败"
	ErrorMsgGenerateToken      = "生成access_token失败"
	ErrorMsgCreateToken        = "创建access_token失败"
	ErrorMsgTokenNotExists     = "access_token为空"
	ErrorMsgCheckToken         = "access_token失效"
	ErrorMsgNoPermission       = "用户无访问权限"
	ErrorMsgUpgradeWebSocket   = "websocket升级失败"

	ErrorMsgDataExists          = "记录已经存在"
	ErrorMsgDataNotExists       = "记录不存在"
	ErrorMsgTransactionOpen     = "事务开启失败"
	ErrorMsgTransactionCommit   = "事务提交失败"
	ErrorMsgTransactionRollback = "事务回滚失败"
	ErrorMsgHttpClientError     = "第三方服务异常"

	ErrorMsgNodeOvertopArea = "节点超出区域范围"
)

var (
	ErrCode = map[string]int{
		ErrorMsgTokenNotExists: 4000,
		ErrorMsgCheckToken:     4001,
		ErrorMsgNoPermission:   4003,

		ErrorMsgMethodNotFound: 5000,
		ErrorMsgHandleNotFound: 5001,

		ErrorMsgLoadParam:     5002,
		ErrorMsgValidateParam: 5003,
		ErrorMsgAtoiParam:     5004,

		ErrorMsgCreateData:     5005,
		ErrorMsgListData:       5006,
		ErrorMsgUpdateData:     5007,
		ErrorMsgDeleteData:     5008,
		ErrorMsgCreateOrUpdate: 5009,

		ErrorMsgUserNameOrPassword: 5010,
		ErrorMsgUserLogin:          5011,
		ErrorMsgGetUserInfo:        5012,
		ErrorMsgUserLoginOut:       5013,

		ErrorMsgGenerateToken:              5014,
		ErrorMsgCreateToken:                5015,
		ErrorMsgTrainTypeNameExists:        5016,
		ErrorMsgRegionCarriageName:         5017,
		ErrorMsgRegionNameExists:           5018,
		ErrorMsgRegionIDExists:             5019,
		ErrorMsgBatchCreate:                5020,
		ErrorMsgTrainRegionNotMatch:        5021,
		ErrorMsgRegionName:                 5022,
		ErrorMsgTrainFrameTypeNotMatch:     5023,
		ErrorMsgFrameMerged:                5024,
		ErrorMsgMergeFrame:                 5025,
		ErrorMsgFileRead:                   5026,
		ErrorMsgFileEmpty:                  5027,
		ErrorFileSave:                      5028,
		ErrorMsgOSMkdir:                    5029,
		ErrorMsgRegionBindPointCLoud:       5030,
		ErrorMsgRegionMatchTrainType:       5031,
		ErrorMsgRegionMatchRegion:          5032,
		ErrorMsgTrainFrameMatch:            5033,
		ErrorMsgStopComputing:              5034,
		ErrorMsgDataUpload:                 5035,
		ErrorMsgSocketConn:                 5036,
		ErrorMsgUploadToken:                5037,
		ErrorMsgUploadSeek:                 5038,
		ErrorMsgUploadWrite:                5039,
		ErrorMsgUploadStat:                 5040,
		ErrorMsgOccupiedModify:             5041,
		ErrorMsgRedisSetExpire:             5042,
		ErrorMsgNerfOccupiedDelete:         5043,
		ErrorMsgNerfModelExists:            5044,
		ErrorMsgNerfDataGroupProcessing:    5046,
		ErrorMsgNerfStudioOccupied:         5047,
		ErrorMsgCameraNotModify:            5048,
		ErrorMsgCameraOccupied:             5049,
		ErrorMsgPullProgress:               5050,
		ErrorMsgGetProgress:                5051,
		ErrorMsgCancel:                     5052,
		ErrorMsgNerfModelNotUploading:      5053,
		ErrorMsgNerfModelCreate:            5054,
		ErrorMsgNerfModelCancel:            5055,
		ErrorMsgStopComputeCancel:          5056,
		ErrorMsgStopComputeProgress:        5057,
		ErrorMsgRegionNoRelatedFrame:       5058,
		ErrorMsgPhotoIsShield:              5059,
		ErrorMsgPhotoIsNoShield:            5060,
		ErrorMsgTrainTypeNotMatch:          5061,
		ErrorMsgDefaultCorrectionNotModify: 5062,
		ErrorMsgNerfModelOccupied:          5063,
		ErrorMsgPlanOccupied:               5064,
		ErrorMsgNerfPhotoPointDelete:       5065,
		ErrorMsgPlanNerfViewerStart:        5066,
		ErrorMsgPlanNerfViewerStop:         5067,
		ErrorMsgPlanNerfViewerNoStart:      5068,
		ErrorMsgGetSocketAddress:           5069,
		ErrorMsgScaleNoRegionFrame:         5070,
		ErrorMsgGenScalePointCloud:         5071,
		ErrorMsgNerfModelStatus:            5072,
		ErrorMsgNerfDataStatus:             5073,
		ErrorMsgCarriageGroupName:          5074,
		ErrorMsgCarriageInfo:               5075,
		ErrorMsgBogieTypeRef:               5076,
		ErrorMsgRegionVerifyState:          5077,
		ErrorMsgIssuedPlan:                 5078,
		ErrorMsgApplyPlan:                  5079,
		ErrorMsgVerifyTask:                 5080,
		ErrorMsgVerifyTaskStateRepOpt:      5081,
		ErrorMsgVerifyTaskStateOpt:         5082,
		ErrorMsgNerfDataOccupied:           5083,
		ErrorMsgNerfDataImage:              5084,
		ErrorMsgOccupyRobot:                5085,
		ErrorMsgNodeOvertopArea:            5086,

		ErrorMsgDataExists:          6000,
		ErrorMsgDataNotExists:       6001,
		ErrorMsgTransactionOpen:     6002,
		ErrorMsgTransactionCommit:   6003,
		ErrorMsgTransactionRollback: 6004,
		ErrorMsgHttpClientError:     6005,
		ErrorMsgUpgradeWebSocket:    6006,
	}

	// CommonErrorMsg 通用错误信息
	CommonErrorMsg = []string{
		ErrorMsgSuffixParamExists,
		ErrorMsgSuffixParamNotExists,
		ErrorMsgPrefixInvalidParameter,
	}

	// PostProcessingMsg 通用的错误处理信息后置处理
	PostProcessingMsg = map[string]string{
		ErrorMsgSuffixParamExists:      ErrorMsgDataExists,
		ErrorMsgSuffixParamNotExists:   ErrorMsgDataNotExists,
		ErrorMsgPrefixInvalidParameter: ErrorMsgValidateParam,
	}
)

// GetErrorCode 从已经记录在案的code中查询，如果有则返回，没有返回默认 ErrorCodeBusiness
func GetErrorCode(msg string) int {
	code, ok := ErrCode[msg]
	if !ok {
		for _, item := range CommonErrorMsg {
			t := strings.TrimPrefix(item, "%v")
			t = strings.TrimSuffix(t, "%v")
			if strings.Contains(msg, t) {
				return ErrCode[PostProcessingMsg[item]]
			}
		}
		return ErrorCodeBusiness
	}
	return code
}
