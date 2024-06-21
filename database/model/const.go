package model

import (
	"demo-gogo/config"
	"fmt"
)

const (
	//FolderModel             = "model"
	FolderImg               = "img"
	FolderNerf              = "nerf"
	LocationNginxPointCloud = "/nginx_point_cloud"
	LocationMinioPointCloud = "/minio_point_cloud"

	// FileUploadToken 文件上传的token
	FileUploadToken = "upload_token"
	// FileName 上传的文件名
	FileName = "file_name"
	// FileSize 上传文件的总大小
	FileSize = "file_size"
	// FileUrl 文件url
	FileUrl = "file_url"
	// FilePath 文件存储的路径
	FilePath = "file_path"
	// FileUploadSize 文件已经上传的大小
	FileUploadSize = "upload_size"
)

const (
	FileDateFormatLayout = "2006-01-02"

	redisKeyStopLocation    = "%s:stop_compute:%s:%d"
	redisKeyBigFileUpload   = "%s:upload:big_file:%s"
	redisKeyNerfProcessData = "%s:nerf:process_data:%d"
	redisKeyTrainNerfModel  = "%s:nerf:train:%d"
	redisKeyNerfModelViewer = "%s:nerf:model_viewer:%d"
	redisKeyNerfViewer      = "%s:nerf:plan_viewer:%d"
	redisKeyVerifyTask      = "%s:nerf:verify_task"
)

// GetStopProgressContentKey 生成停车点计算进度Redis Key
func GetStopProgressContentKey(computeType string, regionID int) string {
	return fmt.Sprintf(redisKeyStopLocation, config.Conf.APP.Name, computeType, regionID)
}

// GetFileUploadKey 生成大文件上传Redis Key
func GetFileUploadKey(token string) string {
	return fmt.Sprintf(redisKeyBigFileUpload, config.Conf.APP.Name, token)
}

func GetNerfVerifyTaskKey() string {
	return fmt.Sprintf(redisKeyVerifyTask, config.Conf.APP.Name)
}

// GetNerfDataProcessKey 生成Nerf数据处理进度Redis Key
func GetNerfDataProcessKey(nerfDataID int) string {
	return fmt.Sprintf(redisKeyNerfProcessData, config.Conf.APP.Name, nerfDataID)
}

// GetNerfTrainKey 生成Nerf模型训练Redis Key
func GetNerfTrainKey(nerfModelID int) string {
	return fmt.Sprintf(redisKeyTrainNerfModel, config.Conf.APP.Name, nerfModelID)
}

// GetNerfModelViewerKey 生成Nerf模型预览Redis Key
func GetNerfModelViewerKey(nerfModelID int) string {
	return fmt.Sprintf(redisKeyNerfModelViewer, config.Conf.APP.Name, nerfModelID)
}

// GetNerfPlanViewerKey 生成示教方案视图预览Redis Key
func GetNerfPlanViewerKey(planID int) string {
	return fmt.Sprintf(redisKeyNerfViewer, config.Conf.APP.Name, planID)
}

const (
	TableNameTrainType = "train_type"

	TableNameCarriage  = "carriage"
	TableNamePosition  = "carriage_position"
	TableNameBogieType = "bogie_type"
	TableNameBogie     = "bogie"
	TableNameAxle      = "bogie_axle"
	TableNameWheel     = "bogie_wheel"

	TableNameOriginRegion    = "origin_region"
	TableNameRegion          = "region"
	TableNameTrainFrame      = "train_frame"
	TableNameRegionFrame     = "region_frame"
	TableNameStopFrame       = "stop_frame"
	TableNamePointCloud      = "point_cloud"
	TableNamePhotoPoint      = "photo_point"
	TableNameMatchResult     = "match_result"
	TableNameStopLocation    = "stop_location"
	TableNameStopPhoto       = "stop_photo"
	TableNameStopNerfPhoto   = "stop_nerf_photo"
	TableNameRegionStopFrame = "region_stop_frame"
	TableNameCorrection      = "correction"

	TableNameNerfModel            = "nerf_model"
	TableNameNerfModelImage       = "nerf_model_images"
	TableNameCamera               = "camera"
	TableNamePlan                 = "plan"
	TableNameVerifyPlan           = "verify_plan"
	TableNameStopNerfCamera       = "stop_nerf_camera"
	TableNameNerfCamera           = "nerf_camera"
	TableNameVerifyPhotoPoint     = "verify_photo_point"
	TableNameNerfPhotoPoint       = "nerf_photo_point"
	TableNameNerfData             = "nerf_data"
	TableNameNerfDataGroup        = "nerf_data_group"
	TableNameNerfDataImage        = "nerf_images"
	TableNameVerifyTask           = "verify_task"
	TableNameRegionNerfPhotoPoint = "region_nerf_photo_point"
	TableNameMap                  = "map"
	TableNameMapRoutes            = "map_routes"
	TableNameMapRouteNodes        = "map_route_nodes"

	FieldID    = "id"
	FieldName  = "name"
	FieldMapId = "map_id"

	FieldCarriageNumber = "carriage_number"
	FieldTrainTypeID    = "train_type_id"
	FieldCarriageName   = "carriage_name"

	FieldCarriageID            = "carriage_id"
	FieldEndPositionID         = "end_position_id"
	FieldSidePositionID        = "side_position_id"
	FieldControlPositionID     = "control_position_id"
	FieldReuseType             = "reuse_type"
	FieldCopyNumber            = "copy_number"
	FieldSortNum               = "sort_number"
	FieldSide                  = "side"
	FieldCalculateState        = "calculate_state"
	FieldCalculateRes          = "calculate_res"
	FieldMergeTrainFrameID     = "merge_train_frame_id"
	FieldMergeCoordinateSystem = "merge_coordinate_system"
	FieldMergeCoordinate       = "merge_coordinate"
	FieldTrainFrameID          = "train_frame_id"
	FieldTrainFrameName        = "train_frame_name"
	FieldRegionID              = "region_id"
	FieldStopUploadState       = "upload_state"
	FieldOriginRegionCode      = "region_code"
	FieldCoordinateSystem      = "coordinate_system"
	FieldCoordinate            = "coordinate"

	FieldPositionKind = "kind"

	FieldBogieTypeID            = "bogie_type_id"
	FieldBogieID                = "bogie_id"
	FieldAxleID                 = "axle_id"
	FieldWheelID                = "wheel_id"
	FieldGroupName              = "group_name"
	FieldSensorType             = "sensor_type"
	FieldSensorName             = "sensor_name"
	FieldSensorData             = "sensor_data"
	FieldFrameUrl               = "frame_url"
	FieldBindOriginRegion       = "bind_origin_region"
	FieldRegionPointCloudID     = "point_cloud_id"
	FieldUrl                    = "url"
	FieldStopFrameID            = "stop_frame_id"
	FieldFrameID                = "frame_id"
	FieldStopLocationID         = "stop_location_id"
	FieldPose                   = "pose"
	FieldRoboticArm             = "robotic_arm"
	FieldDirectionNumber        = "direction_number"
	FieldTextureImageUrl        = "texture_image_url"
	FieldRGBImageUrl            = "rgb_image_url"
	FieldDepthImageUrl          = "depth_image_url"
	FieldPointCloudImageUrl     = "point_cloud_url"
	FieldPhotoPointID           = "photo_point_id"
	FieldShieldPhotoIDs         = "shield_photo_ids"
	FieldOffset                 = "offset"
	FieldPointType              = "point_type"
	FieldCorrectionTypes        = "correction_types"
	FieldImageTypes             = "image_types"
	FieldJointTrajectory        = "joint_trajectory"
	FieldCustomMessage          = "custom_message"
	FieldNerfPhotoPointArchived = "archived"

	FieldMatchResultSourceID = "source_id"
	FieldMatchResultTargetID = "target_id"
	FieldMatchResultType     = "match_type"
	FieldOriginRegionID      = "origin_region_id"
	FieldOriginRegionName    = "origin_region_name"
	FieldCorrectionID        = "correction_id"
	FieldPlanID              = "plan_id"
	FieldVerifyPlanID        = "verify_plan_id"
	FieldPlanName            = "plan_name"

	FieldOccupied     = "occupied"
	FieldNerfModelUrl = "nerf_model_url"

	FieldCameraModelUrl   = "camera_model_url"
	FieldCameraFov        = "fov"
	FieldLimit            = "limit"
	FieldResolution       = "resolution"
	FieldFocalLength      = "focal_length"
	FieldLightPose        = "light_pose"
	FieldCameraModelName  = "camera_model_name"
	FieldNerfPhotoUrl     = "nerf_photo_url"
	FieldNerfPhotoPointID = "nerf_photo_point_id"
	FieldNerfModelStatus  = "nerf_model_status"

	FieldNerfModelMatchResult = "match_result"
	FieldNerfCameraIssued     = "issued"
	FieldNerfCameraID         = "nerf_camera_id"

	FieldCameraID   = "camera_id"
	FieldNerModelID = "nerf_model_id"

	FieldRegionName = "region_name"

	FieldNerfDataStatus         = "status"
	FieldNerfDataID             = "nerf_data_id"
	FieldNerfDataGroupID        = "nerf_data_group_id"
	FieldNerfDataType           = "data_type"
	FieldNerfDataGroupTool      = "tool"
	FieldNerfDataGroupProcessed = "processed"
	FieldNerfDataImageSize      = "size"

	FieldNerfModelID = "nerf_model_id"
	FieldPhotoPath   = "photo_path"
	FieldState       = "state"

	FieldCreatedTime = "created_at"
	FieldUpdatedTime = "updated_at"
	FieldDeletedTime = "deleted_at"
)
