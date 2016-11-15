package vgoapi

/*
#cgo CFLAGS: -I../CFiles
#cgo LDFLAGS: -L. -lremoteApi
#include "extApi.h"
*/
import "C"
import (
	"unsafe"
	"log"
	"time"
	"math"
)

var (
	// opmode for Vrep functions #see Vrep documentation
	opmodeonshot= C.simxInt(0x000000)
	opmodewait = C.simxInt(0x010000)
	opmodesteaming = C.simxInt(0x020000)
	opmodebuffer = C.simxInt(0x060000)

	// connect to Vrep API
	ClientID = C.simxStart(createSimxChar("127.0.0.1"), 19997, 1, 1, 5000, 5)

	// 2w1a motors and robot name
	WristMotor = "WristMotor"
	ElbowMotor = "ElbowMotor"
	ShoulderMotor = "ShoulderMotor"

	WristResp = "HandRespondable"
	ElbowResp = "ForearmRespondable"
	ShoulderResp = "ArmResponsable"

	r2W1A = "2W1A"
	LeftWheel = "LeftWheelJoint"
	RightWheel = "RightWheelJoint"
)

type API struct {
	wristHandle    C.simxInt
	elbowHandle    C.simxInt
	shoulderHandle C.simxInt

	wristRespHandle    C.simxInt
	elbowRespHandle    C.simxInt
	shoulderRespHandle C.simxInt

	robotHandle    C.simxInt

	leftWheelHandle      C.simxInt
	rightWheelHandle     C.simxInt

	robotOrient    [3]float32
	robotPos [3]float32

	wristPos [3]float32
	elbPos [3]float32
	shlPos [3]float32

	wristRespOr [3]float32
	elbRespOr [3]float32
	shlRespOr [3]float32

	rightWheel [3]float32
	leftWheel [3]float32

}

var manager  *API

func init() {
	manager = new(API)
	if ClientID == -1 {
		log.Print("error")
	}
	getRobotHandle()
}

func createSimxChar(src string) *C.simxChar {
	return (*C.simxChar)(unsafe.Pointer(&[]byte(src)[0]))
}

func createSimxFloat(src *[3]float32) *C.simxFloat {
	return (*C.simxFloat)(unsafe.Pointer(&(*src)[0]))
}

func createSimxInt(i *int) *C.simxInt {
	return (*C.simxInt)(unsafe.Pointer(i))
}

//set error
func getObjectHandle(name string, handle *C.simxInt) int {
	 C.simxGetObjectHandle(ClientID, createSimxChar(name), handle, opmodewait)
	return 0
}

func getRobotHandle() bool {
	e1 := getObjectHandle(WristMotor, &manager.wristHandle)
	e2 := getObjectHandle(ElbowMotor, &manager.elbowHandle)
	e3 := getObjectHandle(ShoulderMotor, &manager.shoulderHandle)

	getObjectHandle(LeftWheel, &manager.leftWheelHandle)
	getObjectHandle(RightWheel, &manager.rightWheelHandle)
	getObjectHandle(r2W1A, &manager.robotHandle)

	getObjectHandle(WristResp, &manager.wristRespHandle)
	getObjectHandle(ElbowResp, &manager.elbowRespHandle)
	getObjectHandle(ShoulderResp, &manager.shoulderRespHandle)

	if e1 + e2 + e3 == 0 {
		return true
	}
	return false
}

// On get les positions des wheels a la fin
func GetWheelsEndPosition() ([3]float32, [3]float32) {
	return manager.leftWheel, manager.rightWheel
}

func GetMotorsPositions() ([3]float32, [3]float32, [3]float32) {
	return manager.wristPos, manager.elbPos, manager.shlPos
}

func GetMotorsOrienation() ([3]float32, [3]float32, [3]float32) {
	return manager.wristRespOr, manager.elbRespOr, manager.shlRespOr
}

// On get les positions des wheels au debut
func GetWheelsStarPosition() ([3]float32, [3]float32) {
	return [3]float32 {0, 0, 0.04}, [3]float32 {0, 0, 0.04}
}

// Play le logiciel
func StartSimulation() {
	C.simxStartSimulation(ClientID, opmodewait)

}

// Stop le logiciel
func FinishSimulation() {
	C.simxStopSimulation(ClientID, (opmodewait))
}

func StartRobotMovement(newPos []float32) ([3]float32, [3]float32)  {
	// on recupere l'id du robot et des motors
	C.simxStartSimulation(ClientID, opmodewait)

	// on recupere l'orientation et la position du robot
	C.simxGetObjectPosition(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotPos), C.simxInt(opmodesteaming))
	C.simxGetObjectOrientation(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotOrient), C.simxInt(opmodesteaming))

	for x := 0; x < 3; x++ {
		var pwrist [3]float32
		var pelbow [3]float32
		var pshoulder [3]float32
		for i := 0; i < len(newPos) / 3; i++ {
			// init le mouvement de chaque moteur
			C.simxSetJointTargetPosition(ClientID, manager.wristHandle, (C.simxFloat(newPos[0 + i * 3]))  * (math.Pi / 180), opmodewait)
			C.simxSetJointTargetPosition(ClientID, manager.elbowHandle, (C.simxFloat(newPos[1 + i * 3]))  * (math.Pi / 180), opmodewait)
			C.simxSetJointTargetPosition(ClientID, manager.shoulderHandle, C.simxFloat(newPos[2 + i *3])  * (math.Pi / 180), opmodewait)

			// start chaque mouvement et on recupere la nouvelle position des moteurs (pour l'instant on s'en sert pas)

			C.simxGetJointPosition(ClientID, manager.wristHandle, createSimxFloat(&pwrist), (opmodewait))
			C.simxGetJointPosition(ClientID, manager.elbowHandle, createSimxFloat(&pelbow), (opmodewait))
			C.simxGetJointPosition(ClientID, manager.shoulderHandle, createSimxFloat(&pshoulder), (opmodewait))
		}
	}

	// on recupere la nouvelle position du robot et son orientation
	C.simxGetObjectPosition(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotPos), (opmodesteaming))
	C.simxGetObjectOrientation(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotOrient), (opmodesteaming))
	C.simxGetObjectPosition(ClientID, manager.leftWheelHandle, -1, createSimxFloat(&manager.leftWheel), opmodesteaming)
	C.simxGetObjectPosition(ClientID, manager.rightWheelHandle, -1, createSimxFloat(&manager.rightWheel), opmodesteaming)

	C.simxGetObjectPosition(ClientID, manager.wristHandle, -1, createSimxFloat(&manager.wristPos), opmodesteaming)
	C.simxGetObjectPosition(ClientID, manager.elbowHandle, -1, createSimxFloat(&manager.elbPos), opmodesteaming)
	C.simxGetObjectPosition(ClientID, manager.shoulderHandle, -1, createSimxFloat(&manager.shlPos), opmodesteaming)

	C.simxGetObjectOrientation(ClientID, manager.wristRespHandle, -1, createSimxFloat(&manager.wristRespOr), opmodesteaming)
	C.simxGetObjectOrientation(ClientID, manager.elbowRespHandle, -1, createSimxFloat(&manager.elbRespOr), opmodesteaming)
	C.simxGetObjectOrientation(ClientID, manager.shoulderHandle, -1, createSimxFloat(&manager.shlRespOr), opmodesteaming)


	C.simxStopSimulation(ClientID, (opmodewait))
	time.Sleep(200 * time.Millisecond)
	return manager.robotPos, manager.robotOrient
}

func MoveWhile(newPos []float32) {
	C.simxStartSimulation(ClientID, opmodewait)

	for i := 0;; i++{
		if i % 5 == 0 {
			C.simxStopSimulation(ClientID, (opmodewait))
			time.Sleep(200 * time.Millisecond)
			C.simxStartSimulation(ClientID, opmodewait)

		}
		var pwrist [3]float32
		var pelbow [3]float32
		var pshoulder [3]float32
		for i := 0; i < len(newPos) / 3; i++ {
			// init le mouvement de chaque moteur
			C.simxSetJointTargetPosition(ClientID, manager.wristHandle, (C.simxFloat(newPos[0 + i * 3])  * (math.Pi / 180)) , opmodewait)
			C.simxSetJointTargetPosition(ClientID, manager.elbowHandle, (C.simxFloat(newPos[1 + i * 3])  * (math.Pi / 180)) , opmodewait)
			C.simxSetJointTargetPosition(ClientID, manager.shoulderHandle, C.simxFloat(newPos[2 + i *3]  * (math.Pi / 180)) , opmodewait)

			// start chaque mouvement et on recupere la nouvelle position des moteurs (pour l'instant on s'en sert pas)

			C.simxGetJointPosition(ClientID, manager.wristHandle, createSimxFloat(&pwrist), (opmodewait))
			C.simxGetJointPosition(ClientID, manager.elbowHandle, createSimxFloat(&pelbow), (opmodewait))
			C.simxGetJointPosition(ClientID, manager.shoulderHandle, createSimxFloat(&pshoulder), (opmodewait))
		}
	}
	C.simxStopSimulation(ClientID, (opmodewait))

}


