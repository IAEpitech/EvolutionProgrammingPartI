package vgoapi

/*
#cgo CFLAGS: -I../CFiles
#cgo LDFLAGS: -L. -lremoteApi
#include "extApi.h"
*/
import "C"
import (
	"unsafe"
	"fmt"
	"time"
	"log"
)

var (
	// opmode for Vrep functions #see Vrep documentation
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
	ClientID       int
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
}

func createSimxChar(src string) *C.simxChar {
	return (*C.simxChar)(unsafe.Pointer(&[]byte(src)[0]))
}

func createSimxFloat(src [3]float32) *C.simxFloat {
	return (*C.simxFloat)(unsafe.Pointer(&src[0]))
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

func StartSimulation(newPos [3]float32) ([3]float32, [3]float32)  {
	C.simxStartSimulation(ClientID, opmodewait)

	C.simxGetObjectPosition(ClientID, manager.robotHandle, -1, createSimxFloat(manager.robotPos), C.simxInt(opmodesteaming))
	C.simxGetObjectOrientation(ClientID, manager.robotHandle, -1, createSimxFloat(manager.robotOrient), C.simxInt(opmodesteaming))

	fmt.Printf("OLD position (x = %0.5f, y = %0.5f)\nangle x : %0.5f\ty: %0.5f\tz = %0.5f\n", manager.robotPos[0], manager.robotPos[1],
		manager.robotOrient[0], manager.robotOrient[1], manager.robotOrient[2])

	C.simxSetJointTargetPosition(ClientID, manager.wristHandle, (C.simxFloat(newPos[0])), opmodewait)
	C.simxSetJointTargetPosition(ClientID, manager.elbowHandle, (C.simxFloat(newPos[1])), opmodewait)
	C.simxSetJointTargetPosition(ClientID, manager.shoulderHandle, C.simxFloat(newPos[2]), opmodewait)

	var pwrist [3]float32
	var pelbow [3]float32
	var pshoulder [3]float32
	C.simxGetJointPosition(ClientID, manager.wristHandle, createSimxFloat(pwrist), (opmodewait))
	C.simxGetJointPosition(ClientID, manager.elbowHandle, createSimxFloat(pelbow), (opmodewait))
	C.simxGetJointPosition(ClientID, manager.shoulderHandle, createSimxFloat(pshoulder), (opmodewait))
	C.simxStopSimulation(ClientID, (opmodewait))
	time.Sleep(1 * time.Second)

	C.simxGetObjectPosition(ClientID, manager.robotHandle, -1, createSimxFloat(manager.robotPos), (opmodebuffer))
	C.simxGetObjectOrientation(ClientID, manager.robotHandle, -1, createSimxFloat(manager.robotOrient), (opmodebuffer))

	fmt.Printf("NEW position (x = %0.5f, y = %0.5f)\nangle x : %0.5f\ty: %0.5f\tz = %0.5f\n", manager.robotPos[0], manager.robotPos[1],
		manager.robotOrient[0], manager.robotOrient[1], manager.robotOrient[2])
	return manager.robotPos, manager.robotOrient
}


