package utils

import (
	"fmt"
	"github.com/lib/pq"
	"gonum.org/v1/gonum/mat"
	"math"
	"strconv"
)

// QuaternionToPosMatrix 四元数转矩阵 tx,ty,tz,x,y,z,w
func QuaternionToPosMatrix(pose []float64) ([]float64, error) {
	if len(pose) != 7 {
		return nil, fmt.Errorf("QuaternionToPosMatrix pose len Must be 7")
	}
	res := []float64{
		0, 0, 0, pose[0],
		0, 0, 0, pose[1],
		0, 0, 0, pose[2],
		0, 0, 0, 1,
	}
	x := pose[3]
	y := pose[4]
	z := pose[5]
	w := pose[6]
	res[0] = 1 - 2*math.Pow(y, 2) - 2*math.Pow(z, 2)
	res[1] = 2 * (x*y - w*z)
	res[2] = 2 * (x*z + w*y)

	res[4] = 2 * (x*y + w*z)
	res[5] = 1 - 2*math.Pow(x, 2) - 2*math.Pow(z, 2)
	res[6] = 2 * (y*z - w*x)

	res[8] = 2 * (x*z - w*y)
	res[9] = 2 * (y*z + w*x)
	res[10] = 1 - 2*math.Pow(x, 2) - 2*math.Pow(y, 2)
	return res, nil
}

// MatrixToQuaternion 矩阵转四元素（tx,ty,tz,x,y,z,w）
func MatrixToQuaternion(matrix []float64) ([]float64, error) {
	if len(matrix) != 16 {
		return nil, fmt.Errorf("MatrixToQuaternion matrix len Must be 16")
	}
	res := make([]float64, 7)
	res[0] = matrix[3]
	res[1] = matrix[7]
	res[2] = matrix[11]

	var w, x, y, z float64
	tr := matrix[0] + matrix[5] + matrix[10]
	if tr > 0 {
		s := math.Sqrt(tr+1.0) * 2 // s = 4 * qw
		w = 0.25 * s
		x = (matrix[9] - matrix[6]) / s
		y = (matrix[2] - matrix[8]) / s
		z = (matrix[4] - matrix[1]) / s
	} else if (matrix[0] > matrix[5]) && (matrix[0] > matrix[10]) {
		s := math.Sqrt(1.0+matrix[0]-matrix[5]-matrix[10]) * 2 // s = 4 * qx
		w = (matrix[9] - matrix[6]) / s
		x = 0.25 * s
		y = (matrix[1] + matrix[4]) / s
		z = (matrix[2] + matrix[8]) / s
	} else if matrix[5] > matrix[10] {
		s := math.Sqrt(1.0+matrix[5]-matrix[0]-matrix[10]) * 2 // s = 4 * qy
		w = (matrix[2] - matrix[8]) / s
		x = (matrix[1] + matrix[4]) / s
		y = 0.25 * s
		z = (matrix[6] + matrix[9]) / s
	} else {
		s := math.Sqrt(1.0+matrix[10]-matrix[0]-matrix[5]) * 2 // s = 4 * qz
		w = (matrix[4] - matrix[1]) / s
		x = (matrix[2] + matrix[8]) / s
		y = (matrix[6] + matrix[9]) / s
		z = 0.25 * s
	}
	res[3] = x
	res[4] = y
	res[5] = z
	res[6] = w
	return res, nil
}

// ComputeNewRegFrameToHead 计算合并车头帧后区域配准帧的坐标
// 输入 矩阵 四元数  输出 四元数
func ComputeNewRegFrameToHead(trainFrameMatchResult pq.Float64Array, oldRegFrameToHead pq.Float64Array) (pq.Float64Array, error) {
	var resultDense mat.Dense
	// 车头帧配准结果
	matchDense := mat.NewDense(4, 4, trainFrameMatchResult)
	// 四元数转矩阵
	oldCoordinate, err := QuaternionToPosMatrix(oldRegFrameToHead)
	if err != nil {
		return nil, err
	}
	regionDense := mat.NewDense(4, 4, oldCoordinate)
	resultDense.Mul(matchDense, regionDense)
	// 矩阵转四元数
	return MatrixToQuaternion(resultDense.RawMatrix().Data)
}

// ComputeRegionToHead 计算区域到车头帧坐标系的坐标
// 输入 矩阵 四元数  输出 四元数
func ComputeRegionToHead(reginFrameMatchResult pq.Float64Array, regFrameToHead pq.Float64Array) (pq.Float64Array, error) {
	var resultDense mat.Dense
	// 区域帧配准结果
	matchDense := mat.NewDense(4, 4, reginFrameMatchResult)
	// 四元数转矩阵
	regFrameCoordinate, err := QuaternionToPosMatrix(regFrameToHead)
	if err != nil {
		return nil, err
	}
	regFrameDense := mat.NewDense(4, 4, regFrameCoordinate)
	resultDense.Mul(regFrameDense, matchDense)
	// 矩阵转四元数
	return MatrixToQuaternion(resultDense.RawMatrix().Data)
}

// ComputePhotoToRegion 计算拍照点到区域的坐标关系
// 输入 矩阵 四元数  输出 四元数
func ComputePhotoToRegion(stopFrameMatchResult pq.Float64Array, pose pq.Float64Array) (pq.Float64Array, error) {
	var resultDense mat.Dense
	// 停车点配准结果
	matchDense := mat.NewDense(4, 4, stopFrameMatchResult)
	poseCoordinate, err := QuaternionToPosMatrix(pose)
	if err != nil {
		return nil, err
	}
	photoDense := mat.NewDense(4, 4, poseCoordinate)
	resultDense.Mul(matchDense, photoDense)
	// 矩阵转四元数
	return MatrixToQuaternion(resultDense.RawMatrix().Data)
}

// ComputeStopLocationToHead 计算停车点到车头的坐标
func ComputeStopLocationToHead(regionPose pq.Float64Array, oldCoordinate pq.Float64Array) (pq.Float64Array, error) {
	var resultDense mat.Dense
	regionCoordinate, err := QuaternionToPosMatrix(regionPose)
	if err != nil {
		return nil, err
	}
	regionDense := mat.NewDense(4, 4, regionCoordinate)
	stopCoordinate, err := QuaternionToPosMatrix(oldCoordinate)
	if err != nil {
		return nil, err
	}
	stopDense := mat.NewDense(4, 4, stopCoordinate)
	resultDense.Mul(regionDense, stopDense)
	return MatrixToQuaternion(resultDense.RawMatrix().Data)
}

// PhotoPointCoordinateMirror 拍照点区域镜像转换
// 输入 四元数  输出 四元数
func PhotoPointCoordinateMirror(oldCoordinate pq.Float64Array) (pq.Float64Array, error) {
	rotationMatrix, err := QuaternionToPosMatrix(oldCoordinate)
	if err != nil {
		return nil, err
	}
	//rotationMirror := []float64{
	//	-rotationMatrix[0], rotationMatrix[1], -rotationMatrix[2],
	//	rotationMatrix[4], -rotationMatrix[5], rotationMatrix[6],
	//	rotationMatrix[8], -rotationMatrix[9], rotationMatrix[10],
	//}
	//tx := []float64{
	//	-1, 0, 0,
	//	0, -1, 0,
	//	0, 0, 1,
	//}
	//tempDense := mat.NewDense(3, 3, rotationMirror)
	//txDense := mat.NewDense(3, 3, tx)
	//var dot mat.Dense
	//dot.Mul(tempDense, txDense)
	//
	//rotationMatrix = []float64{
	//	dot.At(0, 0), dot.At(0, 1), dot.At(0, 2), -oldCoordinate[0],
	//	dot.At(1, 0), dot.At(1, 1), dot.At(1, 2), oldCoordinate[1],
	//	dot.At(2, 0), dot.At(2, 1), dot.At(2, 2), oldCoordinate[2],
	//	0, 0, 0, 1,
	//}
	rotationMatrix = []float64{
		-rotationMatrix[0], rotationMatrix[1], -rotationMatrix[2], -oldCoordinate[0],
		rotationMatrix[4], -rotationMatrix[5], rotationMatrix[6], oldCoordinate[1],
		rotationMatrix[8], -rotationMatrix[9], rotationMatrix[10], oldCoordinate[2],
		0, 0, 0, 1,
	}
	return MatrixToQuaternion(rotationMatrix)
}

func ComputeNerfToPhotoPoint(camera pq.Float64Array, scaleFactor float64, match pq.Float64Array) (pq.Float64Array, error) {
	newCamera := MatrixScale(camera, scaleFactor)
	newDense := mat.NewDense(4, 4, newCamera)
	matchDense := mat.NewDense(4, 4, match)
	var ans mat.Dense
	// nerf -> 区域
	ans.Mul(matchDense, newDense)
	fmt.Printf("%#v", ans.RawMatrix().Data)
	return MatrixToQuaternion(ans.RawMatrix().Data)
}

// PhotoPointCoordinateRotate 拍照点区域旋转复用
// 输入 四元数  输出 四元数
func PhotoPointCoordinateRotate(oldCoordinate pq.Float64Array) (pq.Float64Array, error) {
	matrix, err := QuaternionToPosMatrix(oldCoordinate)
	if err != nil {
		return nil, err
	}
	dense := mat.NewDense(4, 4, matrix)
	tInit := []float64{
		-1, 0, 0, 0,
		0, -1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	initDense := mat.NewDense(4, 4, tInit)
	var res mat.Dense
	res.Mul(initDense, dense)
	return MatrixToQuaternion(res.RawMatrix().Data)
}

// StopLocationCoordinateRotate 停车点旋转复用 停车点旋转复用和拍照点旋转复用一样
// 输入 四元数  输出 四元数
func StopLocationCoordinateRotate(oldCoordinate pq.Float64Array) (pq.Float64Array, error) {
	matrix, err := QuaternionToPosMatrix(oldCoordinate)
	if err != nil {
		return nil, err
	}
	dense := mat.NewDense(4, 4, matrix)
	tInit := []float64{
		-1, 0, 0, 0,
		0, -1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	initDense := mat.NewDense(4, 4, tInit)
	var res mat.Dense
	res.Mul(initDense, dense)
	return MatrixToQuaternion(res.RawMatrix().Data)
}

// StopLocationCoordinateMirror 停车点镜像复用 和拍照点镜像复用一样
func StopLocationCoordinateMirror(oldCoordinate pq.Float64Array) (pq.Float64Array, error) {
	return PhotoPointCoordinateMirror(oldCoordinate)
}

// Decimal float64 保留小数点后位数
// value float64 浮点数
// prec int 需保留小数点后的位数
func Decimal(value float64, prec int) float64 {
	value, _ = strconv.ParseFloat(strconv.FormatFloat(value, 'f', prec, 64), 64)
	return value
}

// TransposeMatrix 矩阵转置
func TransposeMatrix(matrix []float64) []float64 {
	transposed := make([]float64, 16)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			index := i*4 + j
			transposed[index] = matrix[j*4+i]
		}
	}

	return transposed
}

// MatrixScale 矩阵尺度调整
func MatrixScale(matrix []float64, scaleFactor float64) []float64 {
	// 创建原始矩阵
	originalMatrix := mat.NewDense(4, 4, matrix)
	// 进行尺度矫正   尺度---区域的点云/nerf的点云
	tra := mat.NewDense(4, 4, []float64{
		originalMatrix.At(0, 0), originalMatrix.At(0, 1), originalMatrix.At(0, 2), originalMatrix.At(0, 3) * scaleFactor,
		originalMatrix.At(1, 0), originalMatrix.At(1, 1), originalMatrix.At(1, 2), originalMatrix.At(1, 3) * scaleFactor,
		originalMatrix.At(2, 0), originalMatrix.At(2, 1), originalMatrix.At(2, 2), originalMatrix.At(2, 3) * scaleFactor,
		originalMatrix.At(3, 0), originalMatrix.At(3, 1), originalMatrix.At(3, 2), originalMatrix.At(3, 3),
	})

	// 创建3x3的目标矩阵
	rot := mat.NewDense(3, 3, []float64{
		originalMatrix.At(0, 0), -originalMatrix.At(0, 1), -originalMatrix.At(0, 2),
		originalMatrix.At(1, 0), -originalMatrix.At(1, 1), -originalMatrix.At(1, 2),
		originalMatrix.At(2, 0), -originalMatrix.At(2, 1), -originalMatrix.At(2, 2),
	})

	newMatrix := []float64{
		rot.At(0, 0), rot.At(0, 1), rot.At(0, 2), tra.At(0, 3),
		rot.At(1, 0), rot.At(1, 1), rot.At(1, 2), tra.At(1, 3),
		rot.At(2, 0), rot.At(2, 1), rot.At(2, 2), tra.At(2, 3),
		0, 0, 0, 1,
	}
	return newMatrix
}

// MatrixToPose 矩阵获取pose
func MatrixToPose(matrix []float64) ([]float64, error) {
	if len(matrix) != 16 {
		return nil, fmt.Errorf("matrix 长度有误")
	}
	// 矩阵转置
	matrixTransposed := TransposeMatrix(matrix)

	rotationMatrix := []float64{
		matrixTransposed[0], matrixTransposed[1], matrixTransposed[2], 0,
		matrixTransposed[4], matrixTransposed[5], matrixTransposed[6], 0,
		matrixTransposed[8], matrixTransposed[9], matrixTransposed[10], 0,
		0, 0, 0, 1,
	}

	// 输出欧拉角
	euler := rotationMatrixToEulerAngles(rotationMatrix)

	// 输出xyz坐标
	tra := []float64{matrixTransposed[3], matrixTransposed[7], matrixTransposed[11]}
	tra = append(tra, euler...)
	return tra, nil
}

// PoseToMatrix pose获取矩阵
func PoseToMatrix(pose []float64) ([]float64, error) {
	if len(pose) != 6 {
		return nil, fmt.Errorf("pose长度有误")
	}
	euler := pose[3:]
	rotationMatrixRet := eulerAnglesToRotationMatrix(euler)
	rotationMatrixRet[3] = pose[0]
	rotationMatrixRet[7] = pose[1]
	rotationMatrixRet[11] = pose[2]
	return rotationMatrixRet, nil
}

func isRotationMatrix(R []float64) bool {
	Rt := TransposeMatrix(R)
	shouldBeIdentity := dotProduct(Rt, R)
	I := identityMatrix(4)
	n := norm(subtractMatrices(I, shouldBeIdentity))
	return n < 1e-6
}

func eulerAnglesToRotationMatrix(theta []float64) []float64 {
	Rx := []float64{
		1, 0, 0, 0,
		0, math.Cos(theta[0]), -math.Sin(theta[0]), 0,
		0, math.Sin(theta[0]), math.Cos(theta[0]), 0,
		0, 0, 0, 1,
	}

	Ry := []float64{
		math.Cos(theta[1]), 0, math.Sin(theta[1]), 0,
		0, 1, 0, 0,
		-math.Sin(theta[1]), 0, math.Cos(theta[1]), 0,
		0, 0, 0, 1,
	}

	Rz := []float64{
		math.Cos(theta[2]), -math.Sin(theta[2]), 0, 0,
		math.Sin(theta[2]), math.Cos(theta[2]), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	R := dotProduct(Rz, dotProduct(Ry, Rx))
	return R
}

func rotationMatrixToEulerAngles(R []float64) []float64 {
	if !isRotationMatrix(R) {
		panic("Input matrix is not a valid rotation matrix")
	}

	sy := math.Sqrt(R[0]*R[0] + R[4]*R[4])
	singular := sy < 1e-6

	var x, y, z float64

	if !singular {
		x = math.Atan2(R[9], R[10])
		y = math.Atan2(-R[8], sy)
		z = math.Atan2(R[4], R[0])
	} else {
		x = math.Atan2(-R[6], R[5])
		y = math.Atan2(-R[8], sy)
		z = 0
	}

	return []float64{x, y, z}
}

func dotProduct(A []float64, B []float64) []float64 {
	result := make([]float64, 16)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := 0.0
			for k := 0; k < 4; k++ {
				sum += A[i*4+k] * B[k*4+j]
			}
			result[i*4+j] = sum
		}
	}
	return result
}

func subtractMatrices(A []float64, B []float64) []float64 {
	result := make([]float64, 16)
	for i := 0; i < 16; i++ {
		result[i] = A[i] - B[i]
	}
	return result
}

func norm(matrix []float64) float64 {
	sum := 0.0
	for i := 0; i < 16; i++ {
		sum += matrix[i] * matrix[i]
	}
	return math.Sqrt(sum)
}

func identityMatrix(size int) []float64 {
	matrix := make([]float64, 16)
	for i := 0; i < size; i++ {
		matrix[i*5] = 1.0
	}
	return matrix
}
