package main

import (
	"math/rand"
	"time"
	"./genalgo"
	"./logfile"
)


func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// on boucle jusqu'a  atteindre notre nombre de generation max ou jusqu'a ce qu'on trouve le maximum de notre courbe d'evolution
	for i := 0; i < genalgo.NB_GENERATION; i++ {
		genalgo.Evaluate()
		genalgo.PrintPopulation()
		selection := genalgo.Selection()
		genalgo.GeneratePopulation(selection)
		//parent1, parent2 := genalgo.SelectParent()
		//genalgo.CreateChild(parent1, parent2)
		// creation de la nouvelle population
	}
	logfile.End()
}


//
// PYTHON 2w1a.py in go
//
/*func main() {
	//b := []byte("127.0.0.1")
	id := C.simxStart((*C.simxChar)(unsafe.Pointer(&[]byte("127.0.0.1")[0])), 19997, 1, 1, 5000, 5)
	if id != -1 {
		fmt.Printf("Sucess")
		opmodewait := C.simxInt(0x010000)
		opmodesteaming := 0x020000
		opmodebuffer := 0x060000
		var wristHandle C.simxInt
		var elbowHandle C.simxInt
		var shoulderHandle C.simxInt
		var robotHandle C.simxInt


		C.simxGetObjectHandle(id, (*C.simxChar)(unsafe.Pointer(&[]byte("WristMotor")[0])), &wristHandle, opmodewait)
		C.simxGetObjectHandle(id, ((*C.simxChar)(unsafe.Pointer(&[]byte("ElbowMotor")[0]))), &elbowHandle, opmodewait)
		C.simxGetObjectHandle(id, ((*C.simxChar)(unsafe.Pointer(&[]byte("ShoulderMotor")[0]))), &shoulderHandle, opmodewait)
		C.simxGetObjectHandle(id, ((*C.simxChar)(unsafe.Pointer(&[]byte("2W1A")[0]))), &robotHandle, opmodewait)

		fmt.Printf("wrst : %0.5f\telb : %0.5f\tsld : %0.5f\trbt : %0.5f\n", wristHandle, elbowHandle, shoulderHandle, robotHandle)

		if true {
			rand.Seed(int64(time.Now().Nanosecond()))
			for i := 0; i < 100; i++ {

				// start simulation
				start := time.Now()
				C.simxStartSimulation(id, opmodewait)
				elapsed := time.Since(start)
				log.Printf("Binomial took %s", elapsed)


				// get RobotPosition / RobotOrientation
				var robotOrient [3]float32 = new(float32, 3)

				var robotOrient [3]float32 = new(float32, 3)
				C.simxGetObjectPosition(id, robotHandle, -1, (*C.simxFloat)(unsafe.Pointer(&(robotPos[0]))), C.simxInt(opmodesteaming))

				C.simxGetObjectPosition(id, robotHandle, -1, (*C.simxFloat)(unsafe.Pointer(&(robotPos[0]))), C.simxInt(opmodesteaming))
				C.simxGetObjectOrientation(id, robotHandle, -1, (*C.simxFloat)(unsafe.Pointer(&(robotOrient[0]))), C.simxInt(opmodesteaming))

				fmt.Printf("position (x = %0.5f, y = %0.5f)\nangle x : %0.5f\ty: %0.5f\tz = %0.5f\n", robotPos[0], robotPos[1],
					robotOrient[0], robotOrient[1], robotOrient[2])

				for j := 0; j < 5; j++ {
					awrist :=  float32(rand.Int31n(300)) * (float32)(math.Pi / 180.0)
					aelbow :=  float32(rand.Int31n(300)) * (float32)(math.Pi / 180.0)
					ashoulder := float32(rand.Int31n(300)) * (float32)(math.Pi / 180.0)

					fmt.Printf("awrst : %0.5f\taelb :%0.5f\tashdl : %0.5f\n", awrist, aelbow, ashoulder)

					C.simxSetJointTargetPosition(id, wristHandle, (C.simxFloat)(C.simxFloat(awrist)), C.simxInt(opmodewait))
					C.simxSetJointTargetPosition(id, elbowHandle, (C.simxFloat)(C.simxFloat(aelbow)), C.simxInt(opmodewait))
					C.simxSetJointTargetPosition(id, shoulderHandle, (C.simxFloat)(C.simxFloat(ashoulder)), C.simxInt(opmodewait))
//					time.Sleep(5 * time.Second)

					var pwrist [3]float32
					var pelbow [3]float32
					var pshoulder [3]float32
					C.simxGetJointPosition(id, wristHandle, (*C.simxFloat)(unsafe.Pointer(&pwrist[0])), C.simxInt(opmodewait))
					C.simxGetJointPosition(id, elbowHandle, (*C.simxFloat)(unsafe.Pointer(&pelbow[0])), C.simxInt(opmodewait))
					C.simxGetJointPosition(id, shoulderHandle, (*C.simxFloat)(unsafe.Pointer(&pshoulder[0])), C.simxInt(opmodewait))


					// Get the robot position after the movement sequence
					var robotPos [3]float32

					C.simxGetObjectPosition(id, robotHandle, -1, (*C.simxFloat)(unsafe.Pointer(&(robotPos[0]))), C.simxInt(opmodebuffer))
					C.simxGetObjectOrientation(id, robotHandle, -1, (*C.simxFloat)(unsafe.Pointer(&(robotOrient[0]))), C.simxInt(opmodebuffer))

					fmt.Printf("new position (x = %0.5f, y = %0.5f\tangle x : %0.5f\ty: %0.5f\tz = %0.5f\n", robotPos[0], robotPos[1],
						robotOrient[0], robotOrient[1], robotOrient[2])

				}
				C.simxStopSimulation(id, C.simxInt(opmodewait))
				time.Sleep(1 * time.Second)
			}

		}
	} else {
		fmt.Printf("Error")
	}

}*/
