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
	r2W1A = "2W1A"


)

type API struct {
	wristHandle    C.simxInt
	elbowHandle    C.simxInt
	shoulderHandle C.simxInt
	robotHandle    C.simxInt
	robotOrient    [3]float32
	robotPos [3]float32
}

var manager  *API

func init() {
	manager = new(API)
	if ClientID == -1 {
		log.Print("error")
	}
	GetRobotHandle()

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

func GetRobotHandle() bool {
	e1 := getObjectHandle(WristMotor, &manager.wristHandle)
	e2 := getObjectHandle(ElbowMotor, &manager.elbowHandle)
	e3 := getObjectHandle(ShoulderMotor, &manager.shoulderHandle)
	getObjectHandle(r2W1A, &manager.robotHandle)
	if e1 + e2 + e3 == 0 {
		return true
	}
	return false
}

// Play le logiciel
func StartSimulation() {
	C.simxStartSimulation(ClientID, opmodewait)

}

// Stop le logiciel
func FinishSimulation() {
	C.simxStopSimulation(ClientID, (opmodewait))
}

func StartRobotMovement(newPos [9]float32) ([3]float32, [3]float32)  {
	// on recupere l'id du robot et des motors
	C.simxStartSimulation(ClientID, opmodewait)

	// on recupere l'orientation et la position du robot
	C.simxGetObjectPosition(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotPos), C.simxInt(opmodesteaming))
	C.simxGetObjectOrientation(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotOrient), C.simxInt(opmodesteaming))

	var pwrist [3]float32
	var pelbow [3]float32
	var pshoulder [3]float32
	for i := 0; i < 3; i++ {
		// init le mouvement de chaque moteur
		C.simxSetJointTargetPosition(ClientID, manager.wristHandle, (C.simxFloat(newPos[0 + i * 3])), opmodewait)
		C.simxSetJointTargetPosition(ClientID, manager.elbowHandle, (C.simxFloat(newPos[1 + i * 3])), opmodewait)
		C.simxSetJointTargetPosition(ClientID, manager.shoulderHandle, C.simxFloat(newPos[2 + i *3]), opmodewait)

		// start chaque mouvement et on recupere la nouvelle position des moteurs (pour l'instant on s'en sert pas)

		C.simxGetJointPosition(ClientID, manager.wristHandle, createSimxFloat(&pwrist), (opmodewait))
		C.simxGetJointPosition(ClientID, manager.elbowHandle, createSimxFloat(&pelbow), (opmodewait))
		C.simxGetJointPosition(ClientID, manager.shoulderHandle, createSimxFloat(&pshoulder), (opmodewait))
	}


	// on recupere la nouvelle position du robot et son orientation
	C.simxGetObjectPosition(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotPos), (opmodebuffer))
	C.simxGetObjectOrientation(ClientID, manager.robotHandle, -1, createSimxFloat(&manager.robotOrient), (opmodebuffer))
	C.simxStopSimulation(ClientID, (opmodewait))
	time.Sleep(200 * time.Millisecond)
	return manager.robotPos, manager.robotOrient
}


